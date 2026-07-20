<template>
  <div class="flex flex-col gap-4 py-3 px-5 border border-solid rounded-md">
    <!-- Headline verdict -->
    <div class="flex items-start justify-between gap-4 flex-wrap">
      <div>
        <div class="text-xs uppercase tracking-wide text-gray-500">{{ eyebrow }}</div>
        <div class="text-2xl font-semibold flex items-baseline gap-2">
          <span :class="toneClass(overall.tone)">{{ overall.headline }}</span>
        </div>
        <div class="text-xs text-gray-500 mt-1">{{ subtitle }}</div>
      </div>
      <button
        type="button"
        class="pi pi-info-circle text-gray-400 hover:text-gray-600 dark:hover:text-gray-200 cursor-help mt-1"
        title="How this verdict is computed"
        aria-label="How this verdict is computed"
        @click="toggleMethodology"
      />
      <Popover
        ref="methodologyPanel"
        append-to="body"
      >
        <div class="flex flex-col gap-2 text-sm max-w-md">
          <span class="font-semibold">How this verdict is computed</span>
          <hr class="w-full border-gray-200 dark:border-gray-600" />
          <p>
            The overall figure is the geometric mean of each scenario's NEW ÷ LEGACY ratio — the SPEC-style summary that keeps a 180× spread of absolute values from letting slow
            files dominate.
          </p>
          <p>
            A single grand mean can hide a Simpson's-paradox split (NEW winning overall while losing every bucket), so the per-bucket and per-phase means are shown alongside it —
            trust the verdict only when they agree.
          </p>
          <p>
            Improved / regressed / neutral counts apply the same two-gate significance test as the compare table (robust effect size <em>and</em> ≥5% change); everything else is
            neutral.
          </p>
        </div>
      </Popover>
    </div>

    <!-- Per-bucket and per-phase geomeans — the Simpson's-paradox guard -->
    <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
      <div>
        <div class="text-xs uppercase tracking-wide text-gray-500 mb-1">By bucket</div>
        <div class="flex flex-wrap gap-x-6 gap-y-1">
          <div
            v-for="group in aggregates.perBucket"
            :key="group.key"
            class="text-sm"
          >
            <span class="text-gray-500">{{ bucketLabel(group.key) }}:</span>
            <span
              class="ml-1 font-medium"
              :class="toneClass(describeGeomean(group.geomean).tone)"
              >{{ describeGeomean(group.geomean).label }}</span
            >
          </div>
        </div>
      </div>
      <div>
        <div class="text-xs uppercase tracking-wide text-gray-500 mb-1">By phase</div>
        <div class="flex flex-wrap gap-x-6 gap-y-1">
          <div
            v-for="group in aggregates.perPhase"
            :key="group.key"
            class="text-sm"
          >
            <span class="text-gray-500">{{ phaseLabel(group.key) }}:</span>
            <span
              class="ml-1 font-medium"
              :class="toneClass(describeGeomean(group.geomean).tone)"
              >{{ describeGeomean(group.geomean).label }}</span
            >
          </div>
        </div>
      </div>
    </div>

    <!-- Three-way counts + worst/best movers -->
    <div class="flex flex-wrap items-stretch gap-6">
      <div class="flex items-center gap-4">
        <div class="text-center">
          <div class="text-lg font-semibold text-green-600 dark:text-green-400">{{ aggregates.improved }}</div>
          <div class="text-xs text-gray-500">improved</div>
        </div>
        <div class="text-center">
          <div class="text-lg font-semibold text-red-600 dark:text-red-400">{{ aggregates.regressed }}</div>
          <div class="text-xs text-gray-500">regressed</div>
        </div>
        <div class="text-center">
          <div class="text-lg font-semibold text-gray-500">{{ aggregates.neutral }}</div>
          <div class="text-xs text-gray-500">neutral</div>
        </div>
      </div>
      <div class="flex flex-col gap-1 text-sm">
        <div v-if="aggregates.worstRegression">
          <span class="text-gray-500">Worst regression:</span>
          <span class="ml-1 font-medium">{{ aggregates.worstRegression.title }} · {{ metricTypeLabel(aggregates.worstRegression.metricType) }}</span>
          <span class="ml-1 text-red-600 dark:text-red-400">{{ formatSignedPercent(aggregates.worstRegression.diffPercent) }}</span>
        </div>
        <div v-if="aggregates.bestImprovement">
          <span class="text-gray-500">Best improvement:</span>
          <span class="ml-1 font-medium">{{ aggregates.bestImprovement.title }} · {{ metricTypeLabel(aggregates.bestImprovement.metricType) }}</span>
          <span class="ml-1 text-green-600 dark:text-green-400">{{ formatSignedPercent(aggregates.bestImprovement.diffPercent) }}</span>
        </div>
      </div>
    </div>

    <p class="text-xs text-gray-500">
      A synthetic win here maps to revenue only when the interactive First-Code-Analysis APDEX gate (which covers highlighting) moves too — read this as a sensor, not a money
      metric.
    </p>
  </div>
</template>

<script setup lang="ts">
import { computed, useTemplateRef } from "vue"
import { PopoverMethods } from "primevue/popover"
import { EngineAggregates } from "./engineCompareAggregates"
import { EngineComparisonMode } from "./useEngineComparison"
import { bucketLabel, metricTypeLabel, phaseLabel } from "./highlightingMetrics"
import { formatSignedPercent } from "../../charts/compareStats"

const { aggregates, mode, runLabel } = defineProps<{ aggregates: EngineAggregates; mode: EngineComparisonMode; runLabel: string }>()

type Tone = "faster" | "slower" | "same" | "none"

const eyebrow = computed(() => (mode === "single" ? `NEW vs LEGACY — ${runLabel || "single run"}` : "NEW vs LEGACY — overall"))

const subtitle = computed(() =>
  mode === "single"
    ? `single run${runLabel ? ` — ${runLabel}` : ""} · ${aggregates.count} scenarios (NEW ÷ LEGACY)`
    : `geometric mean of ${aggregates.count} per-scenario ratios (NEW ÷ LEGACY)`
)

// A geomean of NEW ÷ LEGACY ratios, phrased as a faster/slower verdict. Below 1% is "about the same"
// so rounding noise does not read as a real move.
function describeGeomean(geomean: number): { label: string; tone: Tone } {
  if (!Number.isFinite(geomean) || geomean <= 0) return { label: "—", tone: "none" }
  const percent = (geomean - 1) * 100
  if (Math.abs(percent) < 1) return { label: "about the same", tone: "same" }
  return percent < 0 ? { label: `${Math.abs(percent).toFixed(0)}% faster`, tone: "faster" } : { label: `${percent.toFixed(0)}% slower`, tone: "slower" }
}

const overall = computed<{ headline: string; tone: Tone }>(() => {
  const { label, tone } = describeGeomean(aggregates.grandGeomean)
  if (tone === "none") return { headline: "Not enough comparable data yet", tone }
  if (tone === "same") return { headline: "NEW is about the same as LEGACY", tone }
  return { headline: `NEW is ${label}`, tone }
})

function toneClass(tone: Tone): string {
  switch (tone) {
    case "faster":
      return "text-green-600 dark:text-green-400"
    case "slower":
      return "text-red-600 dark:text-red-400"
    default:
      return "text-gray-500"
  }
}

const methodologyPanel = useTemplateRef<PopoverMethods>("methodologyPanel")

function toggleMethodology(event: Event): void {
  methodologyPanel.value?.toggle(event)
}
</script>
