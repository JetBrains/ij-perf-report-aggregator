<template>
  <DataTable
    v-model:filters="filters"
    :value="metricsLists"
    striped-rows
    paginator
    :rows="100"
    :table-style="{ 'min-width': '50rem' }"
    class="p-datatable-sm"
    filter-display="row"
  >
    <Column
      field="name"
      header="Name"
      :sortable="true"
      :show-filter-menu="false"
    >
      <template #filter="{ filterModel, filterCallback }">
        <InputText
          v-model="filterModel.value"
          type="text"
          class="p-column-filter"
          placeholder="Search by name"
          @input="filterCallback()"
        />
      </template>
    </Column>
    <Column
      field="isMain"
      header="Main"
      :sortable="true"
      :show-filter-menu="false"
      :style="{ width: '9rem' }"
    >
      <template #body="slotProps">
        <span
          v-if="slotProps.data.isMain"
          class="inline-block rounded bg-blue-100 px-2 py-0.5 text-xs font-medium text-blue-800"
        >
          Main
        </span>
      </template>
      <template #filter="{ filterModel, filterCallback }">
        <div class="flex items-center gap-2">
          <Checkbox
            input-id="main-metrics-filter"
            :model-value="filterModel.value === true"
            binary
            @update:model-value="
              (v) => {
                filterModel.value = v ? true : null
                filterCallback()
              }
            "
          />
          <label
            for="main-metrics-filter"
            class="text-xs text-gray-600"
          >
            Only main
          </label>
        </div>
      </template>
    </Column>
    <Column
      field="description"
      header="Description"
    >
      <template #body="slotProps">
        <span v-if="slotProps.data.description">{{ slotProps.data.description }}</span>
        <span
          v-else
          class="italic text-gray-400"
        >
          (no description yet)
        </span>
      </template>
    </Column>
    <Column
      field="url"
      header="URL"
    >
      <template #body="slotProps">
        <a
          :href="slotProps.data.url"
          target="_blank"
          class="underline decoration-dotted hover:no-underline"
          rel="noopener noreferrer"
        >
          {{ slotProps.data.url }}
        </a>
      </template>
    </Column>
  </DataTable>
</template>

<script setup lang="ts">
import { FilterMatchMode } from "@primevue/core/api"
import { ref } from "vue"
import { metricsDescription } from "../../shared/metricsDescription"
import { MAIN_METRICS } from "../../util/mainMetrics"

const mainMetricsSet = new Set(MAIN_METRICS)

function isMainMetric(name: string): boolean {
  if (mainMetricsSet.has(name)) return true
  const hashIdx = name.indexOf("#")
  return hashIdx > 0 && mainMetricsSet.has(name.slice(0, hashIdx))
}

const filters = ref({
  name: { value: null, matchMode: FilterMatchMode.STARTS_WITH },
  isMain: { value: null, matchMode: FilterMatchMode.EQUALS },
})

interface Metric {
  name: string
  description: string
  url?: string
  isMain: boolean
}
const metricsLists: Metric[] = []
const seenNames = new Set<string>()
for (const [name, value] of metricsDescription) {
  const isMain = isMainMetric(name)
  seenNames.add(name)
  if (typeof value === "string") {
    metricsLists.push({
      name,
      description: value,
      isMain,
    })
  } else {
    metricsLists.push({
      name,
      description: value.description,
      url: value.url,
      isMain,
    })
  }
}
for (const name of mainMetricsSet) {
  if (!seenNames.has(name)) {
    metricsLists.push({
      name,
      description: "",
      isMain: true,
    })
  }
}
</script>
