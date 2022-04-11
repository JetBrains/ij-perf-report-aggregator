<template>
  <Button
    title="Updated automatically, but you can force data reloading"
    icon="pi pi-refresh"
    class="p-button-rounded p-button-text"
    :loading="loading"
    @click="doLoad"
  />
</template>
<script setup lang="ts">
import { take, takeLast } from "rxjs"
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
    trigger.subject.asObservable().pipe(take(1)).subscribe(value => {
      setTimeout(() => {
        loading.value = false
      }, 2000)
    })
    trigger.subject.next(counter++)
  }
})
</script>