<template>
  <Dialog
    v-model:visible="showDialog"
    modal
    header="Run LLM Analysis"
    :style="{ width: '40vw' }"
  >
    <div class="flex flex-col space-y-6 mb-4 mt-4">
      <div class="flex flex-col gap-1 text-sm">
        <div><span class="font-medium">Test:</span> {{ data.projectName }}</div>
        <div><span class="font-medium">Metric:</span> {{ metric }}</div>
        <div v-if="data.formattedPreviousValue != null || data.formattedCurrentValue != null">
          <span class="font-medium">Change:</span> {{ data.formattedPreviousValue ?? "?" }} → {{ data.formattedCurrentValue ?? "?" }}
          <span
            v-if="data.deltaPrevious"
            :class="deltaColor"
            >({{ data.deltaPrevious }})</span
          >
        </div>
        <div><span class="font-medium">Builds:</span> {{ data.buildIdPrevious }} → {{ data.buildId }}</div>
      </div>

      <!-- Intervening builds failed/timed out, so commits they consumed appear on no graph dot
           and sit below this build's own change range. The analysis widens the range back to the
           previous data point to include them; tell the user so the wider range isn't surprising. -->
      <WarningNotice
        v-if="changesGap?.hasGap"
        title="Change range widened"
      >
        <div class="text-sm">
          {{ changesGap.gapCommitCount }} commit{{ changesGap.gapCommitCount === 1 ? "" : "s" }} landed in builds that produced no data point (e.g. failed or timed-out runs) and
          {{ changesGap.gapCommitCount === 1 ? "is" : "are" }} not part of this build's own change range. The analysis range has been widened back to the previous data point on the
          graph so {{ changesGap.gapCommitCount === 1 ? "it is" : "they are" }} included as suspects.
        </div>
      </WarningNotice>

      <!-- Whether the selected point itself looks like the wrong one to analyse
           (unchanged, wrong direction, or a much larger change nearby). Unlike the
           bisect dialog we don't check graph stability here: a metric too noisy to
           bisect can still yield a useful LLM analysis, so it's not worth warning about. -->
      <WarningNotice
        v-if="misclickWarning"
        :title="misclickWarning.title"
      >
        <div class="text-sm">{{ misclickWarning.detail }}</div>
        <div class="flex items-center mt-3">
          <Checkbox
            id="acknowledgeMisclick"
            v-model="acknowledgedMisclick"
            binary
          />
          <label
            for="acknowledgeMisclick"
            class="ml-2 text-sm"
            >I checked the selected point and want to analyse it</label
          >
        </div>
      </WarningNotice>
    </div>
    <div
      v-if="error"
      class="text-red-500 mb-4"
    >
      {{ error }}
    </div>
    <template #footer>
      <div class="flex justify-end space-x-2">
        <Button
          label="Cancel"
          icon="pi pi-times"
          severity="secondary"
          @click="showDialog = false"
        />
        <Button
          v-tooltip.top="
            reasonOfDisabling === ''
              ? null
              : {
                  value: reasonOfDisabling,
                  autoHide: false,
                }
          "
          label="Run analysis"
          icon="pi pi-sparkles"
          autofocus
          :loading="loading"
          :disabled="reasonOfDisabling !== ''"
          @click="runAnalysis"
        />
      </div>
    </template>
  </Dialog>
</template>
<script setup lang="ts">
import { computedAsync } from "@vueuse/core"
import { useToast } from "primevue/usetoast"
import { computed, ref } from "vue"
import { AccidentKind } from "../../../configurators/accidents/AccidentsConfigurator"
import { LlmAnalysesConfigurator } from "../../../configurators/llmAnalyses/LlmAnalysesConfigurator"
import { injectOrError } from "../../../shared/injectionKeys"
import { serverConfiguratorKey } from "../../../shared/keys"
import { getFirstAndLastCommit } from "../../../util/changes"
import WarningNotice from "../WarningNotice.vue"
import { BisectClient } from "../sideBar/BisectClient"
import { InfoData } from "../sideBar/InfoSidebar"
import { detectPossibleMisclick } from "../sideBar/MisclickHeuristic"
import { startLlmAnalysisWithToast } from "./LlmAnalysisUtils"

const { data, llmAnalysesConfigurator } = defineProps<{
  data: InfoData
  llmAnalysesConfigurator: LlmAnalysesConfigurator
}>()

const showDialog = defineModel<boolean>("showDialog")

const toast = useToast()

const metric = data.series[0]?.metricName ?? ""
// Mirror the bisect dialog: a delta containing "-" is treated as a degradation.
const isDegradation = data.deltaPrevious?.includes("-") ?? false
const deltaColor = computed(() => (isDegradation ? "text-red-500" : "text-green-600"))

const misclickWarning = computed(() => detectPossibleMisclick(data, isDegradation ? AccidentKind.Regression : AccidentKind.Improvement))
const acknowledgedMisclick = ref(false)

// Predict the server-side range widening so the dialog can tell the user about it. When earlier
// builds failed/timed out, the commits they consumed exist on no graph dot and fall below this
// build's own change range; the analysis backend widens the range back to the previous data point
// to cover them. We detect the same gap here via the shared changesGap endpoint. Only meaningful
// for source-based configs — installer-based ones carry no VCS changes of their own, so there is
// no range to widen (mirrors the bisect dialog).
const serverConfigurator = injectOrError(serverConfiguratorKey)
const bisectClient = new BisectClient(serverConfigurator)
const changesGap = computedAsync(async () => {
  if (data.buildIdPrevious == null || data.installerId != undefined) return null
  const { firstCommit } = await getFirstAndLastCommit(serverConfigurator.db, data.buildId)
  if (!firstCommit) return null
  return bisectClient.fetchChangesGap(String(data.buildId), String(data.buildIdPrevious), firstCommit)
}, null)

const error = ref<string | null>(null)
const loading = ref(false)

const reasonOfDisabling = computed(() => {
  if (misclickWarning.value && !acknowledgedMisclick.value) {
    return "Please acknowledge the possible wrong-point warning"
  }
  return ""
})

async function runAnalysis() {
  error.value = null
  loading.value = true
  try {
    const started = await startLlmAnalysisWithToast(llmAnalysesConfigurator, data, toast)
    if (started) {
      showDialog.value = false
    }
  } finally {
    loading.value = false
  }
}
</script>
