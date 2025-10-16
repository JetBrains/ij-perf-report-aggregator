import { Observable } from "rxjs"
import { assert, beforeEach, describe, test } from "vitest"
import { DataQueryExecutor } from "../../src/components/common/DataQueryExecutor"
import { BranchConfigurator, createBranchConfigurator } from "../../src/configurators/BranchConfigurator"
import { MachineConfigurator } from "../../src/configurators/MachineConfigurator"
import { awaitMockCallsCount } from "../utils/awaitors"
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
      dataQueryExecutor = new DataQueryExecutor([data.serverConfigurator, machineConfigurator])
    })

    test("Valid query on configurator init", async () => {
      await awaitMockCallsCount(data.fetchMock, 1)
      const expectedValue = `${data.serverUrl}{"db":"test","table":"test","fields":[{"n":"machine","sql":"distinct machine"}],"order":"machine","flat":true}`
      assert.equal(data.fetchMock.mock.calls[0][0], expectedValue)
    })

    test("Valid query with default selected value", async () => {
      dataQueryExecutor.subscribe(() => {})
      await awaitMockCallsCount(data.fetchMock, 2)
      const expectedValue = `${data.serverUrl}[{"db":"test","table":"test","fields":[],"filters":[{"f":"machine","v":["machine"]}]}]`
      assert.equal(data.fetchMock.mock.calls[1][0], expectedValue)
    })

    test("Valid query when select single value", async () => {
      dataQueryExecutor.subscribe(() => {})
      await awaitMockCallsCount(data.fetchMock, 2)

      machineConfigurator.selected.value = ["intellij-linux-hw-blade-test"]
      let expectedValue = `${data.serverUrl}[{"db":"test","table":"test","fields":[],"filters":[{"f":"machine","v":["intellij-linux-hw-blade-test"]}]}]`
      await awaitMockCallsCount(data.fetchMock, 3)
      assert.equal(data.fetchMock.mock.calls[2][0], expectedValue)

      expectedValue = `${data.serverUrl}[{"db":"test","table":"test","fields":[],"filters":[{"f":"machine","v":["intellij-macos-unit-2200-large-test"]}]}]`
      machineConfigurator.selected.value = ["intellij-macos-unit-2200-large-test"]
      await awaitMockCallsCount(data.fetchMock, 4)
      assert.equal(data.fetchMock.mock.calls[3][0], expectedValue)
    })

    test("Valid query when select multiple value", async () => {
      console.log("Valid query when select multiple value start")
      dataQueryExecutor.subscribe(() => {})
      await awaitMockCallsCount(data.fetchMock, 2)

      machineConfigurator.selected.value = ["intellij-linux-hw-blade-test", "intellij-macos-unit-2200-large-test"]
      const expectedValue1 = `${data.serverUrl}[{"db":"test","table":"test","fields":[],"filters":[{"f":"machine","v":["intellij-linux-hw-blade-test"]}]}]`
      const expectedValue2 = `${data.serverUrl}[{"db":"test","table":"test","fields":[],"filters":[{"f":"machine","v":["intellij-macos-unit-2200-large-test"]}]}]`
      await awaitMockCallsCount(data.fetchMock, 4)
      assert.equal(data.fetchMock.mock.calls[2][0], expectedValue1)
      assert.equal(data.fetchMock.mock.calls[3][0], expectedValue2)
      console.log("Valid query when select multiple value end")
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
      dataQueryExecutor = new DataQueryExecutor([data.serverConfigurator, machineConfigurator, branchConfigurator, data.timeRangeConfigurator])
    })

    test("Valid query on configurator init", async () => {
      await awaitMockCallsCount(data.fetchMock, 2)
      const expectedValue = `${data.serverUrl}{"db":"test","table":"test","fields":[{"n":"branch","sql":"distinct branch"}],"filters":[{"f":"generated_time","q":">subtractMonths(now(),1)"}],"order":"branch","flat":true}`
      assert.equal(data.fetchMock.mock.calls[0][0], expectedValue)
    })

    test("Valid query when select single value for machine configurator", async () => {
      branchConfigurator.selected.value = ["branch1"]
      machineConfigurator.selected.value = ["intellij-linux-hw-blade-test"]
      dataQueryExecutor.subscribe(() => {})
      let expectedValue = `${data.serverUrl}[{"db":"test","table":"test","fields":[],"filters":[{"f":"machine","v":["intellij-linux-hw-blade-test"]},{"f":"branch","v":"branch1"},{"f":"generated_time","q":">subtractMonths(now(),1)"}]}]`
      await awaitMockCallsCount(data.fetchMock, 4)
      assert.equal(data.fetchMock.mock.calls[3][0], expectedValue)

      expectedValue = `${data.serverUrl}[{"db":"test","table":"test","fields":[],"filters":[{"f":"machine","v":["intellij-macos-unit-2200-large-test"]},{"f":"branch","v":"branch1"},{"f":"generated_time","q":">subtractMonths(now(),1)"}]}]`
      machineConfigurator.selected.value = ["intellij-macos-unit-2200-large-test"]
      await awaitMockCallsCount(data.fetchMock, 5)
      assert.equal(data.fetchMock.mock.calls[4][0], expectedValue)
    })

    test("Valid query when select single value for branch configurator", async () => {
      branchConfigurator.selected.value = ["branch2"]
      dataQueryExecutor.subscribe(() => {})
      let expectedValue = `${data.serverUrl}{"db":"test","table":"test","fields":[{"n":"machine","sql":"distinct machine"}],"filters":[{"f":"generated_time","q":">subtractMonths(now(),1)"},{"f":"branch","q":" like 'branch2%'"}],"order":"machine","flat":true}`
      await awaitMockCallsCount(data.fetchMock, 3)
      assert.equal(data.fetchMock.mock.calls[2][0], expectedValue)

      expectedValue = `${data.serverUrl}{"db":"test","table":"test","fields":[{"n":"machine","sql":"distinct machine"}],"filters":[{"f":"generated_time","q":">subtractMonths(now(),1)"},{"f":"branch","q":" like 'branch1%'"}],"order":"machine","flat":true}`
      branchConfigurator.selected.value = ["branch1"]
      await awaitMockCallsCount(data.fetchMock, 4)
      assert.equal(data.fetchMock.mock.calls[3][0], expectedValue)
    })

    test("Valid query when select multiple value for machine configurator", async () => {
      branchConfigurator.selected.value = ["branch1"]
      const values = ["intellij-linux-hw-blade-test", "intellij-macos-unit-2200-large-test"]
      machineConfigurator.selected.value = values
      dataQueryExecutor.subscribe(() => {})
      await awaitMockCallsCount(data.fetchMock, 4)
      for (const [index, value] of values.entries()) {
        const expected = `${data.serverUrl}[{"db":"test","table":"test","fields":[],"filters":[{"f":"machine","v":["${value}"]},{"f":"branch","v":"branch1"},{"f":"generated_time","q":">subtractMonths(now(),1)"}]}]`
        assert.equal(data.fetchMock.mock.calls[3 + index][0], expected)
      }
    })

    test("Valid query when select multiple value for branch configurator", async () => {
      branchConfigurator.selected.value = ["bar", "foo"]
      dataQueryExecutor.subscribe(() => {})
      const expected = `${data.serverUrl}{"db":"test","table":"test","fields":[{"n":"machine","sql":"distinct machine"}],"filters":[{"f":"generated_time","q":">subtractMonths(now(),1)"},{"f":"branch","q":" = 'bar' or branch = 'foo'"}],"order":"machine","flat":true}`
      await awaitMockCallsCount(data.fetchMock, 3)
      assert.equal(data.fetchMock.mock.calls[2][0], expected)
    })
  })
})
