<template>
  <el-tabs
    v-model="activeName"
    @tab-click="navigate"
  >
    <el-tab-pane
      label="Timeline"
      name="timeline"
      lazy
    >
      <keep-alive>
        <TimelineChart />
      </keep-alive>
    </el-tab-pane>
    <el-tab-pane
      label="Service Timeline"
      name="serviceTimeline"
      lazy
    >
      <keep-alive>
        <ServiceTimelineChart />
      </keep-alive>
    </el-tab-pane>
    <!--  use v-once because `charts` is not going to be changed  -->
    <el-tab-pane
      v-for="item in charts"
      v-once
      :key="item.id"
      :label="item.label"
      :name="item.id"
      lazy
    >
      <keep-alive>
        <ActivityChart :type="item.id" />
      </keep-alive>
    </el-tab-pane>
    <el-tab-pane
      label="Stats"
      name="stats"
      lazy
    >
      <keep-alive>
        <StatsChart />
      </keep-alive>
    </el-tab-pane>
  </el-tabs>
</template>

<script lang="ts">
import { defineComponent, watch, ref } from "vue"
import { useRouter, useRoute, RouteLocationNormalizedLoaded } from "vue-router"
import ActivityChart from "./ActivityChart.vue"
import StatsChart from "./StatsChart.vue"
import { chartDescriptors } from "./charts/ActivityChartDescriptor"
import ServiceTimelineChart from "./timeline/ServiceTimelineChart.vue"
import TimelineChart from "./timeline/TimelineChart.vue"

export default defineComponent({
  name: "TabbedInfoCharts",
  components: {TimelineChart, ServiceTimelineChart, ActivityChart, StatsChart},
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
      navigate(): void {
        void router.push({
          query: {
            ...route.query,
            infoTab: activeName.value,
          },
        })
      },
      charts,
    }
  },
})
</script>
