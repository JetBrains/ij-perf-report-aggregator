<template>
  <button
    title="Updated automatically, but you can force data reloading"
    type="button"
    class="-m-2 ml-5 p-2 text-gray-400 hover:text-gray-500 sm:ml-7"
    @click="doLoad"
  >
    <ArrowPathIcon
      class="h-5 w-5"
      aria-hidden="true"
    />
  </button>
</template>
<script setup lang="ts">
import { ArrowPathIcon } from "@heroicons/vue/20/solid"
import { take } from "rxjs"
import { computed, inject, PropType, shallowRef } from "vue"
import { ReloadConfigurator } from "../configurators/ReloadConfigurator"
import { configuratorListKey } from "../injectionKeys"

const props = defineProps({
  load: {
    type: Function as PropType<() => void>,
    default: null,
  },
})

const loading = shallowRef(false)

const doLoad = computed(() => {
  const load = props.load
  if (load != null) {
    return load
  }

  const configurators = inject(configuratorListKey)
  const trigger = configurators?.find(it => it instanceof ReloadConfigurator) as ReloadConfigurator
  if (trigger == null) {
    console.error("`Neither \\`load\\` function is set, nor \\`ReloadConfigurator\\` is provided`", configurators)
  }
  // return executor.scheduleLoadIncludingConfiguratorsFunctionReference
  let counter = 0
  return () => {
    loading.value = true
    trigger.subject.asObservable().pipe(take(1)).subscribe(() => {
      setTimeout(() => {
        loading.value = false
      }, 2000)
    })
    trigger.subject.next(counter++)
  }
})
</script>