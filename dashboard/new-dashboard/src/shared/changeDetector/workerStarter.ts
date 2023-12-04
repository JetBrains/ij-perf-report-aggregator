import { ChangePointClassification } from "./algorithm"
import MyWorker from "./worker?worker&inline"

export function detectChanges(seriesData: (string | number)[][]): Promise<Map<string, ChangePointClassification>> {
  const worker = new MyWorker()
  worker.postMessage(seriesData)
  return new Promise((resolve, reject) => {
    worker.addEventListener("error", (event) => {
      console.log(event)
      reject(event.error)
      worker.terminate()
    })
    worker.addEventListener("message", (event: MessageEvent<Map<string, ChangePointClassification>>) => {
      resolve(event.data)
      worker.terminate()
    })
  })
}
