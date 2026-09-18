import { createPinia, setActivePinia } from "pinia"
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest"
import { createApp, defineComponent, h, nextTick } from "vue"
import type ToggleSwitch from "openvue/toggleswitch"

// the switch is auto-imported into the SFC, so it is replaced here instead of being stubbed globally
vi.mock(import("openvue/toggleswitch"), () => {
  const stub = defineComponent({
    props: { modelValue: { type: Boolean, required: true } },
    emits: { "update:modelValue": (_: boolean) => true },
    setup(props, { emit }) {
      return () =>
        h(
          "button",
          {
            onClick: () => {
              emit("update:modelValue", !props.modelValue)
            },
          },
          String(props.modelValue)
        )
    },
  })
  return { default: stub as unknown as typeof ToggleSwitch }
})

const { default: SettingSwitch } = await import("../../../src/components/settings/SettingSwitch.vue")
const { useSettingsStore } = await import("../../../src/components/settings/settingsStore")

let container: HTMLElement | null = null
let unmount: (() => void) | null = null

function mountSwitch(setting: "scaling" | "smoothing", label = "Scaling"): HTMLElement {
  container = document.createElement("div")
  document.body.append(container)
  const app = createApp(SettingSwitch, { setting, label, tooltip: `explains ${label}` })
  app.mount(container)
  unmount = () => {
    app.unmount()
  }
  return container
}

describe("SettingSwitch", () => {
  beforeEach(() => {
    localStorage.clear()
    setActivePinia(createPinia())
  })

  afterEach(() => {
    unmount?.()
    container?.remove()
    unmount = null
    container = null
  })

  it("labels the switch with the name of the setting", () => {
    const element = mountSwitch("scaling", "Remove outliers (Alpha)")
    expect(element.querySelector("span")?.textContent).toBe("Remove outliers (Alpha):")
  })

  it("reflects the current value of the setting", async () => {
    useSettingsStore().scaling = true
    const element = mountSwitch("scaling")
    await nextTick()
    expect(element.querySelector("button")?.textContent).toBe("true")
  })

  it("writes back through the store setter, so scaling and smoothing stay mutually exclusive", async () => {
    const settingsStore = useSettingsStore()
    settingsStore.smoothing = true

    const element = mountSwitch("scaling")
    element.querySelector("button")?.click()
    await nextTick()

    expect([settingsStore.scaling, settingsStore.smoothing]).toStrictEqual([true, false])
  })
})
