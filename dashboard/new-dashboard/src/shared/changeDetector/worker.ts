import { BetterDirection, ChangeCriteria, detectChanges } from "./algorithm"

onmessage = (e: MessageEvent<{ seriesData: (string | number)[][]; betterDirection: BetterDirection; criteria: ChangeCriteria }>) => {
  postMessage(detectChanges(e.data.seriesData, e.data.betterDirection, e.data.criteria))
}
