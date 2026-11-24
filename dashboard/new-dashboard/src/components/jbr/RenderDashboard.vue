<template>
  <JbrBenchmarkDashboard
    benchmark="Render"
    :groups="groups"
    :modes="modes"
    :statistics="statistics"
    :description="description"
    :default-configurations="['macOS15aarch64Metal', 'Ubuntu2404x86_64MM_Wayland', 'Windows11x86_64']"
  />
</template>

<script setup lang="ts">
import JbrBenchmarkDashboard, { Choice, ScenarioGroup, Statistic } from "./JbrBenchmarkDashboard.vue"

/** RenderPerfTest's scenarios by what they draw, what an IDE draws most first. */
const groups: ScenarioGroup[] = [
  { title: "Text", scenarios: ["TextNoAA", "TextLCD", "TextGray", "LargeTextNoAA"] },
  { title: "Images", scenarios: ["Image", "VolImage"] },
  { title: "Mixes", scenarios: ["TextWiredQuadMix", "VolImageFlatBoxMix", "VolImageTextNoAABat"] },
  { title: "Shapes", scenarios: ["FlatBox", "ClipFlatOval", "FlatOval_XOR", "FlatQuadAA", "WiredBubbles", "WiredQuadAA"] },
  { title: "Gradients", scenarios: ["LinGrad3RotatedOval", "RadGrad3RotatedOvalAA"] },
]

/** RenderPerfTest's execution modes, the suffix of every measure name ("<scenario>_<mode>"). */
const modes: Choice[] = [
  { id: "onscreen", label: "Onscreen" },
  { id: "volatile", label: "Volatile" },
  { id: "buffer", label: "Buffer" },
]

/**
 * The frame time statistics that render.py reports for each scenario and mode, as measures named "<scenario>_<mode>"
 * (the median) and "<scenario>_<mode>#<statistic>". The mean's measure comes with the standard deviation, which the
 * charts draw as a band around it.
 */
const statistics: Statistic[] = [
  { id: "median", label: "Median", suffix: "" },
  { id: "mean", label: "Mean ± SD", suffix: "#mean_value" },
  { id: "p90", label: "p90", suffix: "#p90" },
  { id: "p99", label: "p99", suffix: "#p99" },
  { id: "max", label: "Max", suffix: "#max" },
]

const description =
  "Frame times in ms, lower is better. Onscreen: drawn into the window, each frame timed until the Robot sees it on the screen; " +
  "Volatile: drawn into a VolatileImage (the accelerated pipeline, no vsync); Buffer: drawn into a BufferedImage (software loops). " +
  "Median, p90 and p99: each repeat's percentile of its frames, the median over the repeats. Max: the slowest frame of all repeats. " +
  "Mean ± SD: the mean and the standard deviation of all frames of all repeats."
</script>
