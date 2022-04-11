import { Observable } from "rxjs"
import { Ref, watch } from "vue"

export function refToObservable<T>(ref: Ref<T>, deep: boolean = false): Observable<T> {
  return new Observable(context => {
    watch(ref, value => {
      return context.next(value)
    }, {deep})
    context.next(ref.value)
  })
}