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

          <el-form-item label="Scenarios">
            <el-select v-model="selectedProjects" multiple collapse-tags style="min-width: 500px;">
              <el-option v-for="project in projects" :key="project" :label="project" :value="project"/>
            </el-select>
          </el-form-item>

          <el-form-item label="Metrics">
            <el-select v-model="selectedMetrics" multiple collapse-tags style="min-width: 200px;">
              <el-option v-for="metric in metrics" :key="metric.name" :label="metric.name" :value="metric.name"/>
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

            <template v-for='metric in selectedMetrics'>
              <el-row :gutter="5" :key="metric">
                <el-col :span="20">
                  <el-card shadow="never" :body-style="{ padding: '0px' }">
                    <LineChartComponent :type='metric.endsWith(".d") ? "duration" : metric.endsWith(".c") ? "counter" : "duration"'
                                        :order="item.order"
                                        :dataRequest="dataRequest"
                                        :timeRange="timeRange"
                                        :metrics='[metric]'
                                        :seriesManager="seriesManager"
                                        :chartSettings="chartSettings"/>
                  </el-card>
                </el-col>
              </el-row>
            </template>
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
import { DataQueryDimension, DataQueryFilter, DataRequest, encodeQueries, Metrics } from "@/aggregatedStats/model"
import { LineChartManager, SeriesDescriptor } from "@/aggregatedStats/LineChartManager"
import { parseTimeRange } from "@/aggregatedStats/parseDuration"
import { createDataQueryWithoutFields, LineChartSeriesManager } from "@/aggregatedStats/LineChartComponent.vue"

@Component({
  components: {LineChartComponent, ClusteredChartComponent},
})
export default class SharedIndexesDashboard extends AggregatedStatsPage {
  protected projectFilterManager = new MultiValueFilter("selectedProjects")

  readonly seriesManager = new SharedIndexesLineChartSeriesManager()

  protected getDbName(): string {
    return "sharedIndexes"
  }

  selectedMetrics: Array<string> = []
}

class SharedIndexesLineChartSeriesManager implements LineChartSeriesManager {
  constructor() {
    Object.seal(this)
  }

  loadAndRender(component: LineChartComponent, request: DataRequest): void {
    const timeRange = parseTimeRange(component.timeRange)
    const chartSettings = component.chartSettings

    if (component.metrics.length != 1 || component.metricDescriptors.length != 1) {
      throw Error("Only one metric is allowed")
    }
    const metric = component.metrics[0]
    const metricDescriptor = component.metricDescriptors[0]

    const chartDescriptors = request.projects.map((experimentName, index) => {
      const metricKey = metricDescriptor.key + "_" + index
      const seriesDescriptor: SeriesDescriptor = {
        name: metricDescriptor.key + " " + experimentName,
        dataField: metricKey,
        hiddenByDefault: metricDescriptor.hiddenByDefault,
      }
      const filters: Array<DataQueryFilter> = [
        {field: "project", value: experimentName},
        {field: "machine", value: request.machine},
      ]

      if (chartSettings.granularity === "as is") {
        // Exclude reports that miss this metric key. Otherwise clickhouse would select all rows and with 0 as the metric value.
        filters.push({field: metricKey, operator: '>', value: 0})
      }

      const query = createDataQueryWithoutFields(request, filters, component, timeRange, chartSettings)
      const field: DataQueryDimension = {
        name: metric,
        resultKey: metricKey,
      }
      query.fields!!.push(field)
      return {
        query: query,
        seriesDescriptor: seriesDescriptor,
      }
    })

    component.loadData(`${chartSettings.serverUrl}/api/v1/metrics/${encodeQueries(chartDescriptors.map(e => e.query))}`, (rawData: any, chartManager: LineChartManager) => {
      const data = (chartDescriptors.length == 1 ? [rawData as Array<Metrics>] : rawData as Array<Array<Metrics>>).flat()
      chartManager.render(data.sort((a, b) => a.t - b.t), chartDescriptors.map(e => e.seriesDescriptor))
    })
  }
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
