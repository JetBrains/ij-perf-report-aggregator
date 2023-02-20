import { CallbackDataParams } from "echarts/types/src/util/types"
import { inject } from "vue"
import { DataQueryExecutor } from "../DataQueryExecutor"
import { ValueUnit } from "../chart"
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
  private readonly valueUnit: ValueUnit

  constructor(valueUnit: ValueUnit) {
    this.valueUnit = valueUnit
  }

  public dataQueryExecutor!: DataQueryExecutor

  readonly reportInfoProvider = inject(reportInfoProviderKey, null)

  private consumer: ((data: TooltipData  | null, target: Event | null) => void) | null = null

  setConsumer(consumer: ((data: TooltipData | null, target: Event | null) => void) | null) {
    this.consumer = consumer
  }

  showTooltip(params: Array<CallbackDataParams> | null, target: Event | null) {
    const query = this.dataQueryExecutor.lastQuery
    if (query == null) {
      return
    }

    const consumer = this.consumer
    if (consumer == null) {
      return
    }

    if (params == null || params.length === 0) {
      consumer(null, null)
      return
    }

    // same for all series
    const values = params[0].value as Array<number>

    consumer({
      items: params.map(measure => {
        // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
        const measureValue = (measure.value as Array<number>)[measure.encode!["y"][0]] / (this.valueUnit == "ns" ? 1_000_000 : 1)
        return {
          // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
          name: measure.seriesName!,
          // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
          value: measureValue,
          color: measure.color as string,
        }
      }),
      firstSeriesData: query.db === "perfint" ? [...values.slice(0, 2), ...values.slice(3)] : values,
      reportInfoProvider: this.reportInfoProvider,
      query,
    }, target)
  }
}