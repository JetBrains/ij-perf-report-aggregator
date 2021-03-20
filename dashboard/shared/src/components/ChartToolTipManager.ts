import { CallbackDataParams } from "echarts/types/src/util/types"
import { inject, reactive, ref } from "vue"
import { DataQueryExecutor } from "../DataQueryExecutor"
import { timeFormat } from "../chart"
import { DataQuery } from "../dataQuery"
import { getValueFormatterByMeasureName } from "../formatter"
import { reportInfoProviderKey } from "../injectionKeys"
import { debounceSync } from "../util/debounce"

export interface ReportInfoProvider {
  createReportUrl(generatedTime: number, query: DataQuery): string

  readonly infoFields: Array<string>
}

interface TooltipData {
  linkText: string
  linkUrl: string | null
  items: Array<TooltipDataItem>
  firstSeriesData: Array<number>
}

interface TooltipDataItem {
  readonly name: string
  readonly value: string
  readonly color: string
}

export class ChartToolTipManager {
  public dataQueryExecutor!: DataQueryExecutor

  readonly reportInfoProvider = inject(reportInfoProviderKey, null)
  readonly infoIsVisible = ref(false)

  readonly reportTooltipData = reactive<TooltipData>({items: [], linkText: "", linkUrl: null, firstSeriesData: []})

  readonly scheduleTooltipHide = debounceSync(() => {
    this.infoIsVisible.value = false
  }, 2_000)

  formatArrayValue(params: Array<CallbackDataParams>): null {
    const query = this.dataQueryExecutor.lastQuery
    if (query == null) {
      return null
    }

    const reportTooltipData = this.reportTooltipData
    reportTooltipData.items = params.map(function (measure) {
      // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
      const measureValue = (measure.value as Array<number>)[measure.encode!["y"][0]]
      return {
        // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
        name: measure.seriesName!,
        // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
        value: getValueFormatterByMeasureName(measure.seriesId!)(measureValue),
        color: measure.color as string,
      }
    })
    const firstSeriesData = params[0].value as Array<number>
    // same for all series
    const generatedTime = firstSeriesData[0]
    reportTooltipData.linkText = timeFormat.format(generatedTime)
    reportTooltipData.firstSeriesData = firstSeriesData
    if (this.reportInfoProvider == null) {
      reportTooltipData.linkUrl = null
    }
    else {
      reportTooltipData.linkUrl = this.reportInfoProvider.createReportUrl(generatedTime, query)
    }
    this.infoIsVisible.value = true
    this.scheduleTooltipHide()
    return null
  }
}