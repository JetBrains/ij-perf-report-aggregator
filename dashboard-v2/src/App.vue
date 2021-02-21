<template>
  <el-container>
    <el-header>
      <el-menu
        mode="horizontal"
        :router="true"
        default-active="/ij/dashboard"
      >
        <el-submenu index="/ij/dashboard">
          <template #title>
            IJ Dashboard
          </template>
          <el-menu-item index="/ij/dashboard">
            Dashboard
          </el-menu-item>
          <el-menu-item index="/ij/explore">
            Explore
          </el-menu-item>
          <el-menu-item index="/ij/compare">
            Compare
          </el-menu-item>
        </el-submenu>
      </el-menu>
    </el-header>
    <el-main>
      <router-view v-slot="{ Component }">
        <keep-alive max="4">
          <component :is="Component" />
        </keep-alive>
      </router-view>
    </el-main>
  </el-container>
</template>
<script lang="ts">
import { defineComponent } from "vue"

export default defineComponent({
  name: "App",
})
</script>
<style>
.dividerAfterForm {
  margin-top: 0 !important;
}
</style>
