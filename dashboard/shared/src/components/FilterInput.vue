<template>
  <InputText
    id="filter"
    v-model="filter"
    type="text"
    placeholder="Filter"
  />
</template>

<script lang="ts">
import { computed, defineComponent } from "vue"
import { DimensionConfigurator } from "../configurators/DimensionConfigurator"

export default defineComponent({
  name: "FilterInput",
  props: {
    configurator: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const configurator = props.configurator as DimensionConfigurator
    return {
      filter: computed({
        get() {
          return configurator.filter.value
        },
        set(value: string) {
          configurator.filter.value = value
        },
      }),
    }
  },
})
</script>