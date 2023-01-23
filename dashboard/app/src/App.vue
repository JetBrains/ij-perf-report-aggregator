<template>
  <router-view v-slot="{ Component, route }">
    <template v-if="isNewDashboardRoute(route)">
      <PageLayout>
        <keep-alive
          :key="route.path"
          max="4"
        >
          <component :is="Component" />
        </keep-alive>
      </PageLayout>
    </template>
    <template v-else>
      <PrimeToast />
      <Menubar
        :model="items"
        class="px-6 border-b"
      >
        <template #start>
          <!-- eslint-disable vue/no-v-html -->
          <span
            v-html="logoUrl"
          />
        </template>
        <template #end>
          <ServerSelect
            v-show='!activePath.startsWith("/report")'
            v-model="serverUrl"
          />
        </template>
      </Menubar>
      <main class="mx-auto px-6 py-4">
        <keep-alive
          :key="route.path"
          max="4"
        >
          <component :is="Component" />
        </keep-alive>
      </main>
    </template>
  </router-view>
</template>
<script setup lang="ts">
import PageLayout from "new-dashboard/src/PageLayout.vue"
import { filter, shareReplay } from "rxjs"
import { PersistentStateManager } from "shared/src/PersistentStateManager"
import ServerSelect from "shared/src/components/ServerSelect.vue"
import { ServerConfigurator } from "shared/src/configurators/ServerConfigurator"
import { refToObservable } from "shared/src/configurators/rxjs"
import { serverUrlObservableKey } from "shared/src/injectionKeys"
import { provide, ref, shallowRef, watch } from "vue"
import { RouteLocationNormalizedLoaded, useRoute } from "vue-router"
import PrimeToast from "./PrimeToast.vue"
import logoUrl from "./jb_square.svg?raw"
import { getItems } from "./route"

const serverUrl = shallowRef(ServerConfigurator.DEFAULT_SERVER_URL)
// shallow ref doesn't work - items are modified by primevue
const items = ref(getItems())
const serverUrlObservable = refToObservable(serverUrl).pipe(
  filter((it: string | null): it is string => it !== null && it.length > 0),
  shareReplay(1),
)
provide(serverUrlObservableKey, serverUrlObservable)

const activePath = shallowRef("")
const _route = useRoute()
watch(() => _route.path, p => {
  activePath.value = p
})

const persistentStateManager = new PersistentStateManager("common", {serverUrl: ServerConfigurator.DEFAULT_SERVER_URL})
persistentStateManager.add("serverUrl", serverUrl)

function isNewDashboardRoute(route: RouteLocationNormalizedLoaded): boolean {
  return !route.path.startsWith("/old") && !route.path.startsWith("/ij")
}
</script>