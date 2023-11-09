import { OptionDataValue } from "echarts/types/src/util/types"
import { durationAxisPointerFormatter } from "../components/common/formatter"

/**
 * This class is actually just a storage.
 * It's used in popup and the instance is returned by the chart library is naked (can't contain any method).
 */
export class Delta {
  constructor(
    public prev: number | null,
    public next: number | null
  ) {}

  public static calculateDeltas(values: number[]): Delta[] {
    const deltas: Delta[] = []
    deltas.push(new Delta(null, values[1]))
    for (let i = 1; i < values.length - 1; i++) {
      if (i > 0 && i < values.length + 1) {
        deltas.push(new Delta(values[i - 1], values[i + 1]))
      }
    }
    deltas.push(new Delta(values.at(-2) as number, null))
    return deltas
  }
}

function appendValueAndPercent(element: HTMLElement, value: number, otherValue: number | null, text: string, isMs: boolean, type: string) {
  if (otherValue != null) {
    element.append(document.createElement("br"))
    element.append(getDifferenceString(value, otherValue, text, isMs, type))
  }
}

export function getDifferenceString(value: number, otherValue: number, text: string, isMs: boolean, type: string): string {
  const deltaAbs = value - otherValue
  const deltaAbsFormatted = durationAxisPointerFormatter(isMs ? Math.abs(deltaAbs) : Math.abs(deltaAbs) / 1000 / 1000, type)
  let deltaPercentFormatted = ""
  const plus = deltaAbs > 0 ? "-" : deltaAbs < 0 ? "+" : ""
  if (value != 0) {
    const deltaPercent = Math.abs((deltaAbs / value) * 100)
    deltaPercentFormatted = ` (${plus}${deltaPercent.toFixed(1)}%)`
  }
  return `${text}${plus}${deltaAbsFormatted}${deltaPercentFormatted}`
}

export function appendPopupElement(element: HTMLElement, value: number, delta: Delta, isMs: boolean, type: string): void {
  appendValueAndPercent(element, value, delta.next, "Next: ", isMs, type)
  appendValueAndPercent(element, value, delta.prev, "Previous: ", isMs, type)
}

export function findDeltaInData(data: (OptionDataValue | Delta)[]): Delta | undefined {
  return data.find((obj) => typeof obj === "object" && "prev" in obj && "next" in obj) as Delta | undefined
}
