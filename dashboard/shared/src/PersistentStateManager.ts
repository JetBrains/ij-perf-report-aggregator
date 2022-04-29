import { debounceTime, Subject } from "rxjs"
import { Ref, watch } from "vue"
import { LocationQueryRaw, RouteLocationNormalizedLoaded, Router, useRoute } from "vue-router"

declare type State = { [key: string]: number | string | Array<string> | unknown }

export class PersistentStateManager {
  private readonly state: State

  private readonly initializers: Array<() => void> = []

  private readonly saveSubject = new Subject<null>()

  private readonly route: RouteLocationNormalizedLoaded | null

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

    if (this.router == null) {
      this.route = null
    }
    else {
      this.route = useRoute()
      const route = this.route
      const query = route.query
      for (const [name, value] of Object.entries(query)) {
        this.state[name] = value
      }
    }

    this.saveSubject.pipe(
      debounceTime(300),
    )
      .subscribe(() => {
        localStorage.setItem(this.getKey(), JSON.stringify(this.state))

        this.updateUrlQuery()
      })
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
    this.updateUrlQuery()
  }

  private updateUrlQuery() {
    if (this.router == null) {
      return
    }

    // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
    const currentRoute = this.route!
    const query: LocationQueryRaw = {...currentRoute.query}
    let isChanged = false
    for (const [name, value] of Object.entries(this.state)) {
      if (name !== "serverUrl" && typeof value === "string" || Array.isArray(value) || value == null) {
        if (isChanged || query[name] !== value) {
          query[name] = value
          isChanged = true
        }
      }
    }

    if (isChanged) {
      // noinspection JSIgnoredPromiseFromCall
      void this.router.push({query})
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
        this.saveSubject.next(null)
      }
    })

    let existingValue = this.state[name]
    if (existingValue != null && (typeof existingValue !== "string" || existingValue.length !== 0)) {
      // noinspection SpellCheckingInspection
      if (existingValue === "73YWaW9bytiPDGuKvwNIYMK5CKI") {
        existingValue = "simple for IJ"
      }
      // noinspection SpellCheckingInspection
      if (existingValue === "nC4MRRFMVYUSQLNIvPgDt+B3JqA") {
        existingValue = "idea"
      }
      this.initializers.push(() => {
        // console.debug(`[persistentState] set ${name} to ${existingValue}`)
        value.value = existingValue
      })
    }
  }
}
