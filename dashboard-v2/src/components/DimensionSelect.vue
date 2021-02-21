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
import { defineComponent, ref, watch } from "vue"
import { DimensionConfigurator, Item } from "@/configurators/DimensionConfigurator"

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
    const dimension = props.dimension as DimensionConfigurator
    const value = ref(dimension.value.value)

    const valueToLabel = props.valueLabel ?? ((v: string) => v)
    watch(value, it => {
      dimension.value.value = it
    })
    watch(dimension.value, it => {
      value.value = it
    })

    const items = ref<Array<Item>>([])
    watch(dimension.values, rawValues => {
      items.value = rawValues.map(it => {
        return {label: valueToLabel(it), value: it}
      })
    })
    return {
      value,
      items,
      loading: dimension.loading,
    }
  },
})
</script>