<template>
  <DashboardPage
    db-name="fleet"
    table="measure_new"
    persistent-id="fleetPerf_dashboard"
    initial-machine="linux-blade-hetzner"
    :with-installer="false"
  >
    <Divider title="Typing" />
    <section>
      <GroupProjectsChart
        label="Typing latency p99"
        measure="p99"
        :projects="typingLatency.projects"
        :aliases="typingLatency.aliases"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Typing latency p50 (technical-only metric)"
        measure="p50"
        :projects="typingLatency.projects"
        :aliases="typingLatency.aliases"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Typing (total time)"
        measure="fleet.test"
        :projects="typingTotalTime.projects"
        :aliases="typingTotalTime.aliases"
      />
    </section>

    <Divider title="Core" />

    <section>
      <GroupProjectsChart
        label="Open Project"
        measure="fleet.test"
        :projects="openProject.projects"
        :aliases="openProject.aliases"
      />
    </section>

    <section>
      <GroupProjectsChart
        label="Highlighting"
        measure="fleet.test"
        :projects="highlighting.projects"
        :aliases="highlighting.aliases"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Tree"
        measure="fleet.test"
        :projects="tree.projects"
        :aliases="tree.aliases"
      />
    </section>

    <section>
      <GroupProjectsChart
        label="Other"
        measure="fleet.test"
        :projects="otherProjects.projects"
        :aliases="otherProjects.aliases"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import { computed } from "vue"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
import Divider from "../common/Divider.vue"

// Old project names these tests used to report under - either a deliberate rename to a more
// readable name (e.g. "stressTyping" -> "typing") or, for the mPDF/tree tests, escapeTestName()
// accidentally PascalCasing a spaced name (e.g. "wide tree" -> "wideTree") until the ../ultimate
// escaping fix (unsafeTestName + toSafePathSegment) ships. Always merged in: it's the same
// benchmark reported under its old name, not unrelated historical data.
const OLD_NAME_ALIASES: Record<string, string[]> = {
  "typing in mPDF": ["TypingInMPDF", "Typing in mPDF"],
  "line wrapping in mPDF": ["PressingEnterInMPDF", "Pressing Enter in mPDF"],
  "open mPDF": ["OpenMPDF", "Open mPDF"],
  "frontend completion in mPDF": ["FrontendCompletionInMPDF", "Frontend Completion in mPDF"],
  "wide tree": ["wideTree"],
  "deep tree": ["deepTree"],
  typing: ["stressTyping"],
  "typing with robot": ["stressTypingWithRobot"],
  "line wrapping": ["stressEnter"],
  "multi-caret typing": ["multiCaretTyping"],
  "multi-caret brace typing": ["multiCaretBraceTyping"],
  highlighting: ["stressHighlighting"],
  "long file insertion": ["stressLongFileInsertion"],
  "JavaScript retyping": ["javascriptRetyping"],
  "JavaScript retyping 2": ["javascriptRetyping2"],
  "open go-delve project": ["openGoDelveProject"],
  "open jediterm project": ["openJeditermProject"],
  "open rust-simple-server project": ["openRustSimpleServerProject"],
  "open spring-petclinic-java project": ["openSpringPetClinicJavaProject"],
  "open spring-petclinic-kotlin project": ["openSpringPetClinicKotlinProject"],
}

function withOldNameAliases(canonical: string): string[] {
  return [canonical, ...(OLD_NAME_ALIASES[canonical] ?? [])]
}

function projectSeries(canonicalNames: string[]): { projects: string[]; aliases: string[] } {
  const projects: string[] = []
  const aliases: string[] = []
  for (const canonical of canonicalNames) {
    for (const p of withOldNameAliases(canonical)) {
      projects.push(p)
      aliases.push(canonical)
    }
  }
  return { projects, aliases }
}

const typingLatency = computed(() => projectSeries(["typing", "long file insertion", "JavaScript retyping"]))
const typingTotalTime = computed(() => projectSeries(["multi-caret typing", "line wrapping", "typing", "typing in mPDF", "Typing in mPDF With Backend", "line wrapping in mPDF"]))
const openProject = computed(() =>
  projectSeries(["open go-delve project", "open jediterm project", "open rust-simple-server project", "open spring-petclinic-java project", "open spring-petclinic-kotlin project"])
)
const highlighting = computed(() => projectSeries(["highlighting"]))
const tree = computed(() => projectSeries(["wide tree", "deep tree"]))
const otherProjects = computed(() => projectSeries(["open mPDF", "frontend completion in mPDF"]))
</script>
