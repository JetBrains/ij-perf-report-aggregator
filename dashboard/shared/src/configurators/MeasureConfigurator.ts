import { LineSeriesOption } from "echarts/charts"
import { DatasetOption } from "echarts/types/dist/shared"
import { DimensionDefinition } from "echarts/types/src/util/types"
import { Ref, shallowRef, watch } from "vue"
import { DataQueryResult } from "../DataQueryExecutor"
import { PersistentStateManager } from "../PersistentStateManager"
import { ChartConfigurator, ChartOptions } from "../chart"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration, encodeQuery, toMutableArray } from "../dataQuery"
import { DebouncedTask, TaskHandle } from "../util/debounce"
import { loadJson } from "../util/httpUtil"
import { DimensionConfigurator } from "./DimensionConfigurator"
import { ServerConfigurator } from "./ServerConfigurator"

// natural sort of alphanumerical strings
export const collator = new Intl.Collator(undefined, {numeric: true, sensitivity: "base"})

export class MeasureConfigurator implements DataQueryConfigurator, ChartConfigurator {
  public readonly data = shallowRef<Array<string>>([])
  public readonly value = shallowRef<Array<string>>([])

  private readonly debouncedLoadMetadata = new DebouncedTask(taskHandle => this.loadMetadata(taskHandle))

  constructor(private readonly serverConfigurator: ServerConfigurator,
              persistentStateManager: PersistentStateManager,
              private readonly parent: DimensionConfigurator | null = null,
              readonly skipZeroValues: boolean = true) {
    persistentStateManager.add("measure", this.value)

    if (this.parent != null) {
      watch(this.parent.value, this.debouncedLoadMetadata.executeFunctionReference)
    }
  }

  configureQuery(query: DataQuery, configuration: DataQueryExecutorConfiguration): boolean {
    const measureNames = toMutableArray(this.value.value)
    if (measureNames.length === 0) {
      return false
    }

    configureQuery(measureNames, query, this.skipZeroValues)
    configuration.measures = measureNames
    configuration.chartConfigurator = this
    return true
  }

  configureChart(data: DataQueryResult, configuration: DataQueryExecutorConfiguration): ChartOptions {
    return configureChart(configuration, data)
  }

  scheduleLoadMetadata(immediately: boolean): void {
    this.debouncedLoadMetadata.execute(immediately)
  }

  loadMetadata(taskHandle: TaskHandle): Promise<unknown> {
    const parent = this.parent
    if (parent == null) {
      const server = this.serverConfigurator.value.value
      if (server == null || server.length === 0) {
        return Promise.resolve()
      }

      return loadJson<Array<string>>(`${server}/api/v1/meta/measure?db=${this.serverConfigurator.databaseName}`, null, taskHandle, data => {
        this.data.value = data
      })
    }

    const query = new DataQuery()
    const configuration = new DataQueryExecutorConfiguration()
    if (!parent.configureQuery(query, configuration) || !parent.serverConfigurator.configureQuery(query, configuration)) {
      return Promise.resolve()
    }

    // "group by" is equivalent of distinct (https://clickhouse.tech/docs/en/sql-reference/statements/select/distinct/#alternatives)
    query.addDimension({name: "measures", subName: "name"})
    query.order = ["measures.name"]
    query.table = "report"
    query.flat = true

    return loadJson<Array<string>>(`${configuration.serverUrl}/api/v1/load/${encodeQuery(query)}`, null, taskHandle, data => {
      this.data.value = data
    })
  }
}

export class PredefinedMeasureConfigurator implements DataQueryConfigurator, ChartConfigurator {
  constructor(private readonly measures: Array<string>, public skipZeroValues: Ref<boolean> = shallowRef(true)) {
  }

  configureQuery(query: DataQuery, configuration: DataQueryExecutorConfiguration): boolean {
    configureQuery(this.measures, query, this.skipZeroValues.value)
    configuration.chartConfigurator = this
    configuration.measures = this.measures
    return true
  }

  configureChart(data: DataQueryResult, configuration: DataQueryExecutorConfiguration): ChartOptions {
    return configureChart(configuration, data)
  }
}

export function measureNameToLabel(key: string): string {
  const metricPathEndDotIndex = key.indexOf(".")
  if (metricPathEndDotIndex == -1) {
    return key.replace(/_[a-z]$/g, "")
  }
  else {
    let name = key.substring(metricPathEndDotIndex + 1)
    if (name.length > 2 && name[name.length - 2] == ".") {
      name = name.substring(0, name.length - 2)
    }
    return name
  }
}

function configureQuery(measureNames: Array<string>, query: DataQuery, skipZeroValues: boolean): void {
  // stable order of series (UI) and fields in query (caching)
  measureNames.sort((a, b) => collator.compare(a, b))

  query.addField({
    name: "t",
    sql: "toUnixTimestamp(generated_time) * 1000",
  })

  if (query.db !== "ij" && query.db !== "fleet") {
    query.addField({name: "measures", subName: "value"})
    query.addFilter({field: "measures.name", value: measureNames})
    if (skipZeroValues) {
      query.addFilter({field: "measures.value", operator: "!=", value: 0})
    }
  }
  else {
    for (let i = 0; i < measureNames.length; i++) {
      const value = measureNames[i]
      query.addField(value)
      if (skipZeroValues) {
        query.addFilter({field: value, operator: "!=", value: 0})
      }
    }
  }

  if (query.order != null) {
    throw new Error("order must be configured only by MetricLoader")
  }
  query.order = ["t"]
}

function configureChart(configuration: DataQueryExecutorConfiguration, dataSetList: DataQueryResult): ChartOptions {
  const measures = configuration.measures
  const series = new Array<LineSeriesOption>(measures.length * dataSetList.length)

  const dataSets = new Array<DatasetOption>(dataSetList.length)
  let seriesIndex = 0
  const extraQueryProducer = configuration.extraQueryProducer
  for (let dataSetIndex = 0; dataSetIndex < dataSetList.length; dataSetIndex++) {
    const dimensions = new Array<DimensionDefinition>(1 + measures.length)
    dimensions[0] = {name: "time", type: "time"}

    const dataSetLabel = extraQueryProducer == null ? null : extraQueryProducer.getDataSetLabel(dataSetIndex)
    for (let i = 0; i < measures.length; i++) {
      let seriesName: string
      if (dataSetLabel === null) {
        seriesName = measureNameToLabel(measures[i])
      }
      else if (measures.length === 1) {
        seriesName = dataSetLabel
      }
      else {
        seriesName = `${dataSetLabel}-${measureNameToLabel(measures[i])}`
      }
      series[seriesIndex++] = {
        name: seriesName,
        datasetIndex: dataSetIndex,
        type: "line",
        smooth: true,
        showSymbol: false,
        legendHoverLink: true,
        sampling: "lttb",
        encode: {
          // index of time
          x: 0,
          // +1 because time is the 0-dimension
          y: i + 1,
          tooltip: [i + 1],
        },
      }
      dimensions[i + 1] = {
        name: seriesName,
        type: "int",
      }
    }

    dataSets[dataSetIndex] = {
      dimensions,
      // just optimization to avoid auto-detect (https://echarts.apache.org/en/option.html#dataset.sourceHeader)
      sourceHeader: false,
      source: dataSetList[dataSetIndex],
    }
  }

  return {
    dataset: dataSets,
    series,
  }
}