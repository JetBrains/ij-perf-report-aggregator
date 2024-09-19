<template>
  <TabGroup
    :selected-index="activeId"
    @change="navigate"
  >
    <div class="border-b">
      <TabList class="-mb-px flex space-x-8">
        <Tab
          v-for="item in charts"
          :key="item.id"
          v-slot="{ selected }"
          as="template"
        >
          <button :class="[selected ? 'border-indigo-500 text-indigo-600' : 'border-transparent hover:border-gray-300', 'whitespace-nowrap py-4 px-1 border-b-2 font-medium']">
            {{ item.label }}
          </button>
        </Tab>
      </TabList>
    </div>

    <TabPanels>
      <TabPanel
        v-for="item in charts"
        :key="item.id"
        class="p-3"
        :unmount="false"
      >
        <ActivityChart :descriptor="item" />
      </TabPanel>
    </TabPanels>
  </TabGroup>
</template>

<script setup lang="ts">
import { Tab, TabGroup, TabList, TabPanel, TabPanels } from "@headlessui/vue"
import { asyncScheduler, distinctUntilChanged, observeOn, Subject, switchMap } from "rxjs"
import { computed, ref, watch } from "vue"
import { RouteLocationNormalizedLoaded, useRoute, useRouter } from "vue-router"
import ActivityChart from "./ActivityChart.vue"
import { chartDescriptors } from "./ActivityChartDescriptor"

const { isInfoChart = false } = defineProps<{
  isInfoChart: boolean
}>()

const charts = chartDescriptors.filter((it) => it.isInfoChart === isInfoChart || (!isInfoChart && it.isInfoChart === undefined))

const activeId = ref(0)
const activeName = computed(() => {
  return charts[activeId.value].id
})
const queryParamName = isInfoChart ? "infoTab" : "tab"

function updateLocation(location: RouteLocationNormalizedLoaded): void {
  const tab = location.query[queryParamName]
  // do not check `location.path === "/"` because if component displayed, so, active
  activeId.value =
    tab == null
      ? 0
      : Math.max(
          charts.findIndex((it) => it.id === tab),
          0
        )
}

const route = useRoute()
updateLocation(route)
watch(route, (location) => {
  updateLocation(location)
})

const subject = new Subject<string>()
subject
  .pipe(
    distinctUntilChanged(),
    switchMap((activeName) => {
      return router.push({
        query: {
          ...route.query,
          [queryParamName]: activeName,
        },
      })
    }),
    observeOn(asyncScheduler)
  )
  .subscribe(() => {
    // empty
  })

const router = useRouter()

function navigate(index: number) {
  activeId.value = index
  subject.next(activeName.value)
}
</script>
