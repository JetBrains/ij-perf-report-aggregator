import { deepEqual } from "fast-equals"
import pLimit, { LimitFunction } from "p-limit"
import { ToastSeverity } from "primevue/api"
import ToastEventBus from "primevue/toasteventbus"
import { catchError, delay, distinctUntilChanged, EMPTY, mergeMap, Observable, of, retry} from "rxjs"
import { fromPromise } from "rxjs/internal/observable/from"
import { Ref, watch } from "vue"

export function refToObservable<T>(ref: Ref<T>, deep= false): Observable<T> {
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

export const limit: LimitFunction = pLimit(100)

export function fromFetchWithRetryAndErrorHandling<T>(
  request: Request | string,
  unavailableErrorMessage: ({ summary: string; detail: string }) | null = null,
  bodyConsumer: (response: Response) => Promise<T> = it => it.json() as Promise<T>,
  controller: AbortController | null = null,
): Observable<T> {
  const signal = (controller ?? new AbortController()).signal
  return fromPromise(limit(() => fetch(request, {signal})))
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
