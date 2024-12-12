import { parseDuration, TimeRangeConfigurator } from "../../configurators/TimeRangeConfigurator"

function getDateAgoByDuration(s: string): Date {
  const result = parseDuration(s)
  let days = 0
  if (result.days != null) {
    days += result.days
  }
  if (result.months != null) {
    days += result.months * 31
  }
  if (result.weeks != null) {
    days += result.weeks * 7
  }
  if (result.years != null) {
    days += result.years * 365
  }
  const date = new Date()
  date.setDate(date.getDate() - days)
  return date
}

export function getPersistentLink(url: string, timerangeConfigurator: TimeRangeConfigurator): string {
  url = url
    .replace(/&+$/, "")
    .replace(/([?&])customRange=[^&]*&?/, "$1")
    .replace(/([?&])timeRange=[^&]*&?/, "$1")

  if (timerangeConfigurator.value.value != "custom") {
    const now = new Date()
    const ago = getDateAgoByDuration(timerangeConfigurator.value.value)
    const dayFrom = ago.getDate() >= 2 ? ago.getDate() - 1 : ago.getDate()
    const dayTo = now.getDate() < 31 ? now.getDate() + 1 : now.getDate()
    const filter = `${ago.getFullYear()}-${ago.getUTCMonth() + 1}-${dayFrom}:${now.getFullYear()}-${now.getUTCMonth() + 1}-${dayTo}`
    url = url + "&timeRange=custom&customRange=" + filter
  } else {
    url = url + "&timeRange=custom&customRange=" + timerangeConfigurator.customRange.value
  }
  return url.replace(/&+/g, "&")
}
