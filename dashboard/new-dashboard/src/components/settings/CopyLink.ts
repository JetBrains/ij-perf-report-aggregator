import { formatCustomRange, parseDuration, TimeRangeConfigurator } from "../../configurators/TimeRangeConfigurator"

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
    const from = getDateAgoByDuration(timerangeConfigurator.value.value)
    from.setDate(from.getDate() - 1)
    const to = new Date()
    to.setDate(to.getDate() + 1)
    url = url + "&timeRange=custom&customRange=" + formatCustomRange(from, to)
  } else {
    url = url + "&timeRange=custom&customRange=" + timerangeConfigurator.customRange.value
  }
  return url.replaceAll(/&+/g, "&")
}
