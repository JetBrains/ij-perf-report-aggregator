<template>
  <div class="flex pt-5 px-7 items-center space-x-3">
    <span
      v-if="subMenuItems.length > 0"
      class="text-xl font-semibold"
    >
      Tests for
    </span>
    <button
      class="text-blue-400 px-1 py-1 inline-flex text-xl items-center"
      type="button"
      @click="toggle"
    >
      {{ selected.name }}
      <div class="pi pi-chevron-down text-sm ml-1.5" />
    </button>
    <Menu
      ref="menu"
      :model="items"
      :popup="true"
    />
    <span
      v-if="subMenuItems.length > 0"
      class="text-xl font-semibold"
    >
      aggregated by
    </span>
    <button
      v-if="subMenuItems.length > 0"
      class="text-blue-400 px-1 py-1 inline-flex text-xl items-center"
      type="button"
      @click="toggleSubMenu"
    >
      {{ selectedSubMenu.name }}
      <div class="pi pi-chevron-down text-sm ml-1.5" />
    </button>

    <Menu
      v-if="subMenuItems.length > 0"
      ref="subMenu"
      :model="subItems"
      :popup="true"
    />
  </div>
</template>

<script setup lang="ts">
import Menu from "primevue/menu"
import { MenuItem } from "primevue/menuitem"
import { ref } from "vue"
import { useRouter } from "vue-router"
import { getSubMenus, topNavigationItems } from "./routes"

const menuItems = topNavigationItems.map(({ path, name }) => ({
    label: name,
    url: path,
})) as MenuItem[]
const items = ref(menuItems)
const menu = ref<Menu | null>(null)
const router = useRouter()
const currentPath = router.currentRoute.value.path
const selected = topNavigationItems.find(({ key }) => {
  return key ? currentPath.startsWith(key) : false
}) ?? topNavigationItems[0]

function toggle(event: PointerEvent) {
  menu.value?.toggle(event)
}

const subMenuItems = getSubMenus(selected.path).map(({ path, name }) => ({
  label: name,
  url: path,
})) as MenuItem[]
const subItems = ref(subMenuItems)
const subMenu = ref<Menu | null>(null)

const selectedSubMenu = getSubMenus(selected.path).findLast(({ key }) => {
  return key ? currentPath.startsWith(key) : false
}) ?? subMenuItems[0]

function toggleSubMenu(event: PointerEvent) {
  subMenu.value?.toggle(event)
}
</script>