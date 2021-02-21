<template>
  <el-form-item :label="label">
    <el-select
      v-model="value"
      :loading="loading"
      filterable
    >
      <el-option
        v-for="item in items"
        :key="item.value"
        :label="item.label"
        :value="item.value"
      />
    </el-select>
  </el-form-item>
</template>
<script lang="ts">
import { computed, defineComponent } from "vue"
import { DimensionConfigurator } from "../configurators/DimensionConfigurator"

export default defineComponent({
  name: "DimensionSelect",
  props: {
    label: {
      type: String,
      required: true,
    },
    dimension: {
      type: Object,
      required: true,
    },
    valueLabel: {
      type: Function,
      default: null,
    },
  },
  setup(props) {
    // map Array<string> to Array<Item> to be able to customize how value is displayed in UI
    const valueToLabel = props.valueLabel ?? ((v: string) => v)
    return {
      value: computed({
        get: () => (props.dimension as DimensionConfigurator).value.value,
        set: value => {
          if (value.length !== 0) {
            (props.dimension as DimensionConfigurator).value.value = value
          }
        },
      }),
      items: computed(() => {
        return (props.dimension as DimensionConfigurator).values.value.map(it => {
          return {label: valueToLabel(it), value: it}
        })
      }),
      loading: (props.dimension as DimensionConfigurator).loading,
    }
  },
})
</script>