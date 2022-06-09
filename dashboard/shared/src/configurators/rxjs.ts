import { deepEqual } from "fast-equals"
import { ToastSeverity } from "primevue/api"
import ToastEventBus from "primevue/toasteventbus"
import { catchError, delay, distinctUntilChanged, EMPTY, mergeMap, Observable, of, retry, takeUntil, timer } from "rxjs"
import { fromFetch } from "rxjs/fetch"
import { Ref, watch } from "vue"

export interface ServerResponse {
  t: Array<number>
  machine: Array<string>
  generated_time: Array<number>
  tc_build_id: Array<number>
  triggeredBy: Array<string>
  branch: Array<string>

  tc_installer_build_id?: Array<string>
  build_c1?: Array<number>
  build_c2?: Array<number>
  build_c3?: Array<number>
  raw_report?: Array<string>
  project?: Array<string>
  "measures.name" : Array<string>
  "measures.values": Array<number>
  "measures.type": Array<"c"|"d">

  [index: string]: Array<number|string>|undefined
}
export function refToObservable<T>(ref: Ref<T>, deep: boolean = false): Observable<T> {
  return new Observable<T>(context => {
    watch(ref, value => {
      return context.next(value)
    }, {deep})
    context.next(ref.value)
  }).pipe(
    deep ? distinctUntilChanged(deepEqual) : distinctUntilChanged(), 
  )
}

const serverNotAvailableMessage = {
  severity: ToastSeverity.ERROR,
  summary: "Server is not available",
  detail: "Please check that server is running.",
  life: 3_000,
}

export function fromFetchWithRetryAndErrorHandling<T>(
  request: Request | string,
  unavailableErrorMessage: ({ summary: string; detail: string }) | null = null,
  bodyConsumer: (response: Response) => Promise<T> = it => it.json() as Promise<T>
): Observable<T> {
  return fromFetch(request)
    .pipe(
      // promise to result
      mergeMap(response => {
        if (response.ok) {
          return bodyConsumer(response)
        }
        else {
          throw new Error(`cannot load (status=${response.status})`)
        }
      }),
      // timeout
      takeUntil(timer(8_000)),
      // retry at least three times
      retry({
        count: 3,
        delay(error, retryCount) {
          return of(error).pipe(delay(1_000 * retryCount))
        },
      }),
      catchError((error, _caught) => {
        if (error instanceof TypeError) {
          if (unavailableErrorMessage == null) {
            ToastEventBus.emit("remove", serverNotAvailableMessage)
            ToastEventBus.emit("add", serverNotAvailableMessage)
          }
          else {
            ToastEventBus.emit("add", {
              severity: ToastSeverity.ERROR,
              summary: unavailableErrorMessage.summary,
              detail: unavailableErrorMessage.detail,
              life: 3_000,
            })
          }
        }
        else {
          ToastEventBus.emit("add", {
            severity: ToastSeverity.ERROR,
            summary: "Cannot load",
            detail: `${(error as Error).message}`,
            life: 3_000
          })
        }
        console.error("cannot load", request, error)
        return EMPTY
      })
    )
}
