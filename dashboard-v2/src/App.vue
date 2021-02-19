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
import ElHeader from "element-plus/es/el-header"
import ElMain from "element-plus/es/el-main"
import ElMenu from "element-plus/es/el-menu"
import ElSubmenu from "element-plus/es/el-submenu"
import ElMenuItem from "element-plus/es/el-menu-item"

export default defineComponent({
  name: "App",
  components: {
    ElHeader,
    ElMain,
    ElMenu,
    ElMenuItem,
    ElSubmenu,
  }
})
</script>
