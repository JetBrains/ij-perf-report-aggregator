export interface TooltipLineDescriptor {
  readonly name: string
  readonly value: string
  main?: boolean
  selectable?: boolean
  extraStyle?: string
}

export function buildTooltip(lines: TooltipLineDescriptor[]): string {
  let result = ""
  for (const line of lines) {
    result += line.main ? `<span style="user-select: text">${line.name}</span>` : `<br/>${line.name}`
    const valueStyleClass = line.selectable ? "tooltipSelectableValue" : line.main ? "tooltipMainValue" : "tooltipValue"
    result += `<span class="${valueStyleClass}" ${line.extraStyle == null ? "" : `style="${line.extraStyle}"`}">${line.value}</span>`
  }
  return result
}
