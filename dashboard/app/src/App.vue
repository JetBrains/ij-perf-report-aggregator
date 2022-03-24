<template>
  <Menubar :model="items">
    <template #end>
      <ServerSelect
        v-show='!activePath.startsWith("/report")'
        v-model="serverUrl"
      />
    </template>
  </Menubar>
  <router-view v-slot="{ Component, route }">
    <keep-alive
      :key="route.path"
      max="4"
    >
      <component :is="Component" />
    </keep-alive>
  </router-view>
</template>
<script lang="ts">
import { PersistentStateManager } from "shared/src/PersistentStateManager"
import ServerSelect from "shared/src/components/ServerSelect.vue"
import { ServerConfigurator } from "shared/src/configurators/ServerConfigurator"
import { serverUrlKey } from "shared/src/injectionKeys"
import { defineComponent, provide, ref, watch } from "vue"
import { useRoute } from "vue-router"
import { getItems, getRoutes } from "./route"

export default defineComponent({
  name: "App",
  components: {
    ServerSelect,
  },

  setup() {
    const serverUrl = ref("")
    const routes = getRoutes()
    const items = ref(getItems())
    provide(serverUrlKey, serverUrl)

    const activePath = ref("")
    const route = useRoute()
    watch(() => route.path, p => {
      activePath.value = p
    })

    const persistentStateManager = new PersistentStateManager("common", {serverUrl: ServerConfigurator.DEFAULT_SERVER_URL})
    persistentStateManager.add("serverUrl", serverUrl)
    persistentStateManager.init()

    return {
      serverUrl,
      activePath,
      routes,
      items,
    }
  },
})
</script>