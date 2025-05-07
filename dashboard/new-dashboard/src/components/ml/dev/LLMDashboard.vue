<template>
  <DashboardPage
    db-name="perfintDev"
    table="ml"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    persistent-id="llmDashboardDev"
    :with-installer="false"
  >
    <section>
      <GroupProjectsChart
        label="Inline completion"
        measure="callInlineCompletionOnCompletion#mean_value"
        :projects="['gradle-calculator_SimpleInlineCompletionTest/simple inline completion']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Code generation"
        measure="ai-generate-code#mean_value"
        :projects="['gradle-calculator_CodeGenerationPerformanceTest/generate code']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Do generate"
        measure="doGenerate#mean_value"
        :projects="['gradle-calculator_CodeGenerationPerformanceTest/generate code']"
      ></GroupProjectsChart>
    </section>
    <section>
      <GroupProjectsChart
        label="Chat/Context"
        :measure="[
          'SimpleCompletableMessage::Collecting context…',
          'SimpleCompletableMessage::Generating answer…',
          'retrieveContextAsync.Chat Submit(userInput=&quot;test&quot;, chatSession=ChatRetrievalSession(&quot;New Chat&quot;), chatRetrieversType=SLOW).time.max',
          'retrieveContextAsync.Initialize chat.time.max',
          'retrieveContextAsync.[Deprecated] Chat Retrieval(userInput=&quot;test&quot;, chatSession=ChatRetrievalSession(&quot;New Chat&quot;), chatRetrieversType=ALL).time.max',
          'computeContext.PsiFileSearchRetriever.[Deprecated] Chat Retrieval(userInput=&quot;test&quot;, chatSession=ChatRetrievalSession(&quot;New Chat&quot;), chatRetrieversType=ALL).time.max',
        ]"
        :legend-formatter="chatContextMetricsLegendFormatter"
        :projects="['kotlinx_coroutines_k2_dev_ContextPerformanceTest/basic context performance test']"
      ></GroupProjectsChart>
    </section>
  </DashboardPage>
</template>
<script setup lang="ts">
import GroupProjectsChart from "../../charts/GroupProjectsChart.vue"
import DashboardPage from "../../common/DashboardPage.vue"

const chatContextMetricsLegendFormatter = (name: string) => {
  if (name.startsWith("retrieveContextAsync.Chat Submit")) return "retrieveContext.Chat Submit"
  else if (name.startsWith("retrieveContextAsync.Initialize chat")) return "retrieveContext.Initialize chat"
  else if (name.startsWith("retrieveContextAsync.[Deprecated] Chat Retrieval")) return "retrieveContext.Chat Retrieval[Deprecated]"
  else if (name.startsWith("computeContext.PsiFileSearchRetriever")) return "PsiFileSearchRetriever"
  else return name.replace("SimpleCompletableMessage::", "")
}
</script>
