import { Observable } from "rxjs"
import { beforeEach, describe, expect, it } from "vitest"
import { MachineConfigurator } from "../../src/configurators/MachineConfigurator"
import { selectedStartupProjectsFilter, startupProjectConfigurator } from "../../src/configurators/StartupProjectConfigurator"
import { awaitCallbackTrue } from "../utils/awaitors"
import ConfiguratorTest, { ConfigurationTestData } from "./ConfiguratorTest"

const startupPredicate = "project like '%/measureStartup%' or project like '%/warmup%'"
const projects = ["clion/clion/cmake/measureStartup", "clion/clion/cmake/warmup", "radler/radler/cmake/measureStartup"]
const allGroups = [
  { group: "linux-blade", machines: ["intellij-linux-hw-blade-test"], predicate: "like 'intellij-linux-hw-blade-%'" },
  { group: "mac large", machines: ["intellij-macos-unit-2200-large-test"], predicate: "like 'intellij-macos-unit-2200-large-%'" },
]

describe("Startup project configurator", () => {
  let data: ConfigurationTestData
  let machineGroupsUrl: string

  // The project list and the machine list share one fetch spy — answer each with its own shape.
  function serveGroups(groups: unknown[]) {
    data.fetchMock.mockImplementation(
      (url: string) =>
        new Observable((sub) => {
          sub.next(url.startsWith(machineGroupsUrl) ? groups : projects)
        })
    )
  }

  function machineGroupsRequests(): string[] {
    return data.fetchMock.mock.calls.map((call) => call[0] as string).filter((it) => it.startsWith(machineGroupsUrl))
  }

  beforeEach(() => {
    localStorage.clear()
    data = ConfiguratorTest.setupPreconditions(projects)
    machineGroupsUrl = data.serverUrl.replace("/api/q/", "/api/machineGroups/")
    serveGroups(allGroups)
  })

  it("merges the startup suffixes into a single selectable project", async () => {
    const configurator = startupProjectConfigurator(data.serverConfigurator, data.persistenceForDashboard, true)
    await awaitCallbackTrue(() => configurator.values.value.length > 0)
    expect(configurator.values.value).toStrictEqual(["clion/clion/cmake", "radler/radler/cmake"])
  })

  describe("machine list narrowing", () => {
    it("restricts the machine list to any startup project while nothing is selected", async () => {
      const configurator = startupProjectConfigurator(data.serverConfigurator, data.persistenceForDashboard, true)
      expect(new MachineConfigurator(data.serverConfigurator, undefined, [selectedStartupProjectsFilter(configurator)])).toBeDefined()

      await awaitCallbackTrue(() => machineGroupsRequests().length > 0)
      expect(machineGroupsRequests()[0]).toBe(
        `${machineGroupsUrl}{"db":"test","table":"test","fields":[{"n":"machine","sql":"distinct machine"}],"filters":[{"f":"","q":"${startupPredicate}"}],"order":"machine","flat":true}`
      )
    })

    it("expands the selected project stems back to the reported project names", async () => {
      const configurator = startupProjectConfigurator(data.serverConfigurator, data.persistenceForDashboard, true)
      await awaitCallbackTrue(() => configurator.values.value.length > 0)

      expect(new MachineConfigurator(data.serverConfigurator, undefined, [selectedStartupProjectsFilter(configurator)])).toBeDefined()
      configurator.selected.value = ["radler/radler/cmake", "clion/clion/cmake"]

      const expected =
        `${machineGroupsUrl}{"db":"test","table":"test","fields":[{"n":"machine","sql":"distinct machine"}],` +
        `"filters":[{"f":"project","v":["clion/clion/cmake/measureStartup","clion/clion/cmake/warmup","radler/radler/cmake/measureStartup","radler/radler/cmake/warmup"]}],` +
        `"order":"machine","flat":true}`
      await awaitCallbackTrue(() => machineGroupsRequests().includes(expected))
      expect(machineGroupsRequests()).toContain(expected)
    })

    it("drops a group that did not run the selected project", async () => {
      const configurator = startupProjectConfigurator(data.serverConfigurator, data.persistenceForDashboard, true)
      await awaitCallbackTrue(() => configurator.values.value.length > 0)

      const machineConfigurator = new MachineConfigurator(data.serverConfigurator, undefined, [selectedStartupProjectsFilter(configurator)])
      await awaitCallbackTrue(() => machineConfigurator.values.value.length === 2)

      // From here on, only the linux group ran the selected project.
      serveGroups([allGroups[0]])
      configurator.selected.value = ["radler/radler/cmake"]
      await awaitCallbackTrue(() => machineConfigurator.values.value.length === 1)
      expect(machineConfigurator.values.value.map((it) => it.value)).toStrictEqual(["linux-blade"])
    })
  })
})
