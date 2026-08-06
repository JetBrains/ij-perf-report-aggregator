import { Observable } from "rxjs"
import { computed } from "vue"
import { beforeEach, describe, expect, it } from "vitest"
import { CompareSectionsRegistry } from "../../src/components/charts/compareMode"
import { dashboardProjectsFilter } from "../../src/configurators/DashboardProjectsFilter"
import { MachineConfigurator } from "../../src/configurators/MachineConfigurator"
import { defaultModeName, TestModeConfigurator } from "../../src/configurators/TestModeConfigurator"
import { awaitCallbackTrue } from "../utils/awaitors"
import ConfiguratorTest, { ConfigurationTestData } from "./ConfiguratorTest"

const groups = [{ group: "linux-blade", machines: ["intellij-linux-hw-blade-test"], predicate: "like 'intellij-linux-hw-blade-%'" }]

function section(id: string, projects: string[]) {
  return { id, label: id, measure: "m", projects, aliases: null, machines: null, valueUnit: "ms" as const }
}

describe("Dashboard projects filter", () => {
  let data: ConfigurationTestData
  let machineGroupsUrl: string
  let registry: CompareSectionsRegistry

  function machineGroupsRequests(): string[] {
    return data.fetchMock.mock.calls.map((call) => call[0] as string).filter((it) => it.startsWith(machineGroupsUrl))
  }

  // Mirrors how DashboardPage derives the project set, so the test exercises the real wiring.
  function machineConfigurator(chartsProjects: string[] = []) {
    const projects = computed(() => [...new Set([...chartsProjects, ...registry.sections.value.flatMap((it) => it.projects)])].toSorted())
    return new MachineConfigurator(data.serverConfigurator, undefined, [dashboardProjectsFilter(projects)])
  }

  beforeEach(() => {
    localStorage.clear()
    registry = new CompareSectionsRegistry()
    data = ConfiguratorTest.setupPreconditions([])
    machineGroupsUrl = data.serverUrl.replace("/api/q/", "/api/machineGroups/")
    data.fetchMock.mockReturnValue(
      new Observable((sub) => {
        sub.next(groups)
      })
    )
  })

  it("coalesces the mount-time registration burst into a single request", async () => {
    machineConfigurator()
    // Every non-lazy chart mounts in the same tick.
    registry.register(section("a", ["p/one"]))
    registry.register(section("b", ["p/two"]))
    registry.register(section("c", ["p/three", "p/one"]))

    await awaitCallbackTrue(() => machineGroupsRequests().length > 0)
    await new Promise((resolve) => {
      setTimeout(resolve, 250)
    })
    // No unfiltered request slipped out before the sections arrived.
    expect(machineGroupsRequests()).toStrictEqual([
      `${machineGroupsUrl}{"db":"test","table":"test","fields":[{"n":"machine","sql":"distinct machine"}],"filters":[{"f":"project","v":["p/one","p/three","p/two"]}],"order":"machine","flat":true}`,
    ])
  })

  it("issues one more request when a lazy accordion chart registers later", async () => {
    machineConfigurator()
    registry.register(section("a", ["p/one"]))
    await awaitCallbackTrue(() => machineGroupsRequests().length === 1)

    // The user expands a ChartAccordion; its charts mount and register now.
    registry.register(section("lazy", ["p/hidden"]))
    await awaitCallbackTrue(() => machineGroupsRequests().length === 2)
    expect(machineGroupsRequests()[1]).toContain('"v":["p/hidden","p/one"]')
  })

  it("does not filter at all while no chart has registered", async () => {
    machineConfigurator()
    await awaitCallbackTrue(() => machineGroupsRequests().length > 0)
    expect(machineGroupsRequests()[0]).toBe(`${machineGroupsUrl}{"db":"test","table":"test","fields":[{"n":"machine","sql":"distinct machine"}],"order":"machine","flat":true}`)
  })

  it("seeds from the charts prop so a :charts page filters before anything mounts", async () => {
    machineConfigurator(["p/from-charts"])
    await awaitCallbackTrue(() => machineGroupsRequests().length > 0)
    expect(machineGroupsRequests()[0]).toContain('"v":["p/from-charts"]')
  })

  // A mode can run on a single hardware class (goland's `wsl` only on windows-azure), so the
  // machine list has to be filtered by it too.
  describe("mode narrowing", () => {
    function withMode(mode: TestModeConfigurator) {
      const projects = computed(() => registry.sections.value.flatMap((it) => it.projects).toSorted())
      return new MachineConfigurator(data.serverConfigurator, undefined, [dashboardProjectsFilter(projects), mode])
    }

    it("sends the selected mode with the machine list query", async () => {
      const mode = new TestModeConfigurator(true)
      mode.selected.value = ["wsl"]
      withMode(mode)

      await awaitCallbackTrue(() => machineGroupsRequests().length > 0)
      expect(machineGroupsRequests()[0]).toContain('{"f":"mode","v":["wsl"]}')
    })

    it("reloads the machine list when the mode changes", async () => {
      const mode = new TestModeConfigurator(true)
      mode.selected.value = ["wsl"]
      withMode(mode)
      await awaitCallbackTrue(() => machineGroupsRequests().length === 1)

      mode.selected.value = ["split"]
      await awaitCallbackTrue(() => machineGroupsRequests().length === 2)
      expect(machineGroupsRequests()[1]).toContain('{"f":"mode","v":["split"]}')
    })

    it("sends the default mode as the empty string the reports are stored under", async () => {
      const mode = new TestModeConfigurator(true)
      mode.selected.value = [defaultModeName]
      withMode(mode)

      await awaitCallbackTrue(() => machineGroupsRequests().length > 0)
      expect(machineGroupsRequests()[0]).toContain('{"f":"mode","v":""}')
    })
  })
})
