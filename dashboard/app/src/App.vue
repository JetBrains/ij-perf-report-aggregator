<template>
  <router-view v-slot="{ Component, route }">
    <PageLayout>
      <keep-alive
        :key="route.path"
        max="4"
      >
        <component :is="Component" />
      </keep-alive>
    </PageLayout>
  </router-view>
</template>
<script setup lang="ts">
import PageLayout from "new-dashboard/src/PageLayout.vue"
import { filter, shareReplay } from "rxjs"
import { PersistentStateManager } from "shared/src/PersistentStateManager"
import { ServerConfigurator } from "shared/src/configurators/ServerConfigurator"
import { limit, refToObservable } from "shared/src/configurators/rxjs"
import { serverUrlObservableKey } from "shared/src/injectionKeys"
import { provide, shallowRef, watch } from "vue"
import { useRoute } from "vue-router"

const serverUrl = shallowRef(ServerConfigurator.DEFAULT_SERVER_URL)
// shallow ref doesn't work - items are modified by primevue
const serverUrlObservable = refToObservable(serverUrl).pipe(
  filter((it: string | null): it is string => it !== null && it.length > 0),
  shareReplay(1),
)
provide(serverUrlObservableKey, serverUrlObservable)

const activePath = shallowRef("")
const _route = useRoute()
watch(() => _route.path, p => {
  activePath.value = p
  limit.clearQueue()
})

const persistentStateManager = new PersistentStateManager("common", {serverUrl: ServerConfigurator.DEFAULT_SERVER_URL})
persistentStateManager.add("serverUrl", serverUrl)

</script>