import { EChartsType, init as initChart, throttle } from "echarts/core"
import { EChartsOption } from "echarts"
import { Ref, watch } from "vue"
import { DataQueryExecutor } from "./DataQueryExecutor"
import { TplFormatterParam } from "echarts/types/src/util/format"
import { CallbackDataParams, OrdinalRawValue } from "echarts/types/src/util/types"
import { initEcharts } from "./echarts"
import { debounceSync } from "./util/debounce"

initEcharts()

const dataZoomConfig = [
  // https://echarts.apache.org/en/option.html#dataZoom-inside
  // type inside means that mouse maybe used to zoom.
  {type: "inside"},
  {},
]

// https://github.com/apache/echarts/issues/8294
export const numberFormat = new Intl.NumberFormat(undefined, {
  maximumFractionDigits: 0,
})

const numberFormatWithUnit = createNumberFormatWithUnit()

export const timeFormat = new Intl.DateTimeFormat(undefined, {
  year: "numeric",
  month: "short",
  day: "numeric",
  hour: "numeric",
  minute: "numeric",
  second: "numeric",
})

export class ChartManager {
  private readonly chart: EChartsType
  private readonly resizeObserver: ResizeObserver

  constructor(container: HTMLElement,
              private readonly dataQueryExecutor: DataQueryExecutor,
              dataZoom: Ref<boolean>,
              tooltipFormatter: (params: CallbackDataParams | Array<CallbackDataParams>, _ticket: string) => HTMLElement | string | null) {
    this.chart = initChart(container)

    this.resizeObserver = new ResizeObserver(throttle(() => {
      this.chart.resize()
    }, 300))
    this.resizeObserver.observe(this.chart.getDom())

    this.chart.setOption<EChartsOption>({
      legend: {},
      animation: false,
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
        formatter: tooltipFormatter as never,
        triggerOn: "click",
      },
      xAxis: {
        type: "time",
        axisPointer: {
          label: {
            formatter(data: TplFormatterParam) {
              return timeFormat.format(data["value"])
            },
          },
        },
      },
      yAxis: {
        type: "value",
        axisLabel: {
          // just to be sure that axis label will be formatted using language-sensitive number formatting
          formatter(value: OrdinalRawValue, _index: number) {
            return numberFormat.format(value as number)
          },
        },
        axisPointer: {
          label: {
            formatter(data: TplFormatterParam) {
              return numberFormatWithUnit.format(data["value"])
            },
          },
        },
      },
      dataZoom: dataZoom.value ? dataZoomConfig : undefined,
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
          // just optimization to avoid auto-detect (https://echarts.apache.org/en/option.html#dataset.sourceHeader)
          sourceHeader: false,
          source: data,
        },
        series: configuration.series,
      }, {replaceMerge: ["dataset", "series"]})
    })

    watch(dataZoom, debounceSync(() => {
      this.chart.setOption({
        dataZoom: dataZoom.value ? dataZoomConfig : undefined,
      })
    }))
  }

  dispose(): void {
    this.resizeObserver.unobserve(this.chart.getDom())
    this.chart.dispose()
    this.dataQueryExecutor.setListener(null)
  }
}

function createNumberFormatWithUnit(): Intl.NumberFormat {
  try {
    return new Intl.NumberFormat(undefined, {
      style: "unit",
      unit: "millisecond",
      maximumFractionDigits: 0,
    })
  }
  catch (e) {
    // Safari doesn't support `unit` (https://caniuse.com/mdn-javascript_builtins_intl_numberformat_numberformat_unit)
    return new Intl.NumberFormat(undefined, {
      style: "decimal",
      maximumFractionDigits: 0,
    })
  }
}