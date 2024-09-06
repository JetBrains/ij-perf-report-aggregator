<template>
  <DashboardPage
    db-name="fleet"
    table="measure_new"
    persistent-id="fleetPerf_dashboard"
    initial-machine="linux-blade-hetzner"
    :with-installer="false"
  >
    <Divider title="Typing" />
    <section>
      <GroupProjectsChart
        label="Typing latency p99"
        measure="p99"
        :projects="['stressTyping', 'stressLongFileInsertion', 'javascriptRetyping']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Typing latency p50 (technical-only metric)"
        measure="p50"
        :projects="['stressTyping', 'stressLongFileInsertion', 'javascriptRetyping']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Typing (max delay, technical-only metric)"
        measure="max.awt.delay"
        :projects="['multiCaretTyping', 'stressEnter', 'stressTyping', 'Typing in mPDF', 'Typing in mPDF With Backend', 'Pressing Enter in mPDF']"
      />
    </section>

    <section>
      <GroupProjectsChart
        label="Typing (total time)"
        measure="fleet.test"
        :projects="['multiCaretTyping', 'stressEnter', 'stressTyping', 'Typing in mPDF', 'Typing in mPDF With Backend', 'Pressing Enter in mPDF']"
      />
    </section>

    <Divider title="Core" />

    <section>
      <GroupProjectsChart
        label="Open Project"
        measure="fleet.test"
        :projects="['openGoDelveProject', 'openJeditermProject', 'openRustSimpleServerProject', 'openSpringPetClinicJavaProject', 'openSpringPetClinicKotlinProject']"
      />
    </section>


    <section>
      <GroupProjectsChart
        label="Highlighting"
        measure="fleet.test"
        :projects="['stressHighlighting']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Tree"
        measure="fleet.test"
        :projects="['wideTree', 'deepTree']"
      />
    </section>

    <section>
      <GroupProjectsChart
        label="Other"
        measure="fleet.test"
        :projects="['Open mPDF', 'Frontend Completion in mPDF']"
      />
    </section>

    <Divider title="Deprecated / Useless" />
    <section>
      <GroupProjectsChart
        label="Typing (average delay)"
        measure="awt.delay"
        :projects="['multiCaretTyping', 'stressEnter', 'stressTyping']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Typing (average delay)"
        measure="awt.delay"
        :projects="['Typing in mPDF', 'Typing in mPDF With Backend', 'Pressing Enter in mPDF']"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
import Divider from "../common/Divider.vue"
</script>
