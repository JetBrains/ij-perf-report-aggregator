function getIcon(pathString: string): Element {
  const svg = document.createElementNS("http://www.w3.org/2000/svg", "svg")
  svg.setAttribute("xmlns", "http://www.w3.org/2000/svg")
  svg.setAttribute("fill", "none")
  svg.setAttribute("viewBox", "0 0 24 24")
  svg.setAttribute("stroke-width", "1.5")
  svg.setAttribute("stroke", "currentColor")
  svg.setAttribute("class", "w-3 h-3")
  const path = document.createElementNS("http://www.w3.org/2000/svg", "path")
  path.setAttribute("stroke-linecap", "round")
  path.setAttribute("stroke-linejoin", "round")
  path.setAttribute("d", pathString)
  svg.append(path)
  return svg
}

export function getWarningIcon() {
  return getIcon(
    "M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 " +
      "0L2.697 16.126zM12 15.75h.007v.008H12v-.008z"
  )
}

export function getLeftArrow() {
  return getIcon("M10.5 19.5L3 12m0 0l7.5-7.5M3 12h18")
}

export function getRightArrow() {
  return getIcon("M13.5 4.5L21 12m0 0l-7.5 7.5M21 12H3")
}

export function appendLineWithIcon(element: HTMLElement, icon: Element, text: string): void {
  const line = document.createElement("span")
  line.setAttribute("class", "flex gap-1.5 items-center")
  line.append(icon)
  line.append(text)
  element.append(line)
}
