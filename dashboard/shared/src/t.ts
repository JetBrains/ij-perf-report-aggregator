/* eslint-disable */
import { encodeRison, makeUrlSafe } from "./rison"
import { crush } from "./json-crush"

const c = require("../j")

function test(_v: any) {
  const r = "(fields:!((n:product,sql:'distinct product')),order:product,table:report,flat:!t,db:ij)"
  console.log(r)
  console.log(r.length)
  console.log("\n\n")

  const rr = crush(r)
  console.log("our encode " + rr.length)
  console.log(rr)

  const r2 = makeUrlSafe(c.crush(r))
  console.log("original encode " + r2.length)
  console.log(r2)

  const r3 = c.uncrush(decodeURIComponent(r2))
  console.log(r3.length)
  console.log(r3)
}

test([{
  "db": "sharedIndexes",
  "fields": [{"n": "t", "sql": "toUnixTimestamp(generated_time) * 1000"}, {
    "n": "measures",
    "subName": "value",
  }, "tc_installer_build_id", "tc_build_id", "build_c1", "build_c2", "build_c3", "machine"],
  "filters": [{"f": "project", "v": "intellij_sources/indexing"}, {"f": "branch", "v": "master"}, {
    "f": "machine",
    "v": ["intellij-linux-hw-blade-023", "intellij-linux-hw-blade-024", "intellij-linux-hw-blade-025", "intellij-linux-hw-blade-026", "intellij-linux-hw-blade-027", "intellij-linux-hw-blade-028", "intellij-linux-hw-blade-029", "intellij-linux-hw-blade-030", "intellij-linux-hw-blade-031", "intellij-linux-hw-blade-032", "intellij-linux-hw-blade-033", "intellij-linux-hw-blade-034", "intellij-linux-hw-blade-035", "intellij-linux-hw-blade-036", "intellij-linux-hw-blade-037", "intellij-linux-hw-blade-038", "intellij-linux-hw-blade-039", "intellij-linux-hw-blade-040", "intellij-linux-hw-blade-041", "intellij-linux-hw-blade-042", "intellij-linux-hw-blade-043", "intellij-linux-hw-blade-044", "intellij-linux-hw-blade-045", "intellij-linux-hw-blade-046", "intellij-linux-hw-blade-047", "intellij-linux-hw-blade-048", "intellij-linux-hw-blade-049"],
  }, {"f": "generated_time", "q": "> subtractMonths(now(), 3)"}, {"f": "measures.name", "v": "indexing"}, {"f": "measures.value", "o": "!=", "v": 0}],
  "order": "t",
}, {
  "db": "sharedIndexes",
  "fields": [{"n": "t", "sql": "toUnixTimestamp(generated_time) * 1000"}, {
    "n": "measures",
    "subName": "value",
  }, "tc_installer_build_id", "tc_build_id", "build_c1", "build_c2", "build_c3", "machine"],
  "filters": [{"f": "project", "v": "community/indexing"}, {"f": "branch", "v": "master"}, {
    "f": "machine",
    "v": ["intellij-linux-hw-blade-023", "intellij-linux-hw-blade-024", "intellij-linux-hw-blade-025", "intellij-linux-hw-blade-026", "intellij-linux-hw-blade-027", "intellij-linux-hw-blade-028", "intellij-linux-hw-blade-029", "intellij-linux-hw-blade-030", "intellij-linux-hw-blade-031", "intellij-linux-hw-blade-032", "intellij-linux-hw-blade-033", "intellij-linux-hw-blade-034", "intellij-linux-hw-blade-035", "intellij-linux-hw-blade-036", "intellij-linux-hw-blade-037", "intellij-linux-hw-blade-038", "intellij-linux-hw-blade-039", "intellij-linux-hw-blade-040", "intellij-linux-hw-blade-041", "intellij-linux-hw-blade-042", "intellij-linux-hw-blade-043", "intellij-linux-hw-blade-044", "intellij-linux-hw-blade-045", "intellij-linux-hw-blade-046", "intellij-linux-hw-blade-047", "intellij-linux-hw-blade-048", "intellij-linux-hw-blade-049"],
  }, {"f": "generated_time", "q": "> subtractMonths(now(), 3)"}, {"f": "measures.name", "v": "indexing"}, {"f": "measures.value", "o": "!=", "v": 0}],
  "order": "t",
}])