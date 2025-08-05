import { LineSeriesOption, ScatterSeriesOption } from "echarts/charts"
import { DatasetOption, ECBasicOption, ZRColor } from "echarts/types/dist/shared"
import type { DefaultLabelFormatterCallbackParams as CallbackDataParams } from "echarts"
import { deepEqual } from "fast-equals"
import { combineLatest, debounceTime, distinctUntilChanged, forkJoin, map, Observable, of, switchMap } from "rxjs"
import { ref, Ref, shallowRef } from "vue"
import { DataQueryResult } from "../components/common/DataQueryExecutor"
import { PersistentStateManager } from "../components/common/PersistentStateManager"
import { ChartConfigurator, ChartType, collator, SymbolOptions, ValueUnit } from "../components/common/chart"
import {
  DataQuery,
  DataQueryConfigurator,
  DataQueryDimension,
  DataQueryExecutorConfiguration,
  DataQueryFilter,
  ServerConfigurator,
  toMutableArray,
} from "../components/common/dataQuery"
import { LineChartOptions, ScatterChartOptions } from "../components/common/echarts"
import { durationAxisPointerFormatter, isDurationFormatterApplicable, nsToMs, numberAxisLabelFormatter } from "../components/common/formatter"
import { DBType } from "../components/common/sideBar/InfoSidebar"
import { useSettingsStore } from "../components/settings/settingsStore"
import { ChangePointClassification } from "../shared/changeDetector/algorithm"
import { detectChanges } from "../shared/changeDetector/workerStarter"
import { dbTypeStore } from "../shared/dbTypes"
import { measureNameToLabel } from "../shared/metricsMapping"
import { Delta } from "../util/Delta"
import { toColor } from "../util/colors"
import { MAIN_METRICS, MAIN_METRICS_SET } from "../util/mainMetrics"
import { Accident, AccidentKind, AccidentsConfigurator } from "./accidents/AccidentsConfigurator"
import { scaleToMedian } from "../components/settings/configurators/ScalingConfigurator"
import { exponentialSmoothingWithAlphaInference } from "../components/settings/configurators/SmoothingConfigurator"
import { createComponentState, updateComponentState } from "./componentState"
import { configureQueryFilters, createFilterObservable, FilterConfigurator } from "./filter"
import { fromFetchWithRetryAndErrorHandling, refToObservable } from "./rxjs"
import { removeOutliers } from "../components/settings/configurators/RemoveOutliersConfigurator"
import { getBasicInfo, getBuildId } from "../components/common/sideBar/InfoSidebarPerformance"
import { useDarkModeStore } from "../shared/useDarkModeStore"
import { useSelectedPointStore } from "../shared/selectedPointStore"

export type TooltipTrigger = "item" | "axis" | "none"

export class MeasureConfigurator implements DataQueryConfigurator, ChartConfigurator, FilterConfigurator {
  readonly data = shallowRef<string[]>([])
  private readonly _selected = shallowRef<string[] | string | null>(null)
  readonly state = createComponentState()

  readonly showAllMetrics = ref(false)

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

