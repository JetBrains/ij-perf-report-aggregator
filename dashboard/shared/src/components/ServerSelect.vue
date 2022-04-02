<template>
  <!-- extra wrapping div is used, otherwise warning https://forum.vuejs.org/t/vue-3-how-do-i-use-custom-directive-on-component/105227/3-->
  <div>
    Server:
    <AutoComplete
      v-model="value"
      class="small"
      placeholder="The stats server URL"
      dropdown
      auto-highlight
      :suggestions="filteredServer"
      @complete="searchServer($event)"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, defineEmits, defineProps, ref } from "vue"

import { ServerConfigurator } from "../configurators/ServerConfigurator"

interface Event{
  query: string
}

const props = defineProps({
  modelValue: {
    type: String,
    default: "",
  },
})

const emit = defineEmits(["update:modelValue"])

const suggestedServers: Array<string> = [ServerConfigurator.DEFAULT_SERVER_URL, "http://localhost:9044", "https://ij-perf-api.labs.jb.gg"]
const filteredServer = ref<Array<string>>([])
const value = computed({
  get() {
    return props.modelValue
  },
  set(value: string) {
    return emit("update:modelValue", value)
  },
})

function searchServer(event: Event): void {
  const queryString = event.query
  if (queryString == null || queryString.length === 0) {
    filteredServer.value = [...suggestedServers]
  }
  else {
    filteredServer.value = [...suggestedServers.filter(it => it.replace(RegExp("http(s)?://"), "").startsWith(queryString.toLowerCase()) && it !== queryString)]
  }
}
</script>