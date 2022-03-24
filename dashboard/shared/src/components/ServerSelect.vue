<template>
  Server:
  <AutoComplete
    v-model="value"
    placeholder="The stats server URL"
    dropdown
    auto-highlight
    :suggestions="filteredServer"
    @complete="searchServer($event)"
  />
</template>
<script lang="ts">
import { computed, defineComponent, ref } from "vue"

import { ServerConfigurator } from "../configurators/ServerConfigurator"

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
    const suggestedServers: Array<string> = [ServerConfigurator.DEFAULT_SERVER_URL, "http://localhost:9044", "https://ij-perf-api.labs.jb.gg"]
    const filteredServer = ref()
    return {
      value: computed({
        get() {
          return props.modelValue
        },
        set(value: string) {
          return emit("update:modelValue", value)
        },
      }),
      filteredServer,
      searchServer(event: Event): void {
        /* eslint-disable @typescript-eslint/no-unsafe-assignment */
        const queryString: string = event.query
        if (queryString == null || queryString.length === 0) {
          filteredServer.value = [...suggestedServers]
        }
        else {
          filteredServer.value = [...suggestedServers.filter(it => it.replace(RegExp("http(s)?://"), "").startsWith(queryString.toLowerCase()) && it !== queryString)]
        }
      },
    }
  },
})
</script>