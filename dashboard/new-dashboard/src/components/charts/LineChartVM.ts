import { ECElementEvent } from "echarts/core"
import type { DefaultLabelFormatterCallbackParams as CallbackDataParams } from "echarts"
import { watch } from "vue"
import type { OptionDataValue } from "../../shared/echarts-types"
import { Accident, AccidentKind, AccidentsConfigurator } from "../../configurators/accidents/AccidentsConfigurator"
import { measureNameToLabel } from "../../shared/metricsMapping"
import { appendLineWithIcon, getDiffIcon, getLeftArrow, getRightArrow, getWarningIcon } from "../../shared/popupIcons"
import { Delta, findDeltaInData, getDifferenceString } from "../../util/Delta"
import { DataQueryExecutor, DataQueryResult } from "../common/DataQueryExecutor"
import { timeFormat, ValueUnit } from "../common/chart"
import { DataQueryExecutorConfiguration } from "../common/dataQuery"
import { LineChartOptions } from "../common/echarts"
import { formatMeasureValue, MeasureUnit, reduceToAxisUnit, resolveMeasureUnit, timeFormatWithoutSeconds } from "../common/formatter"
import { InfoSidebar } from "../common/sideBar/InfoSidebar"
import { getFullBuildId, getInfoDataFrom } from "../common/sideBar/InfoSidebarPerformance"
import { consumeMatchedSelectedPoints } from "../../shared/selectedPointStore"
import { useSettingsStore } from "../settings/settingsStore"
import { ChartManager } from "./ChartManager"
import { useDarkModeStore } from "../../shared/useDarkModeStore"
import { HoverFadeController } from "./hoverFade"
import { exportChartMetricsAsYaml } from "./chartExport"

class ClickedValue {
  constructor(
    public readonly timestamp: number,
    public readonly value: number
  ) {}
}

interface AutoOpenConfig {
  sidebarVm: InfoSidebar
  valueUnit: ValueUnit
  accidentsConfigurator: AccidentsConfigurator | null
}

interface SeriesOption {
  datasetIndex?: number
}

export class LineChartVM {
  private readonly settings = useSettingsStore()
  private lastParams: CallbackDataParams[] | CallbackDataParams | null = null
  private readonly lastClickedValue = new Map<string, ClickedValue>()
  private readonly hasDataCallback?: (hasData: boolean) => void
  // Track if data has been loaded (for accident marker refresh)
  private lastData: DataQueryResult | null = null
  private hoverFade?: HoverFadeController
  private autoOpenConfig: AutoOpenConfig | null = null
  private hasAutoOpened = false
  private getFormatter() {
    return (params: CallbackDataParams[] | CallbackDataParams) => {
      this.lastParams = params
      const paramsArray = Array.isArray(params) ? params : [params]
      return paramsArray.length == 1 ? this.getElementForSingleSerie(paramsArray[0]) : this.getElementForMultipleSeries(paramsArray)
    }
  }

  public getOnClickHandler(sidebarVm: InfoSidebar, chartManager: ChartManager, valueUnit: ValueUnit, accidentsConfigurator: AccidentsConfigurator | null) {
    return (params: ECElementEvent) => {
      const useMetaKey = this.isMacOS() ? params.event?.event.metaKey : params.event?.event.ctrlKey
      if (useMetaKey) {
        chartManager.chart.dispatchAction({
          type: "legendUnSelect",
          name: params.seriesName,
        })
      } else {
        if (params.seriesName && Array.isArray(params.value)) {
          // if the same value is clicked again remove
          if (this.lastClickedValue.get(params.seriesName)?.timestamp == (params.value[0] as number)) {
            this.lastClickedValue.delete(params.seriesName)
            sidebarVm.close()
          } else {
            this.lastClickedValue.set(params.seriesName, new ClickedValue(params.value[0] as number, params.value[1] as number))
            const seriesContext = this.extractSeriesContext(chartManager, params)
            const infoData = getInfoDataFrom(this.lastParams ?? params, valueUnit, accidentsConfigurator, chartManager.chart.getDataURL({ type: "png" }), seriesContext)
            sidebarVm.show(infoData)
          }
        }
      }
    }
  }

  private extractSeriesContext(chartManager: ChartManager, params: ECElementEvent): { seriesValues: number[] | undefined; pointIndex: number | undefined } {
    try {
      const option = chartManager.chart.getOption() as { series?: unknown[]; dataset?: unknown[] }
      const seriesList = (option.series ?? []) as SeriesOption[]
      const datasets = (option.dataset ?? []) as { source?: unknown[] }[]
      // When scaling is on, MeasureConfigurator replaces source[1] with scale-to-median values
      // and appends the original unscaled values at the end. Use those so the heuristic stays
      // on the same scale as InfoData.rawValue / previousValue (both unscaled).
      return { seriesValues: this.getValuesColumn(seriesList, datasets, params.seriesIndex ?? 0), pointIndex: params.dataIndex }
    } catch {
      return { seriesValues: undefined, pointIndex: params.dataIndex }
    }
  }

