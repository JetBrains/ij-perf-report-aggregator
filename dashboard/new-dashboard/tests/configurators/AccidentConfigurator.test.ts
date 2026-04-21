import { createPinia, setActivePinia } from "pinia"
import { expect, beforeAll, afterEach, describe, it, vi } from "vitest"
import { ref } from "vue"
import { useRouter } from "vue-router"
import { PersistentStateManager } from "../../src/components/common/PersistentStateManager"
import { AccidentKind } from "../../src/configurators/accidents/AccidentsConfigurator"
import { TimeRangeConfigurator } from "../../src/configurators/TimeRangeConfigurator"
import { AccidentsConfiguratorForStartup } from "../../src/configurators/accidents/AccidentsConfiguratorForStartup"

describe("Branch configurator", () => {
  let timeRangeConfigurator: TimeRangeConfigurator
  const serverUrl = "http://localhost:7474"

  beforeAll(() => {
    setActivePinia(createPinia())
    const persistence = new PersistentStateManager("test-dashboard", {}, useRouter())
    timeRangeConfigurator = new TimeRangeConfigurator(persistence)
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  it("valid query to create accident for startup", async () => {
    const fetchSpy = vi.spyOn(globalThis, "fetch").mockImplementation(() => Promise.resolve(new Response("1", { status: 200 })))

    const configurator = new AccidentsConfiguratorForStartup(serverUrl, ref("RM"), ref(null), ref(null), timeRangeConfigurator)
    await configurator.writeAccidentToMetaDb("Dec 17, 2023, 5:53 AM", "diaspora", "test", "241.120", AccidentKind.Regression)

    expect(fetchSpy).toHaveBeenCalledWith(serverUrl + "/api/meta/accidents/", expect.objectContaining({ method: "POST" }))
    const accidentCall = fetchSpy.mock.calls.find(([url]) => url === serverUrl + "/api/meta/accidents/")
    const [, init] = accidentCall!
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
    const fetchSpy = vi.spyOn(globalThis, "fetch").mockImplementation(() => Promise.resolve(new Response("1", { status: 200 })))

    const configurator = new AccidentsConfiguratorForStartup(serverUrl, ref("RM"), ref(null), ref(null), timeRangeConfigurator)
    await configurator.writeAccidentToMetaDb("Dec 17, 2023, 5:53 AM", "diaspora", "test", "241.120", AccidentKind.Exception, "some trace")

    const accidentCall = fetchSpy.mock.calls.find(([url]) => url === serverUrl + "/api/meta/accidents/")
    const [, init] = accidentCall!
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
