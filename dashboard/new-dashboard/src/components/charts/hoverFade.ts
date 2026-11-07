import { ECElementEvent } from "echarts/core"
import { watch } from "vue"
import { useSettingsStore } from "../settings/settingsStore"
import { ChartManager } from "./ChartManager"
import { setSeriesHoverCursor } from "./intervalSelection"
import { parseSeriesId } from "./seriesId"
import { STD_DEV_BAND_OPACITY } from "./stdDevInterval"

/**
 * Series option key that names the group a series fades with; series without it fade with their own derived series.
 */
export const HOVER_GROUP_KEY = "hoverGroup"

const FADED_SERIES_OPACITY = 0.2
const VISIBLE_SERIES_OPACITY = 1

interface SeriesWithGroup {
  id?: string
  [HOVER_GROUP_KEY]?: string
}

export class HoverFadeController {
  private hoveredGroupId: string | null = null
  private appliedGroupId: string | null = null
  private resetFrameId?: number
  private readonly settings = useSettingsStore()
  private readonly unwatchSetting: () => void

  constructor(private readonly chartManager: ChartManager) {
    chartManager.chart.on("mouseover", this.onMouseOver)
    chartManager.chart.on("mouseout", this.onMouseOut)
    chartManager.chart.on("globalout", this.onGlobalOut)

    this.unwatchSetting = watch(
      () => this.settings.fadeOnHover,
      () => {
        this.clearPendingReset()
        this.apply(this.activeGroupId, true)
      }
    )
  }

  reapply() {
    this.apply(this.activeGroupId, true)
  }

  dispose() {
    this.clearPendingReset()
    this.unwatchSetting()
    const chart = this.chartManager.chart
    chart.off("mouseover", this.onMouseOver)
    chart.off("mouseout", this.onMouseOut)
    chart.off("globalout", this.onGlobalOut)
    setSeriesHoverCursor(this.chartManager, false)
  }

  private get activeGroupId(): string | null {
    return this.settings.fadeOnHover ? this.hoveredGroupId : null
  }

  private readonly onMouseOver = (params: ECElementEvent) => {
    const series = (this.chartManager.chart.getOption() as { series?: SeriesWithGroup[] }).series
    this.setHoveredGroup(HoverFadeController.groupId(series?.find((it) => it.id === params.seriesId) ?? { id: params.seriesId }))
  }

  private readonly onMouseOut = () => {
    this.scheduleReset()
  }

  private readonly onGlobalOut = () => {
    this.setHoveredGroup(null)
  }

  private setHoveredGroup(groupId: string | null) {
    this.clearPendingReset()
    this.hoveredGroupId = groupId
    setSeriesHoverCursor(this.chartManager, groupId != null)
    if (this.settings.fadeOnHover) {
      this.apply(groupId)
    }
  }

  private apply(groupId: string | null, force = false) {
    if (this.appliedGroupId === groupId && (!force || groupId == null)) {
      return
    }

    const option = this.chartManager.chart.getOption() as { series?: SeriesWithGroup[] }
    if (option.series == null) {
      this.appliedGroupId = groupId
      return
    }

    const series = option.series.map((s) => {
      const opacity = groupId == null || HoverFadeController.groupId(s) === groupId ? VISIBLE_SERIES_OPACITY : FADED_SERIES_OPACITY
      switch (s.id == undefined ? null : parseSeriesId(s.id).role) {
        case "stdDevBandFill": {
          // A band has no stroke of its own, so fading it means fading its fill.
          return { id: s.id, areaStyle: { opacity: STD_DEV_BAND_OPACITY * opacity } }
        }
        case "stdDevBandLowerEdge": {
          // The lower edge is already invisible and has nothing to fade - and an areaStyle it was never given
          // would fill it from the axis up, so it is left exactly as it is. It still has to be patched, or the
          // series below it shift into its slot.
          return { id: s.id }
        }
        default: {
          return { id: s.id, lineStyle: { opacity }, itemStyle: { opacity } }
        }
      }
    })
    this.chartManager.chart.setOption({ series })
    this.appliedGroupId = groupId
  }

  private static groupId(series: SeriesWithGroup): string | null {
    const group = series[HOVER_GROUP_KEY]
    if (group != undefined) {
      return HOVER_GROUP_KEY + ":" + group
    }
    if (series.id == null || series.id.length === 0) {
      return null
    }
    return parseSeriesId(series.id).ownerSeriesId
  }

  private clearPendingReset() {
    if (this.resetFrameId != undefined) {
      cancelAnimationFrame(this.resetFrameId)
      this.resetFrameId = undefined
    }
  }

  private scheduleReset() {
    this.clearPendingReset()
    this.resetFrameId = requestAnimationFrame(() => {
      this.resetFrameId = undefined
      this.setHoveredGroup(null)
    })
  }
}
