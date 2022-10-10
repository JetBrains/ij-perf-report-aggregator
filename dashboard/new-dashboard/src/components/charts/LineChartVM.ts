import * as ecStat from "echarts-stat"
import { LineChart, ScatterChart } from "echarts/charts"
import { DatasetComponent, DataZoomInsideComponent, DataZoomSliderComponent, GridComponent, LegendComponent, ToolboxComponent, TooltipComponent } from "echarts/components"
import { registerTransform, use } from "echarts/core"
import { CanvasRenderer } from "echarts/renderers"
import { CallbackDataParams, OptionDataValue } from "echarts/types/src/util/types"
import { ChartManagerHelper } from "shared/src/ChartManagerHelper"
import { DataQueryExecutor } from "shared/src/DataQueryExecutor"
import { ValueUnit } from "shared/src/chart"
import { LineChartOptions } from "shared/src/echarts"
import { durationAxisPointerFormatter, nsToMs, numberFormat, timeFormatWithoutSeconds } from "shared/src/formatter"

use([
  DatasetComponent,
  ToolboxComponent,
  TooltipComponent,
  GridComponent,
  LineChart,
  ScatterChart,
  LegendComponent,
  CanvasRenderer,
  DataZoomInsideComponent,
  DataZoomSliderComponent,
])

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
// eslint-disable-next-line @typescript-eslint/no-unsafe-argument,@typescript-eslint/no-unsafe-member-access
registerTransform(ecStat["transform"].regression)

export class LineChartVM {
  constructor(
    private readonly eChart: ChartManagerHelper,
    private readonly dataQuery: DataQueryExecutor,
    valueUnit: ValueUnit,
  ) {
    const isMs = valueUnit == "ms"

    this.eChart.chart.setOption<LineChartOptions>({
      legend: {
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
        right: 0,
        bottom: 16,
        containLabel: true,
      },
      // @ts-ignore
      tooltip: {
        show: true,
        trigger: "axis",
        enterable: true,
        axisPointer: {
          type: "cross",
          snap: true
        },
        renderMode: "html",
        position: (pointerCoords, _, tooltipElement) => {
          const [pointerLeft, pointerTop] = pointerCoords
          const element = (tooltipElement as HTMLDivElement)
          const charLeft = this.eChart.chart.getDom().getBoundingClientRect().left
          const isOverflowWindow = (charLeft + pointerLeft + element.offsetWidth) > window.innerWidth

          return [
            isOverflowWindow ? (pointerLeft - element.offsetWidth) : pointerLeft,
            pointerTop - element.clientHeight - 10,
          ]
        },
        // Formatting
        formatter(params: CallbackDataParams[]) {
          const [data] = params
          const element = document.createElement("div")
          const [dateMs, durationMs] = data.value as OptionDataValue[]

          element.append(
            durationAxisPointerFormatter(durationMs as number),
            document.createElement("br"),
            timeFormatWithoutSeconds.format(dateMs as number),
          )

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
        this.eChart.replaceDataSetAndSeries(
          configuration.chartConfigurator.configureChart(data, configuration)
        )
      })
  }

  dispose(): void {
    this.eChart.dispose()
  }
}