<template>
  <Message
    v-show="messageState.isError"
    severity="error"
    @close="closeError"
  >
    {{ messageState.message }}
  </Message>
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
/// <reference types="vite-svg-loader" />

import { PersistentStateManager } from "shared/src/PersistentStateManager"
import ServerSelect from "shared/src/components/ServerSelect.vue"
import { ServerConfigurator } from "shared/src/configurators/ServerConfigurator"
import { serverUrlKey } from "shared/src/injectionKeys"
import { provide, ref, watch } from "vue"
import { useRoute } from "vue-router"
import { useMessageStore } from "../stores/Message"

import logoUrl from "./jb_square.svg?url"
import { getItems } from "./route"

const serverUrl = ref("")
const items = ref(getItems())
provide(serverUrlKey, serverUrl)

const activePath = ref("")
const _route = useRoute()
watch(() => _route.path, p => {
  activePath.value = p
})

const persistentStateManager = new PersistentStateManager("common", {serverUrl: ServerConfigurator.DEFAULT_SERVER_URL})
persistentStateManager.add("serverUrl", serverUrl)
persistentStateManager.init()

const messageState = useMessageStore()

function closeError(): void {
  messageState.isError = false
}
</script>

<style>
.p-toolbar-group-left > * {
  margin: 10px;
}
</style>