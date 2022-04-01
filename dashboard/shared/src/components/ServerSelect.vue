<template>
  <Combobox
    id="combobox"
    v-model="value"
  >
    <div class="relative">
      <div class="flex gap-2 flex-row items-baseline">
        <ComboboxLabel class="block text-sm text-gray-700">
          Server:
        </ComboboxLabel>
        <div class="relative">
          <ComboboxInput
            class="w-full rounded-md border border-gray-300 bg-white py-2 pl-3 pr-12 shadow-sm focus:border-indigo-500
            focus:outline-none focus:ring-1 focus:ring-indigo-500 sm:text-sm"
            :display-value="server => server"
            @change="query = $event.target.value"
          />
          <ComboboxButton
            class="absolute inset-y-0 right-0 flex items-center rounded-r-md px-2 focus:outline-none"
          >
            <ChevronDownIcon
              class="h-5 w-5 text-gray-400"
              aria-hidden="true"
            />
          </ComboboxButton>
        </div>
      </div>
      <TransitionRoot
        leave="transition ease-in duration-100"
        leave-from="opacity-100"
        leave-to="opacity-0"
        @after-leave="query = ''"
      >
        <ComboboxOptions
          class="absolute z-10 mt-1 max-h-60 w-full overflow-auto rounded-md
            bg-white py-1 text-base shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none sm:text-sm"
        >
          <div
            v-if="filteredValues.length === 0 && query !== ''"
            class="cursor-default select-none relative py-2 px-4 text-gray-700"
          >
            Nothing found.
          </div>

          <ComboboxOption
            v-for="server in filteredValues"
            :key="server"
            v-slot="{ selected, active }"
            as="template"
            :value="server"
          >
            <li
              class="relative cursor-default select-none py-2 pl-3 pr-9"
              :class="{
                'text-white bg-indigo-600': active,
                'text-gray-900': !active,
              }"
            >
              <span
                class="block truncate"
                :class="{ 'font-semibold': selected, 'font-normal': !selected }"
              >
                {{ server }}
              </span>
              <span
                v-if="selected"
                class="absolute inset-y-0 right-0 flex items-center pr-4"
                :class="{ 'text-white': active, 'text-indigo-600': !active }"
              >
                <CheckIcon
                  class="w-5 h-5"
                  aria-hidden="true"
                />
              </span>
            </li>
          </ComboboxOption>
        </ComboboxOptions>
      </TransitionRoot>
    </div>
  </Combobox>
</template>

<script setup lang="ts">
import { Combobox, ComboboxButton, ComboboxInput, ComboboxLabel, ComboboxOption, ComboboxOptions } from "@headlessui/vue"
import { computed, ref } from "vue"
import { ServerConfigurator } from "../configurators/ServerConfigurator"

const props = defineProps({
  modelValue: {
    type: String,
    default: "",
  },
})

const emit = defineEmits(["update:modelValue"])

const suggestedServers: Array<string> = [ServerConfigurator.DEFAULT_SERVER_URL, "http://localhost:9044", "https://ij-perf-api.labs.jb.gg"]
const value = computed({
  get() {
    return props.modelValue
  },
  set(value: string) {
    return emit("update:modelValue", value)
  },
})

const query = ref("")
const filteredValues = computed(() => {
  if (query.value === "") {
    return suggestedServers
  }
  else {
    return suggestedServers.filter(server => server.replace(RegExp("http(s)?://"), "").toLowerCase().includes(query.value.toLowerCase()))
  }
})
</script>