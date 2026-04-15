import { expect, beforeEach, describe, it } from "vitest"
import { DataQueryExecutor } from "../../src/components/common/DataQueryExecutor"
import { BranchConfigurator, createBranchConfigurator } from "../../src/configurators/BranchConfigurator"
import { TestMeasureConfigurator } from "../dummy/TestMeasureConfigurator"
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
      const measureConfigurator = new TestMeasureConfigurator()
      dataQueryExecutor = new DataQueryExecutor([data.serverConfigurator, branchConfigurator, measureConfigurator])
    })

    it("valid query on configurator init", () => {
      expect(data.fetchMock).toHaveBeenCalledTimes(1)
      const expectedValue = `${data.serverUrl}{"db":"test","table":"test","fields":[{"n":"branch","sql":"distinct branch"}],"order":"branch","flat":true}`
      expect(data.fetchMock.mock.calls[0][0]).toBe(expectedValue)
    })

    it("valid query with default selected value", async () => {
      dataQueryExecutor.subscribe(() => {})
      await awaitMockCallsCount(data.fetchMock, 2)
      const expectedValue = `${data.serverUrl}[{"db":"test","table":"test","fields":[],"filters":[{"f":"branch","v":"b1"}]}]`
      expect(data.fetchMock.mock.calls[1][0]).toBe(expectedValue)
    })

    it("valid query when edit selected single value", async () => {
      dataQueryExecutor.subscribe(() => {})
      await awaitMockCallsCount(data.fetchMock, 2)

      let expectedValue = `${data.serverUrl}[{"db":"test","table":"test","fields":[],"filters":[{"f":"branch","v":"b2"}]}]`
      branchConfigurator.selected.value = "b2"
      await awaitMockCallsCount(data.fetchMock, 3)
      expect(data.fetchMock.mock.calls[2][0]).toBe(expectedValue)

      expectedValue = `${data.serverUrl}[{"db":"test","table":"test","fields":[],"filters":[{"f":"branch","v":"b1"}]}]`
      branchConfigurator.selected.value = "b1"
      await awaitMockCallsCount(data.fetchMock, 4)
      expect(data.fetchMock.mock.calls[3][0]).toBe(expectedValue)
    })

    it("valid query when edit selected multiple value", async () => {
      dataQueryExecutor.subscribe(() => {})
      await awaitMockCallsCount(data.fetchMock, 2)
      const values = ["foo", "bar"]

      branchConfigurator.selected.value = values
      await awaitMockCallsCount(data.fetchMock, 4)
      const actual = values.map((_, index) => data.fetchMock.mock.calls[index + 2][0] as string)
      const expected = values.map((value) => `${data.serverUrl}[{"db":"test","table":"test","fields":[],"filters":[{"f":"branch","v":"${value}"}]}]`)
      expect(actual).toStrictEqual(expected)
    })
  })

  describe("tests with time range filter", () => {
    beforeEach(async () => {
      branchConfigurator = createBranchConfigurator(data.serverConfigurator, data.persistenceForDashboard, [data.timeRangeConfigurator])
      await awaitMockCallsCount(data.fetchMock, 1)
      const measureConfigurator = new TestMeasureConfigurator()
      dataQueryExecutor = new DataQueryExecutor([data.serverConfigurator, branchConfigurator, data.timeRangeConfigurator, measureConfigurator])
    })

    it("valid query on configurator init", () => {
      expect(data.fetchMock).toHaveBeenCalledTimes(1)
      const expectedValue = `${data.serverUrl}{"db":"test","table":"test","fields":[{"n":"branch","sql":"distinct branch"}],"filters":[{"f":"generated_time","q":">subtractMonths(now(),1)"}],"order":"branch","flat":true}`
      expect(data.fetchMock.mock.calls[0][0]).toBe(expectedValue)
    })

    it("valid query with default selected value", async () => {
      dataQueryExecutor.subscribe(() => {})
      await awaitMockCallsCount(data.fetchMock, 2)
      const expectedValue = `${data.serverUrl}[{"db":"test","table":"test","fields":[],"filters":[{"f":"branch","v":"b1"},{"f":"generated_time","q":">subtractMonths(now(),1)"}]}]`
      expect(data.fetchMock.mock.calls[1][0]).toBe(expectedValue)
    })

    it("valid query when edit selected single value", async () => {
      dataQueryExecutor.subscribe(() => {})
      await awaitMockCallsCount(data.fetchMock, 2)

      let expectedValue = `${data.serverUrl}[{"db":"test","table":"test","fields":[],"filters":[{"f":"branch","v":"b2"},{"f":"generated_time","q":">subtractMonths(now(),1)"}]}]`
      branchConfigurator.selected.value = "b2"
      await awaitMockCallsCount(data.fetchMock, 3)
      expect(data.fetchMock.mock.calls[2][0]).toBe(expectedValue)

      expectedValue = `${data.serverUrl}[{"db":"test","table":"test","fields":[],"filters":[{"f":"branch","v":"b1"},{"f":"generated_time","q":">subtractMonths(now(),1)"}]}]`
      branchConfigurator.selected.value = "b1"
      await awaitMockCallsCount(data.fetchMock, 4)
      expect(data.fetchMock.mock.calls[3][0]).toBe(expectedValue)
    })
  })
})
