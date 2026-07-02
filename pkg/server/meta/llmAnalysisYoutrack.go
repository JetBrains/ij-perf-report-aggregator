package meta

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateIssueByAnalysisRequest struct {
	ProjectId   string `json:"projectId"`
	TicketLabel string `json:"ticketLabel"`
	ChangesLink string `json:"changesLink"`
	Delta       string `json:"delta"`
	ChartPng    []byte `json:"chartPng,omitempty"`
}

func CreatePostCreateIssueByAnalysis(metaDb *pgxpool.Pool) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		response := CreateIssueResponse{}

		id, err := strconv.Atoi(chi.URLParam(request, "id"))
		if err != nil || id <= 0 {
			http.Error(writer, "invalid id", http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(request.Body)
		if err != nil {
			handleError(writer, "cannot read body", err, &response.Exceptions)
			_ = marshalAndWriteIssueResponse(writer, &response)
			return
		}
		defer request.Body.Close()

		var params CreateIssueByAnalysisRequest
		if err = json.Unmarshal(body, &params); err != nil {
			handleError(writer, "cannot unmarshal parameters", err, &response.Exceptions)
			_ = marshalAndWriteIssueResponse(writer, &response)
			return
		}

		details, err := getAnalysisById(request.Context(), metaDb, id)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				http.Error(writer, "analysis not found", http.StatusNotFound)
				return
			}
			handleError(writer, "cannot get analysis", err, &response.Exceptions)
			_ = marshalAndWriteIssueResponse(writer, &response)
			return
		}

		if details.State != string(LlmAnalysisStateSuccess) || details.LlmComment == nil || *details.LlmComment == "" {
			http.Error(writer, "analysis is not in success state or has no comment", http.StatusBadRequest)
			return
		}

		descriptionData := buildAnalysisDescriptionData(request.Context(), metaDb, details, params, &response.Exceptions)

		issue, err := createYoutrackIssueCommon(request.Context(), request, CommonIssueParams{
			ProjectId:   params.ProjectId,
			Summary:     params.TicketLabel,
			Description: generateDescription(descriptionData),
			ExtraTags:   []Tag{analysedByIjPerfTag},
		}, &response.Exceptions)
		if err != nil {
			handleError(writer, "failed to create issue", err, &response.Exceptions)
			_ = marshalAndWriteIssueResponse(writer, &response)
			return
		}

		response.Issue = *issue

		if err := updateLlmAnalysisRun(request.Context(), metaDb, id, LlmAnalysisRunPatch{YtIssueId: &issue.IDReadable}); err != nil {
			slog.Error("failed to persist yt_issue_id on analysis", "error", err, "id", id, "issue", issue.IDReadable)
			logError("failed to link issue to analysis", err, &response.Exceptions)
		}

		if len(params.ChartPng) > 0 {
			if err := youtrackClient.UploadAttachment(request.Context(), issue.ID, bytes.NewReader(params.ChartPng), "dashboard.png", int64(len(params.ChartPng))); err != nil {
				slog.Error("failed to upload chart PNG", "error", err)
				logError("failed to upload chart PNG", err, &response.Exceptions)
			}
		}

		if err := marshalAndWriteIssueResponse(writer, &response); err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

type LinkIssueByAnalysisRequest struct {
	IssueId string `json:"issueId"`
}

