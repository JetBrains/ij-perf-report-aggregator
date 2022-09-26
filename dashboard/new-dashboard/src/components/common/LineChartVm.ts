import { LineSeriesOption } from "echarts"
import * as ecStat from "echarts-stat"
import { LineChart, ScatterChart } from "echarts/charts"
import { DatasetComponent, DataZoomInsideComponent, DataZoomSliderComponent, GridComponent, LegendComponent, ToolboxComponent, TooltipComponent } from "echarts/components"
import { registerTransform, use } from "echarts/core"
import { CanvasRenderer } from "echarts/renderers"
import { CallbackDataParams, OptionDataValue } from "echarts/types/src/util/types"
import { ChartManagerHelper } from "shared/src/ChartManagerHelper"
import { DataQueryExecutor } from "shared/src/DataQueryExecutor"
import { ChartSymbolType, ValueUnit } from "shared/src/chart"
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

export class LineChartVm {
  private isTooltipShown = false

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
        // Formatting
        formatter: (params: CallbackDataParams[]) => {
          if (this.isTooltipShown) {
            const [data] = params
            const element = document.createElement("div")
            const [dateMs, durationMs] = data.value as OptionDataValue[]

            element.append(
              durationAxisPointerFormatter(durationMs as number),
              document.createElement("br"),
              timeFormatWithoutSeconds.format(dateMs as number),
            )

            return element
          }
          
          return ""
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

    this.eChart.chart.on("mouseover", params => {
      // Check if it's hovering a line
      if (params.componentSubType == "line") {
        this.isTooltipShown = true
      }

    })

    this.eChart.chart.on("mouseout", () => {
      this.isTooltipShown = false
    })
  }
  
  private setPointsSymbol(type: ChartSymbolType) {
    const options = this.eChart.chart.getOption()
    const series = (options["series"] as LineSeriesOption[])
      .map(line => Object.assign(line, {
        symbol: type,
      }))

    this.eChart.chart.setOption({ series }, {replaceMerge: ["series"]})
  }

  showPoints() {
    this.setPointsSymbol("circle")
  }

  hidePoints() {
    this.setPointsSymbol("none")
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