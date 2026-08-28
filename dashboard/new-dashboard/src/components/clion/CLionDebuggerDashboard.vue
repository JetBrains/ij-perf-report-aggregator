<template>
  <DashboardPage
    v-slot="{ machineConfigurator, branchConfigurator }"
    db-name="perfintDev"
    table="clion"
    persistent-id="clion_debugger_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    :with-installer="false"
  >
    <template
      v-for="group in visibleGroups(machineConfigurator, branchConfigurator)"
      :key="group.value"
    >
      <Divider :label="group.title" />
      <section>
        <GroupProjectsChart
          v-for="chart in group.charts"
          :key="chart.key"
          :label="`${group.prefix}: ${chart.label}`"
          :measure="chart.measures"
          :projects="group.projects"
        />
      </section>
    </template>
  </DashboardPage>
</template>

<script lang="ts" setup>
import { BranchConfigurator } from "../../configurators/BranchConfigurator"
import { selectedToArray } from "../../configurators/DimensionConfigurator"
import { isMacMachine, MachineConfigurator } from "../../configurators/MachineConfigurator"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
import Divider from "../common/Divider.vue"

interface ChartDef {
  key: string
  label: string
  measures: string[]
}

interface GroupDef {
  value: string
  title: string
  prefix: string
  projects: string[]
  charts: ChartDef[]
  perDebugger: boolean
  availableOnMac: boolean // Mac agents ship LLDB only
}

const debugScenarios = ["radler/fmtlib/debug/args-test/basic", "radler/luau/debug/Analyze.cpp", "radler/opencv/debug/test_arithm.cpp"]

// Each scenario is run once per debugger, and the debugger name is the last project segment.
function projectsFor(debuggerName: string): string[] {
  return debugScenarios.map((scenario) => `${scenario}/${debuggerName}`)
}

const debugCharts: ChartDef[] = [
  { key: "launch", label: "Launch Debug", measures: ["fus_debug_session_initialized_ms"] },
  { key: "frame", label: "Frame Variables Computed", measures: ["fus_frame_variables_computed_ms"] },
  { key: "stepInto", label: "Step Into", measures: ["debugStep_into"] },
  { key: "stepOut", label: "Step Out", measures: ["debugStep_out"] },
  { key: "stepOver", label: "Step Over", measures: ["debugStep_over"] },
  { key: "evaluate", label: "Evaluate Expression Mean", measures: ["evaluateExpression#mean_value"] },
]

const allGroups: GroupDef[] = [
  { value: "gdb", title: "Debug Actions (GDB)", prefix: "GDB", projects: projectsFor("gdb"), charts: debugCharts, perDebugger: true, availableOnMac: false },
  { value: "lldb", title: "Debug Actions (LLDB)", prefix: "LLDB", projects: projectsFor("lldb"), charts: debugCharts, perDebugger: true, availableOnMac: true },
  { value: "combined", title: "Debug Actions (Default)", prefix: "Debug", projects: debugScenarios, charts: debugCharts, perDebugger: false, availableOnMac: true },
]

function isOldBranch(branch: string): boolean {
  return /^26[21]/.test(branch)
}

function visibleGroups(machineConfigurator: MachineConfigurator | undefined, branchConfigurator: BranchConfigurator | null): GroupDef[] {
  const selectedMachines = machineConfigurator?.selected.value ?? []
  const macOnly = selectedMachines.length > 0 && selectedMachines.every((machine) => isMacMachine(machine))

  const selectedBranches = selectedToArray(branchConfigurator?.selected.value)
  const hasSplitBranch = selectedBranches.length === 0 || selectedBranches.some((branch) => !isOldBranch(branch))
  const hasReleaseBranch = selectedBranches.some((branch) => isOldBranch(branch))

  return allGroups.filter((group) => (group.perDebugger ? hasSplitBranch : hasReleaseBranch) && (group.availableOnMac || !macOnly))
}
</script>
