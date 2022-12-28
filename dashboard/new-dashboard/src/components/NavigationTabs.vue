<template>
  <div class="px-7 flex gap-7 items-center border-b border-gray-200 border-solid">
    <router-link
      v-for="item in props.items"
      :key="item.path"
      v-slot="{ href, isActive, navigate }"
      :to="item.path"
      custom
    >
      <a
        :href="href"
        class="py-1 text-sm text-gray-500 relative"
        :class="{
          'text-black font-medium pointer-events-none after:content-[\'\'] after:absolute after:bottom-[-1px] after:left-0 after:h-0.5 after:w-full after:bg-blue-500': isActive,
        }"
        @click="navigate"
      >
        {{ item.name }}
      </a>
    </router-link>
  </div>
</template>
<script setup lang="ts">
import { NavigationItem } from "../routes"

interface Props {
  items: NavigationItem[]
  currentPath: string
}
const props = defineProps<Props>()
</script>