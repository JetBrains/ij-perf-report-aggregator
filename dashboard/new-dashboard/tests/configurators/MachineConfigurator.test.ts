import { Observable } from "rxjs"
import { expect, beforeEach, afterEach, describe, it, vi } from "vitest"
import { DataQueryExecutor } from "../../src/components/common/DataQueryExecutor"
import { BranchConfigurator, createBranchConfigurator } from "../../src/configurators/BranchConfigurator"
import { MachineConfigurator } from "../../src/configurators/MachineConfigurator"
import { TestMeasureConfigurator } from "../dummy/TestMeasureConfigurator"
import { awaitCallbackTrue, awaitMockCallsCount } from "../utils/awaitors"
import ConfiguratorTest, { ConfigurationTestData } from "./ConfiguratorTest"

// The backend returns machines already grouped by hardware class. Each test group has a single
// member, so selecting a group expands to exactly one machine — keeping the query assertions
// identical to selecting that machine directly.
const machineGroupsResponse = [
  { group: "linux-blade", machines: ["intellij-linux-hw-blade-test"] },
  { group: "mac large", machines: ["intellij-macos-unit-2200-large-test"] },
]

describe("Machine configurator", () => {
  let data: ConfigurationTestData
  let machineConfigurator: MachineConfigurator
  let dataQueryExecutor: DataQueryExecutor
  let machineGroupsUrl: string

  beforeEach(() => {
    data = ConfiguratorTest.setupPreconditions(["intellij-linux-hw-blade-test", "intellij-macos-unit-2200-large-test"])
    machineGroupsUrl = data.serverUrl.replace("/api/q/", "/api/machineGroups/")

    // The machine list is fetched from /api/machineGroups (already grouped); everything else
    // (chart data queries) keeps returning the plain list — its response is unused by assertions.
    data.fetchMock.mockImplementation(
      (url: string) =>
        new Observable((sub) => {
          sub.next(url.includes("/api/machineGroups/") ? machineGroupsResponse : [])
        })
    )

    // Raw-agent -> group resolution is a plain fetch; stub it to an unmatched group so it never
    // rewrites the selection in tests.
    vi.stubGlobal(
      "fetch",
      vi.fn(() => Promise.resolve({ json: () => Promise.resolve({ group: "__unmatched__" }) }))
    )
  })

  afterEach(() => {
    vi.unstubAllGlobals()
  })

  describe("tests without filters", () => {
    beforeEach(() => {
      machineConfigurator = new MachineConfigurator(data.serverConfigurator, data.persistenceForDashboard)
      const measureConfigurator = new TestMeasureConfigurator()
      dataQueryExecutor = new DataQueryExecutor([data.serverConfigurator, machineConfigurator, measureConfigurator])
    })

    it("valid query on configurator init", async () => {
      await awaitMockCallsCount(data.fetchMock, 1)
      const expectedValue = `${machineGroupsUrl}{"db":"test","table":"test","fields":[{"n":"machine","sql":"distinct machine"}],"order":"machine","flat":true}`
      expect(data.fetchMock.mock.calls[0][0]).toBe(expectedValue)
    })

    it("valid query with default selected value", async () => {
      dataQueryExecutor.subscribe(() => {})
      await awaitMockCallsCount(data.fetchMock, 2)
      const expectedValue = `${data.serverUrl}[{"db":"test","table":"test","fields":[],"filters":[{"f":"machine","v":["machine"]}]}]`
      expect(data.fetchMock.mock.calls[1][0]).toBe(expectedValue)
    })

    it("valid query when select single value", async () => {
      dataQueryExecutor.subscribe(() => {})
      await awaitMockCallsCount(data.fetchMock, 2)

      machineConfigurator.selected.value = ["linux-blade"]
      let expectedValue = `${data.serverUrl}[{"db":"test","table":"test","fields":[],"filters":[{"f":"machine","v":["intellij-linux-hw-blade-test"]}]}]`
      await awaitMockCallsCount(data.fetchMock, 3)
      expect(data.fetchMock.mock.calls[2][0]).toBe(expectedValue)

      expectedValue = `${data.serverUrl}[{"db":"test","table":"test","fields":[],"filters":[{"f":"machine","v":["intellij-macos-unit-2200-large-test"]}]}]`
      machineConfigurator.selected.value = ["mac large"]
      await awaitMockCallsCount(data.fetchMock, 4)
      expect(data.fetchMock.mock.calls[3][0]).toBe(expectedValue)
    })

    it("valid query when select multiple value", async () => {
      dataQueryExecutor.subscribe(() => {})
      await awaitMockCallsCount(data.fetchMock, 2)

      machineConfigurator.selected.value = ["linux-blade", "mac large"]
      const expectedValue1 = `${data.serverUrl}[{"db":"test","table":"test","fields":[],"filters":[{"f":"machine","v":["intellij-linux-hw-blade-test"]}]}]`
      const expectedValue2 = `${data.serverUrl}[{"db":"test","table":"test","fields":[],"filters":[{"f":"machine","v":["intellij-macos-unit-2200-large-test"]}]}]`
      await awaitMockCallsCount(data.fetchMock, 4)
      expect(data.fetchMock.mock.calls[2][0]).toBe(expectedValue1)
      expect(data.fetchMock.mock.calls[3][0]).toBe(expectedValue2)
    })
  })

  describe("selection normalization", () => {
    it("keeps a group selection absent from the current window instead of rewriting it to Unknown", async () => {
      // A window where the selected group has no agents, but unmapped agents exist: the group
      // list carries an "Unknown" bucket, and the lookup resolves the group's display name to
      // "Unknown" (it matches no raw-agent rule). The selection must survive untouched.
      data.fetchMock.mockReturnValue(
        new Observable((sub) => {
          sub.next([...machineGroupsResponse, { group: "Unknown", machines: ["some-new-agent-1"] }])
        })
      )
      const groupLookup = vi.fn<() => Promise<{ json: () => Promise<{ group: string }> }>>(() => Promise.resolve({ json: () => Promise.resolve({ group: "Unknown" }) }))
      vi.stubGlobal("fetch", groupLookup)

      machineConfigurator = new MachineConfigurator(data.serverConfigurator, data.persistenceForDashboard, [], true, ["linux-blade-hetzner"])
      await awaitMockCallsCount(groupLookup, 1)
      // Let the in-flight normalization settle before asserting it left the selection alone.
      await new Promise((resolve) => {
        setTimeout(resolve, 0)
      })
      expect(machineConfigurator.selected.value).toStrictEqual(["linux-blade-hetzner"])
    })

    it("keeps a deliberately selected single agent instead of expanding it to its group", async () => {
      // "intellij-linux-hw-blade-test" is a live leaf of the "linux-blade" group in the loaded
      // list — a user-picked single agent. It must not be expanded to the whole group, and the
      // backend lookup (which would resolve it to "linux-blade") must not even be consulted.
      const groupLookup = vi.fn<() => Promise<{ json: () => Promise<{ group: string }> }>>(() => Promise.resolve({ json: () => Promise.resolve({ group: "linux-blade" }) }))
      vi.stubGlobal("fetch", groupLookup)

      machineConfigurator = new MachineConfigurator(data.serverConfigurator, data.persistenceForDashboard, [], true, ["intellij-linux-hw-blade-test"])
      await awaitCallbackTrue(() => machineConfigurator.values.value.length > 0)
      await new Promise((resolve) => {
        setTimeout(resolve, 0)
      })
      expect(machineConfigurator.selected.value).toStrictEqual(["intellij-linux-hw-blade-test"])
      expect(groupLookup).not.toHaveBeenCalled()
    })
  })

  describe("tests with branch and time filters", () => {
    let branchConfigurator: BranchConfigurator
    beforeEach(() => {
      data.fetchMock.mockReturnValueOnce(
        new Observable((sub) => {
          sub.next(["branch1", "branch2"])
        })
      )
      branchConfigurator = createBranchConfigurator(data.serverConfigurator, data.persistenceForDashboard, [data.timeRangeConfigurator])
      machineConfigurator = new MachineConfigurator(data.serverConfigurator, data.persistenceForDashboard, [data.timeRangeConfigurator, branchConfigurator])
      const measureConfigurator = new TestMeasureConfigurator()
      dataQueryExecutor = new DataQueryExecutor([data.serverConfigurator, machineConfigurator, branchConfigurator, data.timeRangeConfigurator, measureConfigurator])
    })

    it("valid query on configurator init", async () => {
      await awaitMockCallsCount(data.fetchMock, 2)
      const expectedValue = `${data.serverUrl}{"db":"test","table":"test","fields":[{"n":"branch","sql":"distinct branch"}],"filters":[{"f":"generated_time","q":">subtractMonths(now(),1)"}],"order":"branch","flat":true}`
      expect(data.fetchMock.mock.calls[0][0]).toBe(expectedValue)
    })

    it("valid query when select single value for machine configurator", async () => {
      branchConfigurator.selected.value = ["branch1"]
      machineConfigurator.selected.value = ["linux-blade"]
      dataQueryExecutor.subscribe(() => {})
      let expectedValue = `${data.serverUrl}[{"db":"test","table":"test","fields":[],"filters":[{"f":"machine","v":["intellij-linux-hw-blade-test"]},{"f":"branch","v":"branch1"},{"f":"generated_time","q":">subtractMonths(now(),1)"}]}]`
      await awaitMockCallsCount(data.fetchMock, 4)
      expect(data.fetchMock.mock.calls[3][0]).toBe(expectedValue)

      expectedValue = `${data.serverUrl}[{"db":"test","table":"test","fields":[],"filters":[{"f":"machine","v":["intellij-macos-unit-2200-large-test"]},{"f":"branch","v":"branch1"},{"f":"generated_time","q":">subtractMonths(now(),1)"}]}]`
      machineConfigurator.selected.value = ["mac large"]
      await awaitMockCallsCount(data.fetchMock, 5)
      expect(data.fetchMock.mock.calls[4][0]).toBe(expectedValue)
    })

    it("valid query when select single value for branch configurator", async () => {
      branchConfigurator.selected.value = ["branch2"]
      dataQueryExecutor.subscribe(() => {})
      let expectedValue = `${machineGroupsUrl}{"db":"test","table":"test","fields":[{"n":"machine","sql":"distinct machine"}],"filters":[{"f":"generated_time","q":">subtractMonths(now(),1)"},{"f":"branch","q":" like 'branch2%'"}],"order":"machine","flat":true}`
      await awaitMockCallsCount(data.fetchMock, 3)
      expect(data.fetchMock.mock.calls[2][0]).toBe(expectedValue)

      expectedValue = `${machineGroupsUrl}{"db":"test","table":"test","fields":[{"n":"machine","sql":"distinct machine"}],"filters":[{"f":"generated_time","q":">subtractMonths(now(),1)"},{"f":"branch","q":" like 'branch1%'"}],"order":"machine","flat":true}`
      branchConfigurator.selected.value = ["branch1"]
      await awaitMockCallsCount(data.fetchMock, 4)
      expect(data.fetchMock.mock.calls[3][0]).toBe(expectedValue)
    })

    it("valid query when select multiple value for machine configurator", async () => {
      branchConfigurator.selected.value = ["branch1"]
      machineConfigurator.selected.value = ["linux-blade", "mac large"]
      dataQueryExecutor.subscribe(() => {})
      await awaitMockCallsCount(data.fetchMock, 4)
      const expectedMachines = ["intellij-linux-hw-blade-test", "intellij-macos-unit-2200-large-test"]
      const actual = expectedMachines.map((_, index) => data.fetchMock.mock.calls[3 + index][0] as string)
      const expected = expectedMachines.map(
        (machine) =>
          `${data.serverUrl}[{"db":"test","table":"test","fields":[],"filters":[{"f":"machine","v":["${machine}"]},{"f":"branch","v":"branch1"},{"f":"generated_time","q":">subtractMonths(now(),1)"}]}]`
      )
      expect(actual).toStrictEqual(expected)
    })

    it("valid query when select multiple value for branch configurator", async () => {
      branchConfigurator.selected.value = ["bar", "foo"]
      dataQueryExecutor.subscribe(() => {})
      const expected = `${machineGroupsUrl}{"db":"test","table":"test","fields":[{"n":"machine","sql":"distinct machine"}],"filters":[{"f":"generated_time","q":">subtractMonths(now(),1)"},{"f":"branch","q":" = 'bar' or branch = 'foo'"}],"order":"machine","flat":true}`
      await awaitMockCallsCount(data.fetchMock, 3)
      expect(data.fetchMock.mock.calls[2][0]).toBe(expected)
    })
  })
})
