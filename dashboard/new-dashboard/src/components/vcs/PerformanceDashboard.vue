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
      <Chip><a href="#filter">Filter Log tab</a></Chip>
      <Chip><a href="#commit">Commit</a></Chip>
      <Chip><a href="#widget">Branch widget</a></Chip>
    </div>

    <Accordion
      :multiple="true"
      :active-index="[0, 1, 2, 3, 4, 5]"
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

      <AccordionTab header="Show file history">
        <a name="history" />
        <section>
          <GroupProjectsChart
            label="Show file history (test metric)"
            measure="showFileHistory"
            :projects="showFileHistoryProjects"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Show file history - showing first pack of data (test metric)"
            measure="showFirstPack"
            :projects="showFileHistoryProjects"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Computing - time spent on computing a peace of history.
          If index - time of computing before the first rename. If git - time of computing before timeout of operation occurred"
            measure="file-history-computing"
            :projects="showFileHistoryProjects"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Loading full VCS Log (all commits and references)"
            measure="vcs-log-loading-full-log"
            :projects="showFileHistoryProjects"
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
            label="Filter Vcs Log tab by name"
            measure="vcs-log-filtering"
            :projects="filterByNameProjects"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Filter Vcs Log tab by path"
            measure="vcs-log-filtering"
            :projects="filterByPathProjects"
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
            label="Refreshing VCS Log when repositories change (on commit, rebase, checkout branch, etc.) - vcs-log-refreshing"
            measure="vcs-log-refreshing"
            :projects="commitProjects"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Partial refresh of the VCS Log, building of SmallDataPack (on commit, rebase, checkout branch, etc.) - vcs-log-partial-refreshing"
            measure="vcs-log-partial-refreshing"
            :projects="commitProjects"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Building a [com.intellij.vcs.log.graph.PermanentGraph] for the list of commits - vcs-log-building-graph"
            measure="vcs-log-building-graph"
            :projects="commitProjects"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Loading full VCS Log (all commits and references) - vcs-log-loading-full-log"
            measure="vcs-log-loading-full-log"
            :projects="commitProjects"
          />
        </section>
      </AccordionTab>
      <AccordionTab header="Show branch widget">
        <a name="widget" />
        <section>
          <GroupProjectsChart
            label="Duration of expanding the whole branch tree - GitBranchesTreePopup::waitTreeExpand"
            measure="gitShowBranchWidget"
            :projects="widgetProjects"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Duration of initializing tree - git-branches-popup-building-tree"
            measure="git-branches-popup-building-tree"
            :projects="widgetProjects"
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
const indexingProjects = [intellijSpecificProject]

const showFileHistory = "intellij_clone_specific_commit/EditorImpl-"
const showFileHistoryProjects = [showFileHistory + "phm", showFileHistory + "noindex", "intellij_sources/showFileHistory/EditorImpl"]

const checkoutProjects = ["intellij_clone_specific_commit/git-checkout"]

const vcsLogFilterProject = "intellij_clone_specific_commit/filterVcsLogTab-"
const filterByNameProjects = [vcsLogFilterProject + "phm", vcsLogFilterProject + "noindex"]
const filterByPathProjects = [vcsLogFilterProject + "path-phm", vcsLogFilterProject + "path-noindex"]

const commitProjects = ["intellij_clone_specific_commit/git-commit", "intellij_clone_specific_commit/git-commit-smallDataPack"]

const widgetProjects = ["intellij_clone_specific_commit/git-branch-widget", "vcs-100k-branches/git-branch-widget"]
</script>
