<template>
  <div class="new-dashboard">
    <header class="flex flex-col gap-y-0.5">
      <PageHeader />
      <NavigationTabs
        :items="tabs"
        :current-path="currentPath"
      />
    </header>
    <div class="px-7 py-5">
      <slot />
    </div>
  </div>
</template>
<script setup lang="ts">
import { provide } from "vue"
import { useRouter } from "vue-router"
import PageHeader from "./PageHeader.vue"
import { InfoSidebarVmImpl } from "./components/InfoSidebarVm"
import NavigationTabs from "./components/NavigationTabs.vue"
import { getNavigationTabs } from "./routes"
import {sidebarVmKey } from "./shared/keys"

import "./shared/overrides.css"

const sidebarVm = new InfoSidebarVmImpl()
const router = useRouter()
const currentPath = router.currentRoute.value.path
const tabs = getNavigationTabs(currentPath)

provide(sidebarVmKey, sidebarVm)
</script>