import { CallbackDataParams, OptionDataValue } from "echarts/types/src/util/types"
import { Accident, AccidentKind, AccidentsConfigurator } from "../../configurators/AccidentsConfigurator"
import { appendLineWithIcon, getLeftArrow, getRightArrow, getWarningIcon } from "../../shared/popupIcons"
import { Delta, findDeltaInData, getDifferenceString } from "../../util/Delta"
import { DataQueryExecutor, DataQueryResult } from "../common/DataQueryExecutor"
import { timeFormat, ValueUnit } from "../common/chart"
import { DataQueryExecutorConfiguration } from "../common/dataQuery"
import { LineChartOptions } from "../common/echarts"
import { durationAxisPointerFormatter, isDurationFormatterApplicable, nsToMs, numberFormat, timeFormatWithoutSeconds } from "../common/formatter"
import { useSettingsStore } from "../settings/settingsStore"
import { PerformanceChartManager } from "./PerformanceChartManager"

export class PerformanceLineChartVM {
  private settings = useSettingsStore()
  private getFormatter(isMs: boolean) {
    return (params: CallbackDataParams) => {
      const element = document.createElement("div")
      const data = params.value as (OptionDataValue | Delta)[]
      const dateMs = data[0]
      let type = data[3]
      if (type != "c" && type != "d") {
        type = isDurationFormatterApplicable(data[2] as string) ? "d" : "c"
      }
      const durationMs = this.settings.scaling ? data.at(-1) : data[1]
      element.append(
        durationAxisPointerFormatter(isMs ? (durationMs as number) : (durationMs as number) / 1000 / 1000, type),
        document.createElement("br"),
        timeFormatWithoutSeconds.format(dateMs as number)
      )

      element.append(document.createElement("br"))
      element.append(`${params.seriesName}`)
      const accidents = this.accidentsConfigurator?.getAccidents(data as string[]) ?? null
      if (accidents != null) {
        for (const accident of accidents) {
          appendLineWithIcon(element, getWarningIcon(), this.getAccidentMessage(accident))
        }
      }

      const delta = findDeltaInData(data)
      if (delta != undefined) {
        if (delta.prev != null) {
          appendLineWithIcon(element, getLeftArrow(), getDifferenceString(durationMs as number, delta.prev, isMs, type))
        }
        if (delta.next != null) {
          appendLineWithIcon(element, getRightArrow(), getDifferenceString(durationMs as number, delta.next, isMs, type))
        }
      }
      return element
    }
  }
  private getAccidentMessage(accident: Accident): string {
    return accident.kind == AccidentKind.InferredRegression || accident.kind == AccidentKind.InferredImprovement
      ? accident.reason
      : "Known " + accident.kind.toLowerCase() + ": " + accident.reason
  }

  private accidentsConfigurator: AccidentsConfigurator | null
  constructor(
    private readonly eChart: PerformanceChartManager,
    private readonly dataQuery: DataQueryExecutor,
    valueUnit: ValueUnit,
    accidentsConfigurator: AccidentsConfigurator | null,
    private readonly legendFormatter: (name: string) => string
  ) {
    this.legendFormatter = legendFormatter
    this.accidentsConfigurator = accidentsConfigurator
    const isMs = valueUnit == "ms"
    this.eChart.chart.showLoading("default", { showSpinner: false })
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
        },
      },
      animation: false,
      grid: {
        left: 8,
        right: 8,
        bottom: 16,
        containLabel: true,
      },
      // @ts-expect-error bug in echarts types
      tooltip: {
        show: true,
        trigger: "item",
        enterable: true,
        axisPointer: {
          type: "cross",
        },
        renderMode: "html",
        position: (pointerCoords, _, tooltipElement) => {
          const [pointerLeft, pointerTop] = pointerCoords
          const element = tooltipElement as HTMLDivElement
          const chartRect = this.eChart.chart.getDom().getBoundingClientRect()
          const isOverflowWindow = chartRect.left + pointerLeft + element.offsetWidth > chartRect.right

          return [isOverflowWindow ? pointerLeft - element.offsetWidth : pointerLeft, pointerTop - element.clientHeight - 10]
        },
        // Formatting
        formatter: this.getFormatter(isMs),
        valueFormatter(it) {
          return numberFormat.format(isMs ? (it as number) : nsToMs(it as number)) + " ms"
        },
        // Styling
        extraCssText: "user-select: text",
        borderColor: "#E5E7EB",
        padding: [6, 8],
        textStyle: {
          fontSize: 12,
        },
      },
      xAxis: {
        type: "time",
        max: new Date(),
        axisPointer: {
          snap: false,
          label: {
            formatter(data) {
              return timeFormat.format(data.value as number)
            },
          },
        },
      },
      yAxis: {
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
        this.eChart.chart.showLoading("default", { showSpinner: false })
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
            ],
            formatter(name: string): string {
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
          .catch((error) => {
            console.error(error)
          })
      }
    })
  }

  dispose(): void {
    this.eChart.dispose()
  }
}
