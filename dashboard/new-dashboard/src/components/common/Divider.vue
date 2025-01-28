<template>
  <div
    v-if="label == null"
    class="relative flex py-5 items-center"
  >
    <div class="grow border-t" />
    <span class="shrink mx-4 text-lg">{{ title }}</span>
    <div class="grow border-t" />
  </div>
  <!-- https://tailwindui.com/components/application-ui/layout/dividers -->
  <div class="v-else relative">
    <div
      class="absolute inset-0 flex items-center"
      aria-hidden="true"
    >
      <div class="w-full border-t" />
    </div>
    <div class="relative flex justify-center">
      <span class="bg-white dark:bg-[#121212] px-2">{{ label }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
export interface Props {
  title?: string
  label?: string
}

const { title, label } = defineProps<Props>()
</script>
