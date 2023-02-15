<template>
  <div class="relative flex py-5 items-center">
    <div class="flex-grow border-t border-gray-400" />
    <span class="flex-shrink mx-4 text-gray-400 text-lg">{{ props.title }}</span>
    <div class="flex-grow border-t border-gray-400" />
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  title: string
}>()
</script>