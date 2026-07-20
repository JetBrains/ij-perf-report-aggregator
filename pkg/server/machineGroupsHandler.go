package server

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"

	data_query "github.com/JetBrains/ij-perf-report-aggregator/pkg/data-query"
	"github.com/JetBrains/ij-perf-report-aggregator/pkg/machine"
	"github.com/valyala/bytebufferpool"
	"github.com/valyala/quicktemplate"
)

const machineGroupsPathPrefix = "/api/machineGroups/"

// handleMachineGroupLookup resolves a single raw machine name to its hardware-class group. It
// is a pure lookup (no DB), used by the frontend to map a drilldown's raw agent name to the
// group to preselect — the target page's own (differently filtered) machine list may not
// contain that exact ephemeral instance.
func (t *StatsServer) handleMachineGroupLookup(w http.ResponseWriter, request *http.Request) {
	name := request.URL.Query().Get("machine")
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"group": machine.GroupName(name)}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type machineGroupResponseItem struct {
	Group    string   `json:"group"`
	Machines []string `json:"machines"`
	// Predicate is the /api/q filter suffix ({f: "machine", q: <predicate>}) selecting exactly
	// this hardware class — rendered from the grouping rule itself, so the frontend never has
	// to infer one from the members. Empty for groups without a prefix/name rule (e.g.
	// "Unknown"); callers then filter by the member list.
	Predicate string `json:"predicate,omitempty"`
}

// handleMachineGroups runs the distinct-machine query the caller supplies (same DataQuery
// payload as /api/q, so it honours the same table + filters), then returns the agents already
// bucketed into hardware-class groups. Grouping lives only here on the backend (AT-4930), so
// the frontend never sees raw agents ungrouped.
func (t *StatsServer) handleMachineGroups(request *http.Request) (*bytebufferpool.ByteBuffer, bool, error) {
	dataQueries, _, err := data_query.ReadQueryV2WithPrefix(request, machineGroupsPathPrefix)
	if err != nil {
		return nil, false, err
	}
	if len(dataQueries) == 0 {
		// The cache manager dereferences the returned buffer unconditionally, so it must never
		// be nil on success — mirror /api/q and answer an empty query list with an empty result.
		buffer, err := toJSONBuffer([]machineGroupResponseItem{})
		return buffer, true, err
	}

	// Reuse the standard query pipeline to fetch the flat list of distinct machine names.
	raw := bytebufferpool.Get()
	defer bytebufferpool.Put(raw)
	writer := quicktemplate.AcquireWriter(raw)
	err = t.computeMeasureResponse(request.Context(), dataQueries[0], writer.N())
	quicktemplate.ReleaseWriter(writer)
	if err != nil {
		return nil, false, err
	}

	var machines []string
	if len(raw.B) != 0 {
		if err := json.Unmarshal(raw.B, &machines); err != nil {
			return nil, false, err
		}
	}

	byGroup := make(map[string][]string)
	for _, m := range machines {
		group := machine.GroupName(m)
		byGroup[group] = append(byGroup[group], m)
	}

	response := make([]machineGroupResponseItem, 0, len(byGroup))
	for group, members := range byGroup {
		slices.Sort(members)
		response = append(response, machineGroupResponseItem{Group: group, Machines: members, Predicate: machine.GroupPredicate("machine", group)})
	}
	slices.SortFunc(response, func(a, b machineGroupResponseItem) int {
		return strings.Compare(a.Group, b.Group)
	})

	buffer, err := toJSONBuffer(response)
	return buffer, true, err
}
