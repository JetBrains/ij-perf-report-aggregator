import { deepEqual } from "fast-equals"
import pLimit, { LimitFunction } from "p-limit"
import { catchError, delay, distinctUntilChanged, EMPTY, mergeMap, Observable, of, retry } from "rxjs"
import { fromPromise } from "rxjs/internal/observable/from"
import { Ref, watch } from "vue"

export function refToObservable<T>(ref: Ref<T>, deep = false): Observable<T> {
  return new Observable<T>((context) => {
    watch(
      ref,
      (value) => {
        return context.next(value)
      },
      { deep }
    )
    context.next(ref.value)
  }).pipe(deep ? distinctUntilChanged(deepEqual) : distinctUntilChanged())
}
export const limit: LimitFunction = pLimit(25)

export function defaultBodyConsumer<T>(response: Response): Promise<T> {
  return response
    .clone()
    .text()
    .then((text) => {
      try {
        return JSON.parse(text) as T
      } catch {
        throw new Error("Invalid JSON")
      }
    })
}

export function fromFetchWithRetryAndErrorHandling<T>(
  request: Request | string,
  bodyConsumer: (response: Response) => Promise<T> = defaultBodyConsumer,
  controller: AbortController | null = null
): Observable<T> {
  const signal = (controller ?? new AbortController()).signal
  return fromPromise(limit(() => fetch(request, { signal }))).pipe(
    // promise to result
    mergeMap((response) => {
      if (response.ok) {
        return bodyConsumer(response)
      } else {
        throw new Error(`cannot load (status=${response.status})`)
      }
    }),
    retry({
      count: 30,
      delay(error) {
        return of(error).pipe(delay(2000))
      },
    }),
    catchError((error, _caught) => {
      console.error("cannot load", request, error)
      return EMPTY
    })
  )
}
