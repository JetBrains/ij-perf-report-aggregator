<template>
  <div class="flex justify-between pt-5 px-7">
    <button
      class="px-1 py-1 inline-flex text-xl items-center"
      type="button"
      @click="toggle"
    >
      {{ menuItems[0]?.label }}
      <div class="pi pi-chevron-down text-sm ml-1.5" />
    </button>
    <Menu
      ref="menu"
      :model="items"
      :popup="true"
    />

    <ServerSelect v-model="serverUrl" />
  </div>
</template>

<script setup lang="ts">
import Menu from "primevue/menu"
import ServerSelect from "shared/src/components/ServerSelect.vue"
import { ServerConfigurator } from "shared/src/configurators/ServerConfigurator"
import { ref, shallowRef } from "vue"
import { topNavigationItems } from "./routes"

const serverUrl = shallowRef(ServerConfigurator.DEFAULT_SERVER_URL)
const menuItems = topNavigationItems.map(({ path, name }) => ({
    label: name,
    url: path,
}))
const items = ref(menuItems)
const menu = ref<Menu | null>(null)
function toggle(event: PointerEvent) {
    menu.value?.toggle(event)
}
</script>