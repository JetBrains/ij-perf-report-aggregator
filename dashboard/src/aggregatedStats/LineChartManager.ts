// Copyright 2000-2020 JetBrains s.r.o. Use of this source code is governed by the Apache 2.0 license that can be found in the LICENSE file.
import * as am4charts from "@amcharts/amcharts4/charts"
import * as am4core from "@amcharts/amcharts4/core"
import { ChartSettings } from "@/aggregatedStats/ChartSettings"
import { addExportMenu, StatChartManager } from "@/charts/ChartManager"
import { Metrics } from "@/aggregatedStats/model"
import { ChartConfigurator } from "@/aggregatedStats/ChartConfigurator"
import * as am4plugins_annotation from "@amcharts/amcharts4/plugins/annotation"

export class LineChartManager implements StatChartManager {
  readonly chart: am4charts.XYChart

  constructor(container: HTMLElement,
              private chartSettings: ChartSettings,
              private readonly isInstantEvents: boolean,
              private readonly configurator: ChartConfigurator) {
    this.chart = am4core.create(container, am4charts.XYChart)

    const chart = this.chart

    chart.legend = new am4charts.Legend()
    chart.legend.itemContainers.template.events.on("over", event => {
      const dataItem = event.target.dataItem
      if (dataItem == null) {
        return
      }

      for (const series of this.chart.series) {
        if (dataItem.dataContext === series) {
          continue
        }

        for (const segment of (series as am4charts.LineSeries).segments) {
          segment.setState("inactive")
        }
      }
    })
    chart.legend.itemContainers.template.events.on("out", event => {
      const dataItem: any = event.target.dataItem
      if (dataItem == null) {
        return
      }

      for (const series of this.chart.series) {
        if (dataItem.dataContext === series) {
          continue
        }

        for (const segment of (series as am4charts.LineSeries).segments) {
          segment.setState("default")
        }
      }
    })

    chart.colors.step = 4
    addExportMenu(chart)

    const cursor = new am4charts.XYCursor()
    cursor.behavior = "zoomX"
    chart.cursor = cursor

    // const dateAxis = chart.xAxes.push(new am4charts.DateAxis())
    configurator.configureXAxis(chart)
    // xAxis.groupData = true
    // DurationAxis doesn't work due to some unclear bug
    if (chartSettings.selectedProduct == "sharedIndexes") {
      const durationAxis = chart.yAxes.push(new am4charts.DurationAxis())
      durationAxis.baseUnit = "millisecond"
      durationAxis.durationFormatter.baseUnit = "millisecond"
      durationAxis.durationFormatter.durationFormat = "hh 'h' mm 'm' ss 's'"
    } else {
      const valueAxis = chart.yAxes.push(new am4charts.ValueAxis())

      // do not use logarithmic scale for line chart of duration events - better looking and more clear charts, if height will be a problem, then chart height can be increased
      valueAxis.logarithmic = this.isInstantEvents
      valueAxis.durationFormatter.baseUnit = "millisecond"
      valueAxis.durationFormatter.durationFormat = "S"
    }

    if (this.chartSettings.showScrollbarXPreview) {
      this.configureScrollbarXWithPreview()
    }

    chart.plugins.push(new am4plugins_annotation.Annotation())

    // prevent Vue reactivity
    Object.seal(this)
  }

  private configureScrollbarXWithPreview(): am4charts.XYChartScrollbar {
    const scrollbarX = new am4charts.XYChartScrollbar()
    const chart = this.chart
    chart.scrollbarX = scrollbarX
    scrollbarX.parent = chart.bottomAxesContainer
    return scrollbarX
  }

  scrollbarXPreviewOptionChanged(chartSettings: ChartSettings) {
    this.chartSettings = chartSettings

    // no need to dispose old scrollbar explicit - will be disposed automatically on set
    const chart = this.chart
    if ((chart.scrollbarX != null) === chartSettings.showScrollbarXPreview) {
      return
    }

    if (chartSettings.showScrollbarXPreview) {
      const scrollbarX = this.configureScrollbarXWithPreview()
      chart.series.each(it => {
        scrollbarX.series.push(it)
      })
    }
    else {
      chart.scrollbarX = null as any
    }
  }

  render(data: Array<Metrics>, seriesDescriptors: Array<SeriesDescriptor>): void {
    console.log("render lines")
    const chart = this.chart

    const oldSeries = new Map<string, am4charts.XYSeries>()
    for (const series of chart.series) {
      oldSeries.set(series.name, series)
    }

    const scrollbarX = chart.scrollbarX as am4charts.XYChartScrollbar
    for (const descriptor of seriesDescriptors) {
      const seriesName = descriptor.name
      let series = oldSeries.get(seriesName)
      if (series == null) {
        const series = new am4charts.LineSeries()
        this.configureLineSeries(descriptor, series, seriesName)
        chart.series.push(series)
        if (this.chartSettings.showScrollbarXPreview) {
          scrollbarX.series.push(series)
        }
      }
      else {
        oldSeries.delete(seriesName)
      }
    }

    if (oldSeries.size > 0) {
      console.info("dispose series", Array.from(oldSeries.keys()))
      for (const series of oldSeries.values()) {
        chart.series.removeIndex(chart.series.indexOf(series))
        series.dispose()
      }
    }

    // else {
    //   // https://github.com/amcharts/amcharts4/issues/1908
    //   for (const series of chart.series) {
    //     series.bulletsContainer.disposeChildren()
    //   }
    // }

    chart.data = data
  }

  private configureLineSeries(metric: SeriesDescriptor, series: am4charts.LineSeries, seriesName: string) {
    series.name = seriesName

    this.configurator.configureSeries(series)
    // duration
    series.dataFields.valueY = metric.dataField

    if (metric.hiddenByDefault) {
      series.hidden = true
    }

    series.strokeWidth = 2
    series.smoothing = "monotoneX"
    const segmentState = series.segments.template.states.create("inactive")
    segmentState.properties.strokeOpacity = 0.4
  }

  dispose(): void {
    this.chart.dispose()
  }
}

export interface SeriesDescriptor {
  name: string
  dataField: string
  hiddenByDefault: boolean
}