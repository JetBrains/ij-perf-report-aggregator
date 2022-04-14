<template>
  <MultiSelect
    v-model="value"
    :options="data"
    placeholder="Metrics"
    :max-selected-labels="1"
  />
</template>
<script lang="ts">
import { computed, defineComponent } from "vue"

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
      value: computed({
        get(): Array<string>|null {
          if(!configurator.data.value.some(it => configurator.value.value != null && configurator.value.value.indexOf(it) > -1)){
            return null
          }
          return configurator.data.value.length == 0 ? null : configurator.value.value
        },
        set(value: Array<string>|null) {
          configurator.value.value = value
        },
      }),
      data: configurator.data,
    }
  },
})
</script>