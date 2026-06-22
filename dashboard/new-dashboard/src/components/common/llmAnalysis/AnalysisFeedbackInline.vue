<template>
  <div
    class="flex flex-col items-start gap-1 rounded-lg border px-3 py-1.5 text-sm font-normal transition-colors"
    :class="
      myFeedback
        ? 'border-emerald-300 bg-emerald-50 dark:border-emerald-500/40 dark:bg-emerald-500/10'
        : 'border-amber-300 bg-amber-50 dark:border-amber-500/40 dark:bg-amber-500/10'
    "
  >
    <span
      v-if="feedbacksLoading"
      class="text-xs text-gray-400"
    >
      Loading…
    </span>
    <template v-else>
      <div class="flex items-center gap-2">
        <i
          v-if="myFeedback"
          class="pi pi-check-circle text-emerald-500"
        />
        <span
          class="font-semibold whitespace-nowrap"
          :class="myFeedback ? 'text-emerald-700 dark:text-emerald-300' : 'text-amber-700 dark:text-amber-300'"
        >
          {{ myFeedback ? "Your rating" : "Was this analysis helpful?" }}
        </span>
        <div
          v-tooltip.top="starsDisabledReason"
          class="flex items-center"
          @mouseleave="hoveredStar = 0"
        >
          <button
            v-for="star in 5"
            :key="star"
            v-tooltip.top="ratingLabels[star - 1]"
            type="button"
            :aria-label="ratingLabels[star - 1]"
            :aria-pressed="myFeedback?.rate === star"
            class="cursor-pointer px-0.5 text-xl leading-none transition-transform hover:scale-125 disabled:cursor-not-allowed disabled:opacity-50 disabled:hover:scale-100"
            :class="star <= displayRating ? 'text-amber-400' : 'text-amber-400/40 hover:text-amber-400'"
            :disabled="starsDisabledReason != null || submitting"
            @mouseenter="hoveredStar = star"
            @click="rate(star)"
          >
            <i :class="star <= displayRating ? 'pi pi-star-fill' : 'pi pi-star'" />
          </button>
        </div>

        <Button
          v-if="myFeedback"
          :label="myFeedback.feedback ? 'edit comment' : 'add a comment'"
          :icon="myFeedback.feedback ? 'pi pi-pencil' : 'pi pi-comment'"
          severity="secondary"
          text
          size="small"
          @click="toggleComment"
        />

        <template v-if="otherFeedbacks.length > 0">
          <span class="text-gray-300">·</span>
          <Button
            severity="secondary"
            text
            size="small"
            @click="toggleOthers"
          >
            <span class="flex items-center gap-1 text-xs whitespace-nowrap">
              {{ otherFeedbacks.length === 1 ? "1 other rating" : `${otherFeedbacks.length} others` }}
              <i class="pi pi-star-fill text-amber-400" />
              {{ otherAverage }}
            </span>
          </Button>
        </template>
      </div>

      <div
        class="max-w-md text-xs leading-snug transition-colors"
        :class="activeLabel ? 'text-gray-700 dark:text-gray-200' : 'text-gray-400 dark:text-gray-500'"
      >
        {{ activeLabel ?? "Hover a star to see what each rating means" }}
      </div>

      <span
        v-if="feedbacksError"
        class="text-xs text-red-600"
      >
        {{ feedbacksError }}
      </span>
    </template>

    <Popover ref="commentPopover">
      <div class="flex w-80 flex-col gap-2">
        <Textarea
          v-model="commentText"
          rows="4"
          class="w-full"
          placeholder="What could work better, ideally the actually guilty commit(s)"
        />
        <div class="flex items-center gap-2">
          <Button
            label="Save"
            size="small"
            :loading="submitting"
            :disabled="submitting"
            @click="saveComment"
          />
          <Button
            label="Cancel"
            severity="secondary"
            size="small"
            :disabled="submitting"
            @click="closeComment"
          />
        </div>
      </div>
    </Popover>

    <Popover ref="othersPopover">
      <ul class="flex max-w-sm flex-col gap-1 text-xs">
        <li
          v-for="fb in otherFeedbacks"
          :key="fb.id"
          class="flex items-start gap-2"
        >
          <span class="font-medium whitespace-nowrap text-amber-500">★ {{ fb.rate }}</span>
          <span
            v-if="fb.userEmail"
            class="text-gray-500 whitespace-nowrap"
          >
            {{ userLocalPart(fb.userEmail) }}
          </span>
          <span
            v-if="fb.feedback"
            class="text-gray-700"
          >
            {{ fb.feedback }}
          </span>
        </li>
      </ul>
    </Popover>
  </div>
