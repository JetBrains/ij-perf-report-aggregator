import { CallbackDataParams } from "echarts/types/src/util/types"
import { inject, shallowReactive } from "vue"
import { DataQueryExecutor } from "../DataQueryExecutor"
import { DataQuery } from "../dataQuery"
import { reportInfoProviderKey } from "../injectionKeys"

export interface ReportInfoProvider {
  createReportUrl(generatedTime: number, query: DataQuery): string

  readonly infoFields: Array<string>
}

export interface TooltipData {
  items: Array<TooltipDataItem>
  firstSeriesData: Array<number>
  reportInfoProvider: ReportInfoProvider | null
  query: DataQuery | null
}

interface TooltipDataItem {
  readonly name: string
  readonly value: number
  readonly color: string
}

export class ChartToolTipManager {
  public dataQueryExecutor!: DataQueryExecutor

  readonly reportInfoProvider = inject(reportInfoProviderKey, null)
  readonly reportTooltipData = shallowReactive<TooltipData>({items: [], firstSeriesData: [], reportInfoProvider: null, query: null})

  paused = false

  formatArrayValue(params: Array<CallbackDataParams>): null {
    if (this.paused) {
      console.log("paused")
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
    // same for all series
    data.firstSeriesData = params[0].value as Array<number>
    data.reportInfoProvider = this.reportInfoProvider
    data.query = query
    return null
  }
}