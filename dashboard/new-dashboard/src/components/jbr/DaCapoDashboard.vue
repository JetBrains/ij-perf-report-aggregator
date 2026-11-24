<template>
  <JbrBenchmarkDashboard
    benchmark="Dacapo"
    :groups="[{ title: null, scenarios }]"
    description="Milliseconds per iteration of each DaCapo benchmark, lower is better."
    :default-configurations="['macOS15aarch64', 'Ubuntu2404x86_64', 'Windows11x86_64']"
  />
</template>

<script setup lang="ts">
import JbrBenchmarkDashboard from "./JbrBenchmarkDashboard.vue"

const scenarios = [
  "avrora",
  "batik",
  "biojava",
  "eclipse",
  "fop",
  "graphchi",
  "h2",
  "jme",
  "jython",
  "kafka",
  "luindex",
  "lusearch",
  "pmd",
  "spring",
  "sunflow",
  "tomcat",
  "xalan",
  "zxing",
]
</script>
