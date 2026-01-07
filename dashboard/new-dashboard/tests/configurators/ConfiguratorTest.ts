import { Observable } from "rxjs"
import { MockInstance, vi } from "vitest"
import { useRouter } from "vue-router"
import { PersistentStateManager } from "../../src/components/common/PersistentStateManager"
import { ServerConfigurator } from "../../src/components/common/dataQuery"
import { TimeRangeConfigurator } from "../../src/configurators/TimeRangeConfigurator"
import * as rxjs from "../../src/configurators/rxjs"
import { TestServerConfigurator } from "../dummy/TestServerConfigurator"

export interface ConfigurationTestData {
  serverConfigurator: ServerConfigurator
  persistenceForDashboard: PersistentStateManager
  timeRangeConfigurator: TimeRangeConfigurator
  fetchMock: MockInstance
  serverUrl: string
}

export default {
  setupPreconditions(mockValue: string[]): ConfigurationTestData {
    const fetchMock: MockInstance = vi.spyOn(rxjs, "fromFetchWithRetryAndErrorHandling").mockClear().mockReset()

    fetchMock.mockReturnValue(
      new Observable((sub) => {
        sub.next(mockValue)
      })
    )

    const serverConfigurator = new TestServerConfigurator("test", "test")
    const persistenceForDashboard = new PersistentStateManager(
      "test-dashboard",
      {
        machine: "machine",
        project: [],
        branch: "b1",
      },
      useRouter()
    )
    const timeRangeConfigurator = new TimeRangeConfigurator(persistenceForDashboard)

    const serverUrl = `${serverConfigurator.serverUrl}/api/q/`

    return { serverConfigurator, persistenceForDashboard, timeRangeConfigurator, fetchMock, serverUrl }
  },
}
