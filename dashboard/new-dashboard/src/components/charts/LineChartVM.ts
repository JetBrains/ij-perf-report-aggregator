import { CallbackDataParams, OptionDataValue } from "echarts/types/src/util/types"
import { Accident, convertAccidentsToMap, getAccident } from "../../util/meta"
import { DataQueryExecutor, DataQueryResult } from "../common/DataQueryExecutor"
import { timeFormat, ValueUnit } from "../common/chart"
import { DataQueryExecutorConfiguration } from "../common/dataQuery"
import { LineChartOptions } from "../common/echarts"
import { durationAxisPointerFormatter, nsToMs, numberFormat, timeFormatWithoutSeconds } from "../common/formatter"
import { ChartManager } from "./ChartManager"

function getWarningIcon() {
  const svg = document.createElementNS("http://www.w3.org/2000/svg", "svg")
  svg.setAttribute("xmlns", "http://www.w3.org/2000/svg")
  svg.setAttribute("fill", "none")
  svg.setAttribute("viewBox", "0 0 24 24")
  svg.setAttribute("stroke-width", "1.5")
  svg.setAttribute("stroke", "currentColor")
  svg.setAttribute("class", "w-6 h-6")
  const path = document.createElementNS("http://www.w3.org/2000/svg", "path")
  path.setAttribute("stroke-linecap", "round")
  path.setAttribute("stroke-linejoin", "round")
  path.setAttribute(
    "d",
    "M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 " +
      "0L2.697 16.126zM12 15.75h.007v.008H12v-.008z"
  )
  svg.append(path)

  const div = document.createElement("div")
  div.setAttribute("class", "w-4 h-4")
  div.append(svg)
  return div
}

export class LineChartVM {
  constructor(private readonly eChart: ChartManager, private readonly dataQuery: DataQueryExecutor, valueUnit: ValueUnit, accidents: Accident[] | undefined) {
    const accidentsMap = convertAccidentsToMap(accidents)
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
      // @ts-expect-error
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
        formatter(params: CallbackDataParams) {
          const element = document.createElement("div")
          const data = params.value as OptionDataValue[]
          const [dateMs, durationMs, _, type] = data

          element.append(
            type == "c" ? durationMs.toString() : durationAxisPointerFormatter(isMs ? (durationMs as number) : (durationMs as number) / 1000 / 1000),
            document.createElement("br"),
            timeFormatWithoutSeconds.format(dateMs as number)
          )

          element.append(document.createElement("br"))
          element.append(`${params.seriesName}`)
          const accident = getAccident(accidentsMap, data as string[])
          if (accident != null) {
            //<ExclamationTriangleIcon class="w-4 h-4 text-red-500" /> Known degradation:
            element.append(document.createElement("br"))
            const accidentHtml = document.createElement("span")
            accidentHtml.setAttribute("class", "flex gap-1.5 items-center")
            const div = getWarningIcon()
            accidentHtml.append(div)
            accidentHtml.append("Known " + accident.kind.toLowerCase() + ": " + accident.reason)
            element.append(accidentHtml)
          }

          return element
        },
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
      this.eChart.chart.setOption(
        {
          legend: { type: "scroll" },
          toolbox: { top: 20 },
        },
        {
          replaceMerge: ["legend"],
        }
      )
      this.eChart.updateChart(configuration.chartConfigurator.configureChart(data, configuration))
    })
  }

  dispose(): void {
    this.eChart.dispose()
  }
}
