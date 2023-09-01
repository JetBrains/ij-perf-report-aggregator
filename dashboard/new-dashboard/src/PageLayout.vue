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
import { useRouter } from "vue-router"
import PageHeader from "./PageHeader.vue"
import NavigationTabs from "./components/NavigationTabs.vue"
import { getNavigationElement, Tab } from "./routes"

import "./shared/overrides.scss"

const router = useRouter()
const currentPath = router.currentRoute.value.path
const product = getNavigationElement(currentPath)
const tabs: Tab[] =
  product.children.find((child) => {
    // eslint-disable-next-line @typescript-eslint/prefer-string-starts-ends-with
    return currentPath.slice(0, Math.max(0, currentPath.lastIndexOf("/"))) == child.url
  })?.tabs ?? product.children[0].tabs
</script>
