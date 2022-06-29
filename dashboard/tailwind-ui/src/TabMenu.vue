<!-- https://tailwindui.com/components/application-ui/navigation/tabs -->
<!-- Tabs with underline -->
<!-- Customization of PrimeVue TabMenu was a quite painful - in addition to theme CSS, component itself sets some style -->
<template>
  <div class="border-b border-gray-200">
    <nav
      class="-mb-px flex space-x-8"
      aria-label="Tabs"
    >
      <router-link
        v-for="tab in props.model"
        :key="tab.to"
        v-slot="{ href, isActive, navigate }"
        :to="tab.to"
        custom
      >
        <a
          :href="href"
          :class="[
            isActive ? 'border-primary-500 text-primary-600' : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300',
            'whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm'
          ]"
          :aria-current="isActive ? 'page' : undefined"
          @click="navigate"
        >
          {{ tab.label }}
        </a>
      </router-link>
    </nav>
  </div>
</template>
<script setup lang="ts">

import { TabItem } from "./model"

const props = defineProps<{
  model: Array<TabItem>
}>()
</script>