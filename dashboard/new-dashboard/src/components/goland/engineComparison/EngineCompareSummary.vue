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
          <p>The geometric mean of NEW ÷ LEGACY over all cells, equally weighted. A cell's ratio is the center of its per-build ratios.</p>
          <p>A cell is significant when it moves by ≥5% and lies outside the build-to-build spread of its ratio.</p>
          <p>
            Faster or slower needs a ≥5% move, no bucket or phase moving the other way, and fewer significant cells against the move than with it; otherwise the verdict is mixed.
            In a single run, a missing result blocks the verdict.
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

    <div
      v-for="line in gapLines"
      :key="line.engine"
      class="text-sm"
    >
      <span class="font-medium text-red-600 dark:text-red-400">{{ line.engine }} missing:</span>
      <span class="ml-1">{{ line.cells.join(", ") }}</span>
    </div>

    <p class="text-xs text-gray-500">
      A synthetic win here maps to revenue only when the interactive First-Code-Analysis APDEX gate (which covers highlighting) moves too — read this as a sensor, not a money
      metric.
    </p>
  </div>
</template>

<script setup lang="ts">
import { computed, useTemplateRef } from "vue"
import { PopoverMethods } from "openvue/popover"
import { EngineAggregates, EngineGap, geomeanDirection, resolveVerdict } from "./engineCompareAggregates"
import { EngineComparisonMode } from "./useEngineComparison"
import { bucketLabel, metricTypeLabel, phaseLabel, Quantity } from "./highlightingMetrics"
import { formatSignedPercent } from "../../charts/compareStats"

const { aggregates, mode, runLabel, quantity } = defineProps<{ aggregates: EngineAggregates; mode: EngineComparisonMode; runLabel: string; quantity: Quantity }>()

type Tone = "better" | "worse" | "same" | "none"

const eyebrow = computed(() => {
  const what = quantity.suffix === "" ? "NEW vs LEGACY" : `NEW vs LEGACY — ${quantity.label}`
  return mode === "single" ? `${what} — ${runLabel || "single run"}` : `${what} — overall`
})

const subtitle = computed(() => `${aggregates.count} of ${aggregates.cells} cells (project × phase × file)`)

function gapText(gap: EngineGap): string {
  const builds = mode === "aggregated" ? ` (${gap.builds} ${gap.builds === 1 ? "build" : "builds"})` : ""
  return `${gap.title} · ${metricTypeLabel(gap.metricType)}${builds}`
}

const gapLines = computed(() =>
  (["NEW", "LEGACY"] as const)
    .map((engine) => ({ engine, cells: aggregates.gaps.filter((gap) => gap.engine === engine).map((gap) => gapText(gap)) }))
    .filter((line) => line.cells.length > 0)
)

// Lower is better for every quantity, but only a duration is "faster": less allocation is "lower".
const words = computed(() => (quantity.suffix === "" ? { lower: "faster", higher: "slower" } : { lower: "lower", higher: "higher" }))

// A geomean of NEW ÷ LEGACY ratios as a signed move, or "about the same" inside the practical threshold.
function describeGeomean(geomean: number): { label: string; tone: Tone } {
  if (!Number.isFinite(geomean) || geomean <= 0) return { label: "—", tone: "none" }
  const percent = (geomean - 1) * 100
  switch (geomeanDirection(geomean)) {
    case -1:
      return { label: `${Math.abs(percent).toFixed(0)}% ${words.value.lower}`, tone: "better" }
    case 1:
      return { label: `${percent.toFixed(0)}% ${words.value.higher}`, tone: "worse" }
    default:
      return { label: "about the same", tone: "same" }
  }
}

function cellCount(count: number): string {
  return `${count} ${count === 1 ? "cell" : "cells"}`
}

const overall = computed<{ headline: string; tone: Tone }>(() => {
  const verdict = resolveVerdict(aggregates, mode === "single")
  const { label } = describeGeomean(aggregates.grandGeomean)
  switch (verdict.kind) {
    case "noData":
      return { headline: "No data", tone: "none" }
    case "incomplete": {
      const parts = [verdict.newGaps > 0 ? `NEW missing in ${cellCount(verdict.newGaps)}` : "", verdict.legacyGaps > 0 ? `LEGACY missing in ${cellCount(verdict.legacyGaps)}` : ""]
      return { headline: `No verdict: ${parts.filter((part) => part !== "").join(", ")}`, tone: verdict.newGaps > 0 ? "worse" : "none" }
    }
    case "same":
      return { headline: "About the same", tone: "same" }
    case "mixed":
      return { headline: `Mixed: ${label} overall`, tone: "same" }
  }
  return { headline: `NEW is ${label}`, tone: verdict.kind }
})

function toneClass(tone: Tone): string {
  switch (tone) {
    case "better":
      return "text-green-600 dark:text-green-400"
    case "worse":
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
