<template>
  <div class="px-7 flex gap-7 items-center border-b border-solid">
    <router-link
      v-for="item in items"
      :key="item.url"
      v-slot="{ href, isActive, navigate }"
      :to="item.url"
      custom
    >
      <a
        :href="href"
        class="py-1 text-gray-500 dark:text-gray-300 relative"
        :class="{
          'font-medium pointer-events-none after:content-[\'\'] after:absolute after:bottom-[-1px] after:left-0 after:h-0.5 after:w-full after:bg-primary dark:after:bg-primary-dark':
            isActive,
        }"
        @click="navigate"
      >
        {{ item.label }}
      </a>
    </router-link>
  </div>
</template>
<script setup lang="ts">
import { Tab } from "../routes"

interface Props {
  items: Tab[]
}
const { items } = defineProps<Props>()
</script>