  setShowAllMetrics(value: boolean) {
    this.showAllMetrics.value = value
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

    const isIj = dbTypeStore().isIJStartup()

    combineLatest([createFilterObservable(serverConfigurator, filters), refToObservable(this.showAllMetrics)])
      .pipe(
        debounceTime(100),
        distinctUntilChanged(deepEqual),
        switchMap(() => {
          const loadMeasureListUrl = getLoadMeasureListUrl(serverConfigurator, filters)
          if (loadMeasureListUrl == null) {
            return of(null)
          }

          const loadMetricsUrl = getLoadMetricsListUrl(serverConfigurator, filters)
          if (loadMetricsUrl == null) {
            return of(null)
          }

          this.state.loading = true
          return isIj
            ? forkJoin([
                fromFetchWithRetryAndErrorHandling<string[]>(`${serverConfigurator.serverUrl}/api/v1/meta/measure?db=${serverConfigurator.db}`),
                fromFetchWithRetryAndErrorHandling<string[]>(loadMeasureListUrl),
                fromFetchWithRetryAndErrorHandling<string[]>(loadMetricsUrl),
              ]).pipe(
                map((data) => {
                  data[2] = data[2].map((it) => "metrics." + it)
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
          data = data.filter(
            (it) =>
              !/^c\.i\.ide\.[A-Za-z]\.[A-Za-z]: scheduled$/.test(it) &&
              !/^c\.i\.ide\.[A-Za-z]\.$/.test(it) &&
              !/^c\.i\.ide\.[A-Za-z]\.[A-Za-z](\.)?$/.test(it) &&
              !/^ProjectImpl@\d+ container$/.test(it)
          )
          data = [...new Set(data.map((it) => (/^c\.i\.ide\.[A-Za-z]\.[A-Za-z] preloading$/.test(it) ? "com.intellij.ide.misc.EvaluationSupport" : it)))]
        }

        if (dbTypeStore().dbType == DBType.FLEET) {
          data = data.filter((it) => !/.*id=.*/.test(it) && it.length < 120)
          data = data.map((it) => it + ".end")
        }

        let filtered = data.filter(
          (it) =>
            //filter for editor menu
            !/.*#(update|getchildren|getselection)@.*/i.test(it) &&
            //filter out _23 metrics, we need them in DB but not in UI
            (!/.*_\d+(#.*)?$/.test(it) || this.showAllMetrics.value)
        )

        filtered = customSort(filtered, MAIN_METRICS)

        const selectedRef = this.selected

        const selected = selectedRef.value
        if (selected != null && selected.length > 0) {
          const selectedInData = selected.filter((it) => data.includes(it))
          if (selectedInData.length > 0) {
            filtered = [...new Set([...filtered, ...selectedInData])]
          }

          if (selectedInData.length !== selected.length) {
            selectedRef.value = selectedInData
          }
        }
        this.data.value = filtered
        selectedRef.value = [...new Set([...(selectedRef.value as string[]), ...filtered.filter((value) => MAIN_METRICS_SET.has(value))])]
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

  configureChart(data: DataQueryResult, configuration: DataQueryExecutorConfiguration): Promise<ECBasicOption> {
    return configureChart(configuration, data, this.chartType, "ms", this.symbolOptions)
  }

  configureFilter(query: DataQuery): boolean {
    const currentValue = this._selected.value
    const mergedFilter: string[] = []
    if (Array.isArray(currentValue)) {
      for (const metric of currentValue) {
        mergedFilter.push("has(`measures.name`, '" + metric + "')")
      }
    } else {
      mergedFilter.push("has(`measures.name`, '" + (currentValue ?? "") + "')")
    }
    const filterQuery = mergedFilter.join(" or ")
    if (currentValue != undefined && currentValue.length > 0) {
      query.addFilter({ q: filterQuery })
      return true
    } else {
      return false
    }
  }
}

function customSort(tsArray: string[], referenceArray: string[]): string[] {
  const referenceSet = new Set(referenceArray)

  const inReference = tsArray.filter((x) => referenceSet.has(x))
  const notInReference = tsArray.filter((x) => !referenceSet.has(x))

  return [...inReference, ...notInReference]
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
  if (dbTypeStore().isIJStartup()) {
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

function getLoadMetricsListUrl(serverConfigurator: ServerConfigurator, filters: FilterConfigurator[]): string | null {
  const query = new DataQuery()
  const configuration = new DataQueryExecutorConfiguration()
  if (!serverConfigurator.configureQuery(query, configuration)) {
    return null
  }

  if (!configureQueryFilters(query, filters)) {
    return null
  }

  const fieldPrefix = "metrics"
  // "group by" is equivalent of distinct (https://clickhouse.tech/docs/en/sql-reference/statements/select/distinct/#alternatives)
  query.addDimension(fieldPrefix.length === 0 ? { n: "name" } : { n: fieldPrefix, subName: "name" })
  query.order = fieldPrefix.length === 0 ? "name" : `${fieldPrefix}.name`
  query.table = serverConfigurator.table
  query.flat = true
  return serverConfigurator.computeQueryUrl(query)
}

export class PredefinedMeasureConfigurator implements DataQueryConfigurator, ChartConfigurator {
  constructor(
    private readonly measures: Ref<string[]> = shallowRef([]),
    readonly skipZeroValues: Ref<boolean> = shallowRef(true),
    private readonly chartType: ChartType = "line",
    private readonly valueUnit: ValueUnit = "ms",
    readonly symbolOptions: SymbolOptions = {},
    readonly accidentsConfigurator: AccidentsConfigurator | null = null,
    readonly toolTipTrigger: TooltipTrigger
  ) {}

  createObservable(): Observable<unknown> {
    return combineLatest([refToObservable(this.skipZeroValues), refToObservable(this.measures)])
  }

  configureQuery(query: DataQuery, configuration: DataQueryExecutorConfiguration): boolean {
    if (this.measures.value.length === 0) {
      return false
    }
    configureQuery(this.measures.value, query, configuration, this.skipZeroValues.value)
    configuration.addChartConfigurator(this)
    configuration.measures = this.measures.value
    return true
  }

  configureChart(data: DataQueryResult, configuration: DataQueryExecutorConfiguration): Promise<ECBasicOption> {
    return configureChart(configuration, data, this.chartType, this.valueUnit, this.symbolOptions, this.accidentsConfigurator, this.toolTipTrigger)
  }
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
  const isIj = dbTypeStore().isIJStartup()
  const structureName = isIj ? "measure" : "measures"
  const valueName = isIj ? "duration" : "value"
  const field: DataQueryDimension = { n: "" }
  query.insertField(field, 1)

  if (query.db !== "ij" && query.db !== "ijDev" && !(query.db === "fleet" && query.table === "report")) {
    query.addField({ n: "measures", subName: "name" })
    query.addField({ n: "measures", subName: "type" })
  }

  const metricNameField: DataQueryDimension = { n: "" }
  if (dbTypeStore().isStartup()) {
    query.insertField(metricNameField, 2)
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

      if (dbTypeStore().isStartup()) {
        delete metricNameField.sql
        delete metricNameField.subName
        if (measure.startsWith("metrics.")) {
          metricNameField.n = "metrics"
          metricNameField.subName = "name"
        } else {
          metricNameField.n = "metricName"
          metricNameField.sql = `'${measure}'`
        }
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
      return measureNames.length > 1 ? measureNames[index] : ""
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

function getItemStyleForSeries(accidentConfigurator: AccidentsConfigurator | null, detectedChanges = new Map<string, ChangePointClassification>()) {
  return {
    color(seriesIndex: CallbackDataParams): ZRColor {
      if (useSelectedPointStore().selectedPoint != undefined && getBasicInfo(seriesIndex, "ms").buildId.toString() == useSelectedPointStore().selectedPoint) {
        return getSelectedPointColor()
      }
      const accidents = accidentConfigurator?.getAccidents(seriesIndex.value as string[])
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

class MergeResults {
  constructor(
    readonly data: DataQueryResult,
    private readonly idToSeriesName: Map<number, string>,
    private readonly idToMeasureName: Map<number, string>
  ) {}

  getSeriesName(index: number): string {
    return this.idToSeriesName.get(index) as string
  }

  getMeasureName(index: number): string {
    return this.idToMeasureName.get(index) as string
  }
}

function mergeSeries(dataList: (string | number)[][][], configuration: DataQueryExecutorConfiguration) {
  const mergedDataList: DataQueryResult = []
  const seriesIdsToIndex = new Map<string, number>()
  const seriesIdToSeriesName = new Map<number, string>()
  const seriesIdToMeasureName = new Map<number, string>()
  for (const [dataIndex, seriesData] of dataList.entries()) {
    if (seriesData[1]?.length === 0) {
      console.log("Serie is empty and will be hidden: " + configuration.seriesNames[dataIndex])
      continue
    }
    const measureName = configuration.measureNames[dataIndex]
    let seriesName = configuration.seriesNames[dataIndex]
    //fleet
    if (seriesName == "" && (seriesData.length == 6 || seriesData.length == 10)) {
      seriesName = seriesData[4][0] as string
    } else if (seriesName == "" && seriesData.length > 6) {
      // we take only the one project name, there can't be more
      seriesName = seriesData[6][0] as string
    }
    seriesName = measureNameToLabel(seriesName)
    const id = measureName === seriesName ? seriesName : `${measureName}@${seriesName}`
    if (seriesIdsToIndex.has(id)) {
      const seriesIndex = seriesIdsToIndex.get(id) as number
      const values = mergedDataList[seriesIndex]
      for (const [i, seriesDatum] of seriesData.entries()) {
        values[i] = i < values.length ? [...values[i], ...seriesDatum] : [...seriesDatum]
      }
    } else {
      const newId = mergedDataList.push(seriesData) - 1
      seriesIdsToIndex.set(id, newId)
      seriesIdToSeriesName.set(newId, seriesName)
      seriesIdToMeasureName.set(newId, measureName)
    }
  }
  return new MergeResults(mergedDataList, seriesIdToSeriesName, seriesIdToMeasureName)
}

function getSelectedPointColor() {
  return useDarkModeStore().darkMode ? "white" : "black"
}

async function configureChart(
  configuration: DataQueryExecutorConfiguration,
  dataList: DataQueryResult,
  chartType: ChartType,
  valueUnit: ValueUnit = "ms",
  symbolOptions: SymbolOptions = {},
  accidentsConfigurator: AccidentsConfigurator | null = null,
  tooltipTrigger: TooltipTrigger = "item"
): Promise<LineChartOptions | ScatterChartOptions> {
  const series = new Array<LineSeriesOption | ScatterSeriesOption>()
  let useDurationFormatter = true

  const dataset: DatasetOption[] = []

  //merge series with the same name
  const mergeResults = mergeSeries(dataList, configuration)

  const settings = useSettingsStore()
  // eslint-disable-next-line prefer-const
  for (let [dataIndex, seriesData] of mergeResults.data.entries()) {
    // eslint-disable-next-line @typescript-eslint/no-unnecessary-condition
    if (seriesData[1] == undefined) {
      //we need to push even empty dataset otherwise it will be out of sync with series and plot will be empty
      dataset.push({
        source: seriesData,
        sourceHeader: false,
      })
      continue
    }

    if (seriesData.length > 3) {
      // we take only the last type of the metric since it's not clear how to show different types and last type helps to change the type if necessary
      const type = seriesData[3].at(-1)
      if (type === "c") {
        useDurationFormatter = false
      } else if (type === "d") {
        useDurationFormatter = true
      }
    }

    if (settings.removeOutliers) {
      seriesData = removeOutliers(seriesData)
    }

    if (settings.smoothing) {
      const smoothedData = exponentialSmoothingWithAlphaInference(seriesData[1] as number[])
      seriesData.push(smoothedData)
    }

    const deltaValues = Delta.calculateDeltas(seriesData[1] as number[], getBuildId(seriesData) as number[])
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    //@ts-expect-error
    seriesData.push(deltaValues)

    if (settings.scaling) {
      seriesData.push(seriesData[1])
      seriesData[1] = scaleToMedian(seriesData[1] as number[])
      useDurationFormatter = false
    }

    let detectedChanges = new Map<string, ChangePointClassification>()
    if (settings.detectChanges) {
      detectedChanges = await detectChanges(seriesData)
    }

    let isNotEmpty = false
    for (const data of seriesData) {
      isNotEmpty = isNotEmpty || data.length > 0
    }
    if (isNotEmpty) {
      // noinspection SuspiciousTypeOfGuard
      const measureName = mergeResults.getMeasureName(dataIndex)
      const seriesName = mergeResults.getSeriesName(dataIndex)
      const seriesLayoutBy = "row"
      const datasetIndex = dataIndex
      const xAxisName = useDurationFormatter ? "time" : "count"
      series.push({
        selectedMode: "single",
        select: {
          itemStyle: {
            color: getSelectedPointColor(),
          },
        },
        // formatter is detected by measure name - that's why series id is specified (see usages of seriesId)
        id: measureName === seriesName ? seriesName : `${measureName}@${seriesName}`,
        name: seriesName,
        type: settings.smoothing ? "scatter" : chartType,
        // showSymbol: symbolOptions.showSymbol == undefined ? seriesData[0].length < 100 : symbolOptions.showSymbol,
        // 10 is a default value for scatter (  undefined doesn't work to unset)
        symbolSize(value: string[]): number {
          const symbolSize = symbolOptions.symbolSize ?? (chartType === "line" ? Math.min((10 * 1000) / seriesData[0].length, 7) : 10)
          const accidents = accidentsConfigurator?.getAccidents(value) ?? null
          if (isValueShouldBeMarkedWithPin(accidents)) {
            return symbolSize * 4
          }
          if (isChangeDetected(detectedChanges, value)) {
            return symbolSize * 2.5
          }
          if (isValueShouldBeMarkedAsException(accidents)) {
            return symbolSize * 1.2
          }
          return settings.smoothing ? symbolSize / 3 : symbolSize
        },
        symbolRotate(value: string[]): number {
          const detectChange = detectedChanges.get(JSON.stringify(value))
          return detectChange == ChangePointClassification.OPTIMIZATION ? 180 : 0
        },
        symbol(value: string[]) {
          const accidents = accidentsConfigurator?.getAccidents(value) ?? null
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
        itemStyle: getItemStyleForSeries(accidentsConfigurator, detectedChanges),
      })
      if (settings.smoothing) {
        series.push({
          // formatter is detected by measure name - that's why series id is specified (see usages of seriesId)
          id: (measureName === seriesName ? seriesName : `${measureName}@${seriesName}`) + "smoothed",
          name: seriesName,
          type: "line",
          symbol: "none",
          silent: true,
          seriesLayoutBy,
          datasetIndex,
          encode: {
            x: xAxisName,
            y: seriesData.length - 2,
          },
          itemStyle: getItemStyleForSeries(accidentsConfigurator),
        })
      }
    }
    if (useDurationFormatter && !isDurationFormatterApplicable(mergeResults.getMeasureName(dataIndex))) {
      useDurationFormatter = false
    }

    dataset.push({
      source: seriesData,
      sourceHeader: false,
    })
  }
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
    tooltip: {
      trigger: tooltipTrigger == "item" ? "item" : dataset.length > 5 ? "item" : "axis",
    },
    dataZoom: [
      {
        type: "slider",
        showDataShadow: false,
        width: 10,
        yAxisIndex: 0,
        filterMode: "filter",
        brushSelect: false,
        show: true,
        fillerColor: useDarkModeStore().darkMode ? "rgba(90,90,90,0.25)" : "rgba(106,114,128,0.1)",
        borderColor: useDarkModeStore().darkMode ? "#444444" : "#d2dbee",
        handleStyle: {
          color: useDarkModeStore().darkMode ? "#444444" : "#d2dbee",
        },
      },
    ],
    series: series as LineSeriesOption,
  }
}

function isValueShouldBeMarkedWithPin(accidents: Accident[] | null): boolean {
  return accidents?.some((accident) => accident.kind != AccidentKind.Exception) ?? false
}

function isValueShouldBeMarkedAsException(accidents: Accident[] | null): boolean {
  return accidents != null && accidents.length > 0 && accidents.every((accident) => accident.kind == AccidentKind.Exception)
}
