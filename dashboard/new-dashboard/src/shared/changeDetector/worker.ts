import { detectChanges } from "./algorithm"

onmessage = (e: MessageEvent<number[][]>) => {
  postMessage(detectChanges(e.data))
}
