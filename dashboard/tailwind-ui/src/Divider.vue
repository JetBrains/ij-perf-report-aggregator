<template>
  <div class="relative mt-3 mb-2">
    <div
      class="absolute inset-0 flex items-center"
      aria-hidden="true"
    >
      <div class="w-full border-t border-gray-300" />
    </div>
    <div class="relative flex justify-center">
      <span
        v-if="label != null"
        class="px-2 bg-white text-sm text-gray-500 capitalize"
      >{{ props.label }}</span>
      <span
        v-else-if="title != null"
        class="px-3 bg-white text-lg font-medium text-gray-900"
      >{{ props.title }}</span>
    </div>
  </div>
</template>
<script setup lang="ts">
const props = defineProps<{
  label?: string
  title?: string
}>()
</script>