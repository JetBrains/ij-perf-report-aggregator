import { Observable } from "rxjs"
import { shallowRef } from "vue"
import { DataQueryExecutor } from "../components/common/DataQueryExecutor"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration, SimpleQueryProducer } from "../components/common/dataQuery"
import { ServerWithCompressConfigurator } from "../configurators/ServerWithCompressConfigurator"
import { refToObservable } from "../configurators/rxjs"

export function base64ToHex(base64: string): string {
  const decodedArray = new Uint8Array(Array.from(atob(base64), (c) => c.codePointAt(0) ?? 0))
  let hex = ""
  for (const byte of decodedArray) {
    hex += byte.toString(16).padStart(2, "0")
  }
  return hex
}

export function calculateChanges(db: string, id: number): Promise<string[] | null> {
  return new Promise((resolve, _) => {
    const serverUrlObservable = refToObservable(shallowRef(ServerWithCompressConfigurator.DEFAULT_SERVER_URL))
    const separator = ".."
    new DataQueryExecutor([
      new ServerWithCompressConfigurator(db, "installer", serverUrlObservable),
      new (class implements DataQueryConfigurator {
        configureQuery(query: DataQuery, configuration: DataQueryExecutorConfiguration): boolean {
          configuration.queryProducers.push(new SimpleQueryProducer())
          query.addField({ n: "changes", sql: `arrayStringConcat(changes,'${separator}')` })
          query.addFilter({ f: "id", v: id })
          query.order = "changes"
          return true
        }

        createObservable(): Observable<unknown> | null {
          return null
        }
      })(),
    ]).subscribe((data, _configuration, isLoading) => {
      if (isLoading || data == null) {
        return
      }
      const changes = data.flat(3)[0]
      if (typeof changes === "string") {
        //commit has to be decoded as base64 and converted to hex
        const changesDecoded = changes.split(separator).map((it) => base64ToHex(it))

        // Split into chunks of 150 changes each
        // otherwise the URL will be too long and fail
        const result: string[] = []
        for (let i = 0; i < changesDecoded.length; i += 150) {
          const chunk = changesDecoded.slice(i, i + 150)
          result.push(chunk.join("%2C"))
        }

        resolve(result)
      } else {
        resolve(null)
      }
    })
  })
}
