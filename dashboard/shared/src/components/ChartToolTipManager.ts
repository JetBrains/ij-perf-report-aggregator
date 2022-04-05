import { CallbackDataParams } from "echarts/types/src/util/types"
import { inject, reactive } from "vue"
import { DataQueryExecutor } from "../DataQueryExecutor"
import { DataQuery } from "../dataQuery"
import { reportInfoProviderKey } from "../injectionKeys"

export interface ReportInfoProvider {
  createReportUrl(generatedTime: number, query: DataQuery): string

  readonly infoFields: Array<string>
}

export interface TooltipData {
  linkText: string
  linkUrl: string | null
  items: Array<TooltipDataItem>
  firstSeriesData: Array<number>
}

interface TooltipDataItem {
  readonly name: string
  readonly value: number
  readonly color: string
}

const timeFormatWithoutSeconds = new Intl.DateTimeFormat(undefined, {
  year: "numeric",
  month: "short",
  day: "numeric",
  hour: "numeric",
  minute: "numeric",
})

export class ChartToolTipManager {
  public dataQueryExecutor!: DataQueryExecutor

  readonly reportInfoProvider = inject(reportInfoProviderKey, null)
  readonly reportTooltipData = reactive<TooltipData>({items: [], linkText: "", linkUrl: null, firstSeriesData: []})

  paused = false

  formatArrayValue(params: Array<CallbackDataParams>): null {
    if (this.paused) {
      return null
    }

    const query = this.dataQueryExecutor.lastQuery
    if (query == null) {
      return null
    }

    const data = this.reportTooltipData
    data.items = params.map(measure => {
      // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
      const measureValue = (measure.value as Array<number>)[measure.encode!["y"][0]]
      return {
        // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
        name: measure.seriesName!,
        // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
        value: measureValue,
        color: measure.color as string,
      }
    })
    const firstSeriesData = params[0].value as Array<number>
    // same for all series
    const generatedTime = firstSeriesData[0]
    data.linkText = timeFormatWithoutSeconds.format(generatedTime) + ` (${firstSeriesData[4]}.${firstSeriesData[5]}${firstSeriesData[6] === 0 ? "" : `.${firstSeriesData[6]}`})`
    data.firstSeriesData = firstSeriesData
    if (this.reportInfoProvider == null) {
      data.linkUrl = null
    }
    else {
      data.linkUrl = this.reportInfoProvider.createReportUrl(generatedTime, query)
    }
    return null
  }
}