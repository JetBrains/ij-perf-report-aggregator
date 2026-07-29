<template>
  <DashboardPage
    db-name="perfintDev"
    table="clion"
    persistent-id="clion_indexing_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    :with-installer="false"
  >
    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <section>
          <GroupProjectsChart
            :label="`Indexing CMake projects`"
            :measure="['backendIndexingTimeMs']"
            :projects="[
              'radler/llvm/indexing',
              'radler/opencv/indexing',
              'radler/curl/indexing',
              'radler/big_project_50k_10k/indexing',
              'radler/big_project_50k_10k_many_symbols/indexing',
            ]"
          />
        </section>
      </div>
    </section>

    <Divider title="Compile Commands projects" />

    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <section>
          <GroupProjectsChart
            :label="`Indexing Compile Commands projects`"
            :measure="['backendIndexingTimeMs']"
            :projects="['radler/llvm_compdb/indexing']"
          />
        </section>
      </div>
    </section>

    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <section>
          <GroupProjectsChart
            :label="`Indexing project (CMake LLVM vs compdb LLVM)`"
            :measure="['backendIndexingTimeMs']"
            :projects="['radler/llvm_compdb/indexing', 'radler/llvm/indexing']"
          />
        </section>
      </div>
    </section>

    <Divider title="Makefile projects" />

    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <section>
          <GroupProjectsChart
            :label="`Indexing Makefile projects`"
            :measure="['backendIndexingTimeMs']"
            :projects="['radler/git/indexing']"
          />
        </section>
      </div>
    </section>

    <Divider title="West projects" />

    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <section>
          <GroupProjectsChart
            :label="`Indexing West projects`"
            :measure="['backendIndexingTimeMs']"
            :projects="['radler/zephyr_bap_broadcast_sink/indexing']"
          />
        </section>
      </div>
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
import Divider from "../common/Divider.vue"
</script>
