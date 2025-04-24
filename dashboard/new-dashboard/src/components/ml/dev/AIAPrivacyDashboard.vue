<template>
  <DashboardPage
    db-name="perfintDev"
    table="ml"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    persistent-id="aiaPrivacyDev"
    :with-installer="false"
  >
    <section>
      <GroupProjectsChart
        label="Context Privacy Filters mean"
        :measure="[
          'contextPrivacyFilter.AiIgnoreContextPrivacyFilter.mean.ns',
          'contextPrivacyFilter.VcsIgnoredContextPrivacyFilter.mean.ns',
          'contextPrivacyFilter.VcsUnversionedContextPrivacyFilter.mean.ns',
        ]"
        :projects="['kotlinx_coroutines_k2_dev_ContextPrivacyFiltersTest/coroutinesProject-context-privacy-filters-all-files']"
      ></GroupProjectsChart>
    </section>
    <section>
      <GroupProjectsChart
        label="Context Privacy Filters sum"
        :measure="[
          'ai-ignore.sum.ms',
          'contextPrivacyFilter.AiIgnoreContextPrivacyFilter.sum.ms',
          'contextPrivacyFilter.VcsIgnoredContextPrivacyFilter.sum.ms',
          'contextPrivacyFilter.VcsUnversionedContextPrivacyFilter.sum.ms',
        ]"
        :projects="['kotlinx_coroutines_k2_dev_ContextPrivacyFiltersTest/coroutinesProject-context-privacy-filters-all-files']"
      ></GroupProjectsChart>
    </section>
    <section>
      <GroupProjectsChart
        label="Context Privacy Filters ignored count"
        :measure="[
          'contextPrivacyFilter.AiIgnoreContextPrivacyFilter.ignored.count',
          'contextPrivacyFilter.VcsIgnoredContextPrivacyFilter.ignored.count',
          'contextPrivacyFilter.VcsUnversionedContextPrivacyFilter.ignored.count',
        ]"
        :projects="['kotlinx_coroutines_k2_dev_ContextPrivacyFiltersTest/coroutinesProject-context-privacy-filters-all-files']"
      ></GroupProjectsChart>
    </section>
    <section>
      <GroupProjectsChart
        label="AI Ignore Single"
        :measure="['ai-ignore.mean.ns', 'ai-ignore.median.ns', 'ai-ignore.95.percentile.ns']"
        :projects="['kotlinx_coroutines_k2_dev_ContextPrivacyFiltersTest/coroutinesProject-context-privacy-filters-all-files']"
      ></GroupProjectsChart>
    </section>
  </DashboardPage>
</template>
<script setup lang="ts">
import GroupProjectsChart from "../../charts/GroupProjectsChart.vue"
import DashboardPage from "../../common/DashboardPage.vue"
</script>
