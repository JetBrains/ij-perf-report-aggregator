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
import { getNavigationElement, Tab } from "./routes"
import {sidebarVmKey } from "./shared/keys"

import "./shared/overrides.css"

const sidebarVm = new InfoSidebarVmImpl()
const router = useRouter()
const currentPath = router.currentRoute.value.path
const product = getNavigationElement(currentPath)
const tabs: Tab[] = product.children.find(child => {
  // eslint-disable-next-line @typescript-eslint/prefer-string-starts-ends-with
  return currentPath.slice(0, Math.max(0, currentPath.lastIndexOf("/"))) == child.url
})?.tabs ?? product.children[0].tabs

provide(sidebarVmKey, sidebarVm)
</script>