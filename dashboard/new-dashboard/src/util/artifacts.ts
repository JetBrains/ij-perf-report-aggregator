import { Observable } from "rxjs"
import { shallowRef } from "vue"
import { DataQueryExecutor } from "../components/common/DataQueryExecutor"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration, SimpleQueryProducer } from "../components/common/dataQuery"
import { ServerWithCompressConfigurator } from "../configurators/ServerWithCompressConfigurator"
import { refToObservable } from "../configurators/rxjs"

export async function getTeamcityBuildCounter(buildId: number): Promise<string | null> {
  try {
    const response = await fetch(`${ServerWithCompressConfigurator.DEFAULT_SERVER_URL}/api/meta/teamcity/buildCounter?buildId=${encodeURIComponent(buildId.toString())}`)
    if (!response.ok) {
      console.log(`Failed to fetch TeamCity build counter: ${response.status} - ${await response.text()}`)
      return null
    }
    return await response.text()
  } catch (error) {
    console.log("Error fetching TeamCity build counter:", error)
    return null
  }
}

export function getTeamcityBuildType(db: string, table: string, id: number): Promise<string | null> {
  return new Promise((resolve, _) => {
    const serverUrlObservable = refToObservable(shallowRef(ServerWithCompressConfigurator.DEFAULT_SERVER_URL))

    new DataQueryExecutor([
      new ServerWithCompressConfigurator(db, table, serverUrlObservable),
      new (class implements DataQueryConfigurator {
        configureQuery(query: DataQuery, configuration: DataQueryExecutorConfiguration): boolean {
          configuration.queryProducers.push(new SimpleQueryProducer())
          query.addField({ n: "tc_build_type" })
          query.addFilter({ f: "tc_build_id", v: id })
          query.order = "tc_build_id"
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
        resolve(changes)
      } else {
        resolve(null)
      }
    })
  })
}
