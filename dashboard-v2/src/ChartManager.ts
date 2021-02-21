import { EChartsType, init as initChart, throttle } from "echarts/core"
import { EChartsOption } from "echarts"
import { toRaw } from "vue"
import { DataQueryExecutor } from "./DataQueryExecutor"
import { TplFormatterParam } from "echarts/types/src/util/format"
import { OrdinalRawValue } from "echarts/types/src/util/types"
import { initEcharts } from "./echarts"

initEcharts()

export class ChartManager {
  private readonly chart: EChartsType
  private readonly resizeObserver: ResizeObserver

  constructor(container: HTMLElement, private readonly dataQueryExecutor: DataQueryExecutor) {
    this.chart = initChart(container)

    this.resizeObserver = new ResizeObserver(throttle(() => {
      this.chart.resize()
    }, 300))
    this.resizeObserver.observe(this.chart.getDom())

    // https://github.com/apache/echarts/issues/8294
    const numberFormat = new Intl.NumberFormat(undefined, {
      maximumFractionDigits: 0,
    })
    const numberFormatWithUnit = new Intl.NumberFormat(undefined, {
      style: "unit",
      unit: "millisecond",
      maximumFractionDigits: 0,
    })
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
        // formatter: (params: CallbackDataParams | Array<CallbackDataParams>, _asyncTicket: string) => {
        //   const p = params as Array<CallbackDataParams>
        //   const result = p.map(v => {
        //     return `${v.marker} ${v.seriesName}: ${numberFormatWithUnit.format((v.value as any)[1] as number)}`
        //   })
        //   return `${p[0].name}<br />${result.join("<br />")}`
        // },
      },
      xAxis: {
        type: "time",
      },
      yAxis: {
        type: "value",
        axisLabel: {
          // just to be sure that axis label will be formatted using language-sensitive number formatting
          formatter: function (value: OrdinalRawValue, _index: number) {
            return numberFormat.format(value as number)
          },
        },
        axisPointer: {
          label: {
            formatter: function (data: TplFormatterParam) {
              return numberFormatWithUnit.format(data["value"])
            },
          },
        },
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
          dimensions: configuration.dimensions,
          // https://echarts.apache.org/en/option.html#dataset.sourceHeader, just optimization to avoid auto-detect
          sourceHeader: false,
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