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
      <Chip><a href="#history">Show file history</a></Chip>
      <Chip><a href="#annotate">Annotate</a></Chip>
    </div>

    <Accordion
      :multiple="true"
      :value="[0, 1, 2, 3]"
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
        <AccordionHeader>Commit</AccordionHeader>
        <AccordionContent>
          <section>
            <GroupProjectsChart
              id="commit"
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
        </AccordionContent>
      </AccordionPanel>

      <AccordionPanel :value="2">
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
        <AccordionHeader>Annotate</AccordionHeader>
        <AccordionContent>
          <!--<section>-->
          <!--  <GroupProjectsChart-->
          <!--    id="annotate"-->
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

const spaceProject = "space/gitLogIndexing"
const indexingProjects = [spaceProject]

const commitProjects = ["space/git-commit", "space/git-commit-smallDataPack"]

const showFileHistoryProjects = ["space/DmsFacadeImpl-instant-git", "space/DmsFacadeImpl"]

const annotateProjects = ["space/vcs-annotate-instant-git", "space/vcs-annotate"]
</script>
