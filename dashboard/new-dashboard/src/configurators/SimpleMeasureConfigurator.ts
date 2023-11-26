import { Observable } from "rxjs"
import { ref, Ref, shallowRef } from "vue"
import { PersistentStateManager } from "../components/common/PersistentStateManager"
import { DataQueryConfigurator } from "../components/common/dataQuery"
import { refToObservable } from "./rxjs"

export class SimpleMeasureConfigurator implements DataQueryConfigurator {
  readonly data = shallowRef<string[] | null>(null)
  private readonly _selected = shallowRef<string[] | string | null>(null)
  readonly state = {
    loading: true,
    disabled: true,
  }
  configureQuery(): boolean {
    return true
  }

  initData(value: string[]) {
    this.data.value = value
    this._selected.value = value
    this.state.loading = false
    this.state.disabled = false
  }

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

  get selectedSafe(): Ref<string[]> {
    const reference = this._selected
    if (reference.value === null) {
      return ref([])
    }
    if (typeof reference.value === "string") {
      reference.value = [reference.value]
    }
    return reference as Ref<string[]>
  }
  constructor(measureName: string, persistentStateManager: PersistentStateManager | null) {
    persistentStateManager?.add(measureName, this._selected)
  }
}
