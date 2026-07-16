package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleMachineGroupLookup(t *testing.T) {
	t.Parallel()
	cases := map[string]string{
		"intellij-linux-performance-aws-i-08aec6c8ee5a71bba":        "Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)",
		"intellij-linux-performance-tiny-aws-on-demand-i-0abc12345": "Linux EC2 C6id.xlarge (4 vCPU Xeon, 8 GB)",
		"intellij-macos-hw-unit-1550":                               "macMini 2018", // migrated legacy agent
		"some-brand-new-agent-i-0ffff":                              "Unknown",
	}

	server := &StatsServer{}
	for machineName, want := range cases {
		req := httptest.NewRequest(http.MethodGet, "/api/machineGroup?machine="+machineName, http.NoBody)
		rec := httptest.NewRecorder()
		server.handleMachineGroupLookup(rec, req)

		require.Equalf(t, http.StatusOK, rec.Code, "status for %q", machineName)
		var body map[string]string
		require.NoErrorf(t, json.Unmarshal(rec.Body.Bytes(), &body), "decode for %q", machineName)
		assert.Equalf(t, want, body["group"], "group for %q", machineName)
	}
}
