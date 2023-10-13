import { Observable } from "rxjs"
import { Ref, shallowRef } from "vue"
import { refToObservable } from "../../../configurators/rxjs"
import { PersistentStateManager } from "../PersistentStateManager"

export class MeasureConfiguratorForComparing {
  readonly data = shallowRef<string[]>([])
  private readonly _selected = shallowRef<string[] | string | null>(null)
  readonly state = shallowRef({
    loading: false,
    disabled: false,
  })

  createObservable(): Observable<unknown> {
    return refToObservable(this.selected, true)
  }

  setSelected(value: string[] | string | null) {
    this._selected.value = value
  }

  get selected(): Ref<string[] | null> {
    const ref = this._selected
    if (typeof ref.value === "string") {
      ref.value = [ref.value]
    }
    return ref as Ref<string[] | null>
  }

  constructor(metrics: string[], persistentStateManager: PersistentStateManager) {
    persistentStateManager.add("measure", this._selected)

    const selectedRef = this.selected
    this.data.value = metrics
    selectedRef.value = [...metrics]
  }
}
