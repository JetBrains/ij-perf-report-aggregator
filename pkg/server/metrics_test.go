package server

import (
	"fmt"
	dataQuery "github.com/JetBrains/ij-perf-report-aggregator/pkg/data-query"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestClickhouse(_ *testing.T) {
	testUrl, _ := url.Parse("/api/v1/load/(db:perfintDev,table:kotlin,fields:!((n:t,sql:'toUnixTimestamp(generated_time)*1000'),(n:machine),(n:measures,subName:value)),filters:!((f:measures.name,v:localInspections#mean_value),(f:branch,v:master),(f:project,v:kotlin_empty/highlight/Main_with_library_cache_k1),(f:generated_time,q:'>subtractMonths(now(),1)')),order:t)")
	statsServer := &StatsServer{
		dbUrl: DefaultDbUrl,
	}
	request := http.Request{
		Method: "GET",
		URL:    testUrl,
	}
	dataQueries, _, _ := dataQuery.ReadQuery(&request)
	lowStartTime := time.Now()
	for i := 0; i < 100_000_000; i++ {
		statsServer.oldLoad(&request, dataQueries, false)
	}
	lowElapsedTime := time.Since(lowStartTime)

	heightStartTime := time.Now()
	for i := 0; i < 100_000_000; i++ {
		statsServer.load(request.Context(), dataQueries)
	}
	heightElapsedTime := time.Since(heightStartTime)

	fmt.Printf("Height level api executed in %s\n", heightElapsedTime)
	fmt.Printf("Low level api executed in %s\n", lowElapsedTime)
}
