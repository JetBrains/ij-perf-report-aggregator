<template>
  <Divider align="center">
    {{ props.title }}
  </Divider>
  <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
    <slot />
  </div>
</template>
<script setup lang="ts">
const props = defineProps<{
  title: string
}>()
</script>