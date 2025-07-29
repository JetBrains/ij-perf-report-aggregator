import { ChartDefinition } from "../charts/DashboardCharts"

export const MODES = ["intellij-idea", "android-studio"]

const directMapping: Record<string, string> = {
  "Progress: Setting up run configurations...": "setting_run_configuration_android",
  "Progress: Generating Xcode files…": "generating_prebuild_xcode_files",
}

const MEASUREMENTS_MAPPING = Object.entries(directMapping).reduce<Record<string, string>>(
  (acc, [key, value]) => {
    // acc is the accumulator object that gets built up with each iteration
    acc[key] = value

    MODES.forEach((mode) => {
      acc[`${mode} – ${key}`] = `${mode} – ${value}`
    })

    return acc
  },
  {} // initial value of the accumulator
)

export const legendFormatter = (name: string) => {
  if (name in MEASUREMENTS_MAPPING) {
    return MEASUREMENTS_MAPPING[name]
  }
  return name
}

export function filterChartsByProjects(chartDefinitions: ChartDefinition[], selectedProjectRoots: string[]): ChartDefinition[] {
  return chartDefinitions
    .map((chart) => {
      // Filter projects that start with any of the selected project roots
      const filteredProjects = chart.projects.filter((project) => {
        const projectRoot = project.split("/")[0]
        return selectedProjectRoots.includes(projectRoot)
      })

      // Only return charts that have at least one matching project
      if (filteredProjects.length > 0) {
        // Create a new chart definition with filtered projects but same labels and measures
        return {
          ...chart,
          projects: filteredProjects,
        }
      }

      return null
    })
    .filter((chart) => chart !== null) as ChartDefinition[]
}