</template>

<script setup lang="ts">
import { storeToRefs } from "pinia"
import Popover from "primevue/popover"
import { computed, ref, watch } from "vue"
import { AnalysisFeedback, LlmAnalysisClient } from "./LlmAnalysisClient"
import { injectOrNull } from "../../../shared/injectionKeys"
import { serverConfiguratorKey } from "../../../shared/keys"
import { useUserStore } from "../../../shared/useUserStore"

const { analysisId, isTerminalState } = defineProps<{
  analysisId: number | string | null
  isTerminalState: boolean
}>()

const serverConfigurator = injectOrNull(serverConfiguratorKey)
const client = new LlmAnalysisClient(serverConfigurator)

const feedbacks = ref<AnalysisFeedback[]>([])
const feedbacksLoaded = ref(false)
const feedbacksLoading = ref(false)
const feedbacksError = ref<string | null>(null)
const submitting = ref(false)

const hoveredStar = ref(0)
const commentText = ref("")
const commentPopover = ref<InstanceType<typeof Popover> | null>(null)
const othersPopover = ref<InstanceType<typeof Popover> | null>(null)

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

const otherAverage = computed<string>(() => {
  const list = otherFeedbacks.value
  if (list.length === 0) return "0"
  const sum = list.reduce((acc, fb) => acc + fb.rate, 0)
  return (sum / list.length).toFixed(1)
})

const displayRating = computed<number>(() => hoveredStar.value || myFeedback.value?.rate || 0)

const ratingLabels = [
  "1 — Misleading: pointed at the wrong cause",
  "2 — Not useful: no actionable signal",
  "3 — Right direction: helped in investigation, but culprit not found",
  "4 — Close: culprit in list but ranked low or reasoning weak",
  "5 — Spot on: culprit identified with sound reasoning",
]

// Live caption shown under the stars: the meaning of the star currently hovered,
// falling back to the user's own rating, otherwise null (shows a hint).
const activeLabel = computed<string | null>(() => (displayRating.value > 0 ? ratingLabels[displayRating.value - 1] : null))

const starsDisabledReason = computed<string | null>(() => (isTerminalState ? null : "Available once analysis finishes"))

async function loadFeedbacks() {
  if (analysisId == null || feedbacksLoading.value) return
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

async function submit(rateValue: number, feedback: string | undefined): Promise<boolean> {
  if (analysisId == null || submitting.value) return false
  submitting.value = true
  feedbacksError.value = null
  try {
    await client.submitAnalysisFeedback(analysisId, rateValue, feedback)
    feedbacksLoaded.value = false
    await loadFeedbacks()
    return true
  } catch (e) {
    feedbacksError.value = e instanceof Error ? e.message : String(e)
    return false
  } finally {
    submitting.value = false
  }
}

async function rate(star: number) {
  if (starsDisabledReason.value != null) return
  await submit(star, myFeedback.value?.feedback)
}

function toggleComment(event: Event) {
  commentText.value = myFeedback.value?.feedback ?? ""
  commentPopover.value?.toggle(event)
}

function closeComment() {
  commentPopover.value?.hide()
}

async function saveComment() {
  const mine = myFeedback.value
  if (mine == null) return
  const ok = await submit(mine.rate, commentText.value)
  if (ok) commentPopover.value?.hide()
}

function toggleOthers(event: Event) {
  othersPopover.value?.toggle(event)
}

watch(
  () => analysisId,
  () => {
    feedbacks.value = []
    feedbacksLoaded.value = false
    feedbacksError.value = null
    hoveredStar.value = 0
    commentText.value = ""
    void loadFeedbacks()
  },
  { immediate: true }
)

function userLocalPart(email: string): string {
  const at = email.indexOf("@")
  return at === -1 ? email : email.slice(0, at)
}
</script>
