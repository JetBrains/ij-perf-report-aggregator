import { deepEqual } from "fast-equals"
import pLimit, { LimitFunction } from "p-limit"
import { catchError, defer, delay, distinctUntilChanged, EMPTY, from, mergeMap, Observable, of, retry } from "rxjs"
import { ref, Ref, watch } from "vue"

export function refToObservable<T>(ref: Ref<T>, deep = false): Observable<T> {
  return new Observable<T>((context) => {
    watch(
      ref,
      (value) => {
        context.next(value)
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
  const numberOfRetries = ref(0)
  return defer(() => {
    const signal = (controller ?? new AbortController()).signal
    return from(
      limit(() =>
        fetch(request, {
          signal,
          headers: numberOfRetries.value > 0 ? { "Cache-Control": "no-cache" } : undefined,
        })
      )
    )
  }).pipe(
    // promise to result
    mergeMap((response) => {
      if (response.ok) {
        return bodyConsumer(response)
      }
      throw new Error(`cannot load (status=${response.status})`)
    }),
    retry({
      count: 10,
      delay(error) {
        numberOfRetries.value++
        return of(error).pipe(delay(1000))
      },
    }),
    catchError((error, _caught) => {
      console.error("cannot load", request, error)
      return EMPTY
    })
  )
}
