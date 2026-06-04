<template>
  <div class="flex flex-col gap-3 text-sm">
    <div
      v-if="feedbacksLoading"
      class="text-xs text-gray-500"
    >
      Loading feedback…
    </div>
    <div
      v-else-if="feedbacksError"
      class="text-xs text-red-600"
    >
      {{ feedbacksError }}
    </div>
    <template v-else>
      <div
        v-if="myFeedback && !isEditing"
        class="flex flex-col gap-1"
      >
        <div class="font-medium text-gray-500">Your feedback</div>
        <div class="flex items-center gap-2 text-xs">
          <span class="font-medium whitespace-nowrap">★ {{ myFeedback.rate }}/5</span>
          <span class="text-gray-400 whitespace-nowrap">{{ formatCreatedAt(myFeedback.updatedAt) }}</span>
          <span
            v-if="myFeedback.feedback"
            v-tooltip.top="myFeedback.feedback"
            class="truncate text-gray-700"
          >
            {{ myFeedback.feedback }}
          </span>
          <Button
            label="Edit"
            severity="secondary"
            text
            size="small"
            @click="startEdit"
          />
        </div>
      </div>

      <div
        v-if="otherFeedbacks.length > 0"
        class="flex flex-col gap-1"
      >
        <div class="font-medium text-gray-500">From other users</div>
        <ul class="flex flex-col gap-1 text-xs">
          <li
            v-for="fb in otherFeedbacks"
            :key="fb.id"
            class="flex items-center gap-2"
          >
            <span class="font-medium whitespace-nowrap">★ {{ fb.rate }}/5</span>
            <span
              v-if="fb.userEmail"
              class="text-gray-500 whitespace-nowrap"
            >
              {{ userLocalPart(fb.userEmail) }}
            </span>
            <span class="text-gray-400 whitespace-nowrap">{{ formatCreatedAt(fb.updatedAt) }}</span>
            <span
              v-if="fb.feedback"
              v-tooltip.top="fb.feedback"
              class="truncate text-gray-700"
            >
              {{ fb.feedback }}
            </span>
          </li>
        </ul>
      </div>
    </template>

    <template v-if="!feedbacksLoading && !feedbacksError && (!myFeedback || isEditing)">
      <div class="flex flex-col gap-1">
        <label
          for="analysis-feedback-rating"
          class="font-medium text-gray-500"
        >
          Rating
        </label>
        <Select
          id="analysis-feedback-rating"
          v-model="rating"
          :options="ratingOptions"
          option-label="label"
          option-value="value"
          placeholder="Select a rating"
          class="w-80"
        />
      </div>
      <div class="flex flex-col gap-1">
        <label
          for="analysis-feedback-comments"
          class="font-medium text-gray-500"
        >
          Feedback
        </label>
        <Textarea
          id="analysis-feedback-comments"
          v-model="feedbackText"
          rows="5"
          class="w-full"
          placeholder="What could work better, ideally the actually guilty commit(s)"
        />
      </div>
      <div class="flex items-center gap-2">
        <div
          v-tooltip.top="submitDisabledReason"
          class="inline-block w-fit"
        >
          <Button
            :label="isEditing ? 'Save' : 'Send'"
            :disabled="submitDisabledReason != null || submitting"
            :loading="submitting"
            @click="submitFeedback"
          />
        </div>
        <Button
          v-if="isEditing"
          label="Cancel"
          severity="secondary"
          :disabled="submitting"
          @click="cancelEdit"
        />
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { storeToRefs } from "pinia"
import { computed, ref, watch } from "vue"
import { AnalysisFeedback, LlmAnalysisClient } from "./LlmAnalysisClient"
import { injectOrNull } from "../../../shared/injectionKeys"
import { serverConfiguratorKey } from "../../../shared/keys"
import { useUserStore } from "../../../shared/useUserStore"

const { analysisId, isTerminalState, active } = defineProps<{
  analysisId: number | string | null
  isTerminalState: boolean
  active: boolean
}>()

const serverConfigurator = injectOrNull(serverConfiguratorKey)
const client = new LlmAnalysisClient(serverConfigurator)

const rating = ref<number | null>(null)
const feedbackText = ref("")

const feedbacks = ref<AnalysisFeedback[]>([])
const feedbacksLoaded = ref(false)
const feedbacksLoading = ref(false)
const feedbacksError = ref<string | null>(null)
const submitting = ref(false)
const isEditing = ref(false)

const { user } = storeToRefs(useUserStore())
const currentUserEmail = computed(() => user.value?.email ?? null)

const myFeedback = computed<AnalysisFeedback | undefined>(() => {
  const email = currentUserEmail.value
  return email == null ? undefined : feedbacks.value.find((fb) => fb.userEmail === email)
})

const otherFeedbacks = computed<AnalysisFeedback[]>(() => {
  const mine = myFeedback.value
  return mine == null ? feedbacks.value : feedbacks.value.filter((fb) => fb !== mine)
})

const ratingOptions = [
  { value: 1, label: "1 — Misleading: pointed at the wrong cause" },
  { value: 2, label: "2 — Not useful: no actionable signal" },
  { value: 3, label: "3 — Right direction: helped in investigation, but culprit not found" },
  { value: 4, label: "4 — Close: culprit in list but ranked low or reasoning weak" },
  { value: 5, label: "5 — Spot on: culprit identified with sound reasoning" },
]

const submitDisabledReason = computed<string | null>(() => {
  if (!isTerminalState) return "Please wait till analysis is finished"
  if (rating.value == null) return "Select a rating"
  return null
})

function reset() {
  rating.value = null
  feedbackText.value = ""
  feedbacks.value = []
  feedbacksLoaded.value = false
  feedbacksError.value = null
  isEditing.value = false
}

function startEdit() {
  const mine = myFeedback.value
  if (mine == null) return
  rating.value = mine.rate
  feedbackText.value = mine.feedback ?? ""
  isEditing.value = true
}

function cancelEdit() {
  isEditing.value = false
  rating.value = null
  feedbackText.value = ""
}

async function loadFeedbacks() {
  if (analysisId == null || feedbacksLoaded.value || feedbacksLoading.value) return
  feedbacksLoading.value = true
  feedbacksError.value = null
  try {
    feedbacks.value = await client.getAnalysisFeedback(analysisId)
    feedbacksLoaded.value = true
  } catch (e) {
    feedbacksError.value = e instanceof Error ? e.message : String(e)
  } finally {
    feedbacksLoading.value = false
  }
}

async function submitFeedback() {
  if (rating.value == null || analysisId == null || submitting.value) return
  submitting.value = true
  try {
    await client.submitAnalysisFeedback(analysisId, rating.value, feedbackText.value)
    rating.value = null
    feedbackText.value = ""
    isEditing.value = false
    feedbacksLoaded.value = false
    await loadFeedbacks()
  } catch (e) {
    feedbacksError.value = e instanceof Error ? e.message : String(e)
  } finally {
    submitting.value = false
  }
}

watch(
  () => analysisId,
  () => {
    reset()
  }
)

watch(
  () => active,
  (isActive) => {
    if (isActive) void loadFeedbacks()
  },
  { immediate: true }
)

function formatCreatedAt(iso: string): string {
  return new Date(iso).toLocaleString()
}

function userLocalPart(email: string): string {
  const at = email.indexOf("@")
  return at === -1 ? email : email.slice(0, at)
}
</script>
