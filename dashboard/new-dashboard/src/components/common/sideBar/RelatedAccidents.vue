<template>
  <DeferredContent @load="loadEventsAroundDate">
    <Accordion
      v-if="accidentsAroundDate?.length ?? 0 > 0"
      :active-index="0"
    >
      <AccordionTab header="Events around (± 1 day)">
        <ul
          v-if="accidentsAroundDate"
          class="gap-1.5 text-sm break-all"
        >
          <li
            v-for="accident in accidentsAroundDate"
            :key="accident?.reason + accident?.kind"
            ref="circle"
          >
            <span
              v-tooltip.left="{
                value: getTooltipText(accident),
                autoHide: false,
                showDelay: 500,
              }"
              class="flex gap-1.5 text-sm"
            >
              <span v-if="inDialog">
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
              <GlobeAltIcon
                v-if="getAffectedTests(accident) == ''"
                class="w-4 h-4 flex-none"
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

const { data, accidentsConfigurator, inDialog } = defineProps<{
  data: InfoData | null
  accidentsConfigurator: AccidentsConfigurator | null
  inDialog: boolean
}>()

function getTooltipText(accident: AccidentSimple): string {
  if (inDialog) return ""
  return accident.date.split(" ")[0] + "\n" + getAffectedTests(accident) + "\n" + (accident.userName != "" ? "Created by " + accident.userName : "")
}
function getAffectedTests(accident: AccidentSimple): string {
  return accident.affectedTests.filter((value) => value !== "").join("\n")
}

const accidentsAroundDate = ref<AccidentSimple[] | undefined>([])

interface AccidentSimple {
  kind: string
  reason: string
  affectedTests: string[]
  date: string
  userName: string
}

function deduplicateAccidents(accidents: Accident[]): AccidentSimple[] {
  const accidentMap = new Map<string, AccidentSimple>()

  for (const accident of accidents) {
    const key = `${accident.kind}|${accident.reason.trim()}`
    if (accidentMap.has(key)) {
      const existingAccident = accidentMap.get(key) as AccidentSimple
      existingAccident.affectedTests = [...existingAccident.affectedTests, accident.affectedTest]
    } else {
      accidentMap.set(key, { kind: accident.kind, reason: accident.reason.trim(), affectedTests: [accident.affectedTest], date: accident.date, userName: accident.userName })
    }
  }
  return [...accidentMap.values()]
}

function loadEventsAroundDate() {
  computedAsync(async () => {
    if (data) {
      const accidents = (await accidentsConfigurator?.getAccidentsAroundDate(data.date)) ?? []
      accidentsAroundDate.value = deduplicateAccidents(accidents)
    }
  })
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
