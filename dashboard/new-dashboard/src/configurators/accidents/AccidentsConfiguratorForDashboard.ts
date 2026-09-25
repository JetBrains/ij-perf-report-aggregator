import { Chart } from "../../components/charts/DashboardCharts"
import { combineLatest } from "../rxjs"
import { TimeRangeConfigurator } from "../TimeRangeConfigurator"
import { dbTypeStore } from "../../shared/dbTypes"
import { AccidentsConfigurator } from "./AccidentsConfigurator"

export class AccidentsConfiguratorForDashboard extends AccidentsConfigurator {
  constructor(
    private readonly serverUrl: string,
    charts: Chart[] | null,
    timeRangeConfigurator: TimeRangeConfigurator
  ) {
    super()
    this.dbType = dbTypeStore().dbType
    const tests = this.getProjectAndProjectWithMetrics(charts)
    combineLatest([timeRangeConfigurator.createObservable()]).subscribe(([[timeRange, customRange]]) => {
      this.loadAccidents(tests, timeRange, customRange)
    })
  }

  protected getAccidentUrl(): string {
    return this.serverUrl + "/api/meta/"
  }

  private getProjectAndProjectWithMetrics(charts: Chart[] | null): string[] {
    const projectsWithMetrics =
      charts?.flatMap((chart) => {
        const measures = Array.isArray(chart.definition.measure) ? chart.definition.measure : [chart.definition.measure]
        return chart.projects.flatMap((project) => {
          return measures.map((measure) => project + "/" + measure)
        })
      }) ?? []
    const projects = new Set(charts?.map((it) => it.projects).flat(Number.POSITIVE_INFINITY) as string[])
    return [...projectsWithMetrics, ...projects]
  }
}
