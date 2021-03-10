import { LineSeriesOption } from "echarts/charts"
import { DatasetOption } from "echarts/types/dist/shared"
import { DimensionDefinition } from "echarts/types/src/util/types"
import { Ref, shallowRef, watch } from "vue"
import { DataQueryResult } from "../DataQueryExecutor"
import { PersistentStateManager } from "../PersistentStateManager"
import { ChartConfigurator, ChartOptions } from "../chart"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration, DataQueryFilter, encodeQuery, toMutableArray } from "../dataQuery"
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

    configureQuery(measureNames, query, configuration, this.skipZeroValues)
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
    if (this.serverConfigurator.databaseName === "ij") {
      const server = this.serverConfigurator.value.value
      if (server == null || server.length === 0) {
        return Promise.resolve()
      }

      return loadJson<Array<string>>(`${server}/api/v1/meta/measure?db=${this.serverConfigurator.databaseName}`, null, taskHandle, data => {
        this.data.value = data
      })
    }

    const parent = this.parent
    const query = new DataQuery()
    const configuration = new DataQueryExecutorConfiguration()
    if (!this.serverConfigurator.configureQuery(query, configuration)) {
      return Promise.resolve()
    }
    if (parent != null && !parent.configureQuery(query, configuration)) {
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
  private readonly measures: Array<string>

  constructor(measures: Array<string>, public skipZeroValues: Ref<boolean> = shallowRef(true)) {
    // analyzer replaces space to underscore, but configurator supports specifying original metric names
    // render window.end => render_window.end
    this.measures = measures.map(it => it.replaceAll(" ", "_"))
  }

  configureQuery(query: DataQuery, configuration: DataQueryExecutorConfiguration): boolean {
    configureQuery(this.measures, query, configuration, this.skipZeroValues.value)
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
    return key
  }
  else {
    return key
      .substring(0, metricPathEndDotIndex)
      .replaceAll("_",  " ")
  }
}

function configureQuery(measureNames: Array<string>, query: DataQuery, configuration: DataQueryExecutorConfiguration, skipZeroValues: boolean): void {
  // stable order of series (UI) and fields in query (caching)
  measureNames.sort((a, b) => collator.compare(a, b))

  query.addField({
    name: "t",
    sql: "toUnixTimestamp(generated_time) * 1000",
  })

  if (query.db === "ij") {
    for (let i = 0; i < measureNames.length; i++) {
      const value = measureNames[i]
      query.addField(value)
      if (skipZeroValues) {
        query.addFilter({field: value, operator: "!=", value: 0})
      }
    }
  }
  else {
    const pureNames = measureNames.map(it => it.endsWith(".end") ? it.substring(0, it.length - ".end".length) : it)
    const filter = {field: "measures.name", value: pureNames[0]}
    if (measureNames.length > 1) {
      // we cannot request several measures in one SQL query - for each measure separate SQl query with filter by measure name
      configureQueryProducer(configuration, filter, pureNames)
    }

    if (measureNames.some(it => it.endsWith(".end"))) {
      query.addField({name: "measures", subName: "end", sql: "(measures.start + measures.value)"})
    }
    else {
      query.addField({name: "measures", subName: "value"})
    }

    query.addFilter(filter)

    if (skipZeroValues) {
      // for end we also filter by raw value and not by sum of start + duration (that stored under "value" name)
      query.addFilter({field: "measures.value", operator: "!=", value: 0})
    }
  }

  if (query.order != null) {
    throw new Error("order must be configured only by MetricLoader")
  }
  query.order = ["t"]
}

function configureQueryProducer(configuration: DataQueryExecutorConfiguration, filter: DataQueryFilter, values: Array<string>): void {
  let index = 1
  if (configuration.extraQueryProducer != null) {
    throw new Error("extraQueryMutator is already set")
  }

  configuration.extraQueryProducer = {
    mutate() {
      filter.value = values[index++]
      return index !== values.length
    },
    getDataSetLabel(index: number): string {
      return values[index].replaceAll("_", " ")
    },
    getDataSetMeasureNames(_index: number): Array<string> {
      return [values[_index]]
    }
  }
}

function configureChart(configuration: DataQueryExecutorConfiguration, dataSetList: DataQueryResult): ChartOptions {
  const series = new Array<LineSeriesOption>()

  const dataSets = new Array<DatasetOption>(dataSetList.length)
  const extraQueryProducer = configuration.extraQueryProducer
  for (let dataSetIndex = 0; dataSetIndex < dataSetList.length; dataSetIndex++) {
    const measures = extraQueryProducer == null ? configuration.measures : extraQueryProducer.getDataSetMeasureNames(dataSetIndex)
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
      series.push({
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
      })
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