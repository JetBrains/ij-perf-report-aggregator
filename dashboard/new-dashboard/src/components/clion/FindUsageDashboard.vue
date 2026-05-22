<template>
  <DashboardPage
    db-name="perfintDev"
    table="clion"
    persistent-id="clion_find_usage_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    :with-installer="false"
  >
    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Total Execution Time: Find Usages (luau, enumerable)"
          :measure="['%syncAction FindUsages', 'findUsagesInToolWindow']"
          :projects="['clion/luau/findUsages/enumerable (LuauOpcode)', 'radler/luau/findUsages/enumerable (LuauOpcode)']"
          :aliases="['clion/luau/findUsages/enumerable (LuauOpcode)', 'radler/luau/findUsages/enumerable (LuauOpcode)']"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Total Execution Time: Find Usages (luau, enumerator)"
          :measure="['%syncAction FindUsages', 'findUsagesInToolWindow']"
          :projects="['clion/luau/findUsages/enumerator (LOP_NOP)', 'radler/luau/findUsages/enumerator (LOP_NOP)']"
          :aliases="['clion/luau/findUsages/enumerator (LOP_NOP)', 'radler/luau/findUsages/enumerator (LOP_NOP)']"
        />
      </div>
    </section>
    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Total Execution Time: Find Usages (luau, class template)"
          :measure="['%syncAction FindUsages', 'findUsagesInToolWindow']"
          :projects="['clion/luau/findUsages/class template (DenseHashTable)', 'radler/luau/findUsages/class template (DenseHashTable)']"
          :aliases="['clion/luau/findUsages/class template (DenseHashTable)', 'radler/luau/findUsages/class template (DenseHashTable)']"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Total Execution Time: Find Usages (luau, macro)"
          :measure="['%syncAction FindUsages', 'findUsagesInToolWindow']"
          :projects="['clion/luau/findUsages/macro (LUAU_ASSERT)', 'radler/luau/findUsages/macro (LUAU_ASSERT)']"
          :aliases="['clion/luau/findUsages/macro (LUAU_ASSERT)', 'radler/luau/findUsages/macro (LUAU_ASSERT)']"
        />
      </div>
    </section>
    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Total Execution Time: Find Usages (cmake, class)"
          :measure="['%syncAction FindUsages', 'findUsagesInToolWindow']"
          :projects="['clion/cmake/findUsages/class (cmCTestResourceAllocator)', 'radler/cmake/findUsages/class (cmCTestResourceAllocator)']"
          :aliases="['clion/cmake/findUsages/class (cmCTestResourceAllocator)', 'radler/cmake/findUsages/class (cmCTestResourceAllocator)']"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Total Execution Time: Find Usages (cmake, macro)"
          :measure="['%syncAction FindUsages', 'findUsagesInToolWindow']"
          :projects="['clion/cmake/findUsages/macro (SAFEDIV)', 'radler/cmake/findUsages/macro (SAFEDIV)']"
          :aliases="['clion/cmake/findUsages/macro (SAFEDIV)', 'radler/cmake/findUsages/macro (SAFEDIV)']"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Total Execution Time: Find Usages (cmake, member)"
          :measure="['%syncAction FindUsages', 'findUsagesInToolWindow']"
          :projects="['clion/cmake/findUsages/member (SlotsNeeded)', 'radler/cmake/findUsages/member (SlotsNeeded)']"
          :aliases="['clion/cmake/findUsages/member (SlotsNeeded)', 'radler/cmake/findUsages/member (SlotsNeeded)']"
        />
      </div>
    </section>
    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="[Radler] Total Execution Time"
          :measure="['%syncAction FindUsages', 'findUsagesInToolWindow']"
          :projects="[
            'radler/luau/findUsages/enumerable (LuauOpcode)',
            'radler/luau/findUsages/enumerator (LOP_NOP)',
            'radler/luau/findUsages/class template (DenseHashTable)',
            'radler/luau/findUsages/macro (LUAU_ASSERT)',
            'radler/cmake/findUsages/class (cmCTestResourceAllocator)',
            'radler/cmake/findUsages/macro (SAFEDIV)',
            'radler/cmake/findUsages/member (SlotsNeeded)',
          ]"
          :aliases="[
            'radler/luau/findUsages/enumerable (LuauOpcode)',
            'radler/luau/findUsages/enumerator (LOP_NOP)',
            'radler/luau/findUsages/class template (DenseHashTable)',
            'radler/luau/findUsages/macro (LUAU_ASSERT)',
            'radler/cmake/findUsages/class (cmCTestResourceAllocator)',
            'radler/cmake/findUsages/macro (SAFEDIV)',
            'radler/cmake/findUsages/member (SlotsNeeded)',
          ]"
        />
      </div>
    </section>
    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="[Radler] Number of Usages"
          measure="findUsagesInToolWindow#number"
          :projects="[
            'radler/luau/findUsages/enumerable (LuauOpcode)',
            'radler/luau/findUsages/enumerator (LOP_NOP)',
            'radler/luau/findUsages/class template (DenseHashTable)',
            'radler/luau/findUsages/macro (LUAU_ASSERT)',
            'radler/cmake/findUsages/class (cmCTestResourceAllocator)',
            'radler/cmake/findUsages/macro (SAFEDIV)',
            'radler/cmake/findUsages/member (SlotsNeeded)',
          ]"
          :aliases="[
            'radler/luau/findUsages/enumerable (LuauOpcode)',
            'radler/luau/findUsages/enumerator (LOP_NOP)',
            'radler/luau/findUsages/class template (DenseHashTable)',
            'radler/luau/findUsages/macro (LUAU_ASSERT)',
            'radler/cmake/findUsages/class (cmCTestResourceAllocator)',
            'radler/cmake/findUsages/macro (SAFEDIV)',
            'radler/cmake/findUsages/member (SlotsNeeded)',
          ]"
        />
      </div>
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
</script>
