<template>
  <DashboardPage
    db-name="perfint"
    table="idea"
    persistent-id="vcs_idea_ultimate_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
  >
    <div>
      <Chip><a href="#index">Vcs indexing</a></Chip>
      <Chip><a href="#history">Show file history</a></Chip>
      <Chip><a href="#checkout">Checkout</a></Chip>
      <Chip><a href="#filter">Filter Vcs Log tab</a></Chip>
    </div>

    <Accordion
      :multiple="true"
      :active-index="[0, 1, 2, 3]"
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
      </AccordionTab>

      <AccordionTab header="Show file history">
        <a name="history" />
        <section>
          <GroupProjectsChart
            label="Show file history (test metric)"
            measure="showFileHistory"
            :projects="historyProjects"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Show file history - showing first pack of data (test metric)"
            measure="showFirstPack"
            :projects="historyProjects"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Show file history - showing first pack of data (test metric)"
            measure="showFirstPack"
            :projects="historyProjects"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Computing - time spent on computing a peace of history.
          If index - time of computing before the first rename. If git - time of computing before timeout of operation occurred"
            measure="file-history-computing"
            :projects="historyProjects"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Loading full VCS Log (all commits and references)"
            measure="vcs-log-loading-full-log"
            :projects="historyProjects"
          />
        </section>
      </AccordionTab>

      <AccordionTab header="Checkout">
        <a name="checkout" />
        <section>
          <GroupProjectsChart
            label="Checkout time"
            measure="git-checkout"
            :projects="checkoutProjects"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Checkout duration(FUS)"
            measure="git-checkout#fusCheckoutDuration"
            :projects="checkoutProjects"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Checkout VFS refresh duration(FUS)"
            measure="git-checkout#fusVfsRefreshDuration"
            :projects="checkoutProjects"
          />
        </section>
      </AccordionTab>

      <AccordionTab header="Filter Vcs Log tab">
        <a name="filter" />
        <section>
          <GroupProjectsChart
            label="Filter Vcs Log tab"
            measure="vcs-log-filtering"
            :projects="filteringProjects"
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

const intellijSpecificProject = "intellij_clone_specific_commit/gitLogIndexing"
const indexingProjects = [intellijSpecificProject, intellijSpecificProject + "-sql"]

const showFileHistoryEditorProject = "intellij_clone_specific_commit/EditorImpl-"
const historyProjects = [showFileHistoryEditorProject + "phm", showFileHistoryEditorProject + "sql", showFileHistoryEditorProject + "noindex"]

const checkoutProjects = ["intellij_clone_specific_commit/git-checkout"]

const vcsLogFilterProject = "intellij_clone_specific_commit/filterVcsLogTab-"
const filteringProjects = [vcsLogFilterProject + "phm", vcsLogFilterProject + "sql", vcsLogFilterProject + "noindex"]
</script>
