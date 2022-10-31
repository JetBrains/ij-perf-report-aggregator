<template>
  <header class="flex flex-col gap-y-0.5">
    <PageHeader />
    <NavigationTabs
      :items="tabNavigationItems"
      :current-path="currentPath"
    />
  </header>
  <div class="px-7 py-5">
    <slot />
  </div>
</template>
<script setup lang="ts">
import { provide } from "vue"
import { useRouter } from "vue-router"
import PageHeader from "./PageHeader.vue"
import { InfoSidebarVmImpl } from "./components/InfoSidebarVm"
import NavigationTabs from "./components/NavigationTabs.vue"
import { tabNavigationItems } from './routes'
import {sidebarVmKey } from "./shared/keys"

const sidebarVm = new InfoSidebarVmImpl()
const router = useRouter()
const currentPath = router.currentRoute.value.path
console.log("currentPath:", currentPath)

provide(sidebarVmKey, sidebarVm)
</script>