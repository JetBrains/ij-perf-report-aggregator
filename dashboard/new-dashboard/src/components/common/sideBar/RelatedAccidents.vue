<template>
  <DeferredContent @load="loadEventsAroundDate">
    <Accordion
      v-if="accidentsAroundDate?.length ?? 0 > 0"
      :active-index="0"
    >
      <AccordionTab header="Events around the date">
        <ul
          v-if="accidentsAroundDate"
          class="gap-1.5 text-sm overflow-y-auto max-h-80"
        >
          <li
            v-for="accident in accidentsAroundDate"
            :key="accident?.reason + accident?.kind"
            ref="circle"
          >
            <span
              class="flex gap-1.5 text-sm"
              v-tooltip.left="{
                value: getAffectedTests(accident),
                autoHide: false,
                showDelay: 500,
              }"
            >
              <span v-if="props.inDialog">
                <DocumentDuplicateIcon
                  class="w-4 h-4"
                  @click="copy(accident)"
                />
              </span>
              &bull;
              <!-- eslint-disable vue/no-v-html -->
              <span
                class="w-full"
                :class="accident.kind == 'Regression' ? 'text-red-500' : 'text-green-500'"
                v-html="replaceToLink(accident.reason)"
              />
            </span>
            <!-- eslint-enable -->
          </li>
        </ul>
      </AccordionTab>
    </Accordion>
  </DeferredContent>
</template>
<script setup lang="ts">
import { computedAsync } from "@vueuse/core"
import { ref } from "vue"
import { Accident, AccidentsConfigurator } from "../../../configurators/AccidentsConfigurator"
import { replaceToLink } from "../../../util/linkReplacer"
import { InfoData } from "./InfoSidebar"

const props = defineProps<{
  data: InfoData | null
  accidentsConfigurator: AccidentsConfigurator | null
  inDialog: boolean
}>()

function getAffectedTests(accident: AccidentSimple): string {
  if (props.inDialog) return ""
  return accident.affectedTests.filter((value) => value !== "").join("\n")
}

const accidentsAroundDate = ref<AccidentSimple[] | undefined>([])

interface AccidentSimple {
  kind: string
  reason: string
  affectedTests: string[]
}

function deduplicateAccidents(accidents: Accident[]): AccidentSimple[] {
  const accidentMap = new Map<string, AccidentSimple>()

  accidents.forEach((accident) => {
    const key = `${accident.kind}|${accident.reason.trim()}`
    if (accidentMap.has(key)) {
      const existingAccident = accidentMap.get(key) as AccidentSimple
      existingAccident.affectedTests = existingAccident.affectedTests.concat(accident.affectedTest)
    } else {
      accidentMap.set(key, { kind: accident.kind, reason: accident.reason.trim(), affectedTests: [accident.affectedTest] })
    }
  })
  return Array.from(accidentMap.values())
}

function loadEventsAroundDate() {
  computedAsync(async () => {
    if (props.data) {
      const accidents = (await props.accidentsConfigurator?.getAccidentsAroundDate(props.data.date)) ?? []
      accidentsAroundDate.value = deduplicateAccidents(accidents)
    }
  }).value
}

const emit = defineEmits(["copyAccident"])

function copy(accident: AccidentSimple): void {
  emit("copyAccident", accident)
}
</script>
<style #scoped>
.p-accordion .p-accordion-header .p-accordion-header-link {
  padding: 0 0 1rem 0;
}

.p-accordion .p-accordion-content {
  padding: 0;
}
</style>
