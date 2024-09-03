import { defineStore } from "pinia"
import { ref, Ref } from "vue"
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
      dbType.value = DBType.INTELLIJ_DEV
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
  }

  function isStartup(): boolean {
    return isIJStartup() || dbType.value == DBType.FLEET
  }

  function isIJStartup(): boolean {
    return dbType.value == DBType.STARTUP_TESTS || dbType.value == DBType.STARTUP_TESTS_DEV
  }

  return { dbType, setDbType, isStartup, isIJStartup }
})
