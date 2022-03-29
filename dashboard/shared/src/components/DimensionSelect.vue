<template>
  <template v-if="valueToGroup == null">
    <MultiSelect
      v-model="value"
      :loading="loading"
      :options="items"
      :placeholder="label"
      option-label="label"
      option-value="value"
      :max-selected-labels="3"
      :selection-limit="multiple ? null : 1"
      :filter="true"
    />
  </template>
  <template v-else>
    <MultiSelect
      v-model="value"
      :loading="loading"
      :options="items"
      :placeholder="label"
      option-label="label"
      option-value="value"
      option-group-children="options"
      option-group-label="label"
      :selection-limit="multiple ? null : 1"
      :max-selected-labels="3"
      :filter="true"
    />
  </template>
</template>
<script lang="ts">
import { computed, defineComponent, PropType } from "vue"
import { BaseDimensionConfigurator } from "../configurators/DimensionConfigurator"

export default defineComponent({
  name: "DimensionSelect",
  props: {
    label: {
      type: String,
      required: true,
    },
    dimension: {
      type: Object as PropType<BaseDimensionConfigurator>,
      required: true,
    },
    valueToLabel: {
      type: Function as PropType<(v: string) => string>,
      default: null,
    },
    // todo not working correctly for now (if value is set to not existing value, runtime error on select)
    valueToGroup: {
      type: Function as PropType<(v: string) => string>,
      default: null,
    },
  },
  setup(props) {
    // map Array<string> to Array<Item> to be able to customize how value is displayed in UI
    return {
      multiple: props.dimension.multiple,
      value: computed<string | Array<string>|null|Item>({
        get() {
          const value = props.dimension.value.value
          if (props.dimension.multiple && !Array.isArray(value)) {
            return value == null || value === "" ? [] : [value]
          }
          else {
            if(typeof value != "object") {
              return {label: value, value}
            } else {
              return value
            }
          }
        },
        set(value) {
          if (value == null || value.length === 0) {
            // eslint-disable-next-line vue/no-mutating-props
            props.dimension.value.value = null
          }
          else {
            // eslint-disable-next-line vue/no-mutating-props
            props.dimension.value.value = value
          }
        },
      }),
      items: computed(() => {
        const valueToLabel = props.valueToLabel ?? function (v) {
          return v
        }

        const values = props.dimension.values.value
        if (props.valueToGroup != null) {
          return group(values, props.valueToGroup, valueToLabel)
        }
        else {
          return values.map(it => {
            return {label: valueToLabel(it), value: it}
          })
        }
      }),
      loading: props.dimension.loading,
    }
  },
})

interface Item {
  label: string
  value: string
}

interface GroupItem {
  label: string
  options: Array<Item>
}

function group(values: Array<string>, groupFunction: (v: string) => string, valueToLabel: (v: string) => string): Array<GroupItem> {
  const groupNameToGroup = new Map<string, GroupItem>()
  const groups: Array<GroupItem> = []
  for (const value of values) {
    const groupName = groupFunction(value)
    let group = groupNameToGroup.get(groupName)
    if (group === undefined) {
      group = {
        label: groupName,
        options: [],
      }
      groupNameToGroup.set(groupName, group)
      groups.push(group)
    }
    group.options.push({label: valueToLabel(value), value})
  }
  console.log(JSON.stringify(groups, null, 2))
  return groups
}
</script>