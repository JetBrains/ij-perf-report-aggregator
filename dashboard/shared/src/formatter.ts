import humanizeDuration, { HumanizerOptions } from "humanize-duration"

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
      m: () => "m",
      s: () => "s",
      ms: () => "ms",
    },
  },
}

export const durationAxisPointerFormatter: (valueInMs: number) => string = humanizeDuration.humanizer(durationFormatOptions)

export const durationAxisLabelFormatter: (valueInMs: number) => string = humanizeDuration.humanizer({
  ...durationFormatOptions,
  delimiter: " "
})

export function numberAxisLabelFormatter(value: number): string {
  return numberFormat.format(value)
}

export function isDurationFormatterApplicable(measureName: string): boolean {
  return !measureName.startsWith("numberOf") && !measureName.endsWith("Count_i")
}

export function getValueFormatterByMeasureName(measureName: string): (valueInMs: number) => string {
  return isDurationFormatterApplicable(measureName) ? durationAxisLabelFormatter : numberAxisLabelFormatter
}