package server

import (
	"encoding/json"
	"net/http"
	"slices"

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
		return nil, false, nil
	}

	query := dataQueries[0]
	table := query.Table
	if table == "" {
		table = "report"
	}

	// Reuse the standard query pipeline to fetch the flat list of distinct machine names.
	raw := bytebufferpool.Get()
	defer bytebufferpool.Put(raw)
	writer := quicktemplate.AcquireWriter(raw)
	err = data_query.SelectRows(request.Context(), query, table, t, writer.N())
	quicktemplate.ReleaseWriter(writer)
	if err != nil {
		return nil, false, err
	}

	machines := []string{}
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
		response = append(response, machineGroupResponseItem{Group: group, Machines: members})
	}
	slices.SortFunc(response, func(a, b machineGroupResponseItem) int {
		if a.Group < b.Group {
			return -1
		}
		if a.Group > b.Group {
			return 1
		}
		return 0
	})

	jsonData, err := json.Marshal(response)
	if err != nil {
		return nil, false, err
	}
	buffer := bytebufferpool.Get()
	_, _ = buffer.Write(jsonData)
	return buffer, true, nil
}
