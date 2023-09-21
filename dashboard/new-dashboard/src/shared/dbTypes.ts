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
  return DBType.UNKNOWN
}
