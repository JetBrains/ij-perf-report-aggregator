import { InfoData } from "./InfoSidebar"

export function generateDefaultReason(data: InfoData): string {
  const date = data.date ? new Date(data.date).toISOString().slice(0, 10).replaceAll("-", ".") : ""
  const metricName = data.series[0]?.metricName
  return `${date} ${data.projectName ?? ""}${metricName ? ` - ${metricName}` : ""}`
}
