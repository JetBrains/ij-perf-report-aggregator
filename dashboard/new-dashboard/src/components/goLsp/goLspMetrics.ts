import type { MetricInfo } from "../../shared/metricsDescription"

// The requests and edits of the Kubernetes v1.37.1 Code Search test (tests/golsp KubernetesCodeSearch.kt in ultimate).
const requestDescriptions: Record<string, string> = {
  definition: "Go to Definition on the admission.NewHandler call in plugin/pkg/admission/namespace/autoprovision/admission.go; it resolves into the apiserver staging module.",
  references: "Find References of NewHandler in the apiserver admission package (handler.go), without the declaration.",
  implementation: "Go to Implementation of the Handles method of the apiserver admission interface (interfaces.go).",
  workspaceSymbol: "workspace/symbol request for NewSharedInformerFactory.",
}
const scenarioDescriptions: Record<string, string> = {
  first: "The first request after the server is ready and the files are open.",
  warm: "Median of 10 requests after 3 untimed warmup requests.",
  "edit.nearbyBody": "Median of 5 requests, each sent right after an edit inside the body of NewHandler in admission/handler.go. The package API does not change.",
  "edit.nearbyFunction": "Median of 5 requests, each sent right after an edit that adds a function to admission/handler.go, an API change in the package of the requests.",
  "edit.coreLibraryFunction":
    "Median of 5 requests, each sent right after an edit that adds a function to apimachinery meta/v1/time.go, a package that almost every Kubernetes package imports.",
}

export const GO_LSP_REQUESTS = Object.keys(requestDescriptions)
export const GO_LSP_SCENARIOS = Object.keys(scenarioDescriptions)

export const GO_LSP_METRICS: Map<string, string | MetricInfo> = new Map<string, string | MetricInfo>([
  [
    "ready",
    {
      description:
        "Time from initialize until the server is idle. GoLand server: project import and indexing, until intellij/ready-for-test. gopls: package load and diagnosis of all workspace packages.",
      unit: "milliseconds",
    },
  ],
  ...GO_LSP_REQUESTS.flatMap((request) =>
    GO_LSP_SCENARIOS.map((scenario): [string, MetricInfo] => [
      `${request}.${scenario}`,
      { description: `${requestDescriptions[request]} ${scenarioDescriptions[scenario]}`, unit: "milliseconds" },
    ])
  ),
])
