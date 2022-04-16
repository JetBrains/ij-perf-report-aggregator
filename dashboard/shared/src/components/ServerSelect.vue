<template>
  <label>Server:
    <AutoComplete
      id="serverSelector"
      v-model="value"
      class="small"
      placeholder="The stats server URL"
      dropdown
      auto-highlight
      :suggestions="filteredServer"
      @complete="searchServer($event)"
    />
  </label>
</template>

<script setup lang="ts">
import { AutoCompleteCompleteEvent } from "primevue/autocomplete"
import { computed, ref } from "vue"

import { ServerConfigurator } from "../configurators/ServerConfigurator"

const props = defineProps<{
  modelValue: string
}>()

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

function searchServer(event: AutoCompleteCompleteEvent): void {
  const queryString = event.query
  if (queryString == null || queryString.length === 0) {
    filteredServer.value = [...suggestedServers]
  }
  else {
    filteredServer.value = [...suggestedServers.filter(it => it.replace(RegExp("http(s)?://"), "").startsWith(queryString.toLowerCase()) && it !== queryString)]
  }
}
</script>