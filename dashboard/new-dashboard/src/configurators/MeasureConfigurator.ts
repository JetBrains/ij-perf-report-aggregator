import { LineSeriesOption, ScatterSeriesOption } from "echarts/charts"
import { DatasetOption, ECBasicOption, ZRColor } from "echarts/types/dist/shared"
import { CallbackDataParams } from "echarts/types/src/util/types"
import { deepEqual } from "fast-equals"
import { debounceTime, distinctUntilChanged, forkJoin, map, Observable, of, switchMap } from "rxjs"
import { Ref, shallowRef } from "vue"
import { DataQueryResult } from "../components/common/DataQueryExecutor"
import { PersistentStateManager } from "../components/common/PersistentStateManager"
import { ChartConfigurator, ChartType, collator, SymbolOptions, ValueUnit } from "../components/common/chart"
import { DataQuery, DataQueryConfigurator, DataQueryDimension, DataQueryExecutorConfiguration, DataQueryFilter, toMutableArray } from "../components/common/dataQuery"
import { LineChartOptions, ScatterChartOptions } from "../components/common/echarts"
import { durationAxisPointerFormatter, isDurationFormatterApplicable, nsToMs, numberAxisLabelFormatter } from "../components/common/formatter"
import { useSettingsStore } from "../components/settings/settingsStore"
import { toColor } from "../util/colors"
import { MAIN_METRICS_SET } from "../util/mainMetrics"
import { Accident, AccidentKind, AccidentsConfigurator, getAccidents } from "./AccidentsConfigurator"
import { ChangePointClassification, detectChanges } from "./DetectChangesConfigurator"
import { scaleToMedian } from "./ScalingConfigurator"
import { ServerConfigurator } from "./ServerConfigurator"
import { exponentialSmoothingWithAlphaInference } from "./SmoothingConfigurator"
import { createComponentState, updateComponentState } from "./componentState"
import { configureQueryFilters, createFilterObservable, FilterConfigurator } from "./filter"
import { fromFetchWithRetryAndErrorHandling, refToObservable } from "./rxjs"

export class MeasureConfigurator implements DataQueryConfigurator, ChartConfigurator {
  readonly data = shallowRef<string[]>([])
  private readonly _selected = shallowRef<string[] | string | null>(null)
  readonly state = createComponentState()

  createObservable(): Observable<unknown> {
    return refToObservable(this.selected, true)
  }

  setSelected(value: string[] | string | null) {
    this._selected.value = value
  }

  get selected(): Ref<string[] | null> {
    const ref = this._selected
    if (typeof ref.value === "string") {
      ref.value = [ref.value]
    }
    return ref as Ref<string[] | null>
  }

