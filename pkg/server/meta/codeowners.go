package meta

import (
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strings"
	"time"
)

// CodeOwners resolves a code-owner group name (or one of its aliases) to the Slack channel
// configured for that group, via the CodeOwners service REST API. It is used to route
// performance-test notifications by owner instead of by a hardcoded team list.

const codeOwnersDefaultBaseURL = "https://codeowners.labs.jb.gg"

// slackChannelAttributeKeys are the group attributes that may hold the notification channel,
// in preference order: the alerts-specific channel wins over the general team channel.
var slackChannelAttributeKeys = []string{"Slack Alerts Channel", "Slack Channel"}

type CodeOwnersClient struct {
	baseURL string
	token   string
	http    *http.Client
}

func NewCodeOwnersClient() *CodeOwnersClient {
	return &CodeOwnersClient{
		baseURL: cmp.Or(os.Getenv("CODEOWNERS_URL"), codeOwnersDefaultBaseURL),
		// Reuse the backend's existing Space token; allow a dedicated one to override.
		token: cmp.Or(os.Getenv("CODEOWNERS_TOKEN"), os.Getenv("SPACE_TOKEN")),
		http:  &http.Client{Timeout: 60 * time.Second},
	}
}

var codeOwnersClient = NewCodeOwnersClient()

// codeOwners REST DTOs (subset of GroupListResponseDto / GroupListItemDto).
type coGroupsResponse struct {
	Items []coGroupItem `json:"items"`
}

type coGroupItem struct {
	ID         string        `json:"id"`
	ParentID   string        `json:"parentId"`
	Name       string        `json:"name"`
	Aliases    []coAlias     `json:"aliases"`
	Attributes []coAttribute `json:"attributes"`
}

type coAlias struct {
	Name string `json:"name"`
}

type coAttribute struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// FetchOwnerChannels returns a map from every active group name AND alias to its Slack
// channel. Aliases (previous names) map to the same channel as the canonical name, so a
// lookup by either resolves identically. Groups without a slack-channel attribute are omitted.
//
// When a name is both an active group's canonical name and another group's stale alias
// (a previous name reused after a rename), the canonical name wins regardless of the order
// the service returns groups, so routing is deterministic. Genuine conflicts (two groups
// sharing a canonical name or an alias with different channels) are logged.
//
// limit=-1 (REQUEST_ALL_LIMIT) tells the CodeOwners service to return every group in a single
// response, so no pagination is needed.
func (c *CodeOwnersClient) FetchOwnerChannels(ctx context.Context) (map[string]string, error) {
	params := url.Values{}
	params.Set("fields", "aliases,attributes")
	params.Set("limit", "-1")
	items, err := c.fetchGroups(ctx, params)
	if err != nil {
		return nil, err
	}

	// Collect canonical names and aliases separately, then merge with canonical names on top,
	// so an active group's name always wins over a stale alias another group may still carry
	// (a previous name reused after a rename) — independent of the service's response order.
	aliasChannels := make(map[string]string)
	canonicalChannels := make(map[string]string)
	for _, item := range items {
		channel := slackChannelOf(item.Attributes)
		if channel == "" {
			continue
		}
		if existing, ok := canonicalChannels[item.Name]; ok && existing != channel {
			slog.Warn("duplicate code-owner group name resolves to different channels; keeping last",
				"name", item.Name, "previous", existing, "channel", channel)
		}
		canonicalChannels[item.Name] = channel
		for _, alias := range item.Aliases {
			if existing, ok := aliasChannels[alias.Name]; ok && existing != channel {
				slog.Warn("code-owner alias claimed by multiple groups with different channels; keeping last",
					"alias", alias.Name, "previous", existing, "channel", channel)
			}
			aliasChannels[alias.Name] = channel
		}
	}

	channels := make(map[string]string, len(aliasChannels)+len(canonicalChannels))
	maps.Copy(channels, aliasChannels)
	maps.Copy(channels, canonicalChannels)
	return channels, nil
}

// CodeOwnerGroup is a group of the CodeOwners service with its place in the hierarchy.
type CodeOwnerGroup struct {
	Name string `json:"name"`
	// Aliases are the previous names of the group. A test can still have one of them as its owner.
	Aliases []string `json:"aliases"`
	// Parent is the name of the group above this one, empty for a group at the top.
	Parent string `json:"parent"`
}

