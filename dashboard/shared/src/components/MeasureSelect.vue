<template>
  <MultiSelect
    v-model="value"
    :options="data"
    placeholder="Metrics"
  />
</template>
<script lang="ts">
import { defineComponent } from "vue"

import { MeasureConfigurator } from "../configurators/MeasureConfigurator"

export default defineComponent({
  name: "MeasureSelect",
  props: {
    configurator: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const configurator = props.configurator as MeasureConfigurator
    return {
      value: configurator.value,
      data: configurator.data,
    }
  },
})
</script>