  constructor(
    serverConfigurator: ServerConfigurator,
    persistentStateManager: PersistentStateManager,
    filters: FilterConfigurator[] = [],
    readonly skipZeroValues: boolean = true,
    readonly chartType: ChartType = "line",
    readonly symbolOptions: SymbolOptions = {}
  ) {
    persistentStateManager.add("measure", this._selected)

    const isIj = serverConfigurator.db === "ij"

    createFilterObservable(serverConfigurator, filters)
      .pipe(
        debounceTime(100),
        distinctUntilChanged(deepEqual),
        switchMap(() => {
          const loadMeasureListUrl = getLoadMeasureListUrl(serverConfigurator, filters)
          if (loadMeasureListUrl == null) {
            return of(null)
          }

          this.state.loading = true
          return isIj
            ? forkJoin([
                fromFetchWithRetryAndErrorHandling<string[]>(`${serverConfigurator.serverUrl}/api/v1/meta/measure?db=${serverConfigurator.db}`),
                fromFetchWithRetryAndErrorHandling<string[]>(loadMeasureListUrl),
              ]).pipe(
                map((data) => {
                  return data.flat(1)
                })
              )
            : fromFetchWithRetryAndErrorHandling<string[]>(loadMeasureListUrl)
        }),
        updateComponentState(this.state)
      )
      .subscribe((data) => {
        if (data == null) {
          return
        }

        if (isIj) {
          data = data.filter((it) => !/^c\.i\.ide\.[A-Za-z]\.[A-Za-z]: scheduled$/.test(it))
          data = data.filter((it) => !/^c\.i\.ide\.[A-Za-z]\.$/.test(it))
          data = data.filter((it) => !/^c\.i\.ide\.[A-Za-z]\.[A-Za-z](\.)?$/.test(it))
          data = [...new Set(data.map((it) => (/^c\.i\.ide\.[A-Za-z]\.[A-Za-z] preloading$/.test(it) ? "com.intellij.ide.misc.EvaluationSupport" : it)))]
        }

        const selectedRef = this.selected
        //filter out _23 metrics, we need them in DB but not in UI
        data = data.filter((it) => !/.*_\d+(#.*)?$/.test(it))
        this.data.value = data
        const selected = selectedRef.value
        if (selected != null && selected.length > 0) {
          // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
          const filtered = selected.filter((it) => data!.includes(it))
          if (filtered.length !== selected.length) {
            selectedRef.value = filtered
          }
        }
        selectedRef.value = [...new Set([...(selectedRef.value as string[]), ...data.filter((value) => MAIN_METRICS_SET.has(value))])]
      })
  }

  configureQuery(query: DataQuery, configuration: DataQueryExecutorConfiguration): boolean {
    const measureNames = toMutableArray(this.selected.value)
    if (measureNames.length === 0) {
      return false
    }

    configureQuery(measureNames, query, configuration, this.skipZeroValues)
    configuration.measures = measureNames
    configuration.addChartConfigurator(this)
    return true
  }

  configureChart(data: DataQueryResult, configuration: DataQueryExecutorConfiguration): ECBasicOption {
    return configureChart(configuration, data, this.chartType, "ms", this.symbolOptions)
  }
}

function getLoadMeasureListUrl(serverConfigurator: ServerConfigurator, filters: FilterConfigurator[]): string | null {
  const query = new DataQuery()
  const configuration = new DataQueryExecutorConfiguration()
  if (!serverConfigurator.configureQuery(query, configuration)) {
    return null
  }

  if (!configureQueryFilters(query, filters)) {
    return null
  }

  let fieldPrefix: string
  if (serverConfigurator.db === "ij") {
    fieldPrefix = "measure"
  } else {
    fieldPrefix = serverConfigurator.table === "measure" ? "" : "measures"
  }

  // "group by" is equivalent of distinct (https://clickhouse.tech/docs/en/sql-reference/statements/select/distinct/#alternatives)
  query.addDimension(fieldPrefix.length === 0 ? { n: "name" } : { n: fieldPrefix, subName: "name" })
  query.order = fieldPrefix.length === 0 ? "name" : `${fieldPrefix}.name`
  query.table = serverConfigurator.table
  query.flat = true
  return serverConfigurator.computeQueryUrl(query)
}

export class PredefinedMeasureConfigurator implements DataQueryConfigurator, ChartConfigurator {
  constructor(
    private readonly measures: string[],
    readonly skipZeroValues: Ref<boolean> = shallowRef(true),
    private readonly chartType: ChartType = "line",
    private readonly valueUnit: ValueUnit = "ms",
    readonly symbolOptions: SymbolOptions = {},
    readonly accidentsConfigurator: AccidentsConfigurator | null = null
  ) {}

  createObservable(): Observable<unknown> {
    return refToObservable(this.skipZeroValues)
  }

  configureQuery(query: DataQuery, configuration: DataQueryExecutorConfiguration): boolean {
    configureQuery(this.measures, query, configuration, this.skipZeroValues.value)
    configuration.addChartConfigurator(this)
    configuration.measures = this.measures
    return true
  }

  configureChart(data: DataQueryResult, configuration: DataQueryExecutorConfiguration): ECBasicOption {
    return configureChart(configuration, data, this.chartType, this.valueUnit, this.symbolOptions, this.accidentsConfigurator)
  }
}

export function measureNameToLabel(key: string): string {
  return key.includes(".") ? key : /* remove _d or _i suffix */ key.replaceAll(/_[a-z]$/g, "")
}

function configureQuery(measureNames: string[], query: DataQuery, configuration: DataQueryExecutorConfiguration, skipZeroValues: boolean): void {
  // stable order of series (UI) and fields in query (caching)
  measureNames.sort((a, b) => collator.compare(a, b))

  query.insertField(
    {
      n: "t",
      sql: "toUnixTimestamp(generated_time)*1000",
    },
    0
  )

  // we cannot request several measures in one SQL query - for each measure separate SQl query with filter by measure name
  const isIj = query.db === "ij"
  const structureName = isIj ? "measure" : "measures"
  const valueName = isIj ? "duration" : "value"
  const field: DataQueryDimension = { n: "" }
  query.insertField(field, 1)

  if (query.db === "perfint" || query.db === "jbr" || query.db === "perfintDev" || query.db == "bazel" || query.db == "perfUnitTests") {
    query.addField({ n: "measures", subName: "name" })
    query.addField({ n: "measures", subName: "type" })
  }

  const prevFilters: DataQueryFilter[] = []

  const addFilter = (filter: DataQueryFilter): void => {
    prevFilters.push(filter)
    query.addFilter(filter)
  }

  configuration.queryProducers.push({
    size(): number {
      return measureNames.length
    },
    mutate(index: number): void {
      const measure = measureNames[index]

      delete field.sql
      delete field.subName

      if (prevFilters.length > 0) {
        query.removeFilters(prevFilters)
        prevFilters.length = 0
      }

      let valueFieldName: string
      if (query.table === "measure") {
        field.n = "value"
        field.resultKey = measure.replaceAll(".", "_")
        addFilter({ f: "name", v: measure })
        valueFieldName = "value"
      } else if (isIj && measure.includes("metrics.")) {
        field.n = "metrics"
        field.subName = "value"
        addFilter({ f: "metrics.name", v: measure.split("metrics.", 2)[1] })
        valueFieldName = "metrics.value"
      } else if (isIj && !measure.includes(" ") && measure != "elementTypeCount") {
        field.n = measure
        valueFieldName = measure
      } else {
        if (measure.endsWith(".end")) {
          field.n = structureName
          field.subName = "end"
          field.sql = `(${structureName}.start+${structureName}.${valueName})`
        } else {
          field.n = structureName
          field.subName = valueName
        }

        addFilter({ f: `${structureName}.name`, v: measure.endsWith(".end") ? measure.slice(0, measure.length - ".end".length) : measure })
        valueFieldName = `${structureName}.${valueName}`
      }

      if (skipZeroValues) {
        addFilter({ f: valueFieldName, o: "!=", v: 0 })
      }
    },
    getSeriesName(index: number): string {
      return measureNames.length > 1 ? measureNameToLabel(measureNames[index]) : ""
    },
    getMeasureName(index: number): string {
      return measureNames[index]
    },
  })

  if (query.order != null) {
    throw new Error("order must be configured only by MetricLoader")
  }
  query.order = "t"
}

function getItemStyleForSeries(accidentMap: Map<string, Accident[]> | undefined, detectedChanges = new Map<string, ChangePointClassification>()) {
  return {
    color(seriesIndex: CallbackDataParams): ZRColor {
      const accidents = getAccidents(accidentMap, seriesIndex.value as string[])
      if (accidents == null || accidents.length === 0) {
        const detectChange = detectedChanges.get(JSON.stringify(seriesIndex.value as string[]))
        if (detectChange == ChangePointClassification.DEGRADATION) {
          return "#cc0000"
        } else if (detectChange == ChangePointClassification.OPTIMIZATION) {
          return "#009900"
        } else if (detectChange == ChangePointClassification.NO_CHANGE) {
          // return "#b4b3b3"
        }
        return seriesIndex.color as ZRColor
      }
      for (const accident of accidents) {
        switch (accident.kind) {
          case AccidentKind.Regression:
            return "#cc0000"
          case AccidentKind.InferredRegression:
            return "#efa9a9"
          case AccidentKind.Improvement:
            return "#009900"
          case AccidentKind.InferredImprovement:
            return "#acffac"
          case AccidentKind.Investigation:
            return "orange"
        }
      }
      return toColor(accidents[0].reason)
    },
  }
}

function isChangeDetected(detectedChanges: Map<string, ChangePointClassification>, value: string[]) {
  const changePointClassification = detectedChanges.get(JSON.stringify(value))
  return changePointClassification != undefined && changePointClassification != ChangePointClassification.NO_CHANGE
}

function configureChart(
  configuration: DataQueryExecutorConfiguration,
  dataList: DataQueryResult,
  chartType: ChartType,
  valueUnit: ValueUnit = "ms",
  symbolOptions: SymbolOptions = {},
  accidentsConfigurator: AccidentsConfigurator | null = null
): LineChartOptions | ScatterChartOptions {
  const series = new Array<LineSeriesOption | ScatterSeriesOption>()
  let useDurationFormatter = true

  const dataset: DatasetOption[] = []

  for (let dataIndex = 0, n = dataList.length; dataIndex < n; dataIndex++) {
    const measureName = configuration.measureNames[dataIndex]
    let seriesName = configuration.seriesNames[dataIndex]
    const seriesData = dataList[dataIndex]

    if (seriesData.length > 3) {
      // we take only the last type of the metric since it's not clear how to show different types and last type helps to change the type if necessary
      const type = seriesData[3].at(-1)
      if (type === "c") {
        useDurationFormatter = false
      } else if (type === "d") {
        useDurationFormatter = true
      }
    }
    //fleet
    if (seriesName == "" && (seriesData.length == 6 || seriesData.length == 10)) {
      seriesName = seriesData[4][0] as string
    } else if (seriesName == "" && seriesData.length > 6) {
      // we take only the one project name, there can't be more
      seriesName = seriesData[6][0] as string
    }

    const settings = useSettingsStore()
    if (settings.smoothing) {
      const smoothedData = exponentialSmoothingWithAlphaInference(seriesData[1] as number[])
      seriesData.push(smoothedData)
    }

    if (settings.scaling) {
      seriesData[1] = scaleToMedian(seriesData[1] as number[])
      useDurationFormatter = false
    }

    let detectedChanges = new Map<string, ChangePointClassification>()
    if (settings.detectChanges) {
      detectedChanges = detectChanges(seriesData)
    }

    let isNotEmpty = false
    for (const data of seriesData) {
      isNotEmpty = isNotEmpty || data.length > 0
    }

    if (isNotEmpty) {
      const name = seriesName.startsWith("metrics.") ? seriesName.slice("metrics.".length) : seriesName
      const id = measureName === seriesName ? seriesName : `${measureName}@${seriesName}`
      const seriesLayoutBy = "row"
      const datasetIndex = dataIndex
      const xAxisName = useDurationFormatter ? "time" : "count"
      series.push({
        selectedMode: "single",
        select: {
          itemStyle: {
            color: "black",
          },
        },
        // formatter is detected by measure name - that's why series id is specified (see usages of seriesId)
        id,
        name,
        type: settings.smoothing ? "scatter" : chartType,
        // showSymbol: symbolOptions.showSymbol == undefined ? seriesData[0].length < 100 : symbolOptions.showSymbol,
        // 10 is a default value for scatter (  undefined doesn't work to unset)
        symbolSize(value: string[]): number {
          const symbolSize = symbolOptions.symbolSize ?? (chartType === "line" ? Math.min(800 / seriesData[0].length, 9) : 10)
          const accidents = getAccidents(accidentsConfigurator?.value.value, value)
          if (isValueShouldBeMarkedWithPin(accidents)) {
            return symbolSize * 4
          }
          if (isChangeDetected(detectedChanges, value)) {
            return symbolSize * 2.5
          }
          if (isValueShouldBeMarkedAsException(accidents)) {
            return symbolSize * 1.2
          }
          return symbolSize
        },
        symbolRotate(value: string[]): number {
          const detectChange = detectedChanges.get(JSON.stringify(value))
          return detectChange == ChangePointClassification.OPTIMIZATION ? 180 : 0
        },
        symbol(value: string[]) {
          const accidents = getAccidents(accidentsConfigurator?.value.value, value)
          if (isValueShouldBeMarkedWithPin(accidents)) {
            return "pin"
          }
          if (isChangeDetected(detectedChanges, value)) {
            return "arrow"
          }
          if (isValueShouldBeMarkedAsException(accidents)) {
            return "diamond"
          }
          return "circle"
        },
        seriesLayoutBy,
        datasetIndex,
        dimensions: [
          { name: xAxisName, type: "time" },
          { name: seriesName, type: "int" },
        ],
        itemStyle: getItemStyleForSeries(accidentsConfigurator?.value.value, detectedChanges),
      })
      if (settings.smoothing) {
        series.push({
          // formatter is detected by measure name - that's why series id is specified (see usages of seriesId)
          id: id + "smoothed",
          name,
          type: "line",
          symbol: "none",
          silent: true,
          seriesLayoutBy,
          datasetIndex,
          encode: {
            x: xAxisName,
            y: seriesData.length - 1,
          },
          itemStyle: getItemStyleForSeries(accidentsConfigurator?.value.value),
        })
      }
    }
    if (useDurationFormatter && !isDurationFormatterApplicable(measureName)) {
      useDurationFormatter = false
    }

    dataset.push({
      source: seriesData,
      sourceHeader: false,
    })
  }

  // if (chartType == "scatter") {
  //   dataset.push({
  //     fromDatasetIndex: 0,
  //     transform: {
  //       type: "ecStat:regression",
  //       // eslint-disable-next-line @typescript-eslint/ban-ts-comment
  //       // @ts-ignore
  //       config: {
  //         method: "polynomial",
  //         dimensions: [1],
  //       },
  //     },
  //   })
  //   series.push({
  //     silent: true,
  //     type: "line",
  //     smooth: true,
  //     datasetIndex: dataset.length - 1,
  //     // symbolSize: 0.1,
  //     // symbol: "circle",
  //     // label: {show: false, fontSize: 16},
  //     // labelLayout: {dx: -20},
  //     // encode: {label: 2, tooltip: 1},
  //   })
  // }

  const isNs = valueUnit == "ns"
  const valueInMsFormatter = useDurationFormatter ? durationAxisPointerFormatter : numberAxisLabelFormatter
  const formatter: (valueInMs: number) => string = isNs ? (v) => valueInMsFormatter(nsToMs(v)) : valueInMsFormatter
  return {
    dataset,
    yAxis: {
      axisLabel: {
        formatter,
      },
      axisPointer: {
        label: {
          formatter(data): string {
            return formatter(data.value as number)
          },
        },
      },
    },
    series: series as LineSeriesOption,
  }
}

function isValueShouldBeMarkedWithPin(accidents: Accident[] | null): boolean {
  return accidents != null && accidents.some((accident) => accident.kind != AccidentKind.Exception)
}

function isValueShouldBeMarkedAsException(accidents: Accident[] | null): boolean {
  return accidents != null && accidents.every((accident) => accident.kind == AccidentKind.Exception)
}
