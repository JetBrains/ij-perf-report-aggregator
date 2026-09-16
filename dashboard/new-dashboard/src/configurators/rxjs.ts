import { deepEqual } from "fast-equals"
import pLimit, { LimitFunction } from "p-limit"
import { ColdObservable } from "rxjs"
import { catchError } from "rxjs/catch-error"
import { combineLatest as combineLatestSymbol } from "rxjs/combine-latest"
import { delay } from "rxjs/delay"
import { distinctUntilChanged } from "rxjs/distinct-until-changed"
import { mergeMap } from "rxjs/merge-map"
import { retry } from "rxjs/retry"
import { ref, Ref, watch } from "vue"

// Every subscriber gets its own watcher and the current value right away (ColdObservable, not the
// platform Observable, whose concurrent subscribers share one producer and would miss the initial emit).
export function refToObservable<T>(ref: Ref<T>, deep = false): Observable<T> {
  return new ColdObservable<T>((subscriber) => {
    const stop = watch(
      ref,
      (value) => {
        subscriber.next(value)
      },
      { deep }
    )
    subscriber.addTeardown(stop)
    subscriber.next(ref.value)
  })[distinctUntilChanged](deep ? deepEqual : undefined)
}

// The beta's `combineLatest` typing has no `const` type parameter, so an array literal is inferred as
// `Observable<A | B>[]` and the tuple shape is lost. This wrapper restores it, and keeps every combination cold.
export function combineLatest<const Sources extends readonly Observable<unknown>[]>(
  sources: Sources
): Observable<{ [K in keyof Sources]: Sources[K] extends ObservableValue<infer T> ? T : never }> {
  return ColdObservable[combineLatestSymbol](sources)
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
  // the fetch starts on subscribe (and again on every retry), not when this function is called
  return ColdObservable.from([request])
    [mergeMap]((request) => {
      const signal = (controller ?? new AbortController()).signal
      return limit(() =>
        fetch(request, {
          signal,
          headers: numberOfRetries.value > 0 ? { "Cache-Control": "no-cache" } : undefined,
        })
      )
    })
    [mergeMap]((response) => {
      if (response.ok) {
        return bodyConsumer(response)
      }
      throw new Error(`cannot load (status=${response.status})`)
    })
    [retry]({
      count: 10,
      delay(error) {
        numberOfRetries.value++
        return ColdObservable.from([error])[delay](1000)
      },
    })
    [catchError]((error) => {
      console.error("cannot load", request, error)
      return ColdObservable.from<never>([])
    })
}
