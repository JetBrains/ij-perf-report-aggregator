<template>
  <div
    v-if="label == null"
    class="relative flex py-5 items-center"
  >
    <div class="grow border-t" />
    <div class="shrink mx-4 flex items-center">
      <span class="text-lg">{{ title }}</span>
      <i
        v-if="description"
        v-tooltip="{ value: description, pt: { text: { class: 'text-base w-md' } } }"
        class="pi pi-info-circle text-sm cursor-help ml-2"
      />
    </div>
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
  description?: string
}

const { title = null, label = null, description = null } = defineProps<Props>()
</script>
