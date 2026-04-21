package meta

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
)

type UploadAttachmentsRequest struct {
	Targets                []UploadTarget         `json:"targets"`
	IssueId                string                 `json:"issueId"`
	TeamCityAttachmentInfo TeamCityAttachmentInfo `json:"teamcityAttachmentInfo"`
	AffectedTest           string                 `json:"affectedTest"`
	ChartPng               *[]byte                `json:"chartPng"`
	TestType               string                 `json:"testType"`
}

type UploadTarget string

const (
	YouTrack UploadTarget = "youtrack"
	Space    UploadTarget = "space"
)

type UploadAttachmentsResponse struct {
	Uploads    map[UploadTarget][]string `json:"uploads"`
	Exceptions []string                  `json:"exceptions"`
}

type uploadResults struct {
	mu      sync.Mutex
	uploads map[UploadTarget][]string
}

func newUploadResults() *uploadResults {
	return &uploadResults{
		uploads: make(map[UploadTarget][]string),
	}
}

func (ur *uploadResults) addSuccess(target UploadTarget, fileName string) {
	ur.mu.Lock()
	defer ur.mu.Unlock()
	ur.uploads[target] = append(ur.uploads[target], fileName)
}

func (request *UploadAttachmentsRequest) ToUploaders() ([]ArtifactsUploader, error) {
	uploaders := make([]ArtifactsUploader, 0, len(request.Targets))

	for _, target := range request.Targets {
		switch target {
		case YouTrack:
			uploaders = append(uploaders, &youtrackUploader{
				client:  youtrackClient,
				issueId: request.IssueId,
			})
		case Space:
			uploaders = append(uploaders, &spaceUploader{
				client:  spacePackagesClient,
				issueId: request.IssueId,
			})
		default:
			slog.Warn("Unknown upload target, skipping", "target", target)
		}
	}

	if len(uploaders) == 0 {
		if len(request.Targets) == 0 {
			return nil, errors.New("no upload targets specified")
		}
		return nil, errors.New("no valid upload targets found")
	}

	return uploaders, nil
}

type ArtifactsUploader interface {
	Upload(ctx context.Context, file []byte, fileName string) error
	Type() UploadTarget
}

type youtrackUploader struct {
	client  *YoutrackClient
	issueId string
}

func (u *youtrackUploader) Upload(ctx context.Context, file []byte, fileName string) error {
	return u.client.UploadAttachment(ctx, u.issueId, file, fileName)
}

func (u *youtrackUploader) Type() UploadTarget {
	return YouTrack
}

type spaceUploader struct {
	client  *SpacePackagesClient
	issueId string
}

func (u *spaceUploader) Upload(ctx context.Context, file []byte, fileName string) error {
	return u.client.UploadFile(ctx, "platform-test-automation", "performance-regression-llm-analysis", "analyses/"+u.issueId, fileName, file)
}

func (u *spaceUploader) Type() UploadTarget {
	return Space
}

func uploadToAll(ctx context.Context, uploaders []ArtifactsUploader, file []byte, fileName string, handleUploadErrors func(string, error, UploadTarget), uploadWg *sync.WaitGroup, results *uploadResults) {
	for _, uploader := range uploaders {
		uploadWg.Go(func() {
			err := uploader.Upload(ctx, file, fileName)
			if err != nil {
				handleUploadErrors("Failed to upload attachment", err, uploader.Type())
			} else {
				results.addSuccess(uploader.Type(), fileName)
			}
		})
	}
}

func CreatePostUploadAttachments() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var response UploadAttachmentsResponse
		body := request.Body
		all, err := io.ReadAll(body)
		if err != nil {
			handleError(writer, "cannot read body", err, &response.Exceptions)
			_ = marshalAndWriteIssueResponse(writer, response)
			return
		}

		defer body.Close()

		var params UploadAttachmentsRequest
		if err = json.Unmarshal(all, &params); err != nil {
			handleError(writer, "cannot unmarshal parameters", err, &response.Exceptions)
			_ = marshalAndWriteIssueResponse(writer, response)
			return
		}

		uploaders, err := params.ToUploaders()
		if err != nil {
			handleError(writer, "invalid upload targets", err, &response.Exceptions)
			_ = marshalAndWriteIssueResponse(writer, response)
			return
		}

		var exceptionsMu sync.Mutex
		logAndAddException := func(message string, err error) {
			slog.Error(message, "error", err)
			exceptionsMu.Lock()
			defer exceptionsMu.Unlock()
			response.Exceptions = append(response.Exceptions,
				fmt.Sprintf("Message: %s. Error: %s", message, err.Error()))
		}
		handleUploadError := func(message string, err error, target UploadTarget) {
			if target == YouTrack {
				logAndAddException(message, err)
			} else {
				slog.Error(message, "target", target, "error", err)
			}
		}

		var uploadWg sync.WaitGroup
		results := newUploadResults()

		builds := []int{params.TeamCityAttachmentInfo.CurrentBuildId}
		if params.TeamCityAttachmentInfo.PreviousBuildId != nil {
			builds = append(builds, *params.TeamCityAttachmentInfo.PreviousBuildId)
		}

		if params.ChartPng != nil {
			uploadToAll(request.Context(), uploaders, *params.ChartPng, "dashboard.png", handleUploadError, &uploadWg, results)
		}

		collector := getArtifactCollector(params.TestType)

		if collector != nil {
			var wg sync.WaitGroup
			for index, buildId := range builds {
				wg.Go(func() {
					testArtifactPath := collector.getArtifactsPath(params)

					children, err := teamCityClient.getArtifactChildren(request.Context(), buildId, testArtifactPath)
					if err != nil {
						logAndAddException("Failed to get teamcity artifact children", err)
						return
					}

					var filteredChildren []string

					for _, str := range children {
						if collector.checkArtifact(str) {
							filteredChildren = append(filteredChildren, str)
						}
					}

					var attachmentPostfix string
					if index == 0 {
						attachmentPostfix = "current"
					} else {
						attachmentPostfix = "before"
					}

					var childWg sync.WaitGroup
					for _, str := range filteredChildren {
						childWg.Go(func() {
							artifact, err := teamCityClient.downloadArtifact(request.Context(), buildId, testArtifactPath+"/"+str)
							if err != nil {
								logAndAddException("Failed to download artifacts from teamcity", err)
								return
							}

							attachmentName := getAttachmentName(str, attachmentPostfix)
							uploadToAll(request.Context(), uploaders, artifact, attachmentName, handleUploadError, &uploadWg, results)
						})
					}
					childWg.Wait()
				})
			}
			wg.Wait()
		}

		uploadWg.Wait()

		response.Uploads = results.uploads

		if len(response.Exceptions) > 0 {
			writer.WriteHeader(http.StatusInternalServerError)
			_ = marshalAndWriteIssueResponse(writer, response)
			return
		}

		_ = marshalAndWriteIssueResponse(writer, response)
	}
}