  private getValuesColumn(seriesList: SeriesOption[], datasets: { source?: unknown[] }[], seriesIndex: number): number[] | undefined {
    const datasetIndex = seriesList[seriesIndex]?.datasetIndex ?? seriesIndex
    const source = datasets[datasetIndex]?.source
    if (!Array.isArray(source)) return undefined
    const valuesColumn = this.settings.scaling ? source.at(-1) : source[1]
    if (!Array.isArray(valuesColumn)) return undefined
    return (valuesColumn as unknown[]).map((v) => (typeof v === "number" && Number.isFinite(v) ? v : Number.NaN))
  }

  public enableSidebarAutoOpen(config: AutoOpenConfig): void {
    this.autoOpenConfig = config
    this.hasAutoOpened = false
  }

  private tryAutoOpenSidebar(): void {
    const matches = consumeMatchedSelectedPoints()
    const config = this.autoOpenConfig
    if (this.hasAutoOpened || config == null || config.sidebarVm.visible.value || matches.length === 0) return

    const chart = this.eChart.chart
    const option = chart.getOption() as { series?: unknown[]; dataset?: unknown[] }
    const params: CallbackDataParams | CallbackDataParams[] = matches.length === 1 ? matches[0] : matches
    const infoData = getInfoDataFrom(params, config.valueUnit, config.accidentsConfigurator, chart.getDataURL({ type: "png" }), {
      seriesValues: this.getValuesColumn((option.series ?? []) as SeriesOption[], (option.dataset ?? []) as { source?: unknown[] }[], matches[0].seriesIndex ?? 0),
      pointIndex: matches[0].dataIndex,
    })
    config.sidebarVm.show(infoData)
    this.hasAutoOpened = true
  }

  private isMacOS() {
    const userAgent = navigator.userAgent.toLowerCase()
    return userAgent.includes("mac")
  }

  private getElementForMultipleSeries(params: CallbackDataParams[]) {
    const element = document.createElement("div")
    const dateMs = (params[0].value as (OptionDataValue | Delta)[])[0]
    element.append(timeFormatWithoutSeconds.format(dateMs as number), document.createElement("br"))
    if (this.settings.smoothing) params = params.filter((_, index) => index % 2 == 0)
    const buildId = getFullBuildId(params[0])
    if (buildId != undefined) {
      element.append(buildId, document.createElement("br"))
    }
    for (const param of params) {
      const seriesName = document.createElement("b")
      seriesName.append(measureNameToLabel(param.seriesName))
      element.append(seriesName, document.createElement("br"))
      const data = param.value as (OptionDataValue | Delta)[]
      const unit = this.resolveUnit(data)
      const durationMs = this.settings.scaling ? data.at(-1) : data[1]
      element.append(formatMeasureValue(durationMs as number, unit), document.createElement("br"))
      this.appendAccidentInfo(data, element)
      this.appendDelta(data, element, durationMs as number, unit)
    }
    return element
  }

  private resolveUnit(data: (OptionDataValue | Delta)[]): MeasureUnit {
    return resolveMeasureUnit(data[2] as string, { storedType: data[3] as string, valueUnit: this.valueUnit, scaling: this.settings.scaling })
  }

  private getElementForSingleSerie(params: CallbackDataParams) {
    const element = document.createElement("div")
    const data = params.value as (OptionDataValue | Delta)[]
    const dateMs = data[0]
    const unit = this.resolveUnit(data)
    const durationMs = this.settings.scaling ? data.at(-1) : data[1]
    element.append(formatMeasureValue(durationMs as number, unit), document.createElement("br"))
    element.append(timeFormatWithoutSeconds.format(dateMs as number), document.createElement("br"))
    element.append(measureNameToLabel(params.seriesName))
    const buildId = getFullBuildId(params)
    if (buildId != undefined) {
      element.append(document.createElement("br"), buildId)
    }
    this.appendAccidentInfo(data, element)
    this.appendDelta(data, element, durationMs as number, unit)
    if (params.seriesName && this.lastClickedValue.get(params.seriesName)) {
      this.appendDeltaWithLastClicked(durationMs as number, this.lastClickedValue.get(params.seriesName)?.value as number, element, unit)
    }
    return element
  }

  private appendAccidentInfo(data: (OptionDataValue | Delta)[], element: HTMLDivElement) {
    const accidents = this.accidentsConfigurator?.getAccidents(data as string[]) ?? null
    if (accidents != null) {
      for (const accident of accidents) {
        appendLineWithIcon(element, getWarningIcon(), this.getAccidentMessage(accident))
      }
    }
  }

  private appendDelta(data: (OptionDataValue | Delta)[], element: HTMLDivElement, durationMs: number, unit: MeasureUnit) {
    const delta = findDeltaInData(data)
    if (delta != undefined) {
      if (delta.prev != null) {
        appendLineWithIcon(element, getLeftArrow(), getDifferenceString(durationMs, delta.prev, unit))
      }
      if (delta.next != null) {
        appendLineWithIcon(element, getRightArrow(), getDifferenceString(durationMs, delta.next, unit))
      }
    }
  }

  private appendDeltaWithLastClicked(durationMs: number, lastClicked: number, element: HTMLDivElement, unit: MeasureUnit) {
    appendLineWithIcon(element, getDiffIcon(), getDifferenceString(lastClicked, durationMs, unit))
  }

  private getAccidentMessage(accident: Accident): string {
    return accident.kind == AccidentKind.InferredRegression || accident.kind == AccidentKind.InferredImprovement
      ? accident.reason
      : "Known " + accident.kind.toLowerCase() + ": " + accident.reason
  }

  private readonly accidentsConfigurator: AccidentsConfigurator | null

  constructor(
    private readonly eChart: ChartManager,
    private readonly dataQuery: DataQueryExecutor,
    private readonly valueUnit: ValueUnit,
    measures: string[],
    accidentsConfigurator: AccidentsConfigurator | null,
    private readonly legendFormatter: (name: string) => string,
    chartTitle: string,
    hasDataCallback?: (hasData: boolean) => void
  ) {
    this.accidentsConfigurator = accidentsConfigurator
    this.hasDataCallback = hasDataCallback
    // The axis-pointer label has no per-series context, so resolve a chart-level unit from the
    // chart's measures (falling back to the value-unit request alone when no measure is known),
    // matching the y-axis. A throughput chart then shows kB/s on the pointer, not a raw number.
    const axisUnit =
      measures.length > 0
        ? reduceToAxisUnit(measures.map((measure) => resolveMeasureUnit(measure, { valueUnit, scaling: this.settings.scaling })))
        : resolveMeasureUnit("", { valueUnit, scaling: this.settings.scaling })
    this.eChart.chart.showLoading("default", useDarkModeStore().darkMode ? { maskColor: "#121212", showSpinner: false, textColor: "#D1D5DB" } : { showSpinner: false })
    this.eChart.chart.setOption<LineChartOptions>({
      legend: {
        top: 0,
        left: 0,
        itemHeight: 3,
        itemWidth: 15,
        icon: "rect",
      },
      toolbox: {
        feature: {
          dataZoom: {
            yAxisIndex: false,
          },
          myTool: {
            show: true,
            title: "Full screen",
            icon: "path://M3.75 3.75v4.5m0-4.5h4.5m-4.5 0L9 9M3.75 20.25v-4.5m0 4.5h4.5m-4.5 0L9 15M20.25 3.75h-4.5m4.5 0v4.5m0-4.5L15 9m5.25 11.25h-4.5m4.5 0v-4.5m0 4.5L15 15",
            onclick() {
              /* eslint-disable @typescript-eslint/no-floating-promises */
              // noinspection JSIgnoredPromiseFromCall
              if (document.fullscreenElement) {
                document.exitFullscreen()
              } else {
                eChart.chartContainer.requestFullscreen()
              }
            },
          },
          myExportMetrics: {
            show: true,
            title: "Export metrics",
            icon: "path://M3 16.5v2.25A2.25 2.25 0 0 0 5.25 21h13.5A2.25 2.25 0 0 0 21 18.75V16.5M16.5 12 12 16.5m0 0L7.5 12m4.5 4.5V3",
            onclick: () => {
              exportChartMetricsAsYaml(eChart.chart, chartTitle, valueUnit)
            },
          },
        },
      },
      animation: false,
      grid: {
        left: 8,
        top: 65,
        right: 20,
        bottom: 16,
        containLabel: true,
      },
      tooltip: {
        show: true,
        enterable: true,
        axisPointer: {
          type: "cross",
        },
        renderMode: "html",
        position: (pointerCoords, _, tooltipElement) => {
          const [pointerLeft, pointerTop] = pointerCoords
          const element = tooltipElement as HTMLDivElement
          const chartRect = this.eChart.chart.getDom().getBoundingClientRect()

          const tooltipWidth = element.offsetWidth
          const tooltipHeight = element.clientHeight

          // Calculate initial positions
          let left = pointerLeft + 10
          let top = pointerTop - tooltipHeight - 10

          // Handle horizontal overflow
          const isOverflowRight = chartRect.left + left + tooltipWidth > chartRect.right
          if (isOverflowRight) {
            left = pointerLeft - tooltipWidth - 10
          }

          // Handle vertical overflow
          const isOverflowTop = chartRect.top + top < chartRect.top
          const isOverflowBottom = chartRect.top + top + tooltipHeight > chartRect.bottom
          if (isOverflowTop) {
            top = pointerTop + 10 // Position below the pointer if it overflows on top
          } else if (isOverflowBottom) {
            top = chartRect.bottom - tooltipHeight - 10 // Adjust to stay within the bottom edge
          }

          return [left, top]
        },
        // Formatting
        formatter: this.getFormatter(),
        valueFormatter: (it) => formatMeasureValue(it as number, axisUnit),
        // Styling
        padding: [6, 8],
        backgroundColor: useDarkModeStore().darkMode ? "#121212" : "white",
        borderColor: useDarkModeStore().darkMode ? "#4B5563" : "#E5E7EB",
        borderWidth: 0.3,
        textStyle: {
          color: useDarkModeStore().darkMode ? "#D1D5DB" : "#6B7280",
          fontSize: 12,
        },
      },
      xAxis: {
        type: "time",
        axisPointer: {
          snap: false,
          label: {
            backgroundColor: useDarkModeStore().darkMode ? "#121212" : "#6E7079",
            formatter(data) {
              return timeFormat.format(data.value as number)
            },
          },
        },
      },
      yAxis: {
        axisPointer: {
          label: {
            backgroundColor: useDarkModeStore().darkMode ? "#121212" : "#6E7079",
          },
        },
        min(value) {
          return useSettingsStore().flexibleYZero ? value.min * 0.9 : 0
        },
        type: "value",
        splitLine: {
          show: false,
        },
      },
    })
  }

