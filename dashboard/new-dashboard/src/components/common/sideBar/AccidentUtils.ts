import { AccidentKind } from "../../../configurators/accidents/AccidentsConfigurator"
import { InfoData } from "./InfoSidebar"

export function generateDefaultReason(data: InfoData): string {
  const date = data.date ? new Date(data.date).toISOString().slice(0, 10).replaceAll("-", ".") : ""
  const metricName = data.series[0]?.metricName
  return `${date} ${data.projectName ?? ""}${metricName ? ` - ${metricName}` : ""}`
}

export function inferKindFromData(data: InfoData | null): AccidentKind {
  const value = data?.series[0]?.rawValue
  const prev = data?.previousValue
  if (value == null || prev == null || !Number.isFinite(value) || !Number.isFinite(prev)) {
    return AccidentKind.Regression
  }
  return value < prev ? AccidentKind.Improvement : AccidentKind.Regression
}
