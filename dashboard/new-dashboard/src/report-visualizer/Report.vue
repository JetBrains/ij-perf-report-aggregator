<template>
  <InputForm />
  <span class="text-xs">
    For timeline no threshold, for other charts threshold is 10ms. End time equals to last dumb-aware project start-up activity (so, services and others may be out of end time).
  </span>

  <TabbedCharts :is-info-chart="true" />
  <TabbedCharts :is-info-chart="false" />

  <ul>
    <li>
      <span class="text-xs">
        app initialized: end of phase
        <code>app initialized callback</code>
        .
      </span>
    </li>
    <li>
      <span class="text-xs">
        project initialized: end of phase
        <code>module loading</code>
        .
      </span>
    </li>
  </ul>
</template>
<script setup lang="ts">
import InputForm from "./InputForm.vue"
import TabbedCharts from "./TabbedCharts.vue"
</script>
<!--suppress CssUnusedSymbol -->
<style>
.tooltipMainValue {
  float: right;
  margin-left: 20px;
  font-weight: 900;
}

.tooltipValue {
  float: right;
  margin-left: 20px;
}

.tooltipSelectableValue {
  float: right;
  margin-left: 20px;
  user-select: text;
}
</style>
