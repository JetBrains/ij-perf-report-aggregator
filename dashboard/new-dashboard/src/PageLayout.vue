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
import { computed } from "vue"
import { useRouter } from "vue-router"
import PageHeader from "./PageHeader.vue"
import NavigationTabs from "./components/NavigationTabs.vue"
import { getNavigationElement } from "./routes"

const router = useRouter()
const currentPath = computed(() => router.currentRoute.value.path)
const product = computed(() => {
  return getNavigationElement(currentPath.value)
})

const tabs = computed(() => {
  return (
    product.value.children.find((child) => {
      return currentPath.value.slice(0, Math.max(0, currentPath.value.lastIndexOf("/"))) == child.url
    })?.tabs ?? product.value.children[0].tabs
  )
})
</script>
