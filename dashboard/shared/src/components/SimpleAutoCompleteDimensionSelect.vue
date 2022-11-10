<!-- https://tailwindui.com/components/ecommerce/components/category-filters#component-cb75aacd7c636865835476c98f606a38 -->
<template>
  <Popover
    as="div"
    class="relative inline-block text-left"
  >
    <div>
      <PopoverButton class="group inline-flex items-center justify-center text-sm font-medium text-gray-700 hover:text-gray-900">
        <span>{{ value?.label ?? label }}</span>
        <!--<span-->
        <!--  v-if="sectionIdx === 0"-->
        <!--  class="ml-1.5 rounded bg-gray-200 py-0.5 px-1.5 text-xs font-semibold tabular-nums text-gray-700"-->
        <!--&gt;1</span>-->
        <ChevronDownIcon
          class="-mr-1 ml-1 h-5 w-5 flex-shrink-0 text-gray-400 group-hover:text-gray-500"
          aria-hidden="true"
        />
      </PopoverButton>
    </div>

    <transition
      enter-active-class="transition ease-out duration-100"
      enter-from-class="transform opacity-0 scale-95"
      enter-to-class="transform opacity-100 scale-100"
      leave-active-class="transition ease-in duration-75"
      leave-from-class="transform opacity-100 scale-100"
      leave-to-class="transform opacity-0 scale-95"
    >
      <PopoverPanel class="absolute max-h-96 overflow-y-scroll left-0 z-10 mt-2 origin-top-right rounded-md bg-white p-4 shadow-2xl ring-1 ring-black ring-opacity-5 focus:outline-none">
        <form class="space-y-4">
          <div
            v-for="(option, optionIndex) in filteredItems"
            :key="option.value"
            class="flex items-center"
          >
            <input
              :id="`filter-${label}-${optionIndex}`"
              :name="`${label}[]`"
              :value="option.value"
              type="checkbox"
              class="h-4 w-4 rounded border-gray-300 text-indigo-600 focus:ring-indigo-500"
            >
            <label
              :for="`filter-${label}-${optionIndex}`"
              class="ml-3 whitespace-nowrap pr-6 text-sm font-medium text-gray-900"
            >{{ option.label }}</label>
          </div>
        </form>
      </PopoverPanel>
    </transition>
  </Popover>
</template>
<script setup lang="ts">
import { Popover, PopoverButton, PopoverPanel } from "@headlessui/vue"
import { ChevronDownIcon } from "@heroicons/vue/20/solid"
import { Option } from "tailwind-ui/src/tabModel"
import { computed, ref } from "vue"
import { DimensionConfigurator } from "../configurators/DimensionConfigurator"

const props = withDefaults(defineProps<{
  label: string
  dimension: DimensionConfigurator
  valueToLabel?: (v: string) => string
}>(), {
  valueToLabel: (v: string) => v,
  valueToGroup: null,
})

const items = computed<Array<Option>>(() => {
  const valueToLabel = props.valueToLabel
  return props.dimension.values.value.map(it => {
    return {label: valueToLabel(it.toString()), value: it as string}
  })
})

const query = ref("")

const filteredItems = computed(() =>
  query.value === ""
    ? items.value ?? []
    : items.value.filter(it =>
      it.label
        .toLowerCase()
        .replace(/\s+/g, "")
        .includes(query.value.toLowerCase().replace(/\s+/g, "")),
    ),
)

const value = computed<Option | null>({
  get(): Option | null {
    const values = props.dimension.values.value
    if (values == null || values.length === 0) {
      return null
    }

    const value = props.dimension.selected.value
    const normalizedValue = Array.isArray(value) ? value[0] : (value === "" ? null : value)
    return normalizedValue == null ? null : {label: props.valueToLabel(normalizedValue), value: normalizedValue}
  },
  set(value) {
    // eslint-disable-next-line vue/no-mutating-props
    props.dimension.selected.value = value?.value ?? null
  },
})

</script>