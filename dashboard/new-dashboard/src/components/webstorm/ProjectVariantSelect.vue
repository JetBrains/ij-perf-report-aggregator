<template>
  <MultiSelect
    v-model="selectedVariants"
    :options="props.variantOptions"
    title="Variant"
    option-label="label"
    option-value="value"
    placeholder="All variants"
  >
    <template #value="slotProps">
      <div class="group flex items-center gap-1">
        <CubeIcon class="w-4 h-4" />

        <span v-if="!slotProps.value || slotProps.value.length === 0">
          {{ slotProps.placeholder }}
        </span>

        <span v-if="slotProps.value && slotProps.value.length === 1">
          {{ labelByValue.get(slotProps.value[0]) }}
        </span>

        <span v-if="slotProps.value && slotProps.value.length > 1">{{ slotProps.value.length }} variants</span>

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
import { ChevronDownIcon, CubeIcon } from "@heroicons/vue/20/solid/index"
import { ref, watch } from "vue"

interface VariantOption {
  label: string
  value: string
}

interface Props {
  variantOptions: VariantOption[]
}

const props = defineProps<Props>()

const selectedVariants = ref<string[]>([])

const emit = defineEmits<{ "update:selectedVariants": [selectedVariants: string[]] }>()

watch(selectedVariants, (newValue) => {
  emit("update:selectedVariants", newValue)
})

const labelByValue = new Map(props.variantOptions.map((v) => [v.value, v.label]))
</script>
