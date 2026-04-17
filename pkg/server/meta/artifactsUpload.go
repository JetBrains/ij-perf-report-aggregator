package meta

import (
	"context"
	"encoding/json"
	"errors"
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

func uploadToAll(ctx context.Context, uploaders []ArtifactsUploader, file []byte, fileName string, errCh chan<- error, uploadWg *sync.WaitGroup) {
	for _, uploader := range uploaders {
		uploadWg.Go(func() {
			err := uploader.Upload(ctx, file, fileName)
			if err != nil {
				if uploader.Type() == YouTrack {
					slog.Error("Failed to upload attachment to youtrack", "error", err)
					errCh <- err
				} else {
					slog.Error("Failed to upload attachment", "target", uploader.Type(), "error", err)
				}
			}
		})
	}
}

func CreatePostUploadAttachments() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		type Exceptions struct {
			Exceptions []string `json:"exceptions"`
		}

		var exceptions Exceptions
		body := request.Body
		all, err := io.ReadAll(body)
		if err != nil {
			handleError(writer, "cannot read body", err, &exceptions.Exceptions)
			_ = marshalAndWriteIssueResponse(writer, exceptions)
			return
		}

		defer body.Close()

		var params UploadAttachmentsRequest
		if err = json.Unmarshal(all, &params); err != nil {
			handleError(writer, "cannot unmarshal parameters", err, &exceptions.Exceptions)
			_ = marshalAndWriteIssueResponse(writer, exceptions)
			return
		}

		uploaders, err := params.ToUploaders()
		if err != nil {
			handleError(writer, "invalid upload targets", err, &exceptions.Exceptions)
			_ = marshalAndWriteIssueResponse(writer, exceptions)
			return
		}

		errCh := make(chan error, 10)
		var uploadWg sync.WaitGroup

		builds := []int{params.TeamCityAttachmentInfo.CurrentBuildId}
		if params.TeamCityAttachmentInfo.PreviousBuildId != nil {
			builds = append(builds, *params.TeamCityAttachmentInfo.PreviousBuildId)
		}

		if params.ChartPng != nil {
			uploadToAll(request.Context(), uploaders, *params.ChartPng, "dashboard.png", errCh, &uploadWg)
		}

		collector := getArtifactCollector(params.TestType)

		if collector != nil {
			var wg sync.WaitGroup
			for index, buildId := range builds {
				wg.Go(func() {
					testArtifactPath := collector.getArtifactsPath(params)

					children, err := teamCityClient.getArtifactChildren(request.Context(), buildId, testArtifactPath)
					if err != nil {
						slog.Error("Failed to get teamcity artifact children", "error", err)
						errCh <- err
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
								slog.Error("Failed to download artefacts form teamcity", "error", err)
								errCh <- err
								return
							}

							attachmentName := getAttachmentName(str, attachmentPostfix)
							uploadToAll(request.Context(), uploaders, artifact, attachmentName, errCh, &uploadWg)
						})
					}
					childWg.Wait()
				})
			}
			wg.Wait()
		}

		uploadWg.Wait()
		close(errCh)

		for err := range errCh {
			if err != nil {
				exceptions.Exceptions = append(exceptions.Exceptions, err.Error())
			}
		}

		if len(exceptions.Exceptions) > 0 {
			writer.WriteHeader(http.StatusInternalServerError)
			_ = marshalAndWriteIssueResponse(writer, exceptions)
			return
		}

		_ = marshalAndWriteIssueResponse(writer, exceptions)
	}
}
