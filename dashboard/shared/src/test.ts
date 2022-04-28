/* eslint-disable */
import { inspect } from "util"
import { generateQueries } from "./DataQueryExecutor"
import { DimensionConfigurator } from "./configurators/DimensionConfigurator"
import { ServerConfigurator } from "./configurators/ServerConfigurator"
import { DataQuery, DataQueryExecutorConfiguration } from "./dataQuery"
import { decode as risonDecode } from "rison-node"

function test(): void {
  const configuration = new DataQueryExecutorConfiguration()
  configuration.measures = ["foo", "bar"]

  const query = new DataQuery()

  const serverConfigurator = new ServerConfigurator("ij", null, "http://server")

  const project = new DimensionConfigurator("project", serverConfigurator, null, true)
  project.selected.value = ["community/indexing", "java/indexing"]

  const branch = new DimensionConfigurator("branch", serverConfigurator, null, true)
  branch.selected.value = ["master", "221"]

  const machine = new DimensionConfigurator("machine", serverConfigurator, null, true)
  machine.selected.value = ["mac", "linux", "win"]

  const configurators = [project, branch, machine]
  for (const configurator of configurators) {
    if (!configurator.configureQuery(query, configuration)) {
      throw new Error("not expected")
    }
  }

  const raw = generateQueries(query, configuration)
  const result = raw.map(it => risonDecode(it)).map((it: any) => {
    return it.filters.map(it => it.value)
  })
  //   .sort((a, b) => {
  //   let r = a[0].localeCompare(b[0])
  //   if (r != 0) {
  //     return r
  //   }
  //
  //   r = a[1].localeCompare(b[1])
  //   if (r != 0) {
  //     return r
  //   }
  //   return a[2].localeCompare(b[2])
  // }).filter(function(item, pos, ary) {
  //         return !pos || !isEqual(item, ary[pos - 1]);
  //     })
  console.log(inspect(result, {depth: 20, colors: false}))
  console.log(result.length)
}

function isEqual(a: Array<number>, b: Array<number>) {
  return a[0] == b[0] && a[1] == b[1] && a[2] == b[2]
}

test()