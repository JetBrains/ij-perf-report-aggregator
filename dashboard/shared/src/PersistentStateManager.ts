import { Ref, watch } from "vue"

import { LocationQueryRaw, Router, useRoute } from "vue-router"
import { debounceSync } from "./util/debounce"

declare type State = { [key: string]: number | string | Array<string> | unknown }

export class PersistentStateManager {
  private readonly state: State

  private readonly initializers: Array<() => void> = []

  private readonly debouncedSave = debounceSync(() => {
    localStorage.setItem(this.getKey(), JSON.stringify(this.state))
  }, 300)

  constructor(private id: string, defaultState: State | null = null, private readonly router: Router | null = null) {
    const storedState = localStorage.getItem(this.getKey())
    if (storedState == null) {
      this.state = defaultState ?? {}
    }
    else {
      // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment
      this.state = JSON.parse(storedState)
      if (defaultState != null) {
        this.state = {
          defaultState,
          ...this.state,
        }
      }
    }

    if (this.router != null) {
      const query = useRoute().query
      for (const [name, value] of Object.entries(query)) {
        this.state[name] = value
      }
    }
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

    if (this.router != null) {
      const currentRoute = useRoute()
      const query: LocationQueryRaw = {...currentRoute.query}
      let isChanged = false
      for (const [name, value] of Object.entries(this.state)) {
        if (name !== "serverUrl" && typeof value === "string") {
          if (isChanged || query[name] !== value) {
            query[name] = value
            isChanged = true
          }
        }
      }

      if (isChanged) {
        // noinspection JSIgnoredPromiseFromCall
        void this.router.push({query,})
      }
    }
  }

  add(name: string, value: Ref<unknown>): void {
    if (value == null) {
      throw new Error("value must be not null")
    }

    watch(value, value => {
      const oldValue = this.state[name]
      if (value !== oldValue) {
        this.state[name] = value
        this.debouncedSave()
      }
    })

    const existingValue = this.state[name]
    if (existingValue != null && (typeof existingValue !== "string" || existingValue.length !== 0)) {
      this.initializers.push(() => {
        // console.debug(`[persistentState] set ${name} to ${existingValue}`)
        value.value = existingValue
      })
    }
  }
}
