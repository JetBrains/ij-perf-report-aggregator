<template>
  <Dialog
    v-model:visible="showStacktrace"
    modal
    header="Stacktrace"
    :style="{ width: '80vw' }"
  >
    <div class="flex items-center space-x-4 mb-4 mt-4">
      {{ accident?.stacktrace ?? "" }}
    </div>
    <template #footer>
      <div class="flex justify-end space-x-2">
        <Button
          label="Close"
          icon="pi pi-times"
          @click="(showStacktrace = false)"
        />
      </div>
    </template>
  </Dialog>
</template>
<script setup lang="ts">
import { Accident } from "../../../configurators/accidents/AccidentsConfigurator"

const { accident } = defineProps<{
  accident: Accident | null
}>()

const showStacktrace = defineModel<boolean>()
</script>
