<template>
  <Divider :label="props.title" />
  <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
    <slot />
  </div>
</template>
<script setup lang="ts">
import Divider from "../tailwind-ui/Divider.vue"

const props = defineProps<{
  title: string
}>()
</script>