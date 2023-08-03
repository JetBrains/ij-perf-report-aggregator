<template>
  <div class="flex justify-between pt-5 px-7 items-center">
    <div class="flex space-x-0.5 font-semibold">
      <span
        v-if="isSubMenuExists"
        class="text-xl"
      >
        Tests on
      </span>
      <button
        class="text-blue-500 px-1 py-1 inline-flex text-xl items-center"
        type="button"
        @click="toggle"
      >
        {{ product.label }}
        <div class="pi pi-chevron-down text-sm ml-1" />
      </button>
      <Menu
        ref="menu"
        :model="items"
        :popup="true"
      />
      <span
        v-if="isSubMenuExists"
        class="text-xl"
      >
        aggregated for
      </span>
      <button
        v-if="isSubMenuExists"
        class="text-blue-500 px-1 py-1 inline-flex text-xl items-center"
        type="button"
        @click="toggleSubMenu"
      >
        {{ selectedSubMenu.label }}
        <div class="pi pi-chevron-down text-sm ml-1" />
      </button>

      <Menu
        v-if="isSubMenuExists"
        ref="subMenu"
        :model="subItems"
        :popup="true"
      />
    </div>
    <a href="https://youtrack.jetbrains.com/articles/IDEA-A-2100661420/IJ-Perf>">
      <QuestionMarkCircleIcon class="w-7 h-7 text-blue-500" />
    </a>
  </div>
</template>

<script setup lang="ts">
import Menu from "primevue/menu"
import { ref } from "vue"
import { useRouter } from "vue-router"
import { SubProject, getNavigationElement, PRODUCTS } from "./routes"

const currentPath = useRouter().currentRoute.value.path
const products = PRODUCTS.map((product) => ({ ...product, url: product.children[0].tabs[0].url })) //default to the first element in the first subproject
const items = ref(products)
const menu = ref<Menu | null>(null)
const product = getNavigationElement(currentPath)

function toggle(event: MouseEvent) {
  menu.value?.toggle(event)
}

const subMenuItems = product.children.map((child) => ({ ...child, url: child.tabs[0].url }))
const isSubMenuExists = product.children.length > 1
const subItems = ref(subMenuItems)
const subMenu = ref<Menu | null>(null)

const selectedSubMenu: SubProject =
  product.children.find((child) => {
    return child.url == currentPath.slice(0, Math.max(0, currentPath.lastIndexOf("/")))
  }) ?? product.children[0]

function toggleSubMenu(event: MouseEvent) {
  subMenu.value?.toggle(event)
}
</script>
