<template>
  <DashboardPage
    db-name="perfint"
    table="idea"
    persistent-id="vcs_idea_ultimate_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
  >
    <div>
      <Chip><a href="#index">Vcs indexing</a></Chip>
      <Chip><a href="#commit">Commit</a></Chip>
    </div>

    <Accordion
      :multiple="true"
      :active-index="[0, 1]"
    >
      <AccordionTab header="Indexing">
        <a name="index" />
        <section>
          <GroupProjectsChart
            label="Indexing"
            measure="vcs-log-indexing"
            :projects="indexingProjects"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Number of collected commits"
            measure="vcs-log-indexing#numberOfCommits"
            :projects="indexingProjects"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Real number of collected commits through git rev-list --count --all"
            measure="realNumberOfCommits"
            :projects="indexingProjects"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="'vcs-log' directory size in gb"
            measure="vcs-log-directory-size-gb"
            :projects="indexingProjects"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Building of 'git commit-graph write --reachable --changed-paths' in minutes"
            measure="git-build-commit-graph"
            :projects="indexingProjects"
          />
        </section>
      </AccordionTab>

      <AccordionTab header="Commit">
        <a name="commit" />
        <section>
          <GroupProjectsChart
            label="Commit FUS duration"
            measure="git-commit#fusCommitDuration"
            :projects="commitProjects"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Refreshing VCS Log when repositories change (on commit, rebase, checkout branch, etc.)"
            measure="vcs-log-refreshing"
            :projects="commitProjects"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Partial refresh of the VCS Log, building of SmallDataPack (on commit, rebase, checkout branch, etc.)"
            measure="vcs-log-partial-refreshing"
            :projects="commitProjects"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Building a [com.intellij.vcs.log.graph.PermanentGraph] for the list of commits"
            measure="vcs-log-building-graph"
            :projects="commitProjects"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Loading full VCS Log (all commits and references)"
            measure="vcs-log-loading-full-log"
            :projects="commitProjects"
          />
        </section>
      </AccordionTab>
    </Accordion>
  </DashboardPage>
</template>

<script setup lang="ts">
import Accordion from "primevue/accordion"
import Chip from "primevue/chip"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"

const spaceProject = "space/gitLogIndexing"
const indexingProjects = [spaceProject]

const commitProjects = ["space/git-commit", "space/git-commit-smallDataPack"]
</script>
