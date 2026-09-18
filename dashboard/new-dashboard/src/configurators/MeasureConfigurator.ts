import { LineSeriesOption, ScatterSeriesOption } from "echarts/charts"
import { DatasetOption, ECBasicOption, ZRColor } from "echarts/types/dist/shared"
import type { DefaultLabelFormatterCallbackParams as CallbackDataParams } from "echarts"
import { deepEqual } from "fast-equals"
import { ColdObservable } from "rxjs"
import { debounce } from "rxjs/debounce"
import { distinctUntilChanged } from "rxjs/distinct-until-changed"
import { switchMap } from "rxjs/switch-map"
import { ref, Ref, shallowRef, toRef } from "vue"
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
import { formatMeasureValue, MeasureUnit, reduceToAxisUnit } from "../components/common/formatter"
import { DBType } from "../components/common/sideBar/InfoSidebar"
import { useSettingsStore } from "../components/settings/settingsStore"
import { getChartLastTimestamp, getStaleSeriesMarkLine } from "../components/charts/staleSeries"
import { getSeriesId } from "../components/charts/seriesId"
import { buildStdDevBand, getStdDevBandSeries, getStdDevMeasureName } from "../components/charts/stdDevInterval"
import { BetterDirection, ChangePointClassification, DetectedChange } from "../shared/changeDetector/algorithm"
import { detectChanges } from "../shared/changeDetector/workerStarter"
import { dbTypeStore, resolveMeasureUnitForDb } from "../shared/dbTypes"
import { measureNameToLabel } from "../shared/metricsMapping"
import { Delta } from "../util/Delta"
import { toColor } from "../util/colors"
import { MAIN_METRICS, MAIN_METRICS_SET } from "../util/mainMetrics"
import { Accident, AccidentKind, AccidentsConfigurator } from "./accidents/AccidentsConfigurator"
import { getMedianScale, scaleToMedian } from "../components/settings/transforms/scaling"
import { exponentialSmoothingWithAlphaInference } from "../components/settings/transforms/smoothing"
import { createComponentState, updateComponentState } from "./componentState"
import { configureQueryFilters, createFilterObservable, FilterConfigurator } from "./filter"
import { combineLatest, fromFetchWithRetryAndErrorHandling, refToObservable } from "./rxjs"
import { removeOutliers } from "../components/settings/transforms/outliers"
import { getBasicInfo, getBuildId } from "../components/common/sideBar/InfoSidebarPerformance"
import { useDarkModeStore } from "../shared/useDarkModeStore"
import { captureMatchedSelectedPoint, useSelectedPointStore } from "../shared/selectedPointStore"

export type TooltipTrigger = "item" | "axis" | "none"

export class MeasureConfigurator implements DataQueryConfigurator, ChartConfigurator, FilterConfigurator {
  readonly data = shallowRef<string[]>([])
  private readonly _selected = shallowRef<string[] | string | null>(null)
  readonly state = createComponentState()

  readonly showAllMetrics = ref(false)

  createObservable(): Observable<unknown> {
    return combineLatest([refToObservable(this.selected, true), refToObservable(toRef(useSettingsStore(), "stdDevInterval"))])
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

    updateComponentState(
      combineLatest([createFilterObservable(serverConfigurator, filters), refToObservable(this.showAllMetrics)])
        [debounce](100)
        [distinctUntilChanged](deepEqual)
        [switchMap](() => {
          const loadMeasureListUrl = getLoadMeasureListUrl(serverConfigurator, filters)
          if (loadMeasureListUrl == null) {
            return ColdObservable.from([null])
          }

          this.state.loading = true
          return fromFetchWithRetryAndErrorHandling<string[]>(loadMeasureListUrl)
        }),
      this.state
    ).subscribe((data) => {
      if (data == null) {
        return
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

    configureQuery(measureNames, query, configuration, this.skipZeroValues, true)
    configuration.measures = measureNames
    configuration.addChartConfigurator(this)
    return true
  }

  configureChart(data: DataQueryResult, configuration: DataQueryExecutorConfiguration): Promise<ECBasicOption> {
    return configureChart(configuration, data, this.chartType, "auto", this.symbolOptions)
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
    }
    return false
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

  const fieldPrefix = serverConfigurator.table === "measure" ? "" : "measures"

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
    private readonly valueUnit: ValueUnit = "auto",
    readonly symbolOptions: SymbolOptions = {},
    readonly accidentsConfigurator: AccidentsConfigurator | null = null,
    readonly toolTipTrigger: TooltipTrigger,
    private readonly betterDirection: BetterDirection = "lower",
    // Comparison consumers do not need deviation data.
    private readonly withStdDevInterval: boolean = false
  ) {}

  createObservable(): Observable<unknown> {
    const inputs = [refToObservable(this.skipZeroValues), refToObservable(this.measures)]
    if (this.withStdDevInterval) {
      inputs.push(refToObservable(toRef(useSettingsStore(), "stdDevInterval")))
    }
    return combineLatest(inputs)
  }

  configureQuery(query: DataQuery, configuration: DataQueryExecutorConfiguration): boolean {
    if (this.measures.value.length === 0) {
      return false
    }
    configureQuery(this.measures.value, query, configuration, this.skipZeroValues.value, this.withStdDevInterval)
    configuration.addChartConfigurator(this)
    configuration.measures = this.measures.value
    return true
  }

  configureChart(data: DataQueryResult, configuration: DataQueryExecutorConfiguration): Promise<ECBasicOption> {
    return configureChart(configuration, data, this.chartType, this.valueUnit, this.symbolOptions, this.accidentsConfigurator, this.toolTipTrigger, this.betterDirection)
  }
}

// Only fetch implicit companions; explicitly selected deviations remain ordinary lines.
function getStdDevCompanions(measureNames: readonly string[], withStdDevInterval: boolean): Map<string, string> {
  const companions = new Map<string, string>()
  if (withStdDevInterval && useSettingsStore().stdDevInterval) {
    for (const measure of measureNames) {
      const companion = getStdDevMeasureName(measure)
      if (companion != null && !measureNames.includes(companion)) companions.set(companion, measure)
    }
  }
  return companions
}

function configureQuery(measureNames: string[], query: DataQuery, configuration: DataQueryExecutorConfiguration, skipZeroValues: boolean, withStdDevInterval: boolean): void {
  // stable order of series (UI) and fields in query (caching)
  measureNames.sort((a, b) => collator.compare(a, b))

  const stdDevCompanions = getStdDevCompanions(measureNames, withStdDevInterval)
  const queriedMeasureNames = [...measureNames, ...stdDevCompanions.keys()]
  const ownerMeasureNames = queriedMeasureNames.map((it) => stdDevCompanions.get(it) ?? it)

  query.insertField(
    {
      n: "t",
      sql: "toUnixTimestamp(generated_time)*1000",
    },
    0
  )

  // we cannot request several measures in one SQL query - for each measure separate SQl query with filter by measure name
  const field: DataQueryDimension = { n: "" }
  query.insertField(field, 1)

  if (!(query.db === "fleet" && query.table === "report")) {
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
      return queriedMeasureNames.length
    },
    mutate(index: number): void {
      const measure = queriedMeasureNames[index]

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
      } else {
        field.n = "measures"
        if (measure.endsWith(".end")) {
          field.subName = "end"
          field.sql = "(measures.start+measures.value)"
        } else {
          field.subName = "value"
        }

        addFilter({ f: "measures.name", v: measure.endsWith(".end") ? measure.slice(0, measure.length - ".end".length) : measure })
        valueFieldName = "measures.value"
      }

      // Keep zero deviations: dropping one could shift the pairing of repeated measurements.
      if (skipZeroValues && !stdDevCompanions.has(measure)) {
        addFilter({ f: valueFieldName, o: "!=", v: 0 })
      }
    },
    getSeriesName(index: number): string {
      // Companions share the mean's name without changing single-measure chart labels.
      return measureNames.length > 1 ? ownerMeasureNames[index] : ""
    },
    getMeasureName(index: number): string {
      return queriedMeasureNames[index]
    },
    getOwnerMeasureName(index: number): string {
      return ownerMeasureNames[index]
    },
  })

  if (query.order != null) {
    throw new Error("order must be configured only by MetricLoader")
  }
  query.order = "t"
}

function getItemStyleForSeries(accidentConfigurator: AccidentsConfigurator | null, detectedChanges = new Map<string, DetectedChange>()) {
  return {
    color(seriesIndex: CallbackDataParams): ZRColor {
      if (useSelectedPointStore().selectedPoint != undefined && getBasicInfo(seriesIndex, "ms").buildId.toString() == useSelectedPointStore().selectedPoint) {
        captureMatchedSelectedPoint(seriesIndex)
        return getSelectedPointColor()
      }
      const accidents = accidentConfigurator?.getAccidents(seriesIndex.value as string[])
      if (accidents == null || accidents.length === 0) {
        const classification = detectedChanges.get(JSON.stringify(seriesIndex.value))?.classification
        if (classification == ChangePointClassification.DEGRADATION) {
          return "#cc0000"
        } else if (classification == ChangePointClassification.OPTIMIZATION) {
          return "#009900"
        } else if (classification == ChangePointClassification.NEUTRAL) {
          // A detected change with no good/bad direction: draw the arrow in grey, not red/green.
          return "#808080"
        } else if (classification == ChangePointClassification.NO_CHANGE) {
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

function isChangeDetected(detectedChanges: Map<string, DetectedChange>, value: string[]) {
  const classification = detectedChanges.get(JSON.stringify(value))?.classification
  return classification != undefined && classification != ChangePointClassification.NO_CHANGE
}

class MergeResults {
  constructor(
    readonly data: DataQueryResult,
    private readonly idToSeriesName: Map<number, string>,
    private readonly idToMeasureName: Map<number, string>,
    private readonly idToOwnerMeasureName: Map<number, string>
  ) {}

  getSeriesName(index: number): string {
    return this.idToSeriesName.get(index) as string
  }

  getMeasureName(index: number): string {
    return this.idToMeasureName.get(index) as string
  }

  /** The measure the series was fetched for - the same as {@link getMeasureName} unless it is a companion. */
  getOwnerMeasureName(index: number): string {
    return this.idToOwnerMeasureName.get(index) as string
  }

  /** Whether the series was fetched on another measure's behalf rather than as a measurement of its own. */
  isCompanion(index: number): boolean {
    return this.getOwnerMeasureName(index) !== this.getMeasureName(index)
  }
}

// Aliased series (e.g. a project's old and new name) each arrive pre-sorted by time, but covering
// different time ranges. Concatenating their columns end-to-end can leave the combined series out
// of chronological order, and the line chart connects points in array order, not by x value - so an
// unsorted merge draws stray lines jumping between the end of one chunk and the start of another.
function sortColumnsByTime(seriesData: (string | number)[][]): (string | number)[][] {
  const time = seriesData[0]
  if (time == undefined || time.length <= 1) {
    return seriesData
  }
  const order = time.map((_, index) => index).toSorted((a, b) => (time[a] as number) - (time[b] as number))
  return seriesData.map((column) => (column.length === time.length ? order.map((index) => column[index]) : column))
}

export function mergeSeries(dataList: (string | number)[][][], configuration: DataQueryExecutorConfiguration) {
  const mergedDataList: DataQueryResult = []
  const seriesIdsToIndex = new Map<string, number>()
  const seriesIdToSeriesName = new Map<number, string>()
  const seriesIdToMeasureName = new Map<number, string>()
  const seriesIdToOwnerMeasureName = new Map<number, string>()
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
    const id = getSeriesId(measureName, seriesName)
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
      seriesIdToOwnerMeasureName.set(newId, configuration.ownerMeasureNames[dataIndex] ?? measureName)
    }
  }
  for (const [index, seriesData] of mergedDataList.entries()) {
    mergedDataList[index] = sortColumnsByTime(seriesData)
  }
  return new MergeResults(mergedDataList, seriesIdToSeriesName, seriesIdToMeasureName, seriesIdToOwnerMeasureName)
}

type PointKey = string | number

// Several builds can report in the same second. Some databases have no build column.
function getPointKeys(seriesData: (string | number)[][]): PointKey[] {
  const timestamps = (seriesData[0] ?? []) as number[]
  const buildIds = getBuildId(seriesData)
  return buildIds == undefined ? timestamps : timestamps.map((timestamp, index) => `${buildIds[index]}@${timestamp}`)
}

function collectStdDevCompanions(mergeResults: MergeResults): Map<string, (string | number)[][]> {
  const companions = new Map<string, (string | number)[][]>()
  for (const [index, data] of mergeResults.data.entries()) {
    if (mergeResults.isCompanion(index)) {
      companions.set(getSeriesId(mergeResults.getOwnerMeasureName(index), mergeResults.getSeriesName(index)), data)
    }
  }
  return companions
}

// Pair before removing outliers. Repeated build/timestamp keys must have equal counts on both sides,
// otherwise a skipped mean could receive a different measurement's deviation.
function resolveStdDevsPerPoint(seriesData: (string | number)[][], companion: (string | number)[][] | undefined): number[] | undefined {
  if (companion == undefined) return undefined
  const deviations = new Map<PointKey, number[]>()
  for (const [index, key] of getPointKeys(companion).entries()) {
    const value = companion[1]?.[index]
    if (typeof value !== "number" || !Number.isFinite(value)) continue
    const values = deviations.get(key) ?? []
    values.push(value)
    deviations.set(key, values)
  }
  const keys = getPointKeys(seriesData)
  const counts = new Map<PointKey, number>()
  for (const key of keys) counts.set(key, (counts.get(key) ?? 0) + 1)
  const positions = new Map<PointKey, number>()
  return keys.map((key) => {
    const position = positions.get(key) ?? 0
    positions.set(key, position + 1)
    const values = deviations.get(key)
    return values != undefined && values.length === counts.get(key) ? values[position] : 0
  })
}

function getSelectedPointColor() {
  return useDarkModeStore().darkMode ? "white" : "black"
}

async function configureChart(
  configuration: DataQueryExecutorConfiguration,
  dataList: DataQueryResult,
  chartType: ChartType,
  valueUnit: ValueUnit = "auto",
  symbolOptions: SymbolOptions = {},
  accidentsConfigurator: AccidentsConfigurator | null = null,
  tooltipTrigger: TooltipTrigger = "item",
  betterDirection: BetterDirection = "lower"
): Promise<LineChartOptions | ScatterChartOptions> {
  const series = new Array<LineSeriesOption | ScatterSeriesOption>()

  const dataset: DatasetOption[] = []

  //merge series with the same name
  const mergeResults = mergeSeries(dataList, configuration)

  const settings = useSettingsStore()
  // Only measured values determine staleness.
  const chartLastTimestamp = getChartLastTimestamp(mergeResults.data.filter((_, index) => !mergeResults.isCompanion(index)).map((seriesData) => (seriesData[0] ?? []) as number[]))
  const measureUnits: MeasureUnit[] = []
  const stdDevByOwnerSeriesId = collectStdDevCompanions(mergeResults)
  // eslint-disable-next-line prefer-const
  for (let [dataIndex, seriesData] of mergeResults.data.entries()) {
    // A companion is drawn as a band around its mean rather than as a line.
    const isCompanion = mergeResults.isCompanion(dataIndex)
    // eslint-disable-next-line @typescript-eslint/no-unnecessary-condition
    if (isCompanion || seriesData[1] == undefined) {
      //we need to push even empty dataset otherwise it will be out of sync with series and plot will be empty
      dataset.push({
        source: isCompanion ? [] : seriesData,
        sourceHeader: false,
      })
      continue
    }

    // we take only the last type of the metric since it is not clear how to show different types
    const storedType = seriesData.length > 3 ? (seriesData[3].at(-1) as string) : undefined

    // Staleness is judged on what the series reported, not on what is left after filtering: dropping a
    // trailing outlier would otherwise end the line early and fake a stop. `chartLastTimestamp` comes
    // from the unfiltered data too, so the two stay on the same footing.
    const reportedTimestamps = seriesData[0] as number[]

    const measureName = mergeResults.getMeasureName(dataIndex)
    const seriesName = mergeResults.getSeriesName(dataIndex)
    const seriesId = getSeriesId(measureName, seriesName)
    let stdDevs = resolveStdDevsPerPoint(seriesData, stdDevByOwnerSeriesId.get(seriesId))

    if (settings.removeOutliers) {
      // Filter the deviations as one more column, then remove it before adding tooltip columns.
      seriesData = removeOutliers(stdDevs == undefined ? seriesData : [...seriesData, stdDevs])
      if (stdDevs != undefined) stdDevs = seriesData.pop() as number[]
    }

    if (settings.smoothing) {
      const smoothedData = exponentialSmoothingWithAlphaInference(seriesData[1] as number[])
      seriesData.push(smoothedData)
    }

    const deltaValues = Delta.calculateDeltas(seriesData[1] as number[], getBuildId(seriesData) as number[])
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    //@ts-expect-error
    seriesData.push(deltaValues)

    // What the band has to be multiplied by to end up on the same scale as the line it wraps.
    let valueScale = 1
    if (settings.scaling) {
      seriesData.push(seriesData[1])
      valueScale = getMedianScale(seriesData[1] as number[])
      seriesData[1] = scaleToMedian(seriesData[1] as number[], valueScale)
    }

    let detectedChanges = new Map<string, DetectedChange>()
    if (settings.detectChanges) {
      detectedChanges = await detectChanges(seriesData, betterDirection)
    }

    let isNotEmpty = false
    for (const data of seriesData) {
      isNotEmpty ||= data.length > 0
    }
    if (isNotEmpty) {
      const seriesLayoutBy = "row"
      const datasetIndex = dataIndex
      const seriesUnit = resolveMeasureUnitForDb(measureName, { storedType, valueUnit })
      measureUnits.push(seriesUnit)
      const xAxisName = seriesUnit === "milliseconds" || seriesUnit === "nanoseconds" ? "time" : "count"
      series.push({
        selectedMode: "single",
        select: {
          itemStyle: {
            color: getSelectedPointColor(),
          },
        },
        // formatter is detected by measure name - that's why series id is specified (see usages of seriesId)
        id: seriesId,
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
          // The "arrow" symbol points up by default; rotate 180° to point down when the value fell.
          return detectedChanges.get(JSON.stringify(value))?.direction === "down" ? 180 : 0
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
        markLine: getStaleSeriesMarkLine(reportedTimestamps, chartLastTimestamp),
      })
      if (settings.smoothing) {
        series.push({
          // formatter is detected by measure name - that's why series id is specified (see usages of seriesId)
          id: seriesId + "smoothed",
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

      if (stdDevs != undefined) {
        const band = buildStdDevBand({
          timestamps: seriesData[0] as number[],
          values: seriesData[1] as number[],
          stdDevs,
          valueScale,
        })
        if (band != null) {
          series.push(...getStdDevBandSeries(band, seriesId, seriesName))
        }
      }
    }

    dataset.push({
      source: seriesData,
      sourceHeader: false,
    })
  }

  // While scaling, axis values are baseline ratios, so they render as plain numbers.
  const axisUnit: MeasureUnit = settings.scaling ? "counter" : reduceToAxisUnit(measureUnits)
  const formatter: (value: number) => string = (value) => formatMeasureValue(value, axisUnit)
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
        // Clip instead of filtering on the band's width, which is not its y position.
        filterMode: "none",
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
