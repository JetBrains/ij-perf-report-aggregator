<template>
  <PrimeToast />
  <Menubar
    :model="items"
    class="!rounded-none !border-0 !border-b !bg-inherit"
  >
    <template #start>
      <img
        width="70"
        :src="logoUrl"
        alt="JetBrains Black Box Logo logo."
      >
    </template>
    <template #end>
      <ServerSelect
        v-show='!activePath.startsWith("/report")'
        v-model="serverUrl"
      />
    </template>
  </Menubar>
  <main class="mx-auto px-6 py-4">
    <router-view v-slot="{ Component, route }">
      <keep-alive
        :key="route.path"
        max="4"
      >
        <component :is="Component" />
      </keep-alive>
    </router-view>
  </main>
</template>
<script setup lang="ts">
import { shareReplay } from "rxjs"
import { PersistentStateManager } from "shared/src/PersistentStateManager"
import ServerSelect from "shared/src/components/ServerSelect.vue"
import { ServerConfigurator } from "shared/src/configurators/ServerConfigurator"
import { refToObservable } from "shared/src/configurators/rxjs"
import { serverUrlObservableKey } from "shared/src/injectionKeys"
import { provide, ref, shallowRef, watch } from "vue"
import { useRoute } from "vue-router"
import PrimeToast from "./PrimeToast.vue"
import logoUrl from "./jb_square.svg?url"
import { getItems } from "./route"

const serverUrl = shallowRef(ServerConfigurator.DEFAULT_SERVER_URL)
// shallow ref doesn't work - items are modified by primevue
const items = ref(getItems())
provide(serverUrlObservableKey, refToObservable(serverUrl).pipe(
    shareReplay(1),
  ),
)

const activePath = shallowRef("")
const _route = useRoute()
watch(() => _route.path, p => {
  activePath.value = p
})

const persistentStateManager = new PersistentStateManager("common", {serverUrl: ServerConfigurator.DEFAULT_SERVER_URL})
persistentStateManager.add("serverUrl", serverUrl)
</script>