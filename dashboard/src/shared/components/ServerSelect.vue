<template>
  <label class="text-sm">Server:
    <AutoComplete
      v-model="value"
      class="small"
      placeholder="The stats server URL"
      dropdown
      auto-highlight
      :item-select="itemSelected"
      :suggestions="filteredServer"
      @complete="searchServer($event)"
    />
  </label>
</template>

<script setup lang="ts">
import { AutoCompleteCompleteEvent } from "primevue/autocomplete"
import { debounceTime, distinctUntilChanged, Subject } from "rxjs"
import { computed, ref } from "vue"
import { ServerConfigurator } from "../configurators/ServerConfigurator"

const props = defineProps<{
  modelValue: string
}>()

const emit = defineEmits(["update:modelValue"])

const subjectWithDebounce = new Subject<string>()
subjectWithDebounce
  .pipe(
    distinctUntilChanged(),
    // typing of server may take longer time
    debounceTime(1_000),
  )
  .subscribe(value => {
    subject.next(value)
  })

const subject = new Subject<string>()
subject
  .pipe(
    distinctUntilChanged(),
  )
  .subscribe(value => {
    emit("update:modelValue", value)
  })

const suggestedServers: Array<string> = [ServerConfigurator.DEFAULT_SERVER_URL, "http://localhost:9044", "https://ij-perf-api.labs.jb.gg"]
const filteredServer = ref<Array<string>>([])
const value = computed({
  get() {
    return props.modelValue
  },
  set(value: string) {
    subjectWithDebounce.next(value)
  },
})

function itemSelected(event: { value: string }): void {
  // no delay on explicit item selection - update server immediately
  subject.next(event.value)
}

const httpRegex = new RegExp("http(s)?://")

function searchServer(event: AutoCompleteCompleteEvent): void {
  const queryString = event.query
  filteredServer.value = queryString == null || queryString.length === 0 ? [...suggestedServers]
    : [...suggestedServers.filter(it => it.replace(httpRegex, "").startsWith(queryString.toLowerCase()) && it !== queryString)]
}
</script>