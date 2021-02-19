<template>
  <el-autocomplete
    v-model="value"
    class="inline-input"
    data-lpignore="true"
    placeholder="The stats server URL"
    size="small"
    :fetch-suggestions="getSuggestions"
  >
    <template #prepend>
      Server
    </template>
  </el-autocomplete>
</template>
<script lang="ts">
import { computed, defineComponent } from "vue"

import { ServerConfigurator } from "../configurators/ServerConfigurator"

interface Item {
  value: string
}

export default defineComponent({
  name: "ServerSelect",
  props: {
    modelValue: {
      type: String,
      default: "",
    },
  },
  emits: ["update:modelValue"],
  setup(props, {emit}) {
    const suggestedServers: Array<Item> = [{value: ServerConfigurator.DEFAULT_SERVER_URL}, {value: "http://localhost:9044"}]
    return {
      value: computed({
        get() {
          return props.modelValue
        },
        set(value: string) {
          return emit("update:modelValue", value)
        },
      }),
      getSuggestions(queryString: string | null, callback: (result: Array<Item>) => void) {
        if (queryString == null || queryString.length === 0) {
          callback(suggestedServers)
        }
        else {
          const q = queryString.toLowerCase()
          callback(suggestedServers.filter(it => it.value.startsWith(q) && it.value !== queryString))
        }
      }
    }
  },
})
</script>