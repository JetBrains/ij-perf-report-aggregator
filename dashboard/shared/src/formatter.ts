import humanizeDuration, { HumanizerOptions } from "humanize-duration"

export function nsToMs(v: number) {
  return v / 1_000_000
}

// https://github.com/apache/echarts/issues/8294
export const numberFormat = new Intl.NumberFormat(undefined, {
  maximumFractionDigits: 0,
})

const durationFormatOptions: HumanizerOptions = {
  language: "shortEn",
  round: true,
  units: ["y", "mo", "w", "d", "h", "m", "s", "ms"],
  languages: {
    shortEn: {
      y: () => "y",
      mo: () => "mo",
      w: () => "w",
      d: () => "d",
      h: () => "h",
      m: () => "min",
      s: () => "s",
      ms: () => "ms",
    },
  },
}

export const durationAxisPointerFormatter: (valueInMs: number) => string = humanizeDuration.humanizer(durationFormatOptions)

export const timeFormatWithoutSeconds = new Intl.DateTimeFormat(undefined, {
  year: "numeric",
  month: "short",
  day: "numeric",
  hour: "numeric",
  minute: "numeric",
})

export const durationFormatterInOneWord: (valueInMs: number) => string = humanizeDuration.humanizer({
  ...durationFormatOptions,
  largest: 1,
})

export function numberAxisLabelFormatter(value: number): string {
  return numberFormat.format(value)
}

export function isDurationFormatterApplicable(measureName: string): boolean {
  return !(
    measureName.includes("number") || measureName.includes("Number") ||
    measureName.includes("count") || measureName.includes("Count") ||
    measureName.endsWith("_sources")
  )
}

export function getValueFormatterByMeasureName(measureName: string): (valueInMs: number) => string {
  return isDurationFormatterApplicable(measureName) ? durationAxisPointerFormatter : numberAxisLabelFormatter
}