<template>
  <Message
    v-show="messageState.isError"
    severity="error"
    @close="closeError"
  >
    {{ messageState.message }}
  </Message>
  <Menubar :model="items">
    <template #start>
      <img
        width="70"
        src="https://resources.jetbrains.com/storage/products/company/brand/logos/jb_square.svg"
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
import { useMessageStore } from "../stores/Message"
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

    const messageState = useMessageStore()
    return {
      serverUrl,
      activePath,
      routes,
      items,
      messageState,
      closeError(): void {
        messageState.isError = false
      },
    }
  },
})
</script>

<style>
.p-toolbar-group-left > * {
  margin: 10px;
}
</style>