<template>
  <MultiSelect
    v-model="selectedProjectCategories"
    :options="filteredCategories"
    title="Project"
    option-label="label"
    option-value="prefix"
    placeholder="All projects"
  >
    <template #value="slotProps">
      <div class="group flex items-center gap-1">
        <CubeIcon class="w-4 h-4" />

        <span v-if="!slotProps.value || slotProps.value.length === 0">
          {{ slotProps.placeholder }}
        </span>

        <span v-if="slotProps.value && slotProps.value.length === 1">
          {{ labelByPrefix.get(slotProps.value[0]) }}
        </span>

        <span v-if="slotProps.value && slotProps.value.length > 1">{{ slotProps.value.length }} projects</span>

        <ChevronDownIcon
          class="-mr-1 ml-1 h-5 w-5 shrink-0"
          aria-hidden="true"
        />
      </div>
    </template>
    <template #dropdownicon>
      <span />
    </template>
  </MultiSelect>
</template>

<script setup lang="ts">
import { ChevronDownIcon } from "@heroicons/vue/20/solid/index"
import { ref, watch } from "vue"
import { PROJECT_CATEGORIES } from "../projects"

interface Props {
  initialProjectCategories: string[]
}

const { initialProjectCategories } = defineProps<Props>()

const selectedProjectCategories = ref(initialProjectCategories)

const emit = defineEmits<{ "update:selectedProjectCategories": [selectedCategories: string[]] }>()

watch(selectedProjectCategories, (newValue) => {
  emit("update:selectedProjectCategories", newValue)
})
const filteredCategories = Object.values(PROJECT_CATEGORIES).filter((c) => c.prefix.length > 0)
const labelByPrefix = new Map(filteredCategories.map((c) => [c.prefix, c.label]))
</script>
