import { EChartsType, init as initChart, throttle } from "echarts/core"
import { EChartsOption } from "echarts"
import { toRaw } from "vue"
import { DataQueryExecutor } from "./DataQueryExecutor"

export class ChartManager {
  private readonly chart: EChartsType
  private readonly resizeObserver: ResizeObserver

  constructor(container: HTMLElement, private readonly dataQueryExecutor: DataQueryExecutor) {
    this.chart = initChart(container)

    this.resizeObserver = new ResizeObserver(throttle(() => {
      this.chart.resize()
    }, 300))
    this.resizeObserver.observe(this.chart.getDom())

    this.chart.setOption<EChartsOption>({
      legend: {},
      toolbox: {
        feature: {
          dataZoom: {
            yAxisIndex: false,
          },
          saveAsImage: {},
        },
      },
      tooltip: {
        trigger: "axis",
        axisPointer: {
          type: "cross",
          snap: true,
        },
      },
      xAxis: {
        type: "time",
      },
      yAxis: {
        type: "value",
      },
      dataZoom: [
        {},
        {},
      ],
    })

    // https://github.com/apache/echarts/issues/10274
    this.chart.dispatchAction({
      type: "takeGlobalCursor",
      key: "dataZoomSelect",
      dataZoomSelectActive: true,
    })

    dataQueryExecutor.setListener((data, configuration) => {
      // for (const series of configuration.series) {
      //   // https://echarts.apache.org/en/option.html#series-lines.emphasis.focus
      //   // series.emphasis = {
      //   //   focus: "self",
      //   // }
      // }
      this.chart.setOption<EChartsOption>({
        dataset: {
          source: toRaw(data),
        },
        series: configuration.series,
      }, { replaceMerge: ["dataset", "series"] })
    })
  }

  dispose(): void {
    this.resizeObserver.unobserve(this.chart.getDom())
    this.chart.dispose()
    this.dataQueryExecutor.setListener(null)
  }
}