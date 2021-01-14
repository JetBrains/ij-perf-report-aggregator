// Copyright 2000-2020 JetBrains s.r.o. Use of this source code is governed by the Apache 2.0 license that can be found in the LICENSE file.
import { Component, Vue, Watch } from "vue-property-decorator"
import { loadJson } from "@/httpUtil"
import { DataRequest, expandMachineAsFilterValue, InfoResponse, MachineGroup, MetricInfo } from "@/aggregatedStats/model"
import { debounce } from "debounce"
import { timeRanges } from "./parseDuration"
import { Location } from "vue-router"
import { Notification } from "element-ui"
import { ChartSettings } from "@/aggregatedStats/ChartSettings"
import { asArray, ValueFilter } from "@/aggregatedStats/ValueFilter"
import { getOrCreateChartSettings } from "@/state/state"

// @ts-ignore
@Component
export abstract class AggregatedStatsPage extends Vue {
  readonly timeRanges = timeRanges

  private lastInfoResponse: InfoResponse | null = null

  get chartSettings(): ChartSettings {
    if (this.getDbName() == null) {
      throw new Error("dbName is null")
    }
    return getOrCreateChartSettings(this.$store, this.getStateModuleName())
  }

  get selectedProjects() {
    return this.chartSettings.selectedProjects
  }
  set selectedProjects(value: Array<string> | null) {
    const list = asArray(value)
    this.$store.commit(`${this.getStateModuleName()}/updateProject`, list)

    console.log(`project changed (${value})`)
    if (value == null) {
      return
    }

    this.requestDataReloading(this.chartSettings.selectedProduct, this.chartSettings.selectedMachine, this.projectFilterManager.getValue(this.chartSettings, this.projects))

    if (this.getDbName() === "ij" && list.length > 0) {
      const currentQuery = this.$route.query
      if (currentQuery.project !== list[0]) {
        // noinspection JSIgnoredPromiseFromCall
        this.$router.push({
          query: {
            ...currentQuery,
            project: list[0],
          },
        })
      }
    }
  }

  private getStateModuleName() {
    return `${this.getDbName()}-dashboard`
  }

  get selectedProject() {
    const list = this.chartSettings.selectedProjects
    return list == null || list.length === 0 ? "" : list[0]
  }
  set selectedProject(value: string | null) {
    this.selectedProjects = asArray(value)
  }

  products: Array<string> = []
  projects: Array<string> = []
  machines: Array<MachineGroup> = []
  metrics: Array<MetricInfo> = []

  isFetching: boolean = false

  timeRange: String = "3M"

  private loadDataAfterDelay = debounce(() => {
    this.loadData()
  }, 1000)

  private commitChartSettingsAfterDelay = debounce(() => {
    this.commitChartSettings()
  }, 10_000)

  protected abstract getDbName(): string

  protected abstract get projectFilterManager(): ValueFilter

  dataRequest: DataRequest | null = null

  loadData() {
    this.isFetching = true
    loadJson(`${this.chartSettings.serverUrl}/api/v1/info?db=${this.getDbName()}`, null)
      .then((data: InfoResponse | null) => {
        if (data == null) {
          this.isFetching = false
          return
        }

        this.lastInfoResponse = Object.seal(data)
        this.products = data.productNames
        this.metrics = data.metrics

        let selectedProduct = this.chartSettings.selectedProduct
        if (this.products.length === 0) {
          selectedProduct = ""
        }
        else if (selectedProduct == null || selectedProduct.length === 0 || !this.products.includes(selectedProduct)) {
          selectedProduct = this.products[0]
        }
        this.chartSettings.selectedProduct = selectedProduct

        // not called by Vue for some reasons
        console.log("update product on info response", selectedProduct)
        const oldSelectedMachine = this.chartSettings.selectedMachine
        this.applyChangedProduct(selectedProduct, data)
        const newSelectedMachine = this.chartSettings.selectedMachine
        if (!isArrayContentTheSame(oldSelectedMachine, newSelectedMachine)) {
          this.selectedMachineChanged(newSelectedMachine, oldSelectedMachine)
        }

        this.isFetching = false
      })
      .catch(e => {
        this.isFetching = false
        console.error(e)
      })
  }

  @Watch("$route")
  onRouteChanged(location: Location, _oldLocation: Location): void {
    this.setProductFromQuery(location.query?.product)
    this.setProjectFromQuery(location.query?.project)
  }

  @Watch("chartSettings.selectedProduct")
  selectedProductChanged(product: string | null, oldProduct: string | undefined): void {
    console.info(`product changed (${oldProduct} => ${product})`)

    const infoResponse = this.lastInfoResponse
    if (infoResponse != null) {
      this.applyChangedProduct(product, infoResponse)
    }

    const currentQuery = this.$route.query
    if (currentQuery.product !== product && product != null && product.length > 0) {
      // noinspection JSIgnoredPromiseFromCall
      this.$router.push({
        query: {
          ...currentQuery,
          product,
        },
      })
    }
  }

  protected convertProjectNameToTitle(projectName: string): string {
    return projectName
  }

