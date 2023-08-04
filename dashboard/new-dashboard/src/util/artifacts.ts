import { Observable } from "rxjs"
import { shallowRef } from "vue"
import { DataQueryExecutor } from "../components/common/DataQueryExecutor"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration, SimpleQueryProducer } from "../components/common/dataQuery"
import { ServerConfigurator } from "../configurators/ServerConfigurator"
import { refToObservable } from "../configurators/rxjs"

export function getTeamcityBuildType(db: string, table: string, id: number, whenDone: (type: string | null) => void) {
  const serverUrlObservable = refToObservable(shallowRef(ServerConfigurator.DEFAULT_SERVER_URL))
  new DataQueryExecutor([
    new ServerConfigurator(db, table, serverUrlObservable),
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
      whenDone(changes)
    } else {
      whenDone(null)
    }
  })
}
