import { map, Observable, shareReplay } from "rxjs"
import { DataQuery, DataQueryExecutorConfiguration, serializeQuery, ServerConfigurator } from "../../src/components/common/dataQuery"

export class TestServerConfigurator implements ServerConfigurator {
  static readonly DEFAULT_SERVER_URL = "https://ij-perf.labs.jb.gg"

  private readonly observable: Observable<null>
  private _serverUrl: string = TestServerConfigurator.DEFAULT_SERVER_URL

  constructor(
    readonly db: string,
    readonly table: string
  ) {
    this.observable = new Observable<string>((subscriber) => {
      subscriber.next(TestServerConfigurator.DEFAULT_SERVER_URL)
    }).pipe(
      map((url) => {
        this._serverUrl = url
        return null
      }),
      shareReplay(1)
    )
  }

  get serverUrl(): string {
    return this._serverUrl
  }

  compressString(params: string): string {
    return params
  }

  computeQueryUrl(query: DataQuery): string {
    return `${this._serverUrl}/api/q/${this.compressString(serializeQuery(query))}`
  }

  computeSerializedQueryUrl(url: string): string {
    return `${this._serverUrl}/api/q/${this.compressString(url)}`
  }

  createObservable(): Observable<unknown> {
    return this.observable
  }

  configureQuery(query: DataQuery, _configuration: DataQueryExecutorConfiguration): boolean {
    query.db = this.db
    query.table = this.table
    return true
  }
}
