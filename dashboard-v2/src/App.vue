<template>
  <IntelliJDashboard />
</template>
<script lang="ts">
import { defineComponent } from "vue"
import IntelliJDashboard from "./IntelliJDashboard.vue"

export default defineComponent({
  name: "App",
  components: {
    IntelliJDashboard,
  },
})
</script>
