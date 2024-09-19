import { ECElementEvent } from "echarts/core"
import { CallbackDataParams, OptionDataValue } from "echarts/types/src/util/types"
import { Accident, AccidentKind, AccidentsConfigurator } from "../../configurators/accidents/AccidentsConfigurator"
import { measureNameToLabel } from "../../shared/metricsMapping"
import { appendLineWithIcon, getDiffIcon, getLeftArrow, getRightArrow, getWarningIcon } from "../../shared/popupIcons"
import { Delta, findDeltaInData, getDifferenceString } from "../../util/Delta"
import { DataQueryExecutor, DataQueryResult } from "../common/DataQueryExecutor"
import { timeFormat, ValueUnit } from "../common/chart"
import { DataQueryExecutorConfiguration } from "../common/dataQuery"
import { LineChartOptions } from "../common/echarts"
import { durationAxisPointerFormatter, isDurationFormatterApplicable, nsToMs, numberFormat, timeFormatWithoutSeconds } from "../common/formatter"
import { InfoSidebar } from "../common/sideBar/InfoSidebar"
import { getFullBuildId, getInfoDataFrom } from "../common/sideBar/InfoSidebarPerformance"
import { useSettingsStore } from "../settings/settingsStore"
import { ChartManager } from "./ChartManager"
import { useDarkModeStore } from "../../shared/useDarkModeStore"

class ClickedValue {
  constructor(
    public readonly timestamp: number,
    public readonly value: number
  ) {}
}

export class LineChartVM {
  private settings = useSettingsStore()
  private lastParams: CallbackDataParams[] | CallbackDataParams | null = null
  private lastClickedValue = new Map<string, ClickedValue>()

  private getFormatter(isMs: boolean) {
    return (params: CallbackDataParams[] | CallbackDataParams) => {
      this.lastParams = params
      const paramsArray = Array.isArray(params) ? params : [params]
      return paramsArray.length == 1 ? this.getElementForSingleSerie(isMs, paramsArray[0]) : this.getElementForMultipleSeries(isMs, paramsArray)
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
          } else {
            this.lastClickedValue.set(params.seriesName, new ClickedValue(params.value[0] as number, params.value[1] as number))
          }
        }
        const infoData = getInfoDataFrom(this.lastParams ?? params, valueUnit, accidentsConfigurator, chartManager.chart.getDataURL({ type: "png" }))
        sidebarVm.show(infoData)
      }
    }
  }

  private isMacOS() {
    const userAgent = navigator.userAgent.toLowerCase()
    return userAgent.includes("mac")
  }

  private getElementForMultipleSeries(isMs: boolean, params: CallbackDataParams[]) {
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
      seriesName.append(measureNameToLabel(param.seriesName as string))
      element.append(seriesName, document.createElement("br"))
      const data = param.value as (OptionDataValue | Delta)[]
      const type = this.getType(data)
      const durationMs = this.settings.scaling ? data.at(-1) : data[1]
      element.append(durationAxisPointerFormatter(isMs ? (durationMs as number) : (durationMs as number) / 1000 / 1000, type), document.createElement("br"))
      this.appendAccidentInfo(data, element)
      this.appendDelta(data, element, durationMs as number, isMs, type)
    }
    return element
  }

  private getType(data: (OptionDataValue | Delta)[]) {
    let type = data[3]
    if (type != "c" && type != "d") {
      type = isDurationFormatterApplicable(data[2] as string) ? "d" : "c"
    }
    return type
  }

  private getElementForSingleSerie(isMs: boolean, params: CallbackDataParams) {
    const element = document.createElement("div")
    const data = params.value as (OptionDataValue | Delta)[]
    const dateMs = data[0]
    const type = this.getType(data)
    const durationMs = this.settings.scaling ? data.at(-1) : data[1]
    element.append(durationAxisPointerFormatter(isMs ? (durationMs as number) : (durationMs as number) / 1000 / 1000, type), document.createElement("br"))
    element.append(timeFormatWithoutSeconds.format(dateMs as number), document.createElement("br"))
    element.append(measureNameToLabel(params.seriesName as string))
    const buildId = getFullBuildId(params)
    if (buildId != undefined) {
      element.append(document.createElement("br"), buildId)
    }
    this.appendAccidentInfo(data, element)
    this.appendDelta(data, element, durationMs as number, isMs, type)
    if (params.seriesName && this.lastClickedValue.get(params.seriesName)) {
      this.appendDeltaWithLastClicked(durationMs as number, this.lastClickedValue.get(params.seriesName)?.value as number, element, isMs, type)
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

  private appendDelta(data: (OptionDataValue | Delta)[], element: HTMLDivElement, durationMs: number, isMs: boolean, type: string) {
    const delta = findDeltaInData(data)
    if (delta != undefined) {
      if (delta.prev != null) {
        appendLineWithIcon(element, getLeftArrow(), getDifferenceString(durationMs, delta.prev, isMs, type))
      }
      if (delta.next != null) {
        appendLineWithIcon(element, getRightArrow(), getDifferenceString(durationMs, delta.next, isMs, type))
      }
    }
  }

  private appendDeltaWithLastClicked(durationMs: number, lastClicked: number, element: HTMLDivElement, isMs: boolean, type: string) {
    appendLineWithIcon(element, getDiffIcon(), getDifferenceString(lastClicked, durationMs, isMs, type))
  }

  private getAccidentMessage(accident: Accident): string {
    return accident.kind == AccidentKind.InferredRegression || accident.kind == AccidentKind.InferredImprovement
      ? accident.reason
      : "Known " + accident.kind.toLowerCase() + ": " + accident.reason
  }

  private accidentsConfigurator: AccidentsConfigurator | null

  constructor(
    private readonly eChart: ChartManager,
    private readonly dataQuery: DataQueryExecutor,
    valueUnit: ValueUnit,
    accidentsConfigurator: AccidentsConfigurator | null,
    private readonly legendFormatter: (name: string) => string
  ) {
    this.legendFormatter = legendFormatter
    this.accidentsConfigurator = accidentsConfigurator
    const isMs = valueUnit == "ms"
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
        },
      },
      animation: false,
      grid: {
        left: 8,
        right: 8,
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
        formatter: this.getFormatter(isMs),
        valueFormatter(it) {
          return numberFormat.format(isMs ? (it as number) : nsToMs(it as number)) + " ms"
        },
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
    return this.dataQuery.subscribe((data: DataQueryResult | null, configuration: DataQueryExecutorConfiguration, isLoading) => {
      if (isLoading || data == null) {
        this.eChart.chart.showLoading("default", useDarkModeStore().darkMode ? { maskColor: "#121212", showSpinner: false, textColor: "#D1D5DB" } : { showSpinner: false })
        return
      }
      this.eChart.chart.hideLoading()
      const formatter = this.legendFormatter
      this.eChart.chart.setOption(
        {
          title: {
            show: data.flat(3).length === 0,
            text: "No data",
            subtext: "Please check selectors: machine, branch, time range",
            left: "center",
            top: "center",
            textStyle: {
              fontSize: 20,
              fontWeight: "normal",
              color: "#6B7280",
            },
          },
          legend: {
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
      for (const it of configuration.getChartConfigurators()) {
        it.configureChart(data, configuration)
          .then((options) => {
            this.eChart.updateChart(options)
          })
          .catch((error: unknown) => {
            console.error(error)
          })
      }
    })
  }

  dispose(): void {
    this.eChart.dispose()
  }
}
