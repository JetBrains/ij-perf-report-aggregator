// Shared highlighting config for the engine-comparison views: the highlighting phases, the per-project
// scenario files, measure-name building, the selectable quantity type, and the phase/bucket display
// labels. Dashboard-only chart settings (the engine list, the metric-type grid, the quantity choices,
// and project-variant expansion) live in HighlightingDashboard.vue, which imports what it needs here.

// Highlighting phases published per file by GoHighlightingTest (see runGoHighlightingScenario). Each is
// a duration in ms; the file label ("<bucket>_<name>") is appended to build the full metric name.
export const PHASES = ["coldStartHighlighting", "warmStartHighlighting", "typingHighlighting"]

export interface HighlightingProject {
  // Chart title.
  title: string
  // DB project base, before the engine suffix and the "/highlighting" scenario segment.
  base: string
  // The three scenario files (fast -> medium -> slow), as their metric-name labels ("<bucket>_<name>").
  files: string[]
}

// Base project names and per-project file labels mirror the ScenarioFile lists in GoHighlightingTest.kt.
export const projects: HighlightingProject[] = [
  { title: "kubernetes", base: "kubernetes", files: ["fast_clientTest", "medium_validation", "slow_generatedPb"] },
  { title: "mattermost", base: "mattermost-server", files: ["fast_webContext", "medium_userStore", "slow_opentracingLayer"] },
  { title: "cockroach", base: "cockroach", files: ["fast_storageParam", "medium_geoBuiltins", "slow_projNonConstOps"] },
  { title: "milvus", base: "milvus", files: ["fast_milvusClient", "medium_proxyImpl", "slow_dataCoordPb"] },
  { title: "rclone", base: "rclone", files: ["fast_httpBackend", "medium_azureblob", "slow_encoderCasesTest"] },
  { title: "volcano", base: "volcano", files: ["fast_queueCli", "medium_eventHandlers", "slow_allocateTest"] },
  { title: "caddy", base: "caddy", files: ["fast_importGraph", "medium_reverseProxyCaddyfile", "slow_httpType"] },
  { title: "k8sDevice", base: "k8sDevice", files: ["fast_resourceFactory", "medium_pciUtil", "slow_serverTest"] },
  { title: "fakeKub", base: "fake_kub", files: ["fast_clientTest", "medium_validation", "slow_generatedPb"] },
]

// The phase segment of a metric type ("coldStartHighlighting_fast" -> "coldStartHighlighting"), or "".
export function phaseOf(type: string): string {
  return PHASES.find((phase) => type.startsWith(`${phase}_`)) ?? ""
}

// The bucket segment of a metric type ("coldStartHighlighting_fast" -> "fast"), or "".
export function bucketOf(type: string): string {
  const phase = phaseOf(type)
  return phase === "" ? "" : type.slice(phase.length + 1)
}

// Concrete measure for a project and a "<phase>_<bucket>" type: the project's file replaces the bucket
// (e.g. "coldStartHighlighting_fast" -> "coldStartHighlighting_fast_importGraph"), then the quantity
// suffix is appended (empty for duration, "#jvm.alloc.mb" for the allocation sub-metric, and so on).
export function buildMeasure(project: HighlightingProject, type: string, quantitySuffix: string): string {
  const phase = phaseOf(type)
  const bucket = bucketOf(type)
  const file = project.files.find((f) => f.startsWith(`${bucket}_`)) ?? ""
  return `${phase}_${file}${quantitySuffix}`
}

// A quantity the dashboard can plot for each scenario. `suffix` is appended to the duration metric
// name; "" is the base duration, "#jvm.alloc.mb" and friends select a JVM sub-metric. Units and
// better-direction are not carried here — they resolve from the declared metricsDescription entries
// (sub-metrics) and from the stored type (duration).
export interface Quantity {
  id: string
  label: string
  suffix: string
}

// Short display labels for the phase/bucket of a metric type, shared by the comparison views.
const PHASE_LABELS: Record<string, string> = {
  coldStartHighlighting: "Cold",
  warmStartHighlighting: "Warm",
  typingHighlighting: "Typing",
}

export function phaseLabel(phase: string): string {
  return PHASE_LABELS[phase] ?? phase
}

export function bucketLabel(bucket: string): string {
  return bucket === "" ? "" : bucket.charAt(0).toUpperCase() + bucket.slice(1)
}

// "coldStartHighlighting_fast" -> "Cold · Fast".
export function metricTypeLabel(type: string): string {
  return `${phaseLabel(phaseOf(type))} · ${bucketLabel(bucketOf(type))}`
}
