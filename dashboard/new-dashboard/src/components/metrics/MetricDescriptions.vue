<template>
  <DataTable
    v-model:filters="filters"
    :value="metricsLists"
    striped-rows
    paginator
    :rows="30"
    table-style="min-width: 50rem"
    class="p-datatable-sm"
    filter-display="row"
  >
    <Column
      field="name"
      header="Name"
      :sortable="true"
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
      field="description"
      header="Description"
    ></Column>
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
import { FilterMatchMode } from "primevue/api"
import { ref } from "vue"
import { metricsDescription } from "../../shared/metricsDescription"

const filters = ref({
  name: { value: null, matchMode: FilterMatchMode.STARTS_WITH },
})

interface Metric {
  name: string
  description: string
  url?: string
}
const metricsLists: Metric[] = []
for (const [name, value] of metricsDescription) {
  if (typeof value === "string") {
    metricsLists.push({
      name,
      description: value,
    })
  } else {
    metricsLists.push({
      name,
      description: value.description,
      url: value.url,
    })
  }
}
</script>
