import { DBType } from "../components/common/sideBar/InfoSidebar"

export function getDBType(dbName: string, table: string): DBType {
  if (dbName == "perfint") {
    return DBType.INTELLIJ
  }
  if (dbName == "jbr") {
    return DBType.JBR
  }
  if (dbName == "perfintDev") {
    return DBType.INTELLIJ_DEV
  }
  if (dbName == "fleet" && table == "measure") {
    return DBType.DEV_FLEET
  }
  if (dbName == "fleet" && table == "report") {
    return DBType.FLEET
  }
  if (dbName == "qodana") {
    return DBType.QODANA
  }
  if (dbName == "bazel") {
    return DBType.BAZEL
  }
  if (dbName == "perfUnitTests") {
    return DBType.PERF_UNIT_TESTS
  }
  if (dbName == "ij") {
    return DBType.STARTUP_TESTS
  }
  if (dbName == "ijDev") {
    return DBType.STARTUP_TESTS_DEV
  }
  return DBType.UNKNOWN
}

export function isStartup(db: string | undefined, table: string | undefined): boolean {
  if (db == undefined || table == undefined) {
    throw new Error("db or table is undefined")
  }
  return isIJStartup(db, table) || getDBType(db, table) == DBType.FLEET
}

export function isIJStartup(db: string | undefined, table: string | undefined): boolean {
  if (db == undefined || table == undefined) {
    throw new Error("db or table is undefined")
  }
  return getDBType(db, table) == DBType.STARTUP_TESTS || getDBType(db, table) == DBType.STARTUP_TESTS_DEV
}
