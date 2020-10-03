<!-- Copyright 2000-2020 JetBrains s.r.o. Use of this source code is governed by the Apache 2.0 license that can be found in the LICENSE file. -->
<template>
  <el-container>
    <el-header>
      <el-menu mode="horizontal" :router=true>
        <el-menu-item index="/report">Report Analyzer</el-menu-item>
        <el-menu-item index="/aggregatedStats">IJ Dashboard</el-menu-item>
        <el-menu-item index="/fleetStats">Fleet Dashboard</el-menu-item>
        <el-menu-item index="/sharedIndexes">Shared Indexes Dashboard</el-menu-item>
      </el-menu>
    </el-header>
    <el-main>
      <!-- not re-rendered on switch if keep-alive is used - do not keep-alive for now -->
      <router-view/>
    </el-main>
  </el-container>
</template>

<script lang="ts">
import {Component, Vue} from "vue-property-decorator"

@Component
export default class App extends Vue {
}
</script>