import { Ref, watch } from "vue"
import { debounce } from "./util/debounce"

export class PersistentStateManager {
  private readonly state: { [key: string]: number | string | unknown }

  private readonly initializers: Array<() => void> = []

  private readonly debouncedSave = debounce(() => {
    localStorage.setItem(this.getKey(), JSON.stringify(this.state))
  }, 300)

  constructor(private id: string) {
    const storedState = localStorage.getItem(this.getKey())
    this.state = storedState == null ? {} : JSON.parse(storedState)
  }

  private getKey(): string {
    return `${this.id}-state-v2`
  }

  /**
   * Initialization cannot be performed as part of `add`, because at the moment of add,
   * not every other possible user of data maybe yet created and therefore data is not being watched
   * (avoid explicit `scheduleLoad` or `load` calls).
   */
  init(): void {
    for (const initializer of this.initializers) {
      initializer()
    }
    this.initializers.length = 0
  }

  add(name: string, value: Ref): void {
    if (value == null) {
      throw new Error("value must be not null")
    }

    watch(value, (value) => {
      const oldValue = this.state[name]
      if (value !== oldValue) {
        this.state[name] = value
        this.debouncedSave()
      }
    })

    const existingValue = this.state[name]
    if (existingValue != null) {
      this.initializers.push(() => {
        value.value = existingValue
      })
    }
  }
}
