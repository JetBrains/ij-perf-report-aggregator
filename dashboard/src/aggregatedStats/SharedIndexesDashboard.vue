<!-- Copyright 2000-2020 JetBrains s.r.o. Use of this source code is governed by the Apache 2.0 license that can be found in the LICENSE file. -->
<template>
  <div>
    <el-row>
      <el-col :span="18">
        <el-form :inline="true" size="small">
          <el-form-item label="Server">
            <el-input
                data-lpignore="true"
                placeholder="Enter the aggregated stats server URL..."
                v-model="chartSettings.serverUrl"/>
          </el-form-item>

          <el-form-item label="Scenario">
            <el-select v-model="selectedProjects" multiple collapse-tags style="min-width: 500px;">
              <el-option v-for="project in projects" :key="project" :label="project" :value="project"/>
            </el-select>
          </el-form-item>

          <el-form-item label="Machine">
            <el-cascader
                v-model="chartSettings.selectedMachine"
                :show-all-levels="false"
                :props='{"label": "name", value: "name", checkStrictly: true, emitPath: false}'
                :options="machines"/>
          </el-form-item>

          <el-form-item>
            <el-button title="Updated automatically, but you can force data reloading"
                       type="primary"
                       icon="el-icon-refresh"
                       :loading="isFetching"
                       @click="loadData"/>
          </el-form-item>
        </el-form>
      </el-col>
      <el-col :span="6">
        <div style="float: right">
          <el-checkbox size="small" v-model="chartSettings.showScrollbarXPreview">Show horizontal scrollbar preview</el-checkbox>
        </div>
      </el-col>
    </el-row>

    <el-tabs value="date" size="small">
      <el-tab-pane v-for='item in [
          {name: "By Date", order: "date"},
          {name: "By Build", order: "buildNumber"}
          ]' :key="item.order" :label="item.name" :name="item.order" lazy>
        <keep-alive>
          <div>
            <el-form v-if="item.order === 'date'" :inline="true" size="small">
              <el-form-item label="Granularity">
                <el-select v-model="chartSettings.granularity" data-lpignore="true" filterable>
                  <el-option v-for='name in ["as is", "2 hour", "day", "week", "month"]' :key="name" :label="name" :value="name"/>
                </el-select>
              </el-form-item>
              <el-form-item label="Period">
                <el-select v-model="timeRange" data-lpignore="true" filterable>
                  <el-option v-for='item in timeRanges' :key="item.k" :label="item.l" :value="item.k"/>
                </el-select>
              </el-form-item>
            </el-form>

            <el-row :gutter="5">
              <el-col :span="12">
                <el-card shadow="never" :body-style="{ padding: '0px' }">
                  <LineChartComponent type="duration" :order="item.order" :dataRequest="dataRequest" :timeRange="timeRange"
                                      :metrics='["metrics.scanning.d"]'
                                      :seriesManager="seriesManager"
                                      :chartSettings="chartSettings"/>
                </el-card>
              </el-col>
              <el-col :span="12">
                <el-card shadow="never" :body-style="{ padding: '0px' }">
                  <LineChartComponent type="duration" :order="item.order" :dataRequest="dataRequest" :timeRange="timeRange"
                                      :metrics='["metrics.indexing.d"]'
                                      :seriesManager="seriesManager"
                                      :chartSettings="chartSettings"/>
                </el-card>
              </el-col>
            </el-row>
          </div>
        </keep-alive>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script lang="ts">
import { Component } from "vue-property-decorator"
import LineChartComponent from "./LineChartComponent.vue"
import ClusteredChartComponent from "./ClusteredChartComponent.vue"
import { AggregatedStatsPage } from "./AggregatedStatsPage"
import { MultiValueFilter } from "@/aggregatedStats/ValueFilter"
import { DataQueryDimension, DataQueryFilter, DataRequest, encodeQueries, MetricDescriptor, Metrics } from "@/aggregatedStats/model"
import { LineChartManager, SeriesDescriptor } from "@/aggregatedStats/LineChartManager"
import { parseTimeRange } from "@/aggregatedStats/parseDuration"
import { createDataQueryWithoutFields, LineChartSeriesManager } from "@/aggregatedStats/LineChartComponent.vue"

@Component({
  components: {LineChartComponent, ClusteredChartComponent}
})
export default class SharedIndexesDashboard extends AggregatedStatsPage {
  protected projectFilterManager = new MultiValueFilter("selectedProjects")

  readonly seriesManager = new SharedIndexesLineChartSeriesManager()

  protected getDbName(): string {
    return "sharedIndexes"
  }
}

interface SeriesProjectFilterFactory {
  seriesLabel: string
  filter: (projectName: string) => boolean
}

const seriesFactories: Array<SeriesProjectFilterFactory> = [
  {seriesLabel: "usual", filter: name => name.includes("-usual-")},
  {seriesLabel: "shared", filter: name => name.includes("-shared-index")},
]

class SharedIndexesLineChartSeriesManager implements LineChartSeriesManager{
  constructor() {
    Object.seal(this)
  }

  loadAndRender(component: LineChartComponent, request: DataRequest): void {
    const timeRange = parseTimeRange(component.timeRange)
    const chartSettings = component.chartSettings

    const projectFilters: Array<{
      projects: Array<string>
      label: string
    }> = seriesFactories
      .map(it => ({projects: request.projects.filter(it.filter), label: it.seriesLabel}))
      .filter(it => it.projects.length > 0)
    if (projectFilters.length === 0) {
      console.error("project selection is invalid - no candidates")
      return
    }

    const queries = projectFilters.map((projectFilter, index) => {
      const query = createDataQueryWithoutFields(request, getFilters(request, projectFilter.projects), component, timeRange, chartSettings)
      query.fields!!.push(...metricToFields(component.metrics, index))
      return query
    })

    const seriesDescriptors = projectFilters
      .map((projectFilter, index) => SharedIndexesLineChartSeriesManager.getSeries(component.metricDescriptors, index, it => `${it} ${projectFilter.label}`))
      .flat()

    component.loadData(`${chartSettings.serverUrl}/api/v1/metrics/${encodeQueries(queries)}`, (data: Array<Array<Metrics>>, chartManager: LineChartManager) => {
      // data here is array of two data series - merge it
      const result: Map<number, Metrics> = new Map<number, Metrics>()
      for (const items of data) {
        for (const item of items) {
          // time is unique in our case (not across different OS (machine grouped by OS and hardware))
          const key = item.t
          const existingItem = result.get(key)
          if (existingItem === undefined) {
            result.set(key, item)
          }
          else {
            Object.assign(existingItem, item)
          }
        }
      }

      chartManager.render(Array.from(result.values()).sort((a, b) => a.t - b.t), seriesDescriptors)
    })
  }

  private static getSeries(metricDescriptors: Array<MetricDescriptor>, index: number, nameProducer: (name: string) => string): Array<SeriesDescriptor> {
    return metricDescriptors.map(it => {
      return {
        name: nameProducer(it.name),
        dataField: `${it.key}_${index}`,
        hiddenByDefault: it.hiddenByDefault,
      }
    })
  }
}

function metricToFields(metrics: Array<string>, index: number): Array<DataQueryDimension> {
  return metrics.map(it => {
    return {
      name: it,
      resultKey: `${it.substring(it.indexOf(".") + 1, it.lastIndexOf("."))}_${index}`
    }
  })
}

function getFilters(request: DataRequest, projects: Array<string>): Array<DataQueryFilter> {
  const result: Array<DataQueryFilter> = []
  result.push({field: "project", value: projects})
  result.push({field: "machine", value: request.machine})
  return result
}

</script>

<!--suppress CssUnusedSymbol -->
<style>
.aggregatedChart {
  width: 100%;
  height: 300px;
}

.dividerAfterForm {
  margin-top: 0 !important;
}
</style>
