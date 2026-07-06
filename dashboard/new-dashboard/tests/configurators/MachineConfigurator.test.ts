import { Observable } from "rxjs"
import { expect, beforeEach, describe, it } from "vitest"
import { DataQueryExecutor } from "../../src/components/common/DataQueryExecutor"
import { BranchConfigurator, createBranchConfigurator } from "../../src/configurators/BranchConfigurator"
import { getMachineGroupName, MACHINE_GROUP_QODANA_FLEET_HEAVY, MachineConfigurator } from "../../src/configurators/MachineConfigurator"
import { TestMeasureConfigurator } from "../dummy/TestMeasureConfigurator"
import { awaitCallbackTrue, awaitMockCallsCount } from "../utils/awaitors"
import ConfiguratorTest, { ConfigurationTestData } from "./ConfiguratorTest"

describe("Machine configurator", () => {
  let data: ConfigurationTestData
  let machineConfigurator: MachineConfigurator
  let dataQueryExecutor: DataQueryExecutor

  beforeEach(() => {
    data = ConfiguratorTest.setupPreconditions(["intellij-linux-hw-blade-test", "intellij-macos-unit-2200-large-test"])
  })

  describe("tests without filters", () => {
    beforeEach(() => {
      machineConfigurator = new MachineConfigurator(data.serverConfigurator, data.persistenceForDashboard)
      const measureConfigurator = new TestMeasureConfigurator()
      dataQueryExecutor = new DataQueryExecutor([data.serverConfigurator, machineConfigurator, measureConfigurator])
    })

    it("valid query on configurator init", async () => {
      await awaitMockCallsCount(data.fetchMock, 1)
      const expectedValue = `${data.serverUrl}{"db":"test","table":"test","fields":[{"n":"machine","sql":"distinct machine"}],"order":"machine","flat":true}`
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

      machineConfigurator.selected.value = ["intellij-linux-hw-blade-test"]
      let expectedValue = `${data.serverUrl}[{"db":"test","table":"test","fields":[],"filters":[{"f":"machine","v":["intellij-linux-hw-blade-test"]}]}]`
      await awaitMockCallsCount(data.fetchMock, 3)
      expect(data.fetchMock.mock.calls[2][0]).toBe(expectedValue)

      expectedValue = `${data.serverUrl}[{"db":"test","table":"test","fields":[],"filters":[{"f":"machine","v":["intellij-macos-unit-2200-large-test"]}]}]`
      machineConfigurator.selected.value = ["intellij-macos-unit-2200-large-test"]
      await awaitMockCallsCount(data.fetchMock, 4)
      expect(data.fetchMock.mock.calls[3][0]).toBe(expectedValue)
    })

    it("valid query when select multiple value", async () => {
      dataQueryExecutor.subscribe(() => {})
      await awaitMockCallsCount(data.fetchMock, 2)

      machineConfigurator.selected.value = ["intellij-linux-hw-blade-test", "intellij-macos-unit-2200-large-test"]
      const expectedValue1 = `${data.serverUrl}[{"db":"test","table":"test","fields":[],"filters":[{"f":"machine","v":["intellij-linux-hw-blade-test"]}]}]`
      const expectedValue2 = `${data.serverUrl}[{"db":"test","table":"test","fields":[],"filters":[{"f":"machine","v":["intellij-macos-unit-2200-large-test"]}]}]`
      await awaitMockCallsCount(data.fetchMock, 4)
      expect(data.fetchMock.mock.calls[2][0]).toBe(expectedValue1)
      expect(data.fetchMock.mock.calls[3][0]).toBe(expectedValue2)
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
      machineConfigurator.selected.value = ["intellij-linux-hw-blade-test"]
      dataQueryExecutor.subscribe(() => {})
      let expectedValue = `${data.serverUrl}[{"db":"test","table":"test","fields":[],"filters":[{"f":"machine","v":["intellij-linux-hw-blade-test"]},{"f":"branch","v":"branch1"},{"f":"generated_time","q":">subtractMonths(now(),1)"}]}]`
      await awaitMockCallsCount(data.fetchMock, 4)
      expect(data.fetchMock.mock.calls[3][0]).toBe(expectedValue)

      expectedValue = `${data.serverUrl}[{"db":"test","table":"test","fields":[],"filters":[{"f":"machine","v":["intellij-macos-unit-2200-large-test"]},{"f":"branch","v":"branch1"},{"f":"generated_time","q":">subtractMonths(now(),1)"}]}]`
      machineConfigurator.selected.value = ["intellij-macos-unit-2200-large-test"]
      await awaitMockCallsCount(data.fetchMock, 5)
      expect(data.fetchMock.mock.calls[4][0]).toBe(expectedValue)
    })

    it("valid query when select single value for branch configurator", async () => {
      branchConfigurator.selected.value = ["branch2"]
      dataQueryExecutor.subscribe(() => {})
      let expectedValue = `${data.serverUrl}{"db":"test","table":"test","fields":[{"n":"machine","sql":"distinct machine"}],"filters":[{"f":"generated_time","q":">subtractMonths(now(),1)"},{"f":"branch","q":" like 'branch2%'"}],"order":"machine","flat":true}`
      await awaitMockCallsCount(data.fetchMock, 3)
      expect(data.fetchMock.mock.calls[2][0]).toBe(expectedValue)

      expectedValue = `${data.serverUrl}{"db":"test","table":"test","fields":[{"n":"machine","sql":"distinct machine"}],"filters":[{"f":"generated_time","q":">subtractMonths(now(),1)"},{"f":"branch","q":" like 'branch1%'"}],"order":"machine","flat":true}`
      branchConfigurator.selected.value = ["branch1"]
      await awaitMockCallsCount(data.fetchMock, 4)
      expect(data.fetchMock.mock.calls[3][0]).toBe(expectedValue)
    })

    it("valid query when select multiple value for machine configurator", async () => {
      branchConfigurator.selected.value = ["branch1"]
      const values = ["intellij-linux-hw-blade-test", "intellij-macos-unit-2200-large-test"]
      machineConfigurator.selected.value = values
      dataQueryExecutor.subscribe(() => {})
      await awaitMockCallsCount(data.fetchMock, 4)
      const actual = values.map((_, index) => data.fetchMock.mock.calls[3 + index][0] as string)
      const expected = values.map(
        (value) =>
          `${data.serverUrl}[{"db":"test","table":"test","fields":[],"filters":[{"f":"machine","v":["${value}"]},{"f":"branch","v":"branch1"},{"f":"generated_time","q":">subtractMonths(now(),1)"}]}]`
      )
      expect(actual).toStrictEqual(expected)
    })

    it("valid query when select multiple value for branch configurator", async () => {
      branchConfigurator.selected.value = ["bar", "foo"]
      dataQueryExecutor.subscribe(() => {})
      const expected = `${data.serverUrl}{"db":"test","table":"test","fields":[{"n":"machine","sql":"distinct machine"}],"filters":[{"f":"generated_time","q":">subtractMonths(now(),1)"},{"f":"branch","q":" = 'bar' or branch = 'foo'"}],"order":"machine","flat":true}`
      await awaitMockCallsCount(data.fetchMock, 3)
      expect(data.fetchMock.mock.calls[2][0]).toBe(expected)
    })
  })
})

