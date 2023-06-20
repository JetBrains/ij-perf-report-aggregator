import { CallbackDataParams } from "echarts/types/src/util/types"
import { inject } from "vue"
import { reportInfoProviderKey } from "../../shared/injectionKeys"
import { DataQueryExecutor } from "../common/DataQueryExecutor"
import { ValueUnit } from "../common/chart"

export interface ReportInfoProvider {
  readonly infoFields: string[]
}

export interface TooltipData {
  items: TooltipDataItem[]
  firstSeriesData: number[]
  reportInfoProvider: ReportInfoProvider | null
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

  private consumer: ((data: TooltipData | null, target: Event | null) => void) | null = null

  setConsumer(consumer: ((data: TooltipData | null, target: Event | null) => void) | null) {
    this.consumer = consumer
  }

  showTooltip(params: CallbackDataParams[] | null, target: Event | null) {
    const consumer = this.consumer
    if (consumer == null) {
      return
    }

    if (params == null || params.length === 0) {
      consumer(null, null)
      return
    }

    // same for all series
    const values = params[0].value as number[]

    consumer(
      {
        items: params.map((measure) => {
          // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
          const measureValue = (measure.value as number[])[measure.encode!["y"][0]] / (this.valueUnit == "ns" ? 1_000_000 : 1)
          return {
            // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
            name: measure.seriesName!,
            // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
            value: measureValue,
            color: measure.color as string,
          }
        }),
        firstSeriesData: values,
        reportInfoProvider: this.reportInfoProvider,
      },
      target
    )
  }
}
