import { createPinia, setActivePinia } from "pinia"
import { expect, beforeAll, afterEach, describe, it, vi } from "vitest"
import { ref } from "vue"
import { useRouter } from "vue-router"
import { PersistentStateManager } from "../../src/components/common/PersistentStateManager"
import { Accident, AccidentKind, AccidentsConfigurator } from "../../src/configurators/accidents/AccidentsConfigurator"
import { TimeRangeConfigurator } from "../../src/configurators/TimeRangeConfigurator"
import { AccidentsConfiguratorForStartup } from "../../src/configurators/accidents/AccidentsConfiguratorForStartup"
import { useUserStore } from "../../src/shared/useUserStore"

describe("Branch configurator", () => {
  const serverUrl = "http://localhost:7474"
  let configurator: AccidentsConfiguratorForStartup

  beforeAll(() => {
    setActivePinia(createPinia())
    vi.spyOn(AccidentsConfigurator.prototype, "getAccidentsFromMetaDb").mockResolvedValue(new Map<string, Accident[]>())
    useUserStore()
    const persistence = new PersistentStateManager("test-dashboard", {}, useRouter())
    const timeRangeConfigurator = new TimeRangeConfigurator(persistence)
    configurator = new AccidentsConfiguratorForStartup(serverUrl, ref("RM"), ref(null), ref(null), timeRangeConfigurator)
  })

  afterEach(() => {
    vi.mocked(globalThis.fetch).mockRestore?.()
  })

  it("valid query to create accident for startup", async () => {
    const fetchSpy = vi.spyOn(globalThis, "fetch").mockResolvedValue(new Response("1", { status: 200 }))

    await configurator.writeAccidentToMetaDb("Dec 17, 2023, 5:53 AM", "diaspora", "test", "241.120", AccidentKind.Regression)

    expect(fetchSpy).toHaveBeenCalledWith(serverUrl + "/api/meta/accidents/", expect.objectContaining({ method: "POST" }))
    const [, init] = fetchSpy.mock.calls[0]
    expect(JSON.parse(init!.body as string)).toStrictEqual({
      date: "Dec 17, 2023, 5:53 AM",
      affected_test: "RM/diaspora",
      reason: "test",
      build_number: "241.120",
      kind: "Regression",
      stacktrace: "",
      user_name: "",
    })
  })

  it("valid query to create accident with stacktrace for startup", async () => {
    const fetchSpy = vi.spyOn(globalThis, "fetch").mockResolvedValue(new Response("1", { status: 200 }))

    await configurator.writeAccidentToMetaDb("Dec 17, 2023, 5:53 AM", "diaspora", "test", "241.120", AccidentKind.Exception, "some trace")

    const [, init] = fetchSpy.mock.calls[0]
    expect(JSON.parse(init!.body as string)).toStrictEqual({
      date: "Dec 17, 2023, 5:53 AM",
      affected_test: "RM/diaspora",
      reason: "test",
      build_number: "241.120",
      kind: "Exception",
      stacktrace: "some trace",
      user_name: "",
    })
  })
})
