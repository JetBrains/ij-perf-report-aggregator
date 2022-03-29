<template>
  <InputForm />
  <small>
    For timeline no threshold, for other charts threshold is 10ms.
    End time equals to last dumb-aware project start-up activity (so, services and others may be out of end time).
  </small>

  <TabbedCharts :is-info-chart="true" />
  <TabbedCharts :is-info-chart="false" />
  
  <ul>
    <li>
      <small>app initialized: end of phase <code>app initialized callback</code>.</small>
    </li>
    <li>
      <small>project initialized: end of phase <code>module loading</code>.</small>
    </li>
  </ul>
</template>

<script lang="ts">
import { defineComponent } from "vue"
import InputForm from "./InputForm.vue"
import TabbedCharts from "./TabbedCharts.vue"

export default defineComponent({
  name: "IntelliJReport",
  components: {InputForm, TabbedCharts},
})
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
  user-select: text
}
</style>
