import { Ref } from "vue"
import { TimeRangeConfigurator } from "../TimeRangeConfigurator"
import { dbTypeStore } from "../../shared/dbTypes"
import { combineLatest } from "rxjs"
import { refToObservable } from "../rxjs"
import { AccidentsConfigurator } from "./AccidentsConfigurator"

export class AccidentsConfiguratorForTests extends AccidentsConfigurator {
  constructor(
    private serverUrl: string,
    projects: Ref<string | string[] | null>,
    metrics: Ref<string[] | string | null>,
    timeRangeConfigurator: TimeRangeConfigurator
  ) {
    super()
    this.dbType = dbTypeStore().dbType
    combineLatest([refToObservable(projects), refToObservable(metrics), timeRangeConfigurator.createObservable()]).subscribe(([projects, measures, [timeRange, customRange]]) => {
      const projectAndMetrics = this.combineProjectsAndMetrics(projects, measures)
      this.getAccidentsFromMetaDb(projectAndMetrics, timeRange, customRange)
        .then((value) => {
          this.value.value = value
        })
        .catch((error: unknown) => {
          console.error(error)
        })
    })
  }

  protected getAccidentUrl(): string {
    return this.serverUrl + "/api/meta/"
  }
}
