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
        <ActivityChart :type="item.id" />
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
  name: "TabbedInfoCharts",
  components: {ActivityChart},
  setup() {
    const activeName = ref("")
    const charts = chartDescriptors.filter(it => it.isInfoChart === true)

    function updateLocation(location: RouteLocationNormalizedLoaded): void {
      const tab = location.query["infoTab"]
      // do not check `location.path === "/"` because if component displayed, so, active
      activeName.value = tab == null ? "timeline" : tab as string
    }

    const route = useRoute()
    updateLocation(route)
    watch(route, location => {
      updateLocation(location)
    })
    const router = useRouter()
    return {
      activeName,
      navigate: new DebouncedTask(function () {
        return router.push({
          query: {
            ...route.query,
            infoTab: activeName.value,
          },
        })
      }, 0).executeFunctionReference,
      charts,
    }
  },
})
</script>
