<template>
  <DashboardPage
    db-name="fleet"
    table="measure_new"
    persistent-id="fleetPerf_dashboard"
    initial-machine="linux-blade-hetzner"
    :with-installer="false"
  >
    <Divider title="Core" />
    <section>
      <GroupProjectsChart
        label="Typing (time)"
        measure="fleet.test"
        :projects="['multiCaretTyping', 'stressEnter', 'stressTyping']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Open Project"
        measure="fleet.test"
        :projects="['openGoDelveProject', 'openJeditermProject', 'openRustSimpleServerProject', 'openSpringPetClinicJavaProject', 'openSpringPetClinicKotlinProject']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Typing (average delay)"
        measure="awt.delay"
        :projects="['multiCaretTyping', 'stressEnter', 'stressTyping']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Typing (max delay)"
        measure="max.awt.delay"
        :projects="['multiCaretTyping', 'stressEnter', 'stressTyping']"
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
    <Divider title="PHP" />
    <section>
      <GroupProjectsChart
        label="Typing (time)"
        measure="fleet.test"
        :projects="['Typing in mPDF', 'Typing in mPDF With Backend', 'Pressing Enter in mPDF']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Typing (average delay)"
        measure="awt.delay"
        :projects="['Typing in mPDF', 'Typing in mPDF With Backend', 'Pressing Enter in mPDF']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Typing (max delay)"
        measure="max.awt.delay"
        :projects="['Typing in mPDF', 'Typing in mPDF With Backend', 'Pressing Enter in mPDF']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Other"
        measure="fleet.test"
        :projects="['Open mPDF', 'Frontend Completion in mPDF']"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
import Divider from "../common/Divider.vue"
</script>
