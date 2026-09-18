import { createPinia, setActivePinia } from "pinia"
import { beforeEach, describe, expect, it } from "vitest"
import { nextTick } from "vue"
import { DataQuery, DataQueryExecutorConfiguration } from "../../../src/components/common/dataQuery"
import { SettingConfigurator } from "../../../src/components/settings/SettingConfigurator"
import { useSettingsStore } from "../../../src/components/settings/settingsStore"

describe("the generic setting configurator", () => {
  beforeEach(() => {
    localStorage.clear()
    setActivePinia(createPinia())
  })

  it("starts with the current value of the setting", () => {
    useSettingsStore().removeOutliers = true
    const emitted: unknown[] = []
    new SettingConfigurator("removeOutliers").createObservable().subscribe((value) => {
      emitted.push(value)
    })
    expect(emitted).toStrictEqual([true])
  })

  it("emits when the setting is flipped", async () => {
    const configurator = new SettingConfigurator("detectChanges")
    const emitted: unknown[] = []
    configurator.createObservable().subscribe((value) => {
      emitted.push(value)
    })

    useSettingsStore().detectChanges = true
    await nextTick()

    expect(emitted).toStrictEqual([false, true])
  })

  it("leaves the query and the filter untouched", () => {
    const configurator = new SettingConfigurator("scaling")
    const query = new DataQuery()
    expect([configurator.configureQuery(query, new DataQueryExecutorConfiguration()), configurator.configureFilter(query)]).toStrictEqual([true, true])
    expect(query.fields).toStrictEqual([])
  })
})
