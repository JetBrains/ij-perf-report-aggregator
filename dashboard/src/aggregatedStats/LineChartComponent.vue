<!-- Copyright 2000-2020 JetBrains s.r.o. Use of this source code is governed by the Apache 2.0 license that can be found in the LICENSE file. -->
<template>
  <el-popover
    placement="top"
    trigger="manual"
    v-model="infoIsVisible">
    <div>
      <div>
        <!-- cell text has 10px padding - so, add link margin to align text  -->
        <el-link v-if="reportName.length !== 0" style="margin-left: 10px" :href="reportLink" target="_blank" type="text">{{reportName}}
        </el-link>
        <small v-else style="margin-left: 10px">use <code>as is</code> granularity to see report</small>

        <el-link type="default"
                 style="float: right"
                 :underline="false"
                 icon="el-icon-close"
                 @click='infoIsVisible = false'/>
      </div>
      <el-table :data="reportTableData" :show-header="false">
        <el-table-column property="name" class-name="infoMetricName" min-width="180"/>
        <el-table-column property="value" align="right" class-name="infoMetricValue"/>
      </el-table>
    </div>
    <div slot="reference" v-loading="isLoading" class="aggregatedChart" ref="chartContainer"></div>
  </el-popover>
</template>

<style>
.el-table .cell {
  word-break: normal;
}

.infoMetricName {
  white-space: nowrap;
}

.infoMetricValue {
  font-family: monospace;
}
</style>

<script lang="ts">
import { Component, Prop, Watch } from "vue-property-decorator"
import { LineChartManager, SeriesDescriptor } from "@/aggregatedStats/LineChartManager"
import { ChartSettings } from "@/aggregatedStats/ChartSettings"
import { SortedByCategory, SortedByDate } from "@/aggregatedStats/ChartConfigurator"
import { DataQuery, DataQueryFilter, DataRequest, encodeQuery, getFilters, MetricDescriptor, Metrics, MetricType } from "@/aggregatedStats/model"
import { BaseStatChartComponent } from "@/aggregatedStats/BaseStatChartComponent"
import { DurationParseResult, parseTimeRange, toClickhouseSql } from "@/aggregatedStats/parseDuration"

const rison = require("rison-node")

@Component
export default class LineChartComponent extends BaseStatChartComponent<LineChartManager> {
  @Prop({type: String, required: true})
  type!: MetricType

  @Prop(String)
  order!: "date" | "buildNumber"

  @Prop({type: String})
  timeRange!: string

  @Prop()
  seriesManager!: LineChartSeriesManager | null

  metricDescriptors!: Array<MetricDescriptor>

  @Watch("chartSettings.showScrollbarXPreview")
  showScrollbarXPreviewChanged(): void {
    const chartManager = this.chartManager
    if (chartManager != null) {
      chartManager.scrollbarXPreviewOptionChanged(this.chartSettings)
    }
  }

  infoIsVisible: boolean = false
  reportTableData: Array<any> = []
  reportName: string = ""

  reportLink: string | null = null

  @Watch("chartSettings.granularity")
  granularityChanged() {
    this.loadDataAfterDelay()
  }

  @Watch("timeRange")
  timeRangeChanged(value: string, oldValue: string) {
    console.info(`timeRange changed (${oldValue} => ${value})`)
    this.loadDataAfterDelay()
  }

  protected reloadData(request: DataRequest) {
    const manager = this.seriesManager || DefaultLineChartSeriesManager.INSTANCE
    manager.loadAndRender(this, request)
  }

  protected createChartManager() {
    const metricDescriptors: Array<MetricDescriptor> = this.metrics.map(key => {
      const metricPathEndDotIndex = key.indexOf(".")
      let name: string
      let effectiveKey = key
      if (metricPathEndDotIndex == -1) {
        name = keyToName(key)
      }
      else {
        name = key.substring(metricPathEndDotIndex + 1)
        if (name.length > 2 && name[name.length - 2] == ".") {
          name = name.substring(0, name.length - 2)
        }
        effectiveKey = name.replace(/ /g, "_")
      }

      return {
        key: effectiveKey,
        name,
        hiddenByDefault: false,
      }
    })
    Object.seal(metricDescriptors)
    this.metricDescriptors = metricDescriptors

    const configurator = this.order === "date" ? new SortedByDate(data => {
      if (data == null) {
        this.infoIsVisible = false
        return
      }

      const tableData = []
      for (const metricDescriptor of metricDescriptors) {
        tableData.push({
          name: metricDescriptor.name,
          value: data[metricDescriptor.key],
        })
      }

      const request = this.dataRequest!!
      const reportQuery: DataQuery = {
        db: request.db,
        filters: getFilters(request).concat([{field: "generated_time", value: data.t / 1000}]),
      }

      if (this.chartSettings.granularity === "as is") {
        const reportUrl = `/api/v1/report/${encodeQuery(reportQuery)}`
        this.reportLink = `/#/report?reportUrl=${encodeURIComponent(this.chartSettings.serverUrl)}${reportUrl}`
        const generatedTime = new Date(data.t)
        // 18 Oct, 13:01:49
        this.reportName = `${generatedTime.getDate()} ${generatedTime.toLocaleString("default", {month: "short"})}, ${generatedTime.toLocaleTimeString("default", {hour12: false})}`
      }
      else {
        this.reportName = ""
      }

      this.reportTableData = tableData
      this.infoIsVisible = true
    }) : new SortedByCategory()
    const chartManager = new LineChartManager(this.$refs.chartContainer as HTMLElement, this.chartSettings || new ChartSettings(), this.type, configurator)

    chartManager.chart.exporting.menu!!.items[0]!!.menu!!.push({
      label: "Open",
      type: "custom",
      options: {
        callback: () => {
          const configuration = rison.encode({
            chartSettings: this.chartSettings,
            metrics: this.metrics,
            order: this.order,
            dataRequest: this.dataRequest,
          })
          window.open("/#/line-chart/" + configuration, "_blank")
        }
      }
    })

    return chartManager
  }
}

