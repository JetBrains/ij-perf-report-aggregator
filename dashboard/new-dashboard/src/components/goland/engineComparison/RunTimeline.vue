<template>
  <div class="flex flex-col gap-1">
    <div
      v-if="runs.length === 0"
      class="text-xs text-gray-500"
    >
      No runs available.
    </div>
    <template v-else>
      <div class="overflow-x-auto">
        <div class="relative flex items-center gap-4 px-2 py-2 min-w-max">
          <div
            class="absolute left-2 right-2 top-1/2 -translate-y-1/2 border-t border-gray-300 dark:border-gray-600"
            aria-hidden="true"
          />
          <button
            v-for="run in ordered"
            :key="run.t"
            v-tooltip.top="run.label"
            type="button"
            class="relative z-10 rounded-full transition-all"
            :class="
              isSelected(run.t)
                ? 'w-3.5 h-3.5 bg-blue-500 ring-2 ring-blue-300 dark:ring-blue-700'
                : 'w-2.5 h-2.5 bg-gray-300 dark:bg-gray-500 hover:bg-gray-400 dark:hover:bg-gray-400'
            "
            :aria-label="`Show run ${run.label}`"
            :aria-pressed="isSelected(run.t)"
            @click="select(run.t)"
          />
        </div>
      </div>
      <div class="text-xs text-gray-500">
        Run: <span class="font-medium">{{ selectedLabel }}</span>
        <span v-if="isLatest"> (latest)</span>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue"
import { Run } from "./useEngineComparison"

// null means "follow the latest run"; a build timestamp pins a specific run.
const model = defineModel<number | null>({ default: null })
const { runs } = defineProps<{ runs: Run[] }>()

// runs arrive newest-first; a timeline reads oldest -> newest left -> right.
const ordered = computed<Run[]>(() => runs.toReversed())
const newestRun = computed<number | null>(() => runs[0]?.t ?? null)

function isSelected(t: number): boolean {
  return model.value === t || (model.value == null && t === newestRun.value)
}

const isLatest = computed(() => model.value == null || model.value === newestRun.value)

const selectedLabel = computed(() => {
  const t = model.value ?? newestRun.value
  return runs.find((run) => run.t === t)?.label ?? "—"
})

function select(t: number): void {
  model.value = t
}
</script>
