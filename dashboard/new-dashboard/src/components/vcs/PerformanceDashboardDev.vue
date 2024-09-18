<template>
  <DashboardPage
    db-name="perfintDev"
    table="idea"
    persistent-id="vcs_idea_ultimate_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    :with-installer="false"
  >
    <div>
      <Chip><a href="#index">Vcs indexing</a></Chip>
      <Chip><a href="#history">Show file history</a></Chip>
      <Chip><a href="#checkout">Checkout</a></Chip>
      <Chip><a href="#filter">Filter Log tab</a></Chip>
      <Chip><a href="#commit">Commit</a></Chip>
      <Chip><a href="#widget">Branch widget</a></Chip>
      <Chip><a href="#annotate">Annotate</a></Chip>
    </div>

    <Accordion
      :multiple="true"
      :value="[0, 1, 2, 3, 4, 5, 6]"
    >
      <AccordionPanel :value="0">
        <AccordionHeader>Indexing</AccordionHeader>
        <AccordionContent>
          <section>
            <GroupProjectsChart
              id="index"
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
          <!--<section>-->
          <!--  <GroupProjectsChart-->
          <!--    label="Building of 'git commit-graph write &#45;&#45;reachable &#45;&#45;changed-paths' in minutes"-->
          <!--    measure="git-build-commit-graph"-->
          <!--    :projects="indexingProjects"-->
          <!--  />-->
          <!--</section>-->
        </AccordionContent>
      </AccordionPanel>

      <AccordionPanel :value="1">
        <AccordionHeader>Show file history</AccordionHeader>
        <AccordionContent>
          <section>
            <GroupProjectsChart
              id="history"
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
              label="Loading full VCS Log (all commits and references) - vcs-log-loading-full-log"
              measure="vcs-log-loading-full-log"
              :projects="showFileHistoryProjects"
            />
          </section>
          <section>
            <GroupProjectsChart
              label="Collecting revisions from the [com.intellij.vcs.log.VcsLogFileHistoryHandler]"
              measure="file-history-collecting-revisions-from-handler"
              :projects="showFileHistoryProjects"
            />
          </section>
        </AccordionContent>
      </AccordionPanel>

      <AccordionPanel :value="3">
        <AccordionHeader>Checkout</AccordionHeader>
        <AccordionContent>
          <section>
            <GroupProjectsChart
              id="checkout"
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
        </AccordionContent>
      </AccordionPanel>

      <AccordionPanel :value="4">
        <AccordionHeader>Filter Vcs Log tab</AccordionHeader>
        <AccordionContent>
          <section id="filter">
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
          <section>
            <GroupProjectsChart
              label="Filter Vcs Log tab by date - from 2024-02-22 to 2024-03-21"
              measure="vcs-log-filtering"
              :projects="filterByDateProjects"
            />
          </section>
        </AccordionContent>
      </AccordionPanel>
      <AccordionPanel :value="5">
        <AccordionHeader>Commit</AccordionHeader>
        <AccordionContent>
          <section id="commit">
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
          <section>
            <GroupProjectsChart
              label="git-reading-repo-info"
              measure="git-reading-repo-info"
              :projects="commitProjects"
            />
          </section>
        </AccordionContent>
      </AccordionPanel>
      <AccordionPanel :value="6">
        <AccordionHeader>Show branch widget</AccordionHeader>
        <AccordionContent>
          <section id="widget">
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
        </AccordionContent>
      </AccordionPanel>
      <AccordionPanel :value="7">
        <AccordionHeader>Annotate</AccordionHeader>
        <AccordionContent>
          <!--<section id="annotate">-->
          <!--  <GroupProjectsChart-->
          <!--    label="Duration of opening git annotation - showFileAnnotation"-->
          <!--    measure="showFileAnnotation"-->
          <!--    :projects="annotateProjects"-->
          <!--  />-->
          <!--</section>-->
          <section>
            <GroupProjectsChart
              label="Duration of opening git annotation - git-open-annotation"
              measure="git-open-annotation"
              :projects="annotateProjects"
            />
          </section>
        </AccordionContent>
      </AccordionPanel>
    </Accordion>
  </DashboardPage>
</template>

<script setup lang="ts">
import Accordion from "primevue/accordion"
import Chip from "primevue/chip"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"

const indexingProjects = ["intellij_clone_specific_commit/gitLogIndexing"]

const showFileHistoryProjects = ["intellij_clone_specific_commit/EditorImpl-phm", "intellij_clone_specific_commit/EditorImpl-noindex"]

const checkoutProjects = ["intellij_clone_specific_commit/git-checkout"]

const filterByNameProjects = ["intellij_clone_specific_commit/filterVcsLogTab-phm", "intellij_clone_specific_commit/filterVcsLogTab-noindex"]
const filterByPathProjects = ["intellij_clone_specific_commit/filterVcsLogTab-path-phm", "intellij_clone_specific_commit/filterVcsLogTab-path-noindex"]

const filterByDateProjects = ["intellij_clone_specific_commit/filterVcsLogTab-date-phm", "intellij_clone_specific_commit/filterVcsLogTab-date-noindex"]

const commitProjects = ["intellij_clone_specific_commit/git-commit", "intellij_clone_specific_commit/git-commit-smallDataPack"]

const widgetProjects = ["intellij_clone_specific_commit/git-branch-widget", "vcs_100k_branches/git-branch-widget"]
const annotateProjects = ["intellij_clone_specific_commit/vcs-annotate-instant-git", "intellij_clone_specific_commit/vcs-annotate"]
</script>
