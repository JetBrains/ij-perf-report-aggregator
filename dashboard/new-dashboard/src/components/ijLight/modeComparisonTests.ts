import { Chart, ChartDefinition, combineCharts } from "../charts/DashboardCharts"
import { YAxisScale } from "../charts/yAxisScale"
import { SERIES_NAME_SEPARATOR } from "../common/DataQueryExecutor"

// IJ Light "mode comparison" tests: the same scenario runs as a monolith IDE, as IJ Light, and as IJ Light
// upgraded, each under simulated network latency (the -0ms/-50ms/-1000ms project suffix).
// Reports land in perfintDev.idea with the mode in the `mode` column.
export const MODES = ["light", "light-upgraded", "monolith"]

const MODE_LABELS: Record<string, string> = { light: "Light", "light-upgraded": "Light Upgraded", monolith: "Monolith" }

// A series name is producer tokens joined with SERIES_NAME_SEPARATOR, e.g. "light-upgraded – 0 ms" (mode first,
// then project alias; machine/branch/measure tokens appear only when several are selected). Rewrite only the mode
// token. Must never return "" — LineChartVM ignores a formatter whose result for "test" is empty.
export function legendFormatter(name: string): string {
  return name
    .split(SERIES_NAME_SEPARATOR)
    .map((token) => MODE_LABELS[token] ?? token)
    .join(SERIES_NAME_SEPARATOR)
}

// Monolith runs only the 0 ms variants: its 50/1000 ms series stay empty.
const NETWORK_LATENCIES = ["0ms", "50ms", "1000ms"]
const NETWORK_LATENCY_ALIASES = ["0 ms", "50 ms", "1000 ms"]

interface Scenario {
  label: string
  // project = `${family.projectPrefix}/${stem}-${latency}`
  stem: string
}

interface FamilyMetric {
  // Divider title
  title: string
  measure: string | string[]
  description?: string
}

// One test family = one dashboard tab: a Divider section per metric, a chart per scenario.
interface TestFamily {
  // DashboardPage persistent-id: machine / branch / time range / modes are remembered per tab
  persistentId: string
  projectPrefix: string
  scenarios: Scenario[]
  metrics: FamilyMetric[]
  // "log" when the latencies differ by an order of magnitude and would flatten each other on a linear axis
  yAxisScale?: YAxisScale
}

interface ModeComparisonSection {
  key: string
  title: string
  description?: string
  charts: Chart[]
}

interface ModeComparisonPage {
  persistentId: string
  yAxisScale: YAxisScale
  sections: ModeComparisonSection[]
}

// com.jetbrains.rdct.ijLight.performance.IjLightTypingLatencyTest
const typingLatency: TestFamily = {
  persistentId: "ijLight_typing_latency_dashboard",
  projectPrefix: "ij-light-typing-latency",
  scenarios: [
    { label: "Plain text", stem: "plain-text" },
    { label: "Large plain text", stem: "plain-text-large" },
    { label: "Backspace in plain text", stem: "plain-text-backspace" },
    { label: "Paste into plain text", stem: "plain-text-paste" },
    { label: "TextMate", stem: "textmate" },
    { label: "Large TextMate", stem: "textmate-large" },
  ],
  metrics: [
    { title: "Typing latency — mean", measure: "typing#latency#mean_value", description: "Mean key press → letter shown delay over the typing run." },
    { title: "Typing latency — p90", measure: "typing#latency#p90", description: "90th percentile of the key press → letter shown delay." },
    { title: "Typing latency — max", measure: "typing#latency#max", description: "Worst single keystroke in the run; sensitive to GC pauses and freezes." },
  ],
}

// Reported by the same IjLightTypingLatencyTest methods, for the file the test then types in.
const openEditor: TestFamily = {
  persistentId: "ijLight_open_editor_dashboard",
  projectPrefix: "ij-light-open-editor",
  scenarios: [
    { label: "Plain text", stem: "plain-text" },
    { label: "Large plain text", stem: "plain-text-large" },
    { label: "TextMate", stem: "textmate" },
    { label: "Large TextMate", stem: "textmate-large" },
  ],
  metrics: [
    {
      title: "Editor opening time",
      measure: "openEditor",
      description: "Time to open the file in the editor. Logarithmic axis: at 1000 ms network latency opening takes up to 10× longer than at 0 ms.",
    },
  ],
  yAxisScale: "log",
}

function pageFor(family: TestFamily): ModeComparisonPage {
  const sections = family.metrics.map((metric) => {
    const declaration: ChartDefinition[] = family.scenarios.map((scenario) => ({
      labels: [scenario.label],
      measures: [metric.measure],
      projects: NETWORK_LATENCIES.map((latency) => `${family.projectPrefix}/${scenario.stem}-${latency}`),
      aliases: NETWORK_LATENCY_ALIASES,
    }))
    return {
      key: Array.isArray(metric.measure) ? metric.measure.join("+") : metric.measure,
      title: metric.title,
      description: metric.description,
      charts: combineCharts(declaration),
    }
  })
  return { persistentId: family.persistentId, yAxisScale: family.yAxisScale ?? "linear", sections }
}

// Add the next family here, plus a tab and an ijLightModeComparison route in routes.ts.
export const MODE_COMPARISON_PAGES = {
  typingLatency: pageFor(typingLatency),
  openEditor: pageFor(openEditor),
}

export type ModeComparisonFamily = keyof typeof MODE_COMPARISON_PAGES
