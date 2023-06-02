<template>
  <MultiSelect
    v-model="branchValue"
    title="Branch"
    :loading="branchConfigurator.state.loading"
    :disabled="branchConfigurator.state.disabled"
    :options="branchItems"
    :placeholder="placeholder"
    option-label="label"
    option-value="value"
    :show-toggle-all="false"
    panel-class="w-fit"
    panel-style="overflow: visible"
    @hide="clearSubMenu"
  >
    <template #value="slotProps">
      <div class="group flex items-center gap-1">
        <div class="w-4 h-4 text-gray-500">
          <BranchIcon />
        </div>

        <span v-if="!slotProps.value || slotProps.value.length === 0">
          {{ placeholder }}
        </span>

        <span v-if="slotProps.value && slotProps.value.length === 1">
          {{ slotProps.value[0] }}
        </span>

        <span v-if="slotProps.value && slotProps.value.length > 1">
          {{ branchesSelectLabelFormat(slotProps.value) }}
        </span>

        <ChevronDownIcon
          class="-mr-1 ml-1 h-5 w-5 flex-shrink-0 text-gray-400 group-hover:text-gray-500"
          aria-hidden="true"
        />
      </div>
    </template>
    <template #footer>
      <div class="border-t border-solid border-neutral-200 relative">
        <ul class="p-multiselect-items p-component">
          <li
            v-if="versionItems!== undefined && versionItems.length > 0"
            class="p-multiselect-item flex items-center gap-2"
            @click="openVersionSubmenu"
          >
            <span class="flex items-center gap-1 overflow-hidden">
              Version type
              <span
                v-if="versionValue !== null && versionValue.length > 0"
                class="text-gray-500 truncate"
              >
                {{ versionValue?.length < 2 ? versionValue[0] : `Selected ${versionValue?.length }` }}
              </span>
            </span>
            <span class="pi pi-angle-right ml-[auto]" />
          </li>
          <li
            v-if="triggeredItems.length > 0"
            class="p-multiselect-item flex items-center gap-2"
            @click="openTriggeredSubmenu"
          >
            <span class="flex items-center gap-1 overflow-hidden">
              Triggered by
              <span
                v-if="triggeredValue !== null && triggeredValue.length > 0"
                class="text-gray-500 truncate"
              >
                {{ triggeredValue?.length < 2 ? triggeredValue[0] : `Selected ${triggeredValue?.length }` }}
              </span>
            </span>
            <span class="pi pi-angle-right ml-[auto]" />
          </li>
        </ul>

        <div
          v-if="activeSubMenu === SubMenu.VERSION_TYPE && versionItems.length > 0"
          class="absolute p-dropdown-panel p-component w-[270px] max-h-[200px] branch-select-dropdown"
          style="left: 100%;"
        >
          <ul class="p-multiselect-items p-component">
            <li
              v-for="item in versionItems"
              :key="item.label"
            >
              <div class="flex items-center p-multiselect-item p-component">
                <Checkbox
                  v-model="versionValue"
                  :value="item.value"
                  :input-id="item.value"
                  class="field-checkbox"
                />
                <label
                  class="w-full inline-block"
                  :for="item.value"
                >
                  <span>{{ item.label }}</span>
                </label>
              </div>
            </li>
          </ul>
        </div>

        <div
          v-if="activeSubMenu === SubMenu.TRIGGERED_BY && triggeredItems.length > 0"
          class="absolute p-dropdown-panel p-component w-[270px] max-h-[200px] branch-select-dropdown"
          style="left: 100%;"
        >
          <ul class="p-multiselect-items p-component">
            <div v-if="triggeredItems.length === 0">
              No available options
            </div>
            <li
              v-for="item in triggeredItems"
              :key="item.label"
            >
              <div class="flex items-center p-multiselect-item p-component">
                <Checkbox
                  v-model="triggeredValue"
                  :value="item.value"
                  :input-id="item.value"
                  class="field-checkbox"
                />
                <label
                  class="w-full inline-block"
                  :for="item.value"
                >
                  <span>{{ item.label }}</span>
                </label>
              </div>
            </li>
          </ul>
        </div>
      </div>
    </template>
    <template #dropdownicon>
      <span class="hidden" />
    </template>
  </MultiSelect>
</template>
<script setup lang="ts">
import { ChevronDownIcon } from "@heroicons/vue/20/solid"
import { usePlaceholder } from "shared/src/components/placeholder"
import { DimensionConfigurator } from "shared/src/configurators/DimensionConfigurator"
import { computed, ref } from "vue"
import { branchesSelectLabelFormat } from "../../shared/labels"
import BranchIcon from "./BranchIcon.vue"

interface Props {
  branchConfigurator: DimensionConfigurator
  releaseConfigurator?: DimensionConfigurator
  triggeredByConfigurator?: DimensionConfigurator
}

const enum SubMenu {
  TRIGGERED_BY,
  VERSION_TYPE,
}

const activeSubMenu = ref<SubMenu | null>(null)

const openVersionSubmenu = () => {
  activeSubMenu.value = SubMenu.VERSION_TYPE
}

const openTriggeredSubmenu = () => {
  activeSubMenu.value = SubMenu.TRIGGERED_BY
}

const clearSubMenu = () => {
  activeSubMenu.value = null
}

const props = defineProps<Props>()

function createItems(configurator?: DimensionConfigurator) {
  return computed(() => {
    if (configurator == undefined) {
      return []
    }
    const values = configurator.values.value as string[]

    return values.map(it => {
      return {label: it.toString(), value: it}
    })
  })
}

function createValueFrom(configurator?: DimensionConfigurator) {
  return computed<string[] | null>({
    get() {
      if(configurator == null){
        return null
      }
      const values = configurator.values.value
      if (values.length === 0) {
        return null
      }

      const value = configurator.selected.value

      if (Array.isArray(value)) {
        return value
      }

      return value == null || value === "" ? [] : [value]
    },
    set(value) {
      if(configurator == null) return
      // eslint-disable-next-line vue/no-mutating-props
      configurator.selected.value = value == null || value.length === 0 ? null : value
    },
  })
}

const branchValue = createValueFrom(props.branchConfigurator)
const versionValue = createValueFrom(props.releaseConfigurator)
const triggeredValue = createValueFrom(props.triggeredByConfigurator)

const branchItems = createItems(props.branchConfigurator)
const versionItems = createItems(props.releaseConfigurator)
const triggeredItems = createItems(props.triggeredByConfigurator)

const placeholder = usePlaceholder(
  {label: "Branch"},
  () => props.branchConfigurator.values.value,
  () => props.branchConfigurator.selected.value,
)
</script>
<style>
.branch-select-dropdown {
  top: 0;
  margin-top: 0;
  border-top-left-radius: 0;
}
</style>