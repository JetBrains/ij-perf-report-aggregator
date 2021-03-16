<template>
  <el-header>
    <el-menu
      mode="horizontal"
      :router="true"
      :default-active="activePath"
    >
      <template
        v-for="item in routes"
        :key="item.title"
      >
        <el-submenu
          v-if="item.title !== null"
          :index="item.title"
          @click="topLevelClicked(item.children[0].path)"
        >
          <template #title>
            {{ item.title }}
          </template>
          <el-menu-item
            v-for="child in item.children"
            :key="child.path"
            :index="child.path"
          >
            {{ child.meta["menuTitle"] }}
          </el-menu-item>
        </el-submenu>
        <template v-else>
          <el-menu-item
            v-for="child in item.children"
            :key="child.path"
            :index="child.path"
          >
            {{ child.meta["menuTitle"] }}
          </el-menu-item>
        </template>
      </template>
      <el-menu-item
        v-show='!activePath.startsWith("/report")'
        style="padding: 0; vertical-align: center; float: right"
      >
        <ServerSelect v-model="serverUrl" />
      </el-menu-item>
    </el-menu>
  </el-header>
  <el-main>
    <router-view v-slot="{ Component }">
      <keep-alive max="4">
        <component :is="Component" />
      </keep-alive>
    </router-view>
  </el-main>
</template>
<script lang="ts">
import { PersistentStateManager } from "shared/src/PersistentStateManager"
import ServerSelect from "shared/src/components/ServerSelect.vue"
import { ServerConfigurator } from "shared/src/configurators/ServerConfigurator"
import { serverUrlKey } from "shared/src/injectionKeys"
import { watch, defineComponent, provide, ref } from "vue"
import { useRoute, useRouter } from "vue-router"
import { getRoutes } from "./route"

export default defineComponent({
  name: "App",
  components: {
    ServerSelect,
  },

  setup() {
    const serverUrl = ref("")
    const routes = getRoutes()
    provide(serverUrlKey, serverUrl)

    const persistentStateManager = new PersistentStateManager("common", {serverUrl: ServerConfigurator.DEFAULT_SERVER_URL})
    persistentStateManager.add("serverUrl", serverUrl)
    persistentStateManager.init()

    const route = useRoute()
    const activePath = ref("")
    watch(() => route.path, p => {
      activePath.value = p
    })

    const router = useRouter()
    return {
      serverUrl,
      activePath,
      routes,
      topLevelClicked(path: string) {
        void router.push({
          path,
        })
      }
    }
  },
})
</script>