package data_query

import (
	"github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAdvancedFilter(t *testing.T) {
	queries, err := readQuery([]byte(`
{
  "filters": [
    {"f": "generated_time", "sql": "> subtractMonths(now(), 1)"}
  ],
  "order": "generated_time"
}
`))
	require.NoError(t, err)
	assert.NotEmpty(t, queries)

	sql, _, err := buildSql(queries[0], "test")
	require.NoError(t, err)
	assert.Equal(t, "select from test where (generated_time > subtractMonths(now(), 1)) order by generated_time", sql)
}

func TestQueryWithManyFields(t *testing.T) {
	queries, err := readQuery([]byte(`
{"fields":[{"n":"t","sql":"toUnixTimestamp(generated_time)*1000"},{"n":"measures","subName":"value"},{"n":"measures","subName":"name"},{"n":"measures","subName":"type"},"machine","tc_build_id","project","tc_installer_build_id","build_c1","build_c2","build_c3","branch"],"filters":[{"f":"project","v":"intellij_sources/vfsRefresh/with-1-thread(s)"},{"f":"branch","v":"akoehler/vfs-degradation-before-vk-ffd635a390d3%","o":"like"},{"f":"machine","v":"intellij-linux-performance-aws-%","o":"like"},{"f":"generated_time","q":">subtractMonths(now(),3)"},{"f":"triggeredBy","v":""},{"f":"build_c3","v":0,"o":"="},{"f":"measures.name","v":"vfs_initial_refresh"}],"order":"t"}
`))
	require.NoError(t, err)
	assert.NotEmpty(t, queries)

	sql, _, err := buildSql(queries[0], "test")
	require.NoError(t, err)
	assert.Equal(t, "select toUnixTimestamp(generated_time)*1000 as `t`, measures.value, measures.name, measures.type, machine, tc_build_id, project, tc_installer_build_id, build_c1, build_c2, build_c3, branch from test array join measures where (project = 'intellij_sources/vfsRefresh/with-1-thread(s)') and (branch like 'akoehler/vfs-degradation-before-vk-ffd635a390d3%') and (machine like 'intellij-linux-performance-aws-%') and (generated_time >subtractMonths(now(),3)) and (triggeredBy = '') and (build_c3=0) and (measures.name = 'vfs_initial_refresh') order by t", sql)
}

func TestOrInFilterQuery(t *testing.T) {
	queries, err := readQuery([]byte(`
{
  "filters": [
    {"f": "branch", "sql": "master or branch = 223"},
    {"f": "project", "v": "foo"}
  ],
  "order": "generated_time"
}
`))
	require.NoError(t, err)
	assert.NotEmpty(t, queries)

	sql, _, err := buildSql(queries[0], "test")
	require.NoError(t, err)
	assert.Equal(t, "select from test where (branch master or branch = 223) and (project = 'foo') order by generated_time", sql)
}

func TestAverageAggregate(t *testing.T) {
	queries, err := readQuery([]byte(`{
  "db": "perfint",
  "table": "phpstorm",
  "fields": [
    {
      "n": "t",
      "sql": "toYYYYMMDD(generated_time)"
    },
    {
      "n": "measures",
      "subName": "value"
    }
  ],
  "filters": [
    {
      "f": "measures.name",
      "v": [
        "responsiveness_time"
      ]
    },
    {
      "f": "branch",
      "v": "master"
    }
  ],
  "flat": false,
  "order": [
    "t"
  ],
  "aggregator": "avg",
  "dimensions": [
    {
      "n": "t"
    }
  ]
}
`))
	require.NoError(t, err)
	assert.NotEmpty(t, queries)

	sql, _, err := buildSql(queries[0], "test")
	require.NoError(t, err)
	assert.Equal(t, "select toYYYYMMDD(generated_time) as `t`, avg(measures.value) as measure_value from test array join measures where (measures.name in ('responsiveness_time')) and (branch = 'master') group by t order by t", sql)
}

func TestDecode(t *testing.T) {
	query, err := util.DecodeQuery("KLUv_SAMYQAASGVsbG8genN0ZCEh")

	require.NoError(t, err)
	assert.Equal(t, "Hello zstd!!", string(query))
}

func TestSpecialSymbolInSQL(t *testing.T) {
	queries, err := readQuery([]byte(`
[
    {
        "db": "ij",
        "table": "report",
        "fields": [
            {
                "n": "t",
                "sql": "toUnixTimestamp(generated_time)*1000"
            },
            {
                "n": "editorRestoring"
            },
            {
                "n": "metricName",
                "sql": "'editorRestoring'"
            },
            "machine",
            "tc_build_id",
            "project",
            "tc_installer_build_id",
            "build_c1",
            "build_c2",
            "build_c3",
            "branch"
        ],
        "filters": [
            {
                "f": "machine",
                "v": "intellij-linux-hw-munit-0%",
                "o": "like"
            },
            {
                "f": "generated_time",
                "q": ">subtractMonths(now(),3)"
            },
            {
                "f": "product",
                "v": "IU"
            },
            {
                "f": "project",
                "v": "simple for IJ"
            },
            {
                "f": "branch",
                "v": "master"
            },
            {
                "f": "triggeredBy",
                "v": ""
            },
            {
                "f": "editorRestoring",
                "o": "!=",
                "v": 0
            }
        ],
        "order": "t"
    }
]
`))
	require.NoError(t, err)
	assert.NotEmpty(t, queries)

	sql, _, err := buildSql(queries[0], "test")
	require.NoError(t, err)
	assert.Equal(t, "select toUnixTimestamp(generated_time)*1000 as `t`, editorRestoring, 'editorRestoring' as `metricName`, machine, tc_build_id, project, tc_installer_build_id, build_c1, build_c2, build_c3, branch from test where (machine like 'intellij-linux-hw-munit-0%') and (generated_time >subtractMonths(now(),3)) and (product = 'IU') and (project = 'simple for IJ') and (branch = 'master') and (triggeredBy = '') and (editorRestoring!=0) order by t", sql)
}