func CreatePostLinkIssueByAnalysis(metaDb *pgxpool.Pool) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		response := CreateIssueResponse{}

		id, err := strconv.Atoi(chi.URLParam(request, "id"))
		if err != nil || id <= 0 {
			http.Error(writer, "invalid id", http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(request.Body)
		if err != nil {
			handleError(writer, "cannot read body", err, &response.Exceptions)
			_ = marshalAndWriteIssueResponse(writer, &response)
			return
		}
		defer request.Body.Close()

		var params LinkIssueByAnalysisRequest
		if err = json.Unmarshal(body, &params); err != nil {
			handleError(writer, "cannot unmarshal parameters", err, &response.Exceptions)
			_ = marshalAndWriteIssueResponse(writer, &response)
			return
		}

		params.IssueId = strings.TrimSpace(params.IssueId)
		if params.IssueId == "" {
			http.Error(writer, "issueId is required", http.StatusBadRequest)
			return
		}

		if _, err := getAnalysisById(request.Context(), metaDb, id); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				http.Error(writer, "analysis not found", http.StatusNotFound)
				return
			}
			handleError(writer, "cannot get analysis", err, &response.Exceptions)
			_ = marshalAndWriteIssueResponse(writer, &response)
			return
		}

		issue, err := youtrackClient.ResolveIssue(request.Context(), params.IssueId)
		if err != nil {
			handleError(writer, "cannot resolve issue", err, &response.Exceptions)
			_ = marshalAndWriteIssueResponse(writer, &response)
			return
		}
		response.Issue = *issue

		if err := youtrackClient.AddTag(request.Context(), issue.IDReadable, analysedByIjPerfTag); err != nil {
			slog.Error("failed to tag linked issue", "error", err, "issue", issue.IDReadable)
			logError("failed to tag issue", err, &response.Exceptions)
		}

		if err := updateLlmAnalysisRun(request.Context(), metaDb, id, LlmAnalysisRunPatch{YtIssueId: &issue.IDReadable}); err != nil {
			slog.Error("failed to persist yt_issue_id on analysis", "error", err, "id", id, "issue", issue.IDReadable)
			logError("failed to link issue to analysis", err, &response.Exceptions)
		}

		if err := marshalAndWriteIssueResponse(writer, &response); err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

// buildAnalysisDescriptionData maps a stored LLM analysis row + space-uploaded artifacts + the
// FE-provided changesLink into the shared GenerateDescriptionData shape. TestType is left empty so
// the legacy logs/snapshots block is suppressed; AnalysisAttachments replaces it.
func buildAnalysisDescriptionData(ctx context.Context, metaDb *pgxpool.Pool, details *LlmAnalysisDetails, params CreateIssueByAnalysisRequest, exceptions *[]string) GenerateDescriptionData {
	attach := collectAnalysisAttachments(ctx, metaDb, details, exceptions)

	var dashLink string
	if details.DashboardLink != nil && *details.DashboardLink != "" {
		dashLink = buildDashboardLink(*details.DashboardLink, details.CurrentBuildId, details.Id)
	}

	return GenerateDescriptionData{
		Kind:                "",
		AffectedTest:        details.Project,
		AffectedMetric:      details.Metric,
		Delta:               params.Delta,
		CurrentValue:        derefString(details.CurrentValue),
		PreviousValue:       derefString(details.PreviousValue),
		BuildLink:           "https://buildserver.labs.intellij.net/viewLog.html?buildId=" + details.CurrentBuildId,
		Changes:             params.ChangesLink,
		DashboardLink:       dashLink,
		TestMethod:          details.TestMethodName,
		AnalysisAttachments: attach,
		AnalysisResult:      details.LlmComment,
	}
}

func collectAnalysisAttachments(ctx context.Context, metaDb *pgxpool.Pool, details *LlmAnalysisDetails, exceptions *[]string) *AnalysisAttachments {
	attach := &AnalysisAttachments{}
	if buildId, err := strconv.Atoi(details.CurrentBuildId); err == nil {
		if upload, err := getSpaceUploads(ctx, metaDb, buildId, details.Project); err == nil {
			attach.Current = upload.Files
			attach.CurrentURL = buildSpaceAnalysisFolderURL(buildId, details.Project)
		} else if !errors.Is(err, pgx.ErrNoRows) {
			logError("cannot fetch current space uploads", err, exceptions)
		}
	}
	if buildId, err := strconv.Atoi(details.PrevBuildId); err == nil {
		if upload, err := getSpaceUploads(ctx, metaDb, buildId, details.Project); err == nil {
			attach.Previous = upload.Files
			attach.PreviousURL = buildSpaceAnalysisFolderURL(buildId, details.Project)
		} else if !errors.Is(err, pgx.ErrNoRows) {
			logError("cannot fetch previous space uploads", err, exceptions)
		}
	}
	if len(attach.Current) == 0 && len(attach.Previous) == 0 {
		return nil
	}
	return attach
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
