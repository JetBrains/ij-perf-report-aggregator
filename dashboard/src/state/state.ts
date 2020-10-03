// Copyright 2000-2019 JetBrains s.r.o. Use of this source code is governed by the Apache 2.0 license that can be found in the LICENSE file.
import Vue from "vue"
import Vuex, { Store } from "vuex"
import { ChartSettings } from "@/aggregatedStats/ChartSettings"
import createPersistedState from "vuex-persistedstate"

export class ReportVisualizerSettings {
  data: string | null = ""
  recentlyUsedIdePort: number = 63342
}

export const markerNames = ["app initialized callback", "module loading"]

export function getOrCreateChartSettings(store: Store<any>, moduleName: string): ChartSettings {
  function registerModule() {
    store.registerModule(moduleName, {
      namespaced: true,
      state: new ChartSettings(),
      mutations: {
        updateProject(state: ChartSettings, value: Array<string>) {
          state.selectedProjects = value
        },
        updateSettings(state: ChartSettings) {
          Object.assign(state, state)
        },
      }
    })
    // createPersistedState plugin must be not statically registered,
    // because in this case existing state will create the same key in the global state, and dynamically registered module will overwrite existing data by empty initial state
    createPersistedState({key: moduleName, paths: [moduleName]})(store)
  }

  if (!store.hasModule(moduleName)) {
    registerModule()
  }
  return store.state[moduleName]
}

export const reportVisualizerModuleName = "report-visualizer"

export function getOrCreateReportVisualizerSettings(store: Store<any>): ReportVisualizerSettings {
  const moduleName = reportVisualizerModuleName
  if (!store.hasModule(moduleName)) {
    store.registerModule(moduleName, {
      namespaced: true,
      state: new ReportVisualizerSettings(),
      mutations: {
        updateData(state: ReportVisualizerSettings, value: string) {
          state.data = value
        },
        updateRecentlyUsedIdePort(state: ReportVisualizerSettings, value: number) {
          state.recentlyUsedIdePort = value
        },
      }
    })
    createPersistedState({key: moduleName, paths: [moduleName]})(store)
  }
  return store.state[moduleName]
}

// // name here is important because otherwise getModule from vuex-module-decorators will not work
// @Module({name: mainModuleName})
// export class AppStateModule extends VuexModule implements AppState {
//   data: DataManager | null = null
//   recentlyUsedIdePort: number = defaultIdePort
//
//   chartSettings: { [key: string]: ChartSettings; } = {}
//
//   @Mutation
//   updateData(data: InputData | null) {
//     this.data = stateStorageManager.createDataManager(data)
//   }
//
//   @Mutation
//   updateRecentlyUsedIdePort(value: number) {
//     this.recentlyUsedIdePort = value
//   }
// }

Vue.use(Vuex)

// const vuexLocal = new VuexPersistence({
//   storage: window.localStorage,
//   saveState: function (_key: string, moduleNameToState: any, _storage) {
//     stateStorageManager.saveState(moduleNameToState)
//   },
//   restoreState: function (_key: string, _storage) {
//     return stateStorageManager.restoreState()
//   },
// })

export const store = new Vuex.Store({
})
