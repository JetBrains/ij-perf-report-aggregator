<!-- Copyright 2000-2020 JetBrains s.r.o. Use of this source code is governed by the Apache 2.0 license that can be found in the LICENSE file. -->
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

const DEFAULT_ACTIVE_TAB = "timeline"
export default defineComponent({
  name: "TabbedInfoCharts",
  components: {TimelineChart, ServiceTimelineChart, ActivityChart, StatsChart},
  setup() {
    const activeName = ref("")
    const charts = chartDescriptors.filter(it => it.isInfoChart === true)

    function updateLocation(location: RouteLocationNormalizedLoaded): void {
      const tab = location.query == null ? null : location.query["infoTab"]
      // do not check `location.path === "/"` because if component displayed, so, active
      activeName.value = tab == null ? DEFAULT_ACTIVE_TAB : tab as string
    }

    const route = useRoute()
    updateLocation(route)
    watch(route, location => {
      updateLocation(location)
    })
    return {
      activeName,
      navigate(): void {
        useRouter().push({
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
