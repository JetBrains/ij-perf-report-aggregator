import { Router } from "vue-router"
import { AccidentKind } from "../../../configurators/accidents/AccidentsConfigurator"
import { TimeRangeConfigurator } from "../../../configurators/TimeRangeConfigurator"
import { dbTypeStore } from "../../../shared/dbTypes"
import { useUserStore } from "../../../shared/useUserStore"
import { getTeamcityBuildType } from "../../../util/artifacts"
import { getFirstAndLastCommit } from "../../../util/changes"
import { getPersistentLink } from "../../settings/CopyLink"
import { ServerConfigurator } from "../dataQuery"
import { BisectClient } from "./BisectClient"
import { suggestTargetValue } from "./BisectChecks"
import { DBType, getNavigateToTestUrl, InfoData } from "./InfoSidebar"

// Bisect checks out individual source commits, so it only applies to the source-based
// configurations of the dev database; installer-based points can't be mapped to a commit range.
export function isBisectSupported(data: InfoData | null): boolean {
  return dbTypeStore().dbType == DBType.INTELLIJ_DEV && data != null && data.installerId == undefined
}

export function bisectDirectionOf(accidentKind: string): string {
  return accidentKind == AccidentKind.Regression || accidentKind == AccidentKind.InferredRegression ? "DEGRADATION" : "OPTIMIZATION"
}

// JPS compilation was the default for master builds up to the cutoff date; newer builds use Bazel.
export function defaultJpsCompilation(data: InfoData): boolean {
  return data.branch === "master" && new Date(data.date) <= new Date("2025-10-19T23:59:59.999Z")
}

export function bisectTestPatterns(data: InfoData): string {
  const methodName = data.description.value?.methodName ?? ""
  return methodName.slice(0, Math.max(0, methodName.lastIndexOf("#")))
}

// Bisect runs the single-test flavour of the configuration, which is always the _1 suffix.
export function resolveBisectBuildType(serverConfigurator: ServerConfigurator, buildId: number): Promise<string | null> {
  return getTeamcityBuildType(serverConfigurator.db, serverConfigurator.table, buildId).then((buildType) => (buildType ? buildType.replace(/_\d+$/, "_1") : buildType))
}

// Whether a bisect can be started for this point without asking the user for anything:
// everything the bisect needs must be derivable, including a target value that cleanly
// separates the levels before and after the point. When it isn't, the point has to go
// through the bisect dialog so the user can pick the value and acknowledge the warnings.
export function canStartBisectAutomatically(data: InfoData | null, accidentKind: string): boolean {
  return isBisectSupported(data) && data != null && data.series[0]?.metricName != null && suggestTargetValue(data, bisectDirectionOf(accidentKind)) != null
}

export interface StartBisectParams {
  data: InfoData
  serverConfigurator: ServerConfigurator
  router: Router
  timerangeConfigurator: TimeRangeConfigurator
  accidentKind: string
  // The bisect job reports its result back to this issue.
  ytIssueId?: string
}

// Starts a bisect with the parameters the bisect dialog would pre-fill, without showing it.
// Returns the URL of the started TeamCity build.
export async function startBisect(params: StartBisectParams): Promise<string> {
  const { data, serverConfigurator, router, timerangeConfigurator, accidentKind, ytIssueId } = params

  const metric = data.series[0]?.metricName
  if (metric == null) {
    throw new Error("Metric is required to start a bisect")
  }
  const direction = bisectDirectionOf(accidentKind)
  const targetValue = suggestTargetValue(data, direction)
  if (targetValue == null) {
    throw new Error("Cannot derive a target value for this point, start the bisect from the sidebar instead")
  }

  const [{ firstCommit, lastCommit }, buildType] = await Promise.all([
    getFirstAndLastCommit(serverConfigurator.db, data.installerId ?? data.buildId),
    resolveBisectBuildType(serverConfigurator, data.buildId),
  ])
  if (!firstCommit || !lastCommit) {
    throw new Error("Build has no changes to bisect")
  }
  if (!buildType) {
    throw new Error("Cannot resolve the TeamCity build type of the build")
  }

  return new BisectClient(serverConfigurator).sendBisectRequest({
    targetValue: String(targetValue),
    requester: useUserStore().user?.email ?? "",
    changes: `${firstCommit}^..${lastCommit}`,
    buildId: data.buildId.toString(),
    mode: "build",
    direction,
    test: data.projectName,
    metric,
    buildType,
    testPatterns: bisectTestPatterns(data),
    excludedCommits: "",
    jpsCompilation: defaultJpsCompilation(data) ? "true" : "false",
    dashboardLink: window.location.origin + getPersistentLink(getNavigateToTestUrl(data, router), timerangeConfigurator),
    ytIssueId,
  })
}
