<template>
  <div class="flex items-center gap-2 text-sm">
    <ToggleSwitch v-model="model" />
    <span v-tooltip.bottom="'Includes data reported under pre-rename project names (e.g. ijent-import-intellij-Local). May draw long strokes for sparse historical series.'"
      >Pre-rename data</span
    >
  </div>
</template>

<script setup lang="ts">
const model = defineModel<boolean>({ required: true })
</script>

<style scoped></style>
