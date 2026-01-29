<script setup lang="ts">
import { inject, onMounted } from "vue"
import { SimpleMeasureConfigurator } from "../../configurators/SimpleMeasureConfigurator"
import { persistenceForDashboardKey } from "../../shared/keys"

interface Props {
  configurator: SimpleMeasureConfigurator
  data: string[]
}

const props = defineProps<Props>()

// This component is rendered as a child of DashboardPage, so it can inject from DashboardPage
const persistentStateManager = inject(persistenceForDashboardKey)

onMounted(() => {
  if (persistentStateManager) {
    props.configurator.registerWithPersistentStateManager(persistentStateManager)
  }
  // Initialize data after registering, so saved values are preserved
  props.configurator.initData(props.data)
})
</script>

<template>
  <div style="display: none"></div>
</template>