  private applyChangedProduct(product: string | null, info: InfoResponse) {
    if (product != null && product.length > 0) {
      // later maybe will be more info for machine, so, do not use string instead of Machine
      this.machines = info.productToMachine[product] || []
      if (this.machines.length === 0) {
        Notification.error(`No machines for product ${product}. Please check that product code is valid.`)
      }

      const projects = info.productToProjects[product] || []
      projects.sort((a, b) => {
        const t1 = this.convertProjectNameToTitle(a)
        const t2 = this.convertProjectNameToTitle(b)
        if (t1.startsWith("simple ") && !t2.startsWith("simple ")) {
          return -1
        }
        if (t2.startsWith("simple ") && !t1.startsWith("simple ")) {
          return 1
        }
        return t1.localeCompare(t2)
      })
      this.projects = projects
    }
    else {
      console.error(`set machines to empty list because no product (${product})`)
      this.machines = []
      this.projects = []
    }

    let selectedMachine = this.chartSettings.selectedMachine || []
    const machines = this.machines
    if (machines.length === 0) {
      selectedMachine = []
      this.chartSettings.selectedMachine = selectedMachine
    }
    else if (selectedMachine.length === 0 || !machines.find(it => selectedMachine.includes(it.name)) || !machines.find(it => it.children.find(it => selectedMachine.includes(it.name)))) {
      selectedMachine = [machines[0].name]
    }

    const selectedProjects = this.projectFilterManager.getValue(this.chartSettings, this.projects)
    if (product != null && selectedProjects.length > 0) {
      this.requestDataReloading(product, selectedMachine, selectedProjects)
    }

    if (isArrayContentTheSame(this.chartSettings.selectedMachine, selectedMachine)) {
      // data will be reloaded on machine change, but if product changed but machine remain the same, data reloading must be triggered here
      if (product != null && selectedMachine.length > 0 && selectedProjects.length > 0) {
        this.requestDataReloading(product, selectedMachine, selectedProjects)
      }
    }
    else {
      this.chartSettings.selectedMachine = selectedMachine
    }
  }

  private requestDataReloading(product: string, machine: Array<string>, projects: Array<string>) {
    if (product == null || product.length === 0 || projects.length === 0 || machine.length === 0) {
      return
    }

    this.dataRequest = Object.seal({
      db: this.getDbName(),
      product,
      machine: expandMachineAsFilterValue(product, machine, this.lastInfoResponse!!),
      projects
    })
  }

  @Watch("chartSettings.selectedMachine")
  selectedMachineChanged(machine: Array<string> | string, oldMachine: Array<string>): void {
    if (typeof machine === "string") {
      machine = [machine]
      this.chartSettings.selectedMachine = machine
    }

    if (oldMachine === machine) {
      return
    }

    console.log(`machine changed (${Array.isArray(oldMachine) ? oldMachine.join() : oldMachine} => ${Array.isArray(machine) ? machine.join() : machine})`)
    if (machine == null) {
      return
    }

    const product = this.chartSettings.selectedProduct
    const projects = this.projectFilterManager.getValue(this.chartSettings, this.projects)
    if (product == null || product.length === 0 || projects.length === 0) {
      return
    }

    this.requestDataReloading(product, machine, projects)
  }

  @Watch("chartSettings.serverUrl")
  serverUrlChanged(newV: string | null, _oldV: string) {
    if (!isEmpty(newV)) {
      this.loadDataAfterDelay()
    }
  }

  @Watch("chartSettings", {deep: true})
  chartSettingsChanged(_newV: any, _oldV: any) {
    // https://vuex.vuejs.org/guide/forms.html
    // vuex forces you to do everything via mutation, and if state modified directly, persisted store is not notified about change
    // so, watch chartSettings and schedule update of localStorage
    // this.dataModule.updateChartSettings(this.chartSettings, this.getDbName())
    this.commitChartSettingsAfterDelay()
  }

  private commitChartSettings() {
    this.$store.commit(`${this.getStateModuleName()}/updateSettings`)
  }

  beforeMount() {
    const serverUrl = this.chartSettings.serverUrl
    if (!isEmpty(serverUrl)) {
      const query = this.$route.query
      this.setProductFromQuery(query.product)
      this.setProjectFromQuery(query.project)
      this.loadData()
    }
  }

  private setProductFromQuery(product: string | undefined | null | (string | undefined | null)[]) {
    if (product != null && product.length == 2) {
      const value = (product as string).toUpperCase()
      const chartSettings = this.chartSettings
      if (chartSettings.selectedProduct !== value) {
        console.log(`product specified in query: ${product}`)
        chartSettings.selectedProduct = value
      }
    }
  }

  private setProjectFromQuery(project: string | undefined | null | (string | undefined | null)[]) {
    if (project != null && project.length > 0) {
      const chartSettings = this.chartSettings
      const list = chartSettings.selectedProjects
      if (list.length == 0 || list[0] !== project) {
        console.log(`project specified in query: ${project}`)
        chartSettings.selectedProjects = [project as string]
      }
    }
  }
}

function isEmpty(v: string | null): boolean {
  return v == null || v.length === 0
}

function isArrayContentTheSame(a: Array<string>, b: Array<string>): boolean {
  if (a.length !== b.length) {
    return false
  }

  for (const item of a) {
    if (!b.includes(item)) {
      return false
    }
  }
  return true
}