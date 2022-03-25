<template>
  <TabView
    v-model:active-index="activeId"
    lazy
    @tab-click="navigate"
  >
    <TabPanel
      v-for="item in charts"
      :key="item.id"
      :header="item.label"
    >
      <keep-alive>
        <ActivityChart :descriptor="item" />
      </keep-alive>
    </TabPanel>
  </TabView>
</template>

<script lang="ts">
import { DebouncedTask } from "shared/src/util/debounce"
import { computed, defineComponent, ref, watch } from "vue"
import { RouteLocationNormalizedLoaded, useRoute, useRouter } from "vue-router"
import ActivityChart from "./ActivityChart.vue"
import { chartDescriptors } from "./ActivityChartDescriptor"

export default defineComponent({
  name: "TabbedCharts",
  components: {ActivityChart},
  props: {
    isInfoChart: {
      type: Boolean,
      required: true,
    },
  },
  setup(props) {
    const charts = chartDescriptors.filter(it => it.isInfoChart === props.isInfoChart || (!props.isInfoChart && it.isInfoChart === undefined))
    const activeId = ref(0)
    const activeName = computed(() => {
      return charts[activeId.value].id
    })
    const queryParamName = props.isInfoChart ? "infoTab" : "tab"

    function updateLocation(location: RouteLocationNormalizedLoaded): void {
      const tab = location.query[queryParamName]
      // do not check `location.path === "/"` because if component displayed, so, active
      if (tab == null) {
        activeId.value = 0
      }
      else {
        activeId.value = charts.findIndex(it => it.id === tab)
      }
    }

    const route = useRoute()
    updateLocation(route)
    watch(route, location => {
      updateLocation(location)
    })

    const router = useRouter()
    return {
      charts,
      activeName,
      activeId,
      navigate: new DebouncedTask(function () {
        return router.push({
          query: {
            ...route.query,
            [queryParamName]: activeName.value,
          },
        })
      }, 0).executeFunctionReference,
    }
  },
})
</script>
