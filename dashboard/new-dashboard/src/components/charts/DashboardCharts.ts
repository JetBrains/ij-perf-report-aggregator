export interface Definition {
  label: string
  measure: string | string[]
}

export interface Chart {
  definition: Definition
  projects: string[]
  aliases?: string[]
}

export interface ChartDefinition {
  labels: string[]
  measures: (string | string[])[]
  projects: string[]
  aliases?: string[]
  machines?: string[]
}

export function combineCharts(charts: ChartDefinition[]): Chart[] {
  const resultingCharts = new Array<Chart>()
  for (const chart of charts) {
    if (chart.labels.length != chart.measures.length) {
      throw new Error("Chart labels and measures arrays must be of the same length.")
    }
    const labelsAndMeasuresCombined = chart.labels.map((label, index) => ({
      label,
      measure: chart.measures[index],
    }))
    for (const labelAndMeasure of labelsAndMeasuresCombined) {
      resultingCharts.push({
        definition: labelAndMeasure,
        projects: chart.projects,
        aliases: chart.aliases,
      })
    }
  }
  return resultingCharts
}

export function extractUniqueProjects(charts: Chart[]): string[] {
  const allProjects = new Set<string>()
  for (const chart of charts) {
    for (const project of chart.projects) {
      allProjects.add(project)
    }
  }
  return [...allProjects]
}
