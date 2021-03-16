export interface TooltipLineDescriptor {
  readonly name: string
  readonly value: string
  main?: boolean
  selectable?: boolean
  extraStyle?: string
}

export function buildTooltip(lines: Array<TooltipLineDescriptor>): string {
  let result = ""
  for (const line of lines) {
    if (line.main) {
      result += `<span style="user-select: text">${line.name}</span>`
    }
    else {
      result += `<br/>${line.name}`
    }
    const valueStyleClass = line.selectable ? "tooltipSelectableValue" : (line.main ? "tooltipMainValue" : "tooltipValue")
    result += `<span class="${valueStyleClass}" ${line.extraStyle == null ? "" : `style="${line.extraStyle}"`}">${line.value}</span>`
  }
  return result
}