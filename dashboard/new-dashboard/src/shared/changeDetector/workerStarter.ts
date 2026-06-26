import { BetterDirection, ChangePointClassification } from "./algorithm"
import MyWorker from "./worker?worker&inline"

export function detectChanges(seriesData: (string | number)[][], betterDirection: BetterDirection = "lower"): Promise<Map<string, ChangePointClassification>> {
  const worker = new MyWorker()
  worker.postMessage({ seriesData, betterDirection })
  return new Promise((resolve, reject) => {
    worker.addEventListener("error", (event) => {
      reject(event.error as Error)
      worker.terminate()
    })
    worker.addEventListener("message", (event: MessageEvent<Map<string, ChangePointClassification>>) => {
      resolve(event.data)
      worker.terminate()
    })
  })
}
