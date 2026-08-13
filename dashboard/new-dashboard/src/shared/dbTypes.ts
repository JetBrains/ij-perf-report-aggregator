import { defineStore } from "pinia"
import { ref, Ref } from "vue"
import type { ValueUnit } from "../components/common/chart"
import { MeasureUnit, resolveMeasureUnit } from "../components/common/formatter"
import { DBType } from "../components/common/sideBar/InfoSidebar"

export const dbTypeStore = defineStore("dbTypeStore", () => {
  const dbType: Ref<DBType> = ref(DBType.UNKNOWN)

  function setDbType(dbName: string, table: string): void {
    if (dbName == "perfint") {
      dbType.value = DBType.INTELLIJ
    }
    if (dbName == "jbr") {
      dbType.value = DBType.JBR
    }
    if (dbName == "perfintDev" || dbName == "mlEvaluation") {
      dbType.value = DBType.INTELLIJ_DEV
    }
    if (dbName == "fleet" && table == "measure_new") {
      dbType.value = DBType.FLEET_PERF
    }
    if (dbName == "fleet" && table == "report") {
      dbType.value = DBType.FLEET
    }
    if (dbName == "qodana") {
      dbType.value = DBType.QODANA
    }
    if (dbName == "bazel") {
      dbType.value = DBType.BAZEL
    }
    if (dbName == "perfUnitTests") {
      dbType.value = DBType.PERF_UNIT_TESTS
    }
    if (dbName == "ij") {
      dbType.value = DBType.STARTUP_TESTS
    }
    if (dbName == "ijDev") {
      dbType.value = DBType.STARTUP_TESTS_DEV
    }
    if (dbName == "diogen") {
      dbType.value = DBType.DIOGEN
    }
    if (dbName == "toolbox") {
      dbType.value = DBType.TOOLBOX
    }
  }

  function isStartup(): boolean {
    return isIJStartup() || dbType.value == DBType.FLEET
  }

  function isIJStartup(): boolean {
    return dbType.value == DBType.STARTUP_TESTS || dbType.value == DBType.STARTUP_TESTS_DEV
  }

  function isModeSupported(): boolean {
    return (
      dbType.value == DBType.INTELLIJ ||
      dbType.value == DBType.INTELLIJ_DEV ||
      dbType.value == DBType.PERF_UNIT_TESTS ||
      dbType.value == DBType.DIOGEN ||
      dbType.value == DBType.TOOLBOX
    )
  }

  return { dbType, setDbType, isStartup, isIJStartup, isModeSupported }
})

// resolveMeasureUnit for a stored type taken straight from a query row: ignores it for databases
// where it is meaningless — the jbr analyzer hardcodes "c" for every measure, durations included
// (pkg/analyzer/jbrReport.go). Chart and sidebar code must resolve units through this wrapper,
// not resolveMeasureUnit directly, whenever it has a row's stored type at hand.
// TODO: fix the analyzer to write real types, then scope this to pre-fix data (or drop it once old rows age out).
export function resolveMeasureUnitForDb(measureName: string, opts: { storedType?: string; valueUnit?: ValueUnit; scaling?: boolean } = {}): MeasureUnit {
  const storedType = dbTypeStore().dbType == DBType.JBR ? undefined : opts.storedType
  return resolveMeasureUnit(measureName, { ...opts, storedType })
}
