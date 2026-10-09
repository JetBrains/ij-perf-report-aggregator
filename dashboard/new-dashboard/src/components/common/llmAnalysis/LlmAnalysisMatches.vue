<template>
  <Accordion
    v-if="items.length > 0"
    :value="collapsed ? null : '0'"
    class="llm-analysis-matches matches-card matches-card--analysed"
  >
    <AccordionPanel value="0">
      <AccordionHeader>
        <span class="flex flex-wrap items-center gap-2">
          <i class="pi pi-link matches-card__icon" />
          <span class="font-semibold">{{ title }}</span>
          <span class="matches-card__pill">{{ plural(items.length, "chart") }}</span>
        </span>
      </AccordionHeader>
      <AccordionContent>
        <ul class="flex flex-col gap-1.5 break-all">
          <li
            v-for="match in items"
            :key="match.id"
            :class="{ 'opacity-60': !isSameProduct(match) }"
          >
            <a
              v-if="match.dashboardLink"
              :href="analysisLink(match.dashboardLink, match.id)"
              target="_blank"
              class="underline decoration-dotted hover:no-underline"
            >
              {{ match.project }} / {{ match.metric }}
            </a>
            <span v-else>{{ match.project }} / {{ match.metric }}</span>
            <div class="flex flex-wrap gap-x-2 text-sm text-gray-500">
              <span
                v-tooltip.top="match.matchedCommits.join('\n')"
                class="font-mono"
              >
                {{ shortCommits(match.matchedCommits) }}
              </span>
              <span v-if="rangeSize">{{ match.matchedCommits.length }} of {{ rangeSize }} commits in range</span>
              <a
                v-if="match.ytIssueId"
                :href="`https://youtrack.jetbrains.com/issue/${match.ytIssueId}`"
                target="_blank"
                class="underline decoration-dotted hover:no-underline"
              >
                {{ match.ytIssueId }}
              </a>
              <span v-if="!isSameProduct(match)">other product</span>
            </div>
          </li>
        </ul>
      </AccordionContent>
    </AccordionPanel>
  </Accordion>
  <div
    v-if="loading && analysisId != null"
    class="matches-card matches-card--degraded flex items-center gap-2 text-sm"
  >
    <i class="pi pi-spin pi-spinner" />
    <span>Checking detected degradations…</span>
  </div>
  <Accordion
    v-if="degradationGroups.chartCount > 0"
    :value="collapsed ? null : '0'"
    class="llm-analysis-matches matches-card matches-card--degraded"
  >
    <AccordionPanel value="0">
      <AccordionHeader>
        <span class="flex flex-wrap items-center gap-2">
          <i class="pi pi-exclamation-triangle matches-card__icon" />
          <span class="font-semibold">Not analysed degradations at builds with this commit</span>
          <span class="matches-card__pill">{{ plural(degradationGroups.chartCount, "chart") }}</span>
          <span class="matches-card__pill">{{ plural(degradationGroups.groups.length, "metric") }}</span>
          <span
            v-if="degradationGroups.strongest"
            v-tooltip.top="'Smallest commit range among the charts: the fewer commits, the stronger the evidence'"
            class="matches-card__pill matches-card__pill--outline"
            >strongest {{ degradationGroups.strongest.matchedCommits.length }} of {{ degradationGroups.strongest.rangeSize }} commits</span
          >
        </span>
      </AccordionHeader>
      <AccordionContent>
        <div
          v-if="degradationGroups.commonCommits"
          class="mb-2 text-sm text-gray-500"
        >
          Matched commits:
          <span
            v-tooltip.top="degradationGroups.commonCommits.join('\n')"
            class="font-mono"
            >{{ shortCommits(degradationGroups.commonCommits) }}</span
          >
        </div>
        <section
          v-for="group in degradationGroups.groups"
          :key="group.metric"
          class="mb-3 last:mb-0"
        >
          <div class="mb-1 text-sm font-medium">
            {{ group.metric }} <span class="font-normal text-gray-500">· {{ group.charts.length }}</span>
          </div>
          <ul class="flex flex-col gap-1 break-all">
            <li
              v-for="chart in group.charts"
              :key="chart.project"
              class="flex flex-wrap items-baseline gap-x-2"
            >
              <a
                v-if="chartPointLink(chart.builds)"
                :href="chartPointLink(chart.builds) ?? undefined"
                target="_blank"
                class="underline decoration-dotted hover:no-underline"
              >
                {{ chart.project }}
              </a>
              <span
                v-else
                v-tooltip.top="chartTableKnown ? 'Not on the analysed chart page: tested in another database' : undefined"
                >{{ chart.project }}</span
              >
              <span class="text-sm text-gray-500">{{ chart.date }}</span>
              <span class="text-sm text-gray-500">{{ chart.matchedCommits.length }} of {{ chart.rangeSize }} commits</span>
              <span
                v-if="!degradationGroups.commonCommits"
                v-tooltip.top="chart.matchedCommits.join('\n')"
                class="font-mono text-sm text-gray-500"
              >
                {{ shortCommits(chart.matchedCommits) }}
              </span>
              <span class="text-sm text-gray-500">
                <template
                  v-for="(build, index) in chart.builds"
                  :key="build.buildId"
                >
                  <a
                    v-if="pointLink(build)"
                    v-tooltip.top="`${build.machine ?? 'unknown machine'} · ${build.matchedCommits.length} of ${build.rangeSize} commits`"
                    :href="pointLink(build) ?? undefined"
                    target="_blank"
                    class="underline decoration-dotted hover:no-underline"
                    >{{ chart.builds.length === 1 ? "point" : `point ${index + 1}` }}</a
                  >
                  <a
                    v-tooltip.top="'Open the build in TeamCity'"
                    :href="buildUrl(Number(build.buildId))"
                    target="_blank"
                    class="ml-1 hover:text-gray-900 dark:hover:text-gray-100"
                    ><i class="pi pi-external-link text-xs"
                  /></a>
                  {{ index < chart.builds.length - 1 ? ", " : "" }}
                </template>
              </span>
            </li>
          </ul>
        </section>
      </AccordionContent>
    </AccordionPanel>
  </Accordion>
</template>
<script setup lang="ts">
import { storeToRefs } from "pinia"
import { computed, ref, watch } from "vue"
import { useRouter } from "vue-router"
import { injectOrNull } from "../../../shared/injectionKeys"
import { serverConfiguratorKey } from "../../../shared/keys"
import { formatCustomRange } from "../../../configurators/TimeRangeConfigurator"
import { useSettingsStore } from "../../settings/settingsStore"
import { analysisParamName, pointParamName } from "../../../shared/selectedPointStore"
import { buildUrl, majorBranch } from "../sideBar/InfoSidebar"
import { groupDegradations } from "./degradationGroups"
import { LlmAnalysisClient, LlmAnalysisDegradationMatch, LlmAnalysisMatch } from "./LlmAnalysisClient"

// buildId: analyses whose guilty commit is in this build's range (degradation not analysed yet).
// analysisId: other analyses that blamed the same commit as this one.
const {
  title,
  buildId,
  analysisId,
  chart,
  chartLink,
  collapsed = false,
} = defineProps<{
  title: string
  // the lists start collapsed, their headers still show the counts
  collapsed?: boolean
  // dashboard link of the analysed chart, a degradation link is derived from it as it is the same metric
  chartLink?: string
  buildId?: number | null
  analysisId?: number | string | null
  // the point's own chart (buildId mode): its analyses are not matches, the point already lists them as its runs
  chart?: { project: string; metric?: string; currentBuildId: string }
}>()

const client = new LlmAnalysisClient(injectOrNull(serverConfiguratorKey))
const router = useRouter()
const { similarDegradations } = storeToRefs(useSettingsStore())

const matches = ref<LlmAnalysisMatch[]>([])
const rangeSize = ref<number | undefined>()
const degradations = ref<LlmAnalysisDegradationMatch[]>([])
const loading = ref(false)

// first path segment is the product, e.g. /clion/testsDev
const productOf = (link: string) => new URL(link, globalThis.location.origin).pathname.split("/")[1] ?? ""
// the analysed chart's product when known: the page may be /analyses
// compared to the analysed chart, or to the chart page in the sidebar; with neither (an analysis without a link opened
// from /analyses) there is nothing to compare to
const referenceLink = computed(() => chartLink ?? (buildId == null ? null : router.currentRoute.value.path))
const isSameProduct = (match: LlmAnalysisMatch) => match.dashboardLink == null || referenceLink.value == null || productOf(match.dashboardLink) === productOf(referenceLink.value)

const degradationGroups = computed(() => groupDegradations(degradations.value))

const plural = (count: number, noun: string) => `${count} ${noun}${count === 1 ? "" : "s"}`
const shortCommits = (commits: string[]) => commits.map((c) => c.slice(0, 10)).join(", ")

const items = computed(() => matches.value.toSorted((a, b) => Number(isSameProduct(b)) - Number(isSameProduct(a))))

function analysisLink(dashboardLink: string, id: number): string {
  const url = new URL(dashboardLink, globalThis.location.origin)
  url.searchParams.set(analysisParamName, String(id))
  return url.pathname + url.search
}

// set when the backend resolved machines in the analysed chart table: a build without a machine is then in another
// table, and the analysed chart page can't show it; false when the table is unknown (an analysis started before it was
// stored) or the lookup failed, links then keep the analysed machine
const chartTableKnown = ref(false)

function pointLink(build: LlmAnalysisDegradationMatch): string | null {
  if (chartLink == null || (chartTableKnown.value && !build.machine)) return null
  return degradationLink(chartLink, build)
}

// the strongest build the chart page can show
function chartPointLink(builds: LlmAnalysisDegradationMatch[]): string | null {
  for (const build of builds) {
    const link = pointLink(build)
    if (link != null) return link
  }
  return null
}

// the page of the analysed chart, which holds for the same or a related metric, with the machine and branch of the
// degradation build: the same commit is often tested on several machines and merged to other branches, and the chart
// shows only the selected ones
function degradationLink(chartLink: string, degradation: LlmAnalysisDegradationMatch): string {
  const url = new URL(chartLink, globalThis.location.origin)
  if (degradation.machine) url.searchParams.set("machine", degradation.machine)
  if (degradation.branch) url.searchParams.set("branch", majorBranch(degradation.branch))
  url.searchParams.set("project", degradation.project)
  url.searchParams.set("measure", degradation.metric)
  url.searchParams.set(pointParamName, degradation.buildId)
  url.searchParams.delete(analysisParamName)
  const date = new Date(`${degradation.date}T00:00`)
  const start = new Date(date)
  start.setDate(date.getDate() - 30)
  const end = new Date(date)
  end.setDate(date.getDate() + 7)
  url.searchParams.set("timeRange", "custom")
  url.searchParams.set("customRange", formatCustomRange(start, end))
  return url.pathname + url.search
}

watch(
  () => [buildId, analysisId, similarDegradations.value],
  async () => {
    matches.value = []
    rangeSize.value = undefined
    degradations.value = []
    chartTableKnown.value = false
    loading.value = false
    const requested = [buildId, analysisId]
    if (!similarDegradations.value || (buildId == null && analysisId == null)) return
    const isCurrent = () => similarDegradations.value && requested[0] === buildId && requested[1] === analysisId
    // independent requests: degradations take longer, and either one failing must not hide the other list
    const loadMatches = client.getMatches(buildId == null ? { analysisId: String(analysisId) } : { buildId: String(buildId), ...chart }).then((result) => {
      if (!isCurrent()) return
      matches.value = result.matches
      rangeSize.value = result.rangeSize
    })
    const loads = [loadMatches]
    if (buildId == null) {
      loading.value = true
      const loadDegradations = client
        .getMatches({ analysisId: String(analysisId), degradations: "true" })
        .then((result) => {
          if (!isCurrent()) return
          chartTableKnown.value = result.machinesResolved === true
          degradations.value = result.degradations ?? []
        })
        .finally(() => {
          if (isCurrent()) loading.value = false
        })
      loads.push(loadDegradations)
    }
    for (const result of await Promise.allSettled(loads)) {
      if (result.status === "rejected") console.warn(result.reason)
    }
  },
  { immediate: true }
)
</script>
<style scoped>
.llm-analysis-matches :deep(.p-accordionheader) {
  padding: 0 0 1rem 0;
}

.llm-analysis-matches :deep(.p-accordioncontent) {
  padding: 0;
}

/* a tinted card with an accent bar, so the collapsed lists stand out from the details above */
.matches-card {
  --accent: #0284c7;
  --tint: #f0f9ff;
  --pill: #e0f2fe;
  --pill-text: #075985;
  border-left: 4px solid var(--accent);
  border-radius: 0.375rem;
  background: var(--tint);
  padding: 0.5rem 0.75rem;
  margin-bottom: 0.5rem;
}
.matches-card--degraded {
  --accent: #e11d48;
  --tint: #fff1f2;
  --pill: #ffe4e6;
  --pill-text: #9f1239;
}
.dark-mode .matches-card {
  --tint: rgb(12 74 110 / 0.25);
  --pill: rgb(12 74 110 / 0.6);
  --pill-text: #e0f2fe;
}
.dark-mode .matches-card--degraded {
  --tint: rgb(136 19 55 / 0.25);
  --pill: rgb(136 19 55 / 0.6);
  --pill-text: #ffe4e6;
}

.matches-card :deep(.p-accordionpanel) {
  border: 0;
}
.matches-card :deep(.p-accordionheader) {
  padding: 0;
  background: transparent;
}
.matches-card :deep(.p-accordioncontent-content) {
  background: transparent;
  padding: 0.75rem 0 0.25rem;
}

.matches-card__icon {
  color: var(--accent);
}
.matches-card__pill {
  border-radius: 9999px;
  background: var(--pill);
  color: var(--pill-text);
  padding: 0 0.5rem;
  font-size: 0.75rem;
  font-weight: 500;
  line-height: 1.25rem;
}
.matches-card__pill--outline {
  background: transparent;
  border: 1px solid var(--accent);
}
</style>