// FetchGroups returns every group with its aliases and its parent, sorted by name. The parent is given by name,
// so the caller needs no group ids. `recursive` makes the service list the subgroups as well, and limit=-1
// returns all groups in one response.
func (c *CodeOwnersClient) FetchGroups(ctx context.Context) ([]CodeOwnerGroup, error) {
	params := url.Values{}
	params.Set("fields", "aliases")
	params.Set("recursive", "true")
	params.Set("limit", "-1")
	items, err := c.fetchGroups(ctx, params)
	if err != nil {
		return nil, err
	}

	nameByID := make(map[string]string, len(items))
	for _, item := range items {
		if item.ID != "" {
			nameByID[item.ID] = item.Name
		}
	}
	groups := make([]CodeOwnerGroup, 0, len(items))
	for _, item := range items {
		aliases := make([]string, 0, len(item.Aliases))
		for _, alias := range item.Aliases {
			aliases = append(aliases, alias.Name)
		}
		// A parent that the listing does not contain is the same as no parent for the caller.
		groups = append(groups, CodeOwnerGroup{Name: item.Name, Aliases: aliases, Parent: nameByID[item.ParentID]})
	}
	slices.SortFunc(groups, func(a, b CodeOwnerGroup) int { return cmp.Compare(a.Name, b.Name) })
	return groups, nil
}

// fetchGroups lists the groups of the CodeOwners service with the given query parameters.
func (c *CodeOwnersClient) fetchGroups(ctx context.Context, params url.Values) ([]coGroupItem, error) {
	if c.token == "" {
		return nil, errors.New("no code-owners token configured (set CODEOWNERS_TOKEN or SPACE_TOKEN)")
	}

	requestURL := c.baseURL + "/app/rest/v1/groups?" + params.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/json")

	httpResp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer httpResp.Body.Close()
	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("code-owners request failed: %s", httpResp.Status)
	}

	var resp coGroupsResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to decode code-owners response: %w", err)
	}
	return resp.Items, nil
}

// slackChannelOf returns the notification channel for a group, picking the first populated
// attribute in slackChannelAttributeKeys preference order. The value is normalized to a bare
// channel name (without a leading '#') to match the channel format used elsewhere in routing.
func slackChannelOf(attributes []coAttribute) string {
	byKey := make(map[string]string, len(attributes))
	for _, attr := range attributes {
		byKey[attr.Key] = strings.TrimPrefix(strings.TrimSpace(attr.Value), "#")
	}
	for _, key := range slackChannelAttributeKeys {
		if v := byKey[key]; v != "" {
			return v
		}
	}
	return ""
}

// CreateGetCodeOwnerChannelsHandler exposes the owner(+alias) -> Slack channel map.
// The degradation-detector fetches this once per run to route notifications by owner.
func CreateGetCodeOwnerChannelsHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		channels, err := codeOwnersClient.FetchOwnerChannels(request.Context())
		if err != nil {
			slog.Error("unable to fetch code-owner channels", "error", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		jsonBytes, err := json.Marshal(channels)
		if err != nil {
			slog.Error("unable to marshal code-owner channels", "error", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.Header().Set("Content-Type", "application/json")
		if _, err := writer.Write(jsonBytes); err != nil {
			slog.Error("unable to write response", "error", err)
		}
	}
}

// CreateGetCodeOwnerGroupsHandler exposes the groups of the CodeOwners service with their aliases and parents.
// The degradation-detector fetches this once per run to match a test whose owner is a subgroup, or has an old
// name, of a group that a subscriber claims.
func CreateGetCodeOwnerGroupsHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		groups, err := codeOwnersClient.FetchGroups(request.Context())
		if err != nil {
			slog.Error("unable to fetch code-owner groups", "error", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		jsonBytes, err := json.Marshal(groups)
		if err != nil {
			slog.Error("unable to marshal code-owner groups", "error", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.Header().Set("Content-Type", "application/json")
		if _, err := writer.Write(jsonBytes); err != nil {
			slog.Error("unable to write response", "error", err)
		}
	}
}
