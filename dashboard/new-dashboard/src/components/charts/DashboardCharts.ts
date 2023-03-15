export interface Definition {
  label: string
  measure: string
}

export interface Chart {
  definition: Definition
  projects: Array<string>
}

export interface ChartDefinition {
  labels: Array<string>
  measures: Array<string>
  projects: Array<string>
}

export function combineCharts(charts: Array<ChartDefinition>): Array<Chart> {
  const resultingCharts = new Array<Chart>
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
      })
    }
  }
  return resultingCharts
}

export function extractUniqueProjects(charts: Array<ChartDefinition>): Array<string> {
  const allProjects = new Set<string>
  for (const chart of charts) {
    for (const project of chart.projects) {
      allProjects.add(project)
    }
  }
  return [...allProjects]
}