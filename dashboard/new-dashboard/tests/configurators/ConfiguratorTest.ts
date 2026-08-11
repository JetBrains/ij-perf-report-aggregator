import { Observable } from "rxjs"
import { MockInstance, vi } from "vitest"
import { PersistentStateManager } from "../../src/components/common/PersistentStateManager"
import { ServerConfigurator } from "../../src/components/common/dataQuery"
import { TimeRangeConfigurator } from "../../src/configurators/TimeRangeConfigurator"
import * as rxjs from "../../src/configurators/rxjs"
import { TestServerConfigurator } from "../dummy/TestServerConfigurator"
import { withSetup } from "../withSetup"

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
    const { persistenceForDashboard, timeRangeConfigurator } = withSetup(() => {
      const persistenceForDashboard = new PersistentStateManager(
        "test-dashboard",
        {
          machine: "machine",
          project: [],
          branch: "b1",
        },
        null
      )
      return { persistenceForDashboard, timeRangeConfigurator: new TimeRangeConfigurator(persistenceForDashboard) }
    })

    const serverUrl = `${serverConfigurator.serverUrl}/api/q/`

    return { serverConfigurator, persistenceForDashboard, timeRangeConfigurator, fetchMock, serverUrl }
  },
}
