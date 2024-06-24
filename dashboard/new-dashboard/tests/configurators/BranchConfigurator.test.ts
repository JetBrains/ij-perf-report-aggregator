import { assert, beforeEach, describe, test } from "vitest"
import { DataQueryExecutor } from "../../src/components/common/DataQueryExecutor"
import { BranchConfigurator, createBranchConfigurator } from "../../src/configurators/BranchConfigurator"
import { awaitMockCallsCount } from "../utils/awaitors"
import ConfiguratorTest, { ConfigurationTestData } from "./ConfiguratorTest"

describe("Branch configurator", () => {
  let data: ConfigurationTestData
  let branchConfigurator: BranchConfigurator
  let dataQueryExecutor: DataQueryExecutor
  beforeEach(() => {
    data = ConfiguratorTest.setupPreconditions(["b1", "b2"])
  })

  describe("tests without filters", () => {
    beforeEach(async () => {
      branchConfigurator = createBranchConfigurator(data.serverConfigurator, data.persistenceForDashboard)
      await awaitMockCallsCount(data.fetchMock, 1)
      dataQueryExecutor = new DataQueryExecutor([data.serverConfigurator, branchConfigurator])
    })

    test("Valid query on configurator init", () => {
      assert.equal(data.fetchMock.mock.calls.length, 1)
      const expectedValue = `${data.serverUrl}{"db":"test","table":"test","fields":[{"n":"branch","sql":"distinct branch"}],"order":"branch","flat":true}`
      assert.equal(data.fetchMock.mock.calls[0][0], expectedValue)
    })

    test("Valid query with default selected value", async () => {
      dataQueryExecutor.subscribe(() => {})
      await awaitMockCallsCount(data.fetchMock, 2)
      const expectedValue = `${data.serverUrl}[{"db":"test","table":"test","fields":[],"filters":[{"f":"branch","v":"b1"}]}]`
      assert.equal(data.fetchMock.mock.calls[1][0], expectedValue)
    })

    test("Valid query when edit selected single value", async () => {
      dataQueryExecutor.subscribe(() => {})
      await awaitMockCallsCount(data.fetchMock, 2)

      let expectedValue = `${data.serverUrl}[{"db":"test","table":"test","fields":[],"filters":[{"f":"branch","v":"b2"}]}]`
      branchConfigurator.selected.value = "b2"
      await awaitMockCallsCount(data.fetchMock, 3)
      assert.equal(data.fetchMock.mock.calls[2][0], expectedValue)

      expectedValue = `${data.serverUrl}[{"db":"test","table":"test","fields":[],"filters":[{"f":"branch","v":"b1"}]}]`
      branchConfigurator.selected.value = "b1"
      await awaitMockCallsCount(data.fetchMock, 4)
      assert.equal(data.fetchMock.mock.calls[3][0], expectedValue)
    })

    test("Valid query when edit selected multiple value", async () => {
      dataQueryExecutor.subscribe(() => {})
      await awaitMockCallsCount(data.fetchMock, 2)
      const values = ["foo", "bar"]

      branchConfigurator.selected.value = values
      await awaitMockCallsCount(data.fetchMock, 4)
      for (const [index, value] of values.entries()) {
        const expectedValue = `${data.serverUrl}[{"db":"test","table":"test","fields":[],"filters":[{"f":"branch","v":"${value}"}]}]`
        assert.equal(data.fetchMock.mock.calls[index + 2][0], expectedValue)
      }
    })
  })

  describe("tests with time range filter", () => {
    beforeEach(async () => {
      branchConfigurator = createBranchConfigurator(data.serverConfigurator, data.persistenceForDashboard, [data.timeRangeConfigurator])
      await awaitMockCallsCount(data.fetchMock, 1)
      dataQueryExecutor = new DataQueryExecutor([data.serverConfigurator, branchConfigurator, data.timeRangeConfigurator])
    })

    test("Valid query on configurator init", () => {
      assert.equal(data.fetchMock.mock.calls.length, 1)
      const expectedValue = `${data.serverUrl}{"db":"test","table":"test","fields":[{"n":"branch","sql":"distinct branch"}],"filters":[{"f":"generated_time","q":">subtractMonths(now(),1)"}],"order":"branch","flat":true}`
      assert.equal(data.fetchMock.mock.calls[0][0], expectedValue)
    })

    test("Valid query with default selected value", async () => {
      dataQueryExecutor.subscribe(() => {})
      await awaitMockCallsCount(data.fetchMock, 2)
      const expectedValue = `${data.serverUrl}[{"db":"test","table":"test","fields":[],"filters":[{"f":"branch","v":"b1"},{"f":"generated_time","q":">subtractMonths(now(),1)"}]}]`
      assert.equal(data.fetchMock.mock.calls[1][0], expectedValue)
    })

    test("Valid query when edit selected single value", async () => {
      dataQueryExecutor.subscribe(() => {})
      await awaitMockCallsCount(data.fetchMock, 2)

      let expectedValue = `${data.serverUrl}[{"db":"test","table":"test","fields":[],"filters":[{"f":"branch","v":"b2"},{"f":"generated_time","q":">subtractMonths(now(),1)"}]}]`
      branchConfigurator.selected.value = "b2"
      await awaitMockCallsCount(data.fetchMock, 3)
      assert.equal(data.fetchMock.mock.calls[2][0], expectedValue)

      expectedValue = `${data.serverUrl}[{"db":"test","table":"test","fields":[],"filters":[{"f":"branch","v":"b1"},{"f":"generated_time","q":">subtractMonths(now(),1)"}]}]`
      branchConfigurator.selected.value = "b1"
      await awaitMockCallsCount(data.fetchMock, 4)
      assert.equal(data.fetchMock.mock.calls[3][0], expectedValue)
    })

    test.skip("Valid query when edit selected multiple value", async () => {
      dataQueryExecutor.subscribe(() => {})
      console.log(data.fetchMock.mock.calls)
      await awaitMockCallsCount(data.fetchMock, 2)
      const values = ["foo", "bar"]

      branchConfigurator.selected.value = values
      await awaitMockCallsCount(data.fetchMock, 4)
      for (const [index, value] of values.entries()) {
        const expectedValue = `${data.serverUrl}[{"db":"test","table":"test","fields":[],"filters":[{"f":"branch","v":"${value}"},{"f":"generated_time","q":">subtractMonths(now(),1)"}]}]`
        assert.equal(data.fetchMock.mock.calls[index + 2][0], expectedValue)
      }
    })
  })
})
