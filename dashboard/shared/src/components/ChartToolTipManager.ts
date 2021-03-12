import { CallbackDataParams } from "echarts/types/src/util/types"
import { inject, reactive, ref } from "vue"
import { DataQueryExecutor } from "../DataQueryExecutor"
import { timeFormat } from "../chart"
import { DataQuery } from "../dataQuery"
import { getValueFormatterByMeasureName } from "../formatter"
import { tooltipUrlProviderKey } from "../injectionKeys"
import { debounceSync } from "../util/debounce"

export type ChartTooltipLinkProvider = (generatedTime: number, query: DataQuery) => string

interface TooltipData {
  linkText: string
  linkUrl: string | null
  items: Array<TooltipDataItem>
}

interface TooltipDataItem {
  name: string
  value: string
  color: string
}

export class ChartToolTipManager {
  public dataQueryExecutor!: DataQueryExecutor

  private readonly tooltipUrlProvider = inject(tooltipUrlProviderKey, null)
  readonly infoIsVisible = ref(false)

  readonly reportTooltipData = reactive<TooltipData>({items: [], linkText: "", linkUrl: null})

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
    const generatedTime = (params[0].value as Array<number>)[0]
    reportTooltipData.linkText = timeFormat.format(generatedTime)
    if (this.tooltipUrlProvider == null) {
      reportTooltipData.linkUrl = null
    }
    else {
      reportTooltipData.linkUrl = this.tooltipUrlProvider(generatedTime, query)
    }
    this.infoIsVisible.value = true
    this.scheduleTooltipHide()
    return null
  }
}