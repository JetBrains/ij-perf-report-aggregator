import { Ref } from "vue"
import { TimeRangeConfigurator } from "../TimeRangeConfigurator"
import { dbTypeStore } from "../../shared/dbTypes"
import { combineLatest } from "rxjs"
import { refToObservable } from "../rxjs"
import { useUserStore } from "../../shared/useUserStore"
import { Accident, AccidentKind, AccidentsConfigurator } from "./AccidentsConfigurator"

export class AccidentsConfiguratorForStartup extends AccidentsConfigurator {
  constructor(
    private serverUrl: string,
    private product: Ref<string | string[] | null>,
    projects: Ref<string | string[] | null>,
    metrics: Ref<string[] | string | null>,
    timeRangeConfigurator: TimeRangeConfigurator
  ) {
    super()
    this.dbType = dbTypeStore().dbType
    combineLatest([refToObservable(projects), refToObservable(metrics), timeRangeConfigurator.createObservable(), refToObservable(product)]).subscribe(
      ([projects, measures, [timeRange, customRange], product]) => {
        if (product == null) return
        if (Array.isArray(product)) return
        const projectAndMetrics = this.combineProjectsAndMetrics(projects, measures)
        const projectAndMetricsWithProduct = projectAndMetrics.map((it) => `${product}/${it}`)
        this.getAccidentsFromMetaDb(projectAndMetricsWithProduct, timeRange, customRange)
          .then((value) => {
            this.value.value = this.removeProductPrefix(product, value)
          })
          .catch((error: unknown) => {
            console.error(error)
          })
      }
    )
  }

  protected getAccidentUrl(): string {
    return this.serverUrl + "/api/meta/"
  }

  private removeProductPrefix(product: string, response: Map<string, Accident[]>): Map<string, Accident[]> {
    const map = new Map<string, Accident[]>()
    for (const [key, value] of response) {
      const keyWithoutProduct = key.replace(`${product}/`, "")
      map.set(keyWithoutProduct, value)
    }
    return map
  }

  async writeAccidentToMetaDb(date: string, affected_test: string, reason: string, build_number: string, kind: string | undefined, stacktrace: string = "") {
    if (this.product.value == null || Array.isArray(this.product.value)) return
    const test = `${this.product.value}/${affected_test}`
    try {
      const userName = useUserStore().user?.name ?? ""
      const response = await fetch(this.getAccidentUrl() + "accidents/", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ date, affected_test: test, reason, build_number: build_number.toString(), kind, stacktrace, user_name: userName }),
      })

      if (!response.ok) {
        throw new Error("The accident wasn't created")
      }
      const idString: string = await response.text()
      const id = Number(idString)
      this.value.value ??= new Map<string, Accident[]>()
      const updatedMap = new Map(this.value.value)
      updatedMap.set(`${affected_test}_${build_number}`, [
        { id, affectedTest: affected_test, date, reason, buildNumber: build_number, kind: kind as AccidentKind, stacktrace, userName },
      ])
      this.value.value = updatedMap //we need to update value in reference to trigger the change
      return id
    } catch (error) {
      console.error(error)
      return
    }
  }
}