  subscribe(): () => void {
    const chart = this.eChart.chart
    const dataUnsubscribe = this.dataQuery.subscribe((data: DataQueryResult | null, configuration: DataQueryExecutorConfiguration, isLoading) => {
      if (isLoading || data == null) {
        chart.showLoading("default", useDarkModeStore().darkMode ? { maskColor: "#121212", showSpinner: false, textColor: "#D1D5DB" } : { showSpinner: false })
        return
      }
      chart.hideLoading()

      // Track that data has been loaded (for accident marker refresh)
      this.lastData = data

      this.renderChart(data, configuration)
    })

    // Watch for accidents changes to re-render the chart without re-fetching data
    let accidentsUnwatch: (() => void) | null = null
    if (this.accidentsConfigurator != null) {
      accidentsUnwatch = watch(
        () => this.accidentsConfigurator?.value.value,
        () => {
          // Trigger a lightweight re-render by "refreshing" the series.
          // The itemStyle.color callback will read the updated accidents.
          if (this.lastData != null) {
            const option = this.eChart.chart.getOption()
            if (option["series"]) {
              this.eChart.chart.setOption({ series: option["series"] })
            }
          }
        }
      )
    }

    this.hoverFade = new HoverFadeController(this.eChart)

    return () => {
      dataUnsubscribe()
      accidentsUnwatch?.()
      this.hoverFade?.dispose()
    }
  }

  private renderChart(data: DataQueryResult, configuration: DataQueryExecutorConfiguration): void {
    consumeMatchedSelectedPoints()
    const chart = this.eChart.chart
    const hasData = data.flat(3).length > 0
    this.hasDataCallback?.(hasData)

    const formatter = this.legendFormatter
    chart.setOption(
      {
        legend: {
          bottom: null,
          type: "scroll",
          selector: [
            {
              type: "inverse",
              title: "inverse",
            },
            {
              type: "all",
              title: "enable all",
            },
          ],
          formatter(name: string): string {
            name = measureNameToLabel(name)
            if (formatter("test") != "") {
              return formatter(name)
            }
            return name
          },
        },
        toolbox: {
          top: 20,
          feature: {
            saveAsImage: {
              name: "plot",
            },
          },
        },
      },
      {
        replaceMerge: ["legend"],
      }
    )

    for (const configurator of configuration.getChartConfigurators()) {
      configurator
        .configureChart(data, configuration)
        .then((options) => {
          this.eChart.updateChart(options)
          this.hoverFade?.reapply()
          this.tryAutoOpenSidebar()
        })
        .catch((error: unknown) => {
          console.error(error)
        })
    }
  }

  dispose(): void {
    this.autoOpenConfig = null
    this.hoverFade?.dispose()
    this.eChart.dispose()
  }
}
