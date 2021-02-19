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
import { computed, defineComponent, PropType } from "vue"
import { DimensionConfigurator } from "../configurators/DimensionConfigurator"

export default defineComponent({
  name: "DimensionSelect",
  props: {
    label: {
      type: String,
      required: true,
    },
    dimension: {
      type: Object as PropType<DimensionConfigurator>,
      required: true,
    },
    valueLabel: {
      type: Function as PropType<(v: string) => string>,
      default: null,
    },
  },
  setup(props) {
    // map Array<string> to Array<Item> to be able to customize how value is displayed in UI
    const valueToLabel = props.valueLabel ?? function (v) { return v}
    return {
      value: computed<string>({
        get() {
          return props.dimension.value.value
        },
        set(value) {
          if (value.length !== 0) {
            // eslint-disable-next-line vue/no-mutating-props
            props.dimension.value.value = value
          }
        },
      }),
      items: computed(() => {
        return props.dimension.values.value.map(it => {
          return {label: valueToLabel(it), value: it}
        })
      }),
      loading: props.dimension.loading,
    }
  },
})
</script>