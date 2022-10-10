<template>
  <PageHeader />
  <slot />
</template>
<script setup lang="ts">
import { provide, ref } from "vue"
import PageHeader from "./PageHeader.vue"
import { InfoSidebarVmImpl } from "./components/InfoSidebarVm"
import {sidebarVmKey } from "./shared/keys"

const sidebarVm = new InfoSidebarVmImpl()

provide(sidebarVmKey, sidebarVm)
</script>