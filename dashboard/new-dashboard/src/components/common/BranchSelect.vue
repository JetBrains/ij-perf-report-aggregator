<template>
  <MultiSelect
    @hide="clearSubMenu"
    v-model="branchValue"
    title="Branch"
    :loading="branchConfigurator.state.loading"
    :disabled="branchConfigurator.state.disabled"
    :options="branchItems"
    :placeholder="placeholder"
    option-label="label"
    option-value="value"
    :show-toggle-all="false"
    panel-class="w-[270px]"
    panel-style="overflow: visible"
  >
    <template #value="slotProps">
      <span v-if="!slotProps.value || slotProps.value.length === 0" class="flex items-center gap-1 ">
         <div class="w-4 h-4 text-gray-500">
           <BranchIcon />
         </div>
          {{ placeholder }}
      </span>
      <span v-if="slotProps.value && slotProps.value.length === 1" class="flex items-center gap-1">
         <div class="w-4 h-4 text-gray-500">
           <BranchIcon/>
         </div>
         {{ slotProps.value[0] }}
      </span>
      <span v-if="slotProps.value && slotProps.value.length > 1" class="flex items-center gap-1">
        <div class="w-4 h-4 text-gray-500">
          <BranchIcon/>
        </div>
        {{ branchesSelectLabelFormat(slotProps.value) }}
      </span>
    </template>
    <template #footer="slotProps">
      <div class="border-t border-solid border-neutral-200 relative">
        <ul class="p-multiselect-items p-component">
          <li
            class="p-multiselect-item flex items-center"
            @click="openVersionSubmenu"
            v-if="versionItems.length > 0"
          >
            Version type
            <span class="pi pi-angle-right ml-[auto]" />
          </li>
          <li
            class="p-multiselect-item flex items-center"
            @click="openTriggeredSubmenu"
            v-if="triggeredItems.length > 0"
          >
            Triggered by
            <span class="pi pi-angle-right ml-[auto]" />
          </li>
        </ul>

        <div
          v-if="activeSubMenu === SubMenu.VERSION_TYPE && versionItems.length > 0"
          class="absolute p-dropdown-panel p-component w-[270px] max-h-[200px] branch-select-dropdown"
        >
          <ul class="p-multiselect-items p-component">
            <li v-for="item in versionItems">
              <label class="field-checkbox w-full p-multiselect-item p-component" :for="item.value">
                <Checkbox
                  :value="item.value"
                  :input-id="item.value"
                  v-model="versionValue"
                />
                <span>{{ item.label }}</span>
              </label>
            </li>
          </ul>
        </div>

        <div
          v-if="activeSubMenu === SubMenu.TRIGGERED_BY && triggeredItems.length > 0"
          class="absolute p-dropdown-panel p-component w-[270px] max-h-[200px] branch-select-dropdown"
        >
          <ul class="p-multiselect-items p-component">
            <div v-if="triggeredItems.length === 0">
              No available options
            </div>
            <li v-for="item in triggeredItems">
              <label class="field-checkbox w-full p-multiselect-item p-component" :for="item.value">
                <Checkbox
                  :value="item.value"
                  :input-id="item.value"
                  v-model="triggeredValue"
                />
                <span>{{ item.label }}</span>
              </label>
            </li>
          </ul>
        </div>
      </div>
    </template>
  </MultiSelect>
</template>
<style>
.branch-select-dropdown {
  top: 0;
  left: 100%;
  margin-top: 0;
  border-top-left-radius: 0;
}
</style>
<script setup lang="ts">
import { DimensionConfigurator } from "shared/src/configurators/DimensionConfigurator"
import { computed, ref } from "vue"
import { usePlaceholder } from "shared/src/components/placeholder"
import { branchesSelectLabelFormat } from "../../shared/labels"
import BranchIcon from "./BranchIcon.vue"

interface Props {
  branchConfigurator: DimensionConfigurator
  releaseConfigurator: DimensionConfigurator
  triggeredByConfigurator: DimensionConfigurator
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

function createItems(configurator: DimensionConfigurator) {
  return computed(() => {
    const values = configurator.values.value

    return values.map(it => {
      return {label: it.toString(), value: it}
    })
  })
}

function createValueFrom(configurator: DimensionConfigurator) {
  return computed<Array<string> | null>({
    get() {
      const values = configurator.values.value
      if (values == null || values.length === 0) {
        return null
      }

      const value = configurator.selected.value

      if (Array.isArray(value)) {
        return value
      }

      return value == null || value === "" ? [] : [value]
    },
    set(value) {
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