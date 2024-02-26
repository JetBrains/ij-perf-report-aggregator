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
          >
            <span class="flex gap-1.5 text-sm">
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
import { AccidentsConfigurator } from "../../../configurators/AccidentsConfigurator"
import { replaceToLink } from "../../../util/linkReplacer"
import { InfoData } from "./InfoSidebar"

const props = defineProps<{
  data: InfoData | null
  accidentsConfigurator: AccidentsConfigurator | null
  inDialog: boolean
}>()

const accidentsAroundDate = ref<AccidentSimple[] | undefined>([])
interface AccidentSimple {
  kind: string
  reason: string
}
function deduplicateAccidents(accidents: AccidentSimple[]): AccidentSimple[] {
  const uniqueJson = [...new Set(accidents.map((accident) => JSON.stringify(accident)))]
  return uniqueJson.map((json) => JSON.parse(json) as AccidentSimple)
}

function loadEventsAroundDate() {
  computedAsync(async () => {
    if (props.data) {
      const accidents = (await props.accidentsConfigurator?.getAccidentsAroundDate(props.data.date)) ?? []
      const transformedAccidents = accidents.map((accident) => ({ kind: accident.kind, reason: accident.reason.trim() }))
      accidentsAroundDate.value = deduplicateAccidents(transformedAccidents)
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
