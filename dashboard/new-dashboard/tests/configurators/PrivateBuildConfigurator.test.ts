import { assert, beforeEach, describe, test } from "vitest"
import { DataQueryExecutor } from "../../src/components/common/DataQueryExecutor"
import { MeasureConfigurator } from "../../src/configurators/MeasureConfigurator"
import { PrivateBuildConfigurator, privateBuildConfigurator } from "../../src/configurators/PrivateBuildConfigurator"
import { awaitMockCallsCount } from "../utils/awaitors"
import ConfiguratorTest, { ConfigurationTestData } from "./ConfiguratorTest"

describe("Private build configurator", () => {
  let data: ConfigurationTestData
  let configurator: PrivateBuildConfigurator
  let dataQueryExecutor: DataQueryExecutor
  beforeEach(() => {
    data = ConfiguratorTest.setupPreconditions(["b1", "b2"])
  })

  describe("tests without filters", () => {
    beforeEach(async () => {
      configurator = privateBuildConfigurator(data.serverConfigurator, data.persistenceForDashboard)
      await awaitMockCallsCount(data.fetchMock, 1)
      const measureConfigurator = new MeasureConfigurator(data.serverConfigurator, data.persistenceForDashboard, [], false)
      measureConfigurator.setSelected("b1")
      await awaitMockCallsCount(data.fetchMock, 2)

      dataQueryExecutor = new DataQueryExecutor([data.serverConfigurator, configurator, measureConfigurator])
    })

    test("Valid query on configurator init", () => {
      assert.equal(data.fetchMock.mock.calls.length, 2)
      const expectedValue = `${data.serverUrl}{"db":"test","table":"test","fields":[{"n":"triggeredBy","sql":"distinct triggeredBy"}],"order":"triggeredBy","flat":true}`
      assert.equal(data.fetchMock.mock.calls[0][0], expectedValue)
    })

    test("Valid query with default selected value", async () => {
      await new Promise<void>((resolve) => {
        dataQueryExecutor.subscribe(() => {
          resolve()
        })
      })
      const expectedValue = `${data.serverUrl}[{"db":"test","table":"test","fields":[{"n":"t","sql":"toUnixTimestamp(generated_time)*1000"},{"n":"measures","subName":"value"}],"filters":[{"f":"triggeredBy","v":""},{"f":"measures.name","v":"b1"}],"order":"t"}]`
      assert.equal(data.fetchMock.mock.calls[2][0], expectedValue)
    })

    test("Valid query when edit selected single value", async () => {
      configurator.selected.value = "b2"
      await new Promise<void>((resolve) => {
        dataQueryExecutor.subscribe(() => {
          resolve()
        })
      })
      assert.equal(
        data.fetchMock.mock.calls[2][0],
        `${data.serverUrl}[{"db":"test","table":"test","fields":[{"n":"t","sql":"toUnixTimestamp(generated_time)*1000"},{"n":"measures","subName":"value"}],"filters":[{"f":"triggeredBy","v":""},{"f":"measures.name","v":"b1"}],"order":"t"}]`
      )
      assert.equal(
        data.fetchMock.mock.calls[3][0],
        `${data.serverUrl}[{"db":"test","table":"test","fields":[{"n":"t","sql":"toUnixTimestamp(generated_time)*1000"},{"n":"measures","subName":"value"}],"filters":[{"f":"triggeredBy","v":"b2"},{"f":"measures.name","v":"b1"}],"order":"t"}]`
      )
    })
  })
})
