import { CallbackDataParams, OptionDataValue } from "echarts/types/src/util/types"
import { DataQueryExecutor } from "shared/src/DataQueryExecutor"
import { timeFormat, ValueUnit } from "shared/src/chart"
import { LineChartOptions } from "shared/src/echarts"
import { durationAxisPointerFormatter, nsToMs, numberFormat, timeFormatWithoutSeconds } from "shared/src/formatter"
import { ChartManager } from "./ChartManager"

export class LineChartVM {
  constructor(
    private readonly eChart: ChartManager,
    private readonly dataQuery: DataQueryExecutor,
    valueUnit: ValueUnit,
  ) {
    const isMs = valueUnit == "ms"
    this.eChart.chart.showLoading()
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
      // @ts-ignore
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
          const element = (tooltipElement as HTMLDivElement)
          const chartRect = this.eChart.chart.getDom().getBoundingClientRect()
          const isOverflowWindow = (chartRect.left + pointerLeft + element.offsetWidth) > chartRect.right

          return [
            isOverflowWindow ? (pointerLeft - element.offsetWidth) : pointerLeft,
            pointerTop - element.clientHeight - 10,
          ]
        },
        // Formatting
        formatter(params: CallbackDataParams) {
          const element = document.createElement("div")
          const [dateMs, durationMs, type] = params.value as OptionDataValue[]

          element.append(
            type == "c" ? durationMs.toString() : durationAxisPointerFormatter(isMs ? durationMs as number : durationMs as number / 1000 / 1000),
            document.createElement("br"),
            timeFormatWithoutSeconds.format(dateMs as number),
          )

          element.append(document.createElement("br"))
          element.append(`${params.seriesName}`)

          return element
        },
        valueFormatter(it) {
          return numberFormat.format(isMs ? it as number : nsToMs(it as number)) + " ms"
        },
        // Styling
        extraCssText: "user-select: text",
        borderColor: "#E5E7EB",
        padding: [6, 8],
        textStyle: {
          fontSize: 12
        },
      },
      xAxis: {
        type: "time",
        axisPointer: {
          snap: false,
          label: {
            formatter(data) {
              return timeFormat.format(data["value"] as number)
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
    return this.dataQuery.subscribe(
      (data, configuration) => {
        this.eChart.chart.hideLoading()
        this.eChart.chart.setOption(
          {
            legend: {type: "scroll"},
            toolbox: {top: 20},
          },
          {
            replaceMerge: ["legend"],
          },
        )
        this.eChart.updateChart(
          configuration.chartConfigurator.configureChart(data, configuration)
        )
      })
  }

  dispose(): void {
    this.eChart.dispose()
  }
}