describe("getMachineGroupName - qodana agents", () => {
  const c5xlarge = "Linux EC2 c5.xlarge (4 vCPU, 8 GB)"

  it("maps the new qodana fleet heavy agents to their own group", () => {
    expect(getMachineGroupName("qodana-fleet-linux-amd64-heavy-test-i-00375d8f0e61ed132")).toBe(MACHINE_GROUP_QODANA_FLEET_HEAVY)
  })

  it("keeps legacy qodana heavy agents in the c5.xlarge group", () => {
    expect(getMachineGroupName("qodana-linux-amd64-heavy-test-A-i-0a1b2c3d4e5f60718")).toBe(c5xlarge)
  })

  it("classifies fleet and legacy heavy agents into different groups", () => {
    // Same group would collapse the group's common-prefix machine filter to `qodana-`,
    // over-matching every other qodana machine class (see the query-layer test below).
    expect(getMachineGroupName("qodana-fleet-linux-amd64-heavy-test-i-00375d8f0e61ed132")).not.toBe(getMachineGroupName("qodana-linux-amd64-heavy-test-A-i-0a1b2c3d4e5f60718"))
  })

  it("keeps qodana-aws-cpu-x64 agents in the c5a(d).xlarge group", () => {
    expect(getMachineGroupName("qodana-aws-cpu-x64-i-0123456789abcdef0")).toBe("Linux EC2 c5a(d).xlarge (4 vCPU, 8 GB)")
  })

  it("returns Unknown for unmapped agents", () => {
    expect(getMachineGroupName("some-unmapped-agent-i-000")).toBe("Unknown")
  })
})

// Serialized fragment of a `machine LIKE 'qodana-fleet-linux-amd64-heavy…%'` filter. Its presence proves
// the fleet group is queried by a fleet-only prefix; a merged fleet+heavy group would collapse to `qodana-%`.
const FLEET_MACHINE_LIKE_FRAGMENT = '"f":"machine","v":"qodana-fleet-linux-amd64-heavy'
const COLLAPSED_MACHINE_LIKE_FRAGMENT = '"f":"machine","v":"qodana-%"'

describe("Machine configurator - qodana fleet group query", () => {
  let data: ConfigurationTestData
  let machineConfigurator: MachineConfigurator
  let dataQueryExecutor: DataQueryExecutor

  beforeEach(() => {
    // A legacy heavy agent coexists with the fleet agents so the test proves the fleet group's filter
    // stays fleet-specific instead of widening to also match the heavy class.
    data = ConfiguratorTest.setupPreconditions([
      "qodana-fleet-linux-amd64-heavy-test-i-0000aaaa",
      "qodana-fleet-linux-amd64-heavy-test-i-0000bbbb",
      "qodana-linux-amd64-heavy-test-A-i-0000cccc",
    ])
    machineConfigurator = new MachineConfigurator(data.serverConfigurator, data.persistenceForDashboard)
    const measureConfigurator = new TestMeasureConfigurator()
    dataQueryExecutor = new DataQueryExecutor([data.serverConfigurator, machineConfigurator, measureConfigurator])
  })

  it("filters the fleet group by a fleet-only prefix, never collapsing to qodana-", async () => {
    dataQueryExecutor.subscribe(() => {})
    // Wait until the loaded machine list has been grouped, so selecting the group hits the prefix path.
    await awaitCallbackTrue(() => machineConfigurator.values.value.some((g) => g.value === MACHINE_GROUP_QODANA_FLEET_HEAVY))

    machineConfigurator.selected.value = [MACHINE_GROUP_QODANA_FLEET_HEAVY]
    await awaitCallbackTrue(() => data.fetchMock.mock.calls.some((c) => (c[0] as string).includes(FLEET_MACHINE_LIKE_FRAGMENT)))

    const urls = data.fetchMock.mock.calls.map((c) => c[0] as string)
    expect(urls.find((u) => u.includes(FLEET_MACHINE_LIKE_FRAGMENT))).toBeDefined()
    expect(urls.find((u) => u.includes(COLLAPSED_MACHINE_LIKE_FRAGMENT))).toBeUndefined()
  })
})
