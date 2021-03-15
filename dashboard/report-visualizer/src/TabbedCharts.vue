<template>
  <el-tabs
    v-model="activeName"
    @tab-click="navigate"
  >
    <el-tab-pane
      v-for="item in charts"
      :key="item.id"
      :label="item.label"
      :name="item.id"
      lazy
    >
      <keep-alive>
        <ActivityChart :descriptor="item" />
      </keep-alive>
    </el-tab-pane>
  </el-tabs>
</template>

<script lang="ts">
import { DebouncedTask } from "shared/src/util/debounce"
import { defineComponent, ref, watch } from "vue"
import { RouteLocationNormalizedLoaded, useRoute, useRouter } from "vue-router"
import ActivityChart from "./ActivityChart.vue"
import { chartDescriptors } from "./charts/ActivityChartDescriptor"

export default defineComponent({
  name: "TabbedCharts",
  components: {ActivityChart},
  props: {
    isInfoChart: {
      type: Boolean,
      required: true,
    }
  },
  setup(props) {
    const charts = chartDescriptors.filter(it => it.isInfoChart === props.isInfoChart || (!props.isInfoChart && it.isInfoChart === undefined))
    const activeName = ref(charts[0].id)
    const queryParamName = props.isInfoChart ? "infoTab" : "tab"
    function updateLocation(location: RouteLocationNormalizedLoaded): void {
      const tab = location.query[queryParamName]
      // do not check `location.path === "/"` because if component displayed, so, active
      if (tab == null) {
        activeName.value = charts[0].id
      }
      else {
        const descriptor = charts.find(it => it.id === tab)
        activeName.value = descriptor === undefined ? charts[0].id : descriptor.id
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
      navigate: new DebouncedTask(function () {
        return router.push({
          query: {
            ...route.query,
            [queryParamName]: activeName.value,
          },
        })
      }, 0).executeFunctionReference,
    }
  }
})
</script>
