<template>
  <SplitButton
    v-if="testActions.length > 0"
    :label="testActions[0].label"
    :model="testActions.slice(1)"
    link
    icon="pi pi-chart-line"
    @click="testActions[0].command"
  />
</template>
<script setup lang="ts">
import { computed } from "vue"
import { useRouter } from "vue-router"
import { dbTypeStore } from "../../../shared/dbTypes"
import { DBType, InfoData } from "./InfoSidebar"
import { getMachineGroupName } from "../../../configurators/MachineConfigurator"

const router = useRouter()

const props = defineProps<{
  data: InfoData | null
}>()

const testActions = computed(() => getTestActions())

function getTestActions(): {
  label: string
  icon: string
  command: () => void
}[] {
  const actions = []
  if (isNavigateToTestSupported()) {
    actions.push({
      label: "Navigate to test",
      icon: "pi pi-chart-line",
      command() {
        handleNavigateToTest()
      },
    })
  }
  if (props.data?.description != null) {
    const methodName = props.data.description.value?.methodName
    if (methodName && methodName != "") {
      actions.push(
        {
          label: "Open test method",
          icon: "pi pi-folder-open",
          command() {
            openTestInIDE(methodName)
          },
        },
        {
          label: "Copy test method name",
          icon: "pi pi-copy",
          command() {
            copyMethodNameToClipboard(methodName)
          },
        }
      )
    }

    const url = props.data.description.value?.url
    if (url && url != "") {
      actions.push({
        label: "Download test project",
        icon: "pi pi-download",
        command() {
          window.open(url)
        },
      })
    }
  }
  return actions
}

function isNavigateToTestSupported(): boolean {
  const currentRoute = router.currentRoute.value
  const parts = currentRoute.path.split("/")
  const pageName = parts.at(-1)?.toLowerCase()
  return pageName != "testsDev" && pageName != "tests" && pageName != "explore"
}

function copyMethodNameToClipboard(methodName: string) {
  void navigator.clipboard.writeText(methodName)
}

function openTestInIDE(methodName: string) {
  const origin = encodeURIComponent("ssh://git@git.jetbrains.team/ij/intellij.git")
  window.open(`jetbrains://idea/navigate/reference?origin=${origin}&fqn=${methodName}`)
}

function handleNavigateToTest() {
  const currentRoute = router.currentRoute.value
  let parts = currentRoute.path.split("/")
  if (parts.at(-1) == "startup" || parts.at(1) == "ij") {
    parts = ["", "ij", "explore"]
  } else if (parts.at(1) == "fleet" && parts.at(2) == "startupDashboard") {
    parts = ["", "fleet", "startupExplore"]
  } else {
    parts[parts.length - 1] = dbTypeStore().dbType == DBType.INTELLIJ_DEV ? "testsDev" : "tests"
  }
  const branch = props.data?.branch ?? ""
  const machineGroup = getMachineGroupName(props.data?.machineName ?? "")
  const majorBranch = branch.match(/\d+\.\d+/) ? branch.slice(0, branch.indexOf(".")) : branch
  const testURL = parts.join("/")

  const queryParams: string = new URLSearchParams({
    ...currentRoute.query,
    project: props.data?.projectName ?? "",
    branch: majorBranch,
    machine: machineGroup,
  }).toString()

  const measures =
    props.data?.series
      .map((s) => s.metricName)
      .filter((m) => m != undefined)
      .map((m) => (dbTypeStore().isIJStartup() && (m as string).includes("/") ? "metrics." + (m as string) : m))
      .map((m) => encodeURIComponent(m as string))
      .map((m) => "&measure=" + m)
      .join("") ?? ""

  window.open(router.resolve(testURL + "?" + queryParams + measures).href, "_blank")
}
</script>
<style #scoped>
.p-tieredmenu .p-menuitem-content {
  @apply text-sm text-gray-600 font-medium text-left relative;
}
</style>
