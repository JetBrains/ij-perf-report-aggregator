import { http, HttpResponse } from "msw"
import { setupServer } from "msw/node"
import { createPinia, setActivePinia } from "pinia"
import { assert, beforeAll, afterAll, describe, test, afterEach } from "vitest"
import { ref } from "vue"
import { useRouter } from "vue-router"
import { PersistentStateManager } from "../../src/components/common/PersistentStateManager"
import { AccidentKind, AccidentsConfiguratorForStartup } from "../../src/configurators/AccidentsConfigurator"
import { TimeRangeConfigurator } from "../../src/configurators/TimeRangeConfigurator"

describe("Branch configurator", () => {
  let timeRangeConfigurator: TimeRangeConfigurator

  const server = setupServer()

  beforeAll(() => {
    server.listen()
    setActivePinia(createPinia())
    const persistence = new PersistentStateManager("test-dashboard", {}, useRouter())
    timeRangeConfigurator = new TimeRangeConfigurator(persistence)
  })

  afterAll(() => {
    server.close()
  })

  afterEach(() => {
    server.resetHandlers()
  })

  test("Valid query to create accident for startup", async () => {
    const serverUrl = "http://localhost:7474"
    const testPromise = new Promise<void>((resolve, reject) => {
      server.use(
        http.post(serverUrl + "/api/meta/accidents/", async (req) => {
          try {
            const text = await req.request.json()
            assert.deepEqual(text, {
              date: "Dec 17, 2023, 5:53 AM",
              affected_test: "RM/diaspora",
              reason: "test",
              build_number: "241.120",
              kind: "Regression",
            })
            resolve()
          } catch (error) {
            reject(error as Error)
          }
          return HttpResponse.json({})
        })
      )
    })

    const configurator = new AccidentsConfiguratorForStartup(serverUrl, ref("RM"), ref(null), ref(null), timeRangeConfigurator)
    configurator.writeAccidentToMetaDb("Dec 17, 2023, 5:53 AM", "diaspora", "test", "241.120", AccidentKind.Regression)
    return testPromise
  })
})
