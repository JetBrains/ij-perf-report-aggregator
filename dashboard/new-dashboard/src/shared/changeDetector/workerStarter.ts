import { ChangePointClassification } from "./algorithm"

export function detectChanges(seriesData: (string | number)[][]): Promise<Map<string, ChangePointClassification>> {
  const worker = new Worker(new URL("worker.ts", import.meta.url), { type: "module" })
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
