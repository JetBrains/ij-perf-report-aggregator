import { MonoTypeOperatorFunction, tap } from "rxjs"
import { shallowReactive } from "vue"

export class ComponentState {
  loading = false
  // primevue has disabled property, so, we use "disabled" and not "enabled"
  disabled = true
}

export function createComponentState(): ComponentState {
  return shallowReactive(new ComponentState())
}

export function updateComponentState<T>(status: ComponentState): MonoTypeOperatorFunction<T> {
  return tap<T>({
    next(value) {
      status.disabled = value === null
      status.loading = false
    },
    error(_error) {
      status.loading = false
    },
  })
}