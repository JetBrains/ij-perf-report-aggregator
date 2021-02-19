import { TplFormatterParam } from "echarts/types/src/util/format"
import { OrdinalRawValue } from "echarts/types/src/util/types"
import { watch , Ref } from "vue"
import { ChartManagerHelper } from "./ChartManagerHelper"
import { DataQueryExecutor } from "./DataQueryExecutor"
import { adaptToolTipFormatter, ChartOptions, numberFormat, timeFormat, ToolTipFormatter, useLineAndBarCharts } from "./chart"
import { DataQuery } from "./dataQuery"
import { debounceSync } from "./util/debounce"

export type ChartTooltipLinkProvider = (name: string, query: DataQuery) => string

const dataZoomConfig = [
  // https://echarts.apache.org/en/option.html#dataZoom-inside
  // type inside means that mouse maybe used to zoom.
  {type: "inside"},
  {},
]

useLineAndBarCharts()

const numberFormatWithUnit = createNumberFormatWithUnit()

export class LineChartManager {
  private readonly chart: ChartManagerHelper

  constructor(container: HTMLElement,
              private readonly dataQueryExecutor: DataQueryExecutor,
              dataZoom: Ref<boolean>,
              tooltipFormatter: ToolTipFormatter) {
    this.chart = new ChartManagerHelper(container)

    this.chart.chart.setOption<ChartOptions>({
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
        },
        formatter: adaptToolTipFormatter(tooltipFormatter),
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

    this.chart.enableZoomTool()

    dataQueryExecutor.setListener((data, configuration) => {
      this.chart.replaceDataSetAndSeries(configuration.chartConfigurator.configureChart(data, configuration))
      // console.log(JSON.stringify(this.chart.getOption(), null, 2))
    })

    watch(dataZoom, debounceSync(() => {
      this.chart.chart.setOption({
        dataZoom: dataZoom.value ? dataZoomConfig : undefined,
      })
    }))
  }

  dispose(): void {
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