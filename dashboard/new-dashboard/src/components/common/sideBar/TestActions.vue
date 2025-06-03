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
import { getNavigateToTestUrl, InfoData } from "./InfoSidebar"

const router = useRouter()

const { data } = defineProps<{
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
  if (data?.description() != null) {
    const methodName = data.description().value?.methodName
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

    const url = data.description().value?.url
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
  window.open(getNavigateToTestUrl(data, router), "_blank")
}
</script>