export function createDataQueryWithoutFields(request: DataRequest,
                                             filters: Array<DataQueryFilter>,
                                             component: LineChartComponent,
                                             timeRange: DurationParseResult,
                                             chartSettings: ChartSettings): DataQuery {
  const dataQuery: DataQuery = {
    db: request.db,
    filters: filters.concat([
      {field: "generated_time", sql: `> ${toClickhouseSql(timeRange)}`}
    ]),
  }

  let granularity = chartSettings.granularity
  if (granularity == null) {
    granularity = "2 hour"
  }

  if (component.order === "buildNumber") {
    dataQuery.dimensions = [
      {name: "build_c1"},
      {name: "build_c2"},
      {name: "build_c3"},
    ]

    dataQuery.fields = [{name: "t", sql: `toUnixTimestamp(anyHeavy(generated_time)) * 1000`}]
    dataQuery.order = dataQuery.dimensions.map(it => it.name)
    dataQuery.aggregator = "medianTDigest"
  }
  else {
    if (granularity === "as is") {
      dataQuery.fields = [{
        name: "t",
        sql: `toUnixTimestamp(generated_time) * 1000`
      }, "build_c1", "build_c2", "build_c3"]
    }
    else {
      dataQuery.aggregator = "medianTDigest"

      dataQuery.dimensions = [
        {name: "t", sql: getTimeDimension(granularity)}
      ]

      dataQuery.fields = ["build_c1", "build_c2", "build_c3"].map(it => ({name: it, sql: `anyHeavy(${it})`}))
    }

    dataQuery.order = ["t"]
  }
  return dataQuery
}

function getTimeDimension(granularity: "2 hour" | "day" | "week" | "month") {
  let sql = "toStartOfInterval(generated_time, interval "
  // hour - backward compatibility
  if (granularity == null || granularity === "2 hour" || granularity === "hour" as any) {
    sql += "2 hour"
  }
  else if (granularity === "day") {
    sql += "1 day"
  }
  else if (granularity === "week") {
    sql += "1 week"
  }
  else {
    sql += "1 month"
  }
  sql += ")"
  return `toUnixTimestamp(${sql}) * 1000`
}

class DefaultLineChartSeriesManager implements LineChartSeriesManager{
  static INSTANCE = new DefaultLineChartSeriesManager()

  constructor() {
    Object.seal(this)
  }

  loadAndRender(component: LineChartComponent, request: DataRequest): void {
    const timeRange = parseTimeRange(component.timeRange)
    const chartSettings = component.chartSettings
    const dataQuery = createDataQueryWithoutFields(request, getFilters(request), component, timeRange, chartSettings)
    dataQuery.fields!!.push(...component.metrics)

    component.loadData(`${chartSettings.serverUrl}/api/v1/metrics/${encodeQuery(dataQuery)}`, (data: Array<Metrics>, chartManager: LineChartManager) => {
      chartManager.render(data, DefaultLineChartSeriesManager.getSeries(component.metricDescriptors))
    })
  }

  private static getSeries(metricDescriptors: Array<MetricDescriptor>): Array<SeriesDescriptor> {
    return metricDescriptors.map(it => {
      return {
        name: getSeriesName(it),
        dataField: it.key,
        hiddenByDefault: it.hiddenByDefault,
      }
    })
  }
}

export interface LineChartSeriesManager {
  loadAndRender(component: LineChartComponent, request: DataRequest): void
}

function getSeriesName(metric: MetricDescriptor) {
  const metricPathEndDotIndex = metric.name.indexOf(".")
  return metricPathEndDotIndex == -1 ? metric.name : metric.name.substring(metricPathEndDotIndex + 1)
}

function keyToName(key: string) {
  if (key == "pluginDescriptorInitV18_d") {
    return "pluginDescriptorInit"
  }
  else if (key == "appStarter_d") {
    return "licenseCheck"
  }
  else {
    return (key.endsWith("_d") || key.endsWith("_i")) ? key.substring(0, key.length - 2) : key
  }
}
</script>