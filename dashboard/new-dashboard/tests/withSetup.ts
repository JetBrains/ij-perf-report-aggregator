import { createApp, h } from "vue"

// Runs a factory inside a real component setup() so that provide()/inject() calls
// (e.g. in the TimeRangeConfigurator constructor) have a component context and
// don't trigger "[Vue warn]: provide() can only be used inside setup()".
// The app is intentionally left mounted so watchers created by the factory keep running.
export function withSetup<T>(factory: () => T): T {
  let result!: T
  const app = createApp({
    setup() {
      result = factory()
      return () => h("div")
    },
  })
  app.mount(document.createElement("div"))
  return result
}
