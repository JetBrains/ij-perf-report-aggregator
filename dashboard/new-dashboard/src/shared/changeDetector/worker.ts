import { BetterDirection, detectChanges } from "./algorithm"

onmessage = (e: MessageEvent<{ seriesData: (string | number)[][]; betterDirection: BetterDirection }>) => {
  postMessage(detectChanges(e.data.seriesData, e.data.betterDirection))
}
