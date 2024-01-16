import { debounceTime, Subject } from "rxjs"
import { Ref, watch } from "vue"
import { LocationQueryRaw, RouteLocationNormalizedLoaded, Router, useRoute } from "vue-router"

// eslint-disable-next-line @typescript-eslint/no-redundant-type-constituents
type State = Record<string, number | string | string[] | unknown>

export class PersistentStateManager {
  private readonly state: State

  private readonly saveSubject = new Subject<null>()
  private readonly updateUrlSubject = new Subject<null>()

  private MAX_ARRAY_SIZE = 30

  private readonly route: RouteLocationNormalizedLoaded | null

  constructor(
    private readonly id: string,
    defaultState: State | null = null,
    private readonly router: Router | null = null
  ) {
    const storedState = localStorage.getItem(this.getKey())
    if (storedState == null) {
      this.state = defaultState ?? {}
    } else {
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
    } else {
      this.route = useRoute()
      const route = this.route
      const query = route.query
      for (const [name, value] of Object.entries(query)) {
        if (value === "[]") {
          this.state[name] = []
        } else {
          this.state[name] = Array.isArray(value) ? value.map((element) => boolFromString(element)) : boolFromString(value)
        }
      }
    }

    this.saveSubject.pipe(debounceTime(300)).subscribe(() => {
      localStorage.setItem(this.getKey(), JSON.stringify(this.state))

      this.updateUrlQuery()
    })

    this.updateUrlSubject.pipe(debounceTime(300)).subscribe(() => {
      this.updateUrlQuery()
    })
  }

  private getKey(): string {
    return `${this.id}-state-v2`
  }

  private updateUrlQuery() {
    if (this.router == null) {
      return
    }

    // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
    const currentRoute = this.route!
    const query: LocationQueryRaw = { ...currentRoute.query }
    let isChanged = false
    for (const [name, value] of Object.entries(this.state)) {
      if (((name !== "serverUrl" && typeof value === "string") || Array.isArray(value) || value === null) && (isChanged || query[name] !== value)) {
        // Persist empty arrays as `[]` to allow 0-value selections shared via the URL to override a user's local state.
        query[name] = Array.isArray(value) && value.length === 0 ? "[]" : value
        isChanged = true
      }
    }

    if (isChanged) {
      // noinspection JSIgnoredPromiseFromCall
      void this.router.push({ query })
    }
  }

  add(name: string, value: Ref<unknown>, existingValueTransformer: ((v: unknown) => unknown) | null = null): void {
    watch(value, (value) => {
      const oldValue = this.state[name]
      if (Array.isArray(value) && value.length > this.MAX_ARRAY_SIZE) {
        value = value.slice(0, this.MAX_ARRAY_SIZE)
      }
      if (value !== oldValue) {
        this.state[name] = value
        this.saveSubject.next(null)
      }
    })

    const existingValue = this.state[name]
    if (existingValue != null && (typeof existingValue !== "string" || existingValue.length > 0)) {
      // console.debug(`[persistentState] set ${name} to ${existingValue}`)
      value.value = existingValueTransformer == null ? existingValue : existingValueTransformer(existingValue)
      this.updateUrlSubject.next(null)
    }
  }
}

function boolFromString(s: string | null) {
  if (s == "true") return true
  if (s == "false") return false
  return s
}
