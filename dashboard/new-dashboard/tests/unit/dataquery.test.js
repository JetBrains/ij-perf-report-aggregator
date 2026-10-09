import { expect, test } from "vitest"
import { mergeQueries } from "../../src/components/common/DataQueryExecutor"

const MASTER = { f: "branch", v: "master" }
const BRANCH_232 = { f: "branch", v: "232%", o: "like" }
const LAST_3_MONTHS = { f: "generated_time", q: ">subtractMonths(now(),3)" }
const NOT_TRIGGERED = { f: "triggeredBy", v: "" }
const AWS_LINUX = { f: "machine", v: "intellij-linux-performance-aws-%", o: "like" }
const NO_BUILD_C3 = { f: "build_c3", v: 0, o: "=" }

const MEASURE_FIELDS = [
  { n: "t", sql: "toUnixTimestamp(generated_time)*1000" },
  { n: "measures", subName: "value" },
  { n: "measures", subName: "name" },
  { n: "measures", subName: "type" },
  "machine",
  "tc_build_id",
  "project",
]
const BUILD_FIELDS = [...MEASURE_FIELDS, "tc_installer_build_id", "build_c1", "build_c2", "build_c3"]

const project = (v) => ({ f: "project", v })
const measure = (v) => ({ f: "measures.name", v })
// the filter mergeQueries produces when it combines queries that differ only in `f`
const merged = (f, v) => ({ f, v, s: true })

// cloned so that expected values never share objects with the queries passed to mergeQueries
function query({ db = "perfint", fields, filters, ...rest }) {
  return structuredClone({ db, table: "idea", ...(fields === undefined ? {} : { fields }), filters, order: "t", ...rest })
}

// the filters a perfint dashboard chart query is built with
function dashboardFilters(branch, projectFilter, measureFilter) {
  return [branch, LAST_3_MONTHS, NOT_TRIGGERED, AWS_LINUX, NO_BUILD_C3, projectFilter, measureFilter]
}

test("complex queries are not merged", () => {
  const queries = () =>
    [MASTER, BRANCH_232].map((branch) =>
      query({
        fields: [{ n: "measures", subName: "value" }],
        filters: [branch, LAST_3_MONTHS, AWS_LINUX, { f: "measures.name", v: "completion_%", o: "like" }],
        aggregator: "avg",
        dimensions: [{ n: "t", sql: "toYYYYMMDD(generated_time)" }],
      })
    )
  expect(mergeQueries(queries(), null)).toStrictEqual(queries())
})

test("merge by project", () => {
  const projects = ["intellij_sources/vfsRefresh/default", "intellij_sources/vfsRefresh/with-1-thread(s)", "intellij_sources/vfsRefresh/git-status"]
  const vfsQuery = (projectFilter) => query({ fields: ["project"], filters: dashboardFilters(MASTER, projectFilter, measure("vfs_initial_refresh")) })

  const actual = mergeQueries(
    projects.map((it) => vfsQuery(project(it))),
    null
  )
  expect(actual).toStrictEqual([vfsQuery(merged("project", projects))])
})

test("merge in case of metric and project", () => {
  const projects = ["community/indexing", "intellij_sources/indexing", "space/indexing"]
  const measures = ["indexing", "indexingTimeWithoutPauses"]
  const indexingQuery = (projectName, measureFilter) => query({ fields: BUILD_FIELDS, filters: dashboardFilters(MASTER, project(projectName), measureFilter) })

  const actual = mergeQueries(
    projects.flatMap((p) => measures.map((m) => indexingQuery(p, measure(m)))),
    null
  )
  expect(actual).toStrictEqual(projects.map((p) => indexingQuery(p, merged("measures.name", measures))))
})

test("single query as is", () => {
  const queries = () => [query({ fields: BUILD_FIELDS, filters: dashboardFilters(MASTER, project("kotlin_petclinic/debug"), measure("debugRunConfiguration")) })]
  expect(mergeQueries(queries(), null)).toStrictEqual(queries())
})

test("don't merge with operator", () => {
  const queries = () => [MASTER, BRANCH_232].map((branch) => query({ filters: dashboardFilters(branch, project("kotlin_petclinic/debug"), measure("debugRunConfiguration")) }))
  expect(mergeQueries(queries(), null)).toStrictEqual(queries())
})

test("don't merge with but merge with project", () => {
  const projects = ["community/rebuild", "intellij_sources/rebuild"]
  const rebuildQuery = (branch, projectFilter) => query({ fields: ["project"], filters: dashboardFilters(branch, projectFilter, measure("build_compilation_duration")) })

  const actual = mergeQueries(
    [MASTER, BRANCH_232].flatMap((branch) => projects.map((p) => rebuildQuery(branch, project(p)))),
    null
  )
  expect(actual).toStrictEqual([MASTER, BRANCH_232].map((branch) => rebuildQuery(branch, merged("project", projects))))
})

test("don't merge if merging field is not in query", () => {
  const importQuery = (branch) =>
    query({
      db: "perfintDev",
      fields: MEASURE_FIELDS,
      filters: [
        project("project-import-jps-kotlin-50_000-modules/fastInstaller"),
        { f: "branch", v: branch },
        {
          f: "machine",
          v: ["intellij-linux-hw-hetzner-agent-06", "intellij-linux-hw-hetzner-agent-13", "intellij-linux-hw-hetzner-agent-17", "intellij-linux-hw-hetzner-agent-21"],
        },
        { f: "generated_time", q: ">subtractYears(now(),1)" },
        NOT_TRIGGERED,
        measure("workspaceModel.to.snapshot.ms"),
      ],
    })
  const queries = () => ["nikita.kudrin/jps_12_september_regression_on_suspect", "nikita.kudrin/jps_12_september_regression_before"].map(importQuery)
  expect(mergeQueries(queries(), null)).toStrictEqual(queries())
})

test("merging is correct", () => {
  const projects = ["gitlab-project-inspections-test/inspection-app", "gitlab-project-inspections-test/inspection-RubyResolve-app"]
  const inspectionQuery = (projectFilter) =>
    query({ db: "perfintDev", fields: ["project", "machine"], filters: [projectFilter, { f: "machine", v: ["intellij-macos-hw-munit-692"] }] })

  const actual = mergeQueries(
    projects.map((it) => inspectionQuery(project(it))),
    null
  )
  expect(actual).toStrictEqual([inspectionQuery(merged("project", projects))])
})
