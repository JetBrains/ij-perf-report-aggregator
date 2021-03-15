<template>
  <InputForm />

  <el-row>
    <el-col>
      <small>
        For timeline no threshold, for other charts threshold is 10ms.
        End time equals to last dumb-aware project start-up activity (so, services and others may be out of end time).
      </small>
    </el-col>
  </el-row>

  <TabbedCharts :is-info-chart="true" />
  <TabbedCharts :is-info-chart="false" />

  <el-row>
    <el-col>
      <ul>
        <li>
          <small>app initialized: end of phase <code>app initialized callback</code>.</small>
        </li>
        <li>
          <small>project initialized: end of phase <code>module loading</code>.</small>
        </li>
      </ul>
    </el-col>
  </el-row>
</template>

<script lang="ts">
import { defineComponent } from "vue"
import InputForm from "./InputForm.vue"
import TabbedCharts from "./TabbedCharts.vue"

export default defineComponent({
  name: "Report",
  components: {InputForm, TabbedCharts},
})
</script>
