<template>
  <DashboardPage
    v-slot="{ dashboardConfigurators, warnings}"
    db-name="fleet"
    table="measure"
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
        value-unit="ns"
        :configurators="dashboardConfigurators"
        :accidents="warnings"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Open Project"
        measure="fleet.test"
        value-unit="ns"
        :projects="['openGoDelveProject', 'openJeditermProject', 'openRustSimpleServerProject', 'openSpringPetClinicJavaProject', 'openSpringPetClinicKotlinProject']"
        :configurators="dashboardConfigurators"
        :accidents="warnings"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Typing (average delay)"
        measure="awt.delay"
        value-unit="ns"
        :projects="['multiCaretTyping', 'stressEnter', 'stressTyping']"
        :configurators="dashboardConfigurators"
        :accidents="warnings"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Typing (max delay)"
        measure="max.awt.delay"
        value-unit="ns"
        :projects="['multiCaretTyping', 'stressEnter', 'stressTyping']"
        :configurators="dashboardConfigurators"
        :accidents="warnings"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Highlighting"
        value-unit="ns"
        measure="fleet.test"
        :projects="['stressHighlighting']"
        :configurators="dashboardConfigurators"
        :accidents="warnings"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Tree"
        value-unit="ns"
        measure="fleet.test"
        :projects="['wideTree', 'deepTree']"
        :configurators="dashboardConfigurators"
        :accidents="warnings"
      />
    </section>
    <Divider title="PHP" />
    <section>
      <GroupProjectsChart
        label="Typing (time)"
        value-unit="ns"
        measure="fleet.test"
        :projects="['Typing in mPDF', 'Typing in mPDF With Backend', 'Pressing Enter in mPDF']"
        :configurators="dashboardConfigurators"
        :accidents="warnings"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Typing (average delay)"
        value-unit="ns"
        measure="awt.delay"
        :projects="['Typing in mPDF', 'Typing in mPDF With Backend', 'Pressing Enter in mPDF']"
        :configurators="dashboardConfigurators"
        :accidents="warnings"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Typing (max delay)"
        value-unit="ns"
        measure="max.awt.delay"
        :projects="['Typing in mPDF', 'Typing in mPDF With Backend', 'Pressing Enter in mPDF']"
        :configurators="dashboardConfigurators"
        :accidents="warnings"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Other"
        value-unit="ns"
        measure="fleet.test"
        :projects="['Open mPDF', 'Frontend Completion in mPDF']"
        :configurators="dashboardConfigurators"
        :accidents="warnings"
      />
    </section>
  </DashboardPage>>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
import Divider from "../common/Divider.vue"

</script>