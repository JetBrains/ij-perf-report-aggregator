import { createPinia, setActivePinia } from "pinia"
import { expect, beforeEach, describe, test } from "vitest"
import { DataQueryExecutor } from "../../src/components/common/DataQueryExecutor"
import { dimensionConfigurator } from "../../src/configurators/DimensionConfigurator"
import { MeasureConfigurator } from "../../src/configurators/MeasureConfigurator"
import { ReleaseNightlyConfigurator } from "../../src/configurators/ReleaseNightlyConfigurator"
import { awaitMockCallsCount } from "../utils/awaitors"
import ConfiguratorTest, { ConfigurationTestData } from "./ConfiguratorTest"

describe("Release Nightly configurator", () => {
  let data: ConfigurationTestData
  let configurator: ReleaseNightlyConfigurator
  beforeEach(() => {
    setActivePinia(createPinia())
    data = ConfiguratorTest.setupPreconditions(["b1", "b2"])
  })

  describe("tests without filters", () => {
    beforeEach(() => {
      configurator = new ReleaseNightlyConfigurator(data.persistenceForDashboard)
    })

    test("Valid initial configuration", () => {
      expect(configurator.values.value).toStrictEqual(["EAP / Release", "Nightly"])
      expect(configurator.selected.value).toBe("Nightly")
    })

    test("Valid filter for other configurators", () => {
      dimensionConfigurator("project", data.serverConfigurator, data.persistenceForDashboard, false, [configurator])
      expect(data.fetchMock).toHaveBeenCalledTimes(1)
      const expectedValue = `${data.serverUrl}{"db":"test","table":"test","fields":[{"n":"project","sql":"distinct project"}],"filters":[{"f":"build_c3","v":0,"o":"="}],"order":"project","flat":true}`
      expect(data.fetchMock.mock.calls[0][0]).toBe(expectedValue)
    })

    test("Valid filter query with Nightly", async () => {
      const measureConfigurator = new MeasureConfigurator(data.serverConfigurator, data.persistenceForDashboard, [], false)
      measureConfigurator.setSelected("b1")
      const dataQueryExecutor = new DataQueryExecutor([data.serverConfigurator, configurator, measureConfigurator])
      dataQueryExecutor.subscribe(() => {})
      await awaitMockCallsCount(data.fetchMock, 2)
      expect(data.fetchMock).toHaveBeenCalledTimes(2)
      const expectedValue = `${data.serverUrl}[{"db":"test","table":"test","fields":[{"n":"t","sql":"toUnixTimestamp(generated_time)*1000"},{"n":"measures","subName":"value"},{"n":"measures","subName":"name"},{"n":"measures","subName":"type"}],"filters":[{"f":"build_c3","v":0,"o":"="},{"f":"measures.name","v":"b1"}],"order":"t"}]`
      expect(data.fetchMock.mock.calls[1][0]).toBe(expectedValue)
    })

    test("Valid filter query with Release", async () => {
      const measureConfigurator = new MeasureConfigurator(data.serverConfigurator, data.persistenceForDashboard, [], false)
      measureConfigurator.setSelected("b1")
      const dataQueryExecutor = new DataQueryExecutor([data.serverConfigurator, configurator, measureConfigurator])
      dataQueryExecutor.subscribe(() => {})
      configurator.selected.value = "EAP / Release"
      await awaitMockCallsCount(data.fetchMock, 2)
      expect(data.fetchMock).toHaveBeenCalledTimes(2)
      const expectedValue = `${data.serverUrl}[{"db":"test","table":"test","fields":[{"n":"t","sql":"toUnixTimestamp(generated_time)*1000"},{"n":"measures","subName":"value"},{"n":"measures","subName":"name"},{"n":"measures","subName":"type"}],"filters":[{"f":"build_c3","v":0,"o":"!="},{"f":"measures.name","v":"b1"}],"order":"t"}]`
      expect(data.fetchMock.mock.calls[1][0]).toBe(expectedValue)
    })
  })
})
