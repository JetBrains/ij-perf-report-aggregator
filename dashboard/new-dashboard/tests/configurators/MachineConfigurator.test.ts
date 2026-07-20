import { Observable } from "rxjs"
import { expect, beforeEach, afterEach, describe, it, vi } from "vitest"
import { DataQueryExecutor } from "../../src/components/common/DataQueryExecutor"
import { BranchConfigurator, createBranchConfigurator } from "../../src/configurators/BranchConfigurator"
import { MachineConfigurator } from "../../src/configurators/MachineConfigurator"
import { TestMeasureConfigurator } from "../dummy/TestMeasureConfigurator"
import { awaitCallbackTrue, awaitMockCallsCount } from "../utils/awaitors"
import ConfiguratorTest, { ConfigurationTestData } from "./ConfiguratorTest"

// The backend returns machines already grouped by hardware class, each group carrying the
// filter predicate rendered from its grouping rule — selecting a group serializes as that
// predicate ({f: "machine", q: ...}), not as the member list.
const bladePredicate = "like 'intellij-linux-hw-blade-%'"
const macLargePredicate = "like 'intellij-macos-unit-2200-large-%'"
const machineGroupsResponse = [
  { group: "linux-blade", machines: ["intellij-linux-hw-blade-test"], predicate: bladePredicate },
  { group: "mac large", machines: ["intellij-macos-unit-2200-large-test"], predicate: macLargePredicate },
]

describe("Machine configurator", () => {
  let data: ConfigurationTestData
  let machineConfigurator: MachineConfigurator
  let dataQueryExecutor: DataQueryExecutor
  let machineGroupsUrl: string

  beforeEach(() => {
    // PersistentStateManager writes selections to localStorage on a 300ms debounce, so a
    // previous test's selection can leak into the next construction — start clean.
    localStorage.clear()
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
      let expectedValue = `${data.serverUrl}[{"db":"test","table":"test","fields":[],"filters":[{"f":"machine","q":"${bladePredicate}"}]}]`
      await awaitMockCallsCount(data.fetchMock, 3)
      expect(data.fetchMock.mock.calls[2][0]).toBe(expectedValue)

      expectedValue = `${data.serverUrl}[{"db":"test","table":"test","fields":[],"filters":[{"f":"machine","q":"${macLargePredicate}"}]}]`
      machineConfigurator.selected.value = ["mac large"]
      await awaitMockCallsCount(data.fetchMock, 4)
      expect(data.fetchMock.mock.calls[3][0]).toBe(expectedValue)
    })

    it("valid query when select multiple value", async () => {
      dataQueryExecutor.subscribe(() => {})
      await awaitMockCallsCount(data.fetchMock, 2)

      machineConfigurator.selected.value = ["linux-blade", "mac large"]
      const expectedValue1 = `${data.serverUrl}[{"db":"test","table":"test","fields":[],"filters":[{"f":"machine","q":"${bladePredicate}"}]}]`
      const expectedValue2 = `${data.serverUrl}[{"db":"test","table":"test","fields":[],"filters":[{"f":"machine","q":"${macLargePredicate}"}]}]`
      await awaitMockCallsCount(data.fetchMock, 4)
      expect(data.fetchMock.mock.calls[2][0]).toBe(expectedValue1)
      expect(data.fetchMock.mock.calls[3][0]).toBe(expectedValue2)
    })

    it("serializes a group without a predicate (the Unknown bucket) by its member list", async () => {
      // Serve the group list for every URL — chart-query responses are unused by assertions.
      data.fetchMock.mockReturnValue(
        new Observable((sub) => {
          sub.next([...machineGroupsResponse, { group: "Unknown", machines: ["zeta-agent-1", "alpha-agent-2"] }])
        })
      )
      machineConfigurator = new MachineConfigurator(data.serverConfigurator, data.persistenceForDashboard)
      dataQueryExecutor = new DataQueryExecutor([data.serverConfigurator, machineConfigurator, new TestMeasureConfigurator()])
      dataQueryExecutor.subscribe(() => {})
      await awaitMockCallsCount(data.fetchMock, 2)

      machineConfigurator.selected.value = ["Unknown"]
      await awaitMockCallsCount(data.fetchMock, 3)
      const expectedValue = `${data.serverUrl}[{"db":"test","table":"test","fields":[],"filters":[{"f":"machine","v":["alpha-agent-2","zeta-agent-1"]}]}]`
      expect(data.fetchMock.mock.calls[2][0]).toBe(expectedValue)
    })

    it("passes a multi-prefix group predicate through verbatim", async () => {
      const hetznerPredicate = "like 'intellij-linux-hw-hetzner%' or machine like 'intellij-linux-agg-hw-hetzner-agent%'"
      // Serve the group list for every URL — chart-query responses are unused by assertions.
      data.fetchMock.mockReturnValue(
        new Observable((sub) => {
          sub.next([{ group: "linux-blade-hetzner", machines: ["intellij-linux-agg-hw-hetzner-agent-1", "intellij-linux-hw-hetzner-agent-2"], predicate: hetznerPredicate }])
        })
      )
      machineConfigurator = new MachineConfigurator(data.serverConfigurator, data.persistenceForDashboard)
      dataQueryExecutor = new DataQueryExecutor([data.serverConfigurator, machineConfigurator, new TestMeasureConfigurator()])
      dataQueryExecutor.subscribe(() => {})
      await awaitMockCallsCount(data.fetchMock, 2)

      machineConfigurator.selected.value = ["linux-blade-hetzner"]
      await awaitMockCallsCount(data.fetchMock, 3)
      const expectedValue = `${data.serverUrl}[{"db":"test","table":"test","fields":[],"filters":[{"f":"machine","q":"${hetznerPredicate}"}]}]`
      expect(data.fetchMock.mock.calls[2][0]).toBe(expectedValue)
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
      let expectedValue = `${data.serverUrl}[{"db":"test","table":"test","fields":[],"filters":[{"f":"machine","q":"${bladePredicate}"},{"f":"branch","v":"branch1"},{"f":"generated_time","q":">subtractMonths(now(),1)"}]}]`
      await awaitMockCallsCount(data.fetchMock, 4)
      expect(data.fetchMock.mock.calls[3][0]).toBe(expectedValue)

      expectedValue = `${data.serverUrl}[{"db":"test","table":"test","fields":[],"filters":[{"f":"machine","q":"${macLargePredicate}"},{"f":"branch","v":"branch1"},{"f":"generated_time","q":">subtractMonths(now(),1)"}]}]`
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
      // Leftover async emissions from earlier tests can inject stray calls into the shared
      // fetch spy, so assert on content rather than call indices.
      const expected = [bladePredicate, macLargePredicate].map(
        (predicate) =>
          `${data.serverUrl}[{"db":"test","table":"test","fields":[],"filters":[{"f":"machine","q":"${predicate}"},{"f":"branch","v":"branch1"},{"f":"generated_time","q":">subtractMonths(now(),1)"}]}]`
      )
      await awaitCallbackTrue(() => {
        const urls = new Set(data.fetchMock.mock.calls.map((call) => call[0] as string))
        return expected.every((url) => urls.has(url))
      })
      const urls = data.fetchMock.mock.calls.map((call) => call[0] as string)
      expect(urls).toContain(expected[0])
      expect(urls).toContain(expected[1])
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
