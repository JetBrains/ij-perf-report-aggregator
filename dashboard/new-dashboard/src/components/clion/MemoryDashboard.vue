<template>
  <DashboardPage
    db-name="perfintDev"
    table="clion"
    persistent-id="clion_memory_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    :with-installer="false"
  >
    <Divider title="Indexing" />

    <section>
      <CLionVsRadlerMemoryUsageChart
        label="LLVM - After Indexing"
        measure="afterIndexing"
        project="llvm/indexing"
      />
    </section>

    <section>
      <CLionVsRadlerMemoryUsageChart
        label="50k sources, 10k headers - After Indexing"
        measure="afterIndexing"
        project="big_project_50k_10k/indexing"
      />
    </section>

    <section>
      <CLionVsRadlerMemoryUsageChart
        label="OpenCV - After Indexing"
        measure="afterIndexing"
        project="opencv/indexing"
      />
    </section>

    <section>
      <CLionVsRadlerMemoryUsageChart
        label="curl - After Indexing"
        measure="afterIndexing"
        project="curl/indexing"
      />
    </section>

    <Divider title="Resolving References" />

    <section>
      <GroupProjectsChart
        label="[Radler] Resolve All References (cmake)"
        :measure="['clangd_light_modules_total_memory_mb', 'clangd_no_modules_total_memory_mb', 'nova_total_memory_mb']"
        :projects="['radler/cmakeResolving/measureResolve/cmake']"
        value-unit="counter"
      />
    </section>

    <Divider title="Go To Declaration" />

    <section>
      <CLionVsRadlerMemoryUsageChart
        label="After Go to Declaration (ctor) (luau)"
        measure="afterGotoDeclaration"
        project="luau/gotoDeclaration/AstStatDeclareFunction.ctor"
      />
    </section>

    <section>
      <CLionVsRadlerMemoryUsageChart
        label="After Go to Declaration (std::string) (luau)"
        measure="afterGotoDeclaration"
        project="luau/gotoDeclaration/std.string"
      />
    </section>

    <section>
      <CLionVsRadlerMemoryUsageChart
        label="After Go to Declaration (time.h) (luau)"
        measure="afterGotoDeclaration"
        project="luau/gotoDeclaration/time.h"
      />
    </section>

    <Divider title="Find Usages" />

    <section>
      <CLionVsRadlerMemoryUsageChart
        label="After Find Usages (macro) (luau)"
        measure="afterFindUsages"
        project="luau/findUsages/macro (LUAU_ASSERT)"
      />
    </section>

    <section>
      <CLionVsRadlerMemoryUsageChart
        label="After Find Usages (enumerable) (luau)"
        measure="afterFindUsages"
        project="luau/findUsages/enumerable (LuauOpcode)"
      />
    </section>

    <section>
      <CLionVsRadlerMemoryUsageChart
        label="After Find Usages (enumerator) (luau)"
        measure="afterFindUsages"
        project="luau/findUsages/enumerator (LOP_NOP)"
      />
    </section>

    <section>
      <CLionVsRadlerMemoryUsageChart
        label="After Find Usages (class template) (luau)"
        measure="afterFindUsages"
        project="luau/findUsages/class template (DenseHashTable)"
      />
    </section>

    <section>
      <CLionVsRadlerMemoryUsageChart
        label="After Find Usages (class) (cmake)"
        measure="afterFindUsages"
        project="cmake/findUsages/class (cmCTestResourceAllocator)"
      />
    </section>

    <section>
      <CLionVsRadlerMemoryUsageChart
        label="After Find Usages (macro) (cmake)"
        measure="afterFindUsages"
        project="cmake/findUsages/macro (SAFEDIV)"
      />
    </section>

    <section>
      <CLionVsRadlerMemoryUsageChart
        label="After Find Usages (member) (cmake)"
        measure="afterFindUsages"
        project="cmake/findUsages/member (SlotsNeeded)"
      />
    </section>

    <Divider title="Check Test Config" />

    <section>
      <CLionVsRadlerMemoryUsageChart
        label="After Check Test Config (test_houghlines.cpp) (opencv)"
        measure="afterCheckLocalTestConfig"
        project="opencv/checkLocalTestConfig/test.houghlines.cpp.marks"
      />
    </section>

    <section>
      <CLionVsRadlerMemoryUsageChart
        label="After Check Test Config (test.kalman.cpp) (opencv)"
        measure="afterCheckLocalTestConfig"
        project="opencv/checkLocalTestConfig/test.kalman.cpp.marks"
      />
    </section>

    <section>
      <CLionVsRadlerMemoryUsageChart
        label="After Check Test Config (test.kalman.cpp) (luau)"
        measure="afterCheckLocalTestConfig"
        project="luau/checkLocalTestConfig/AstQuery.test.cpp.marks"
      />
    </section>

    <section>
      <CLionVsRadlerMemoryUsageChart
        label="After Check Test Config (Linter.test.cpp) (luau)"
        measure="afterCheckLocalTestConfig"
        project="luau/checkLocalTestConfig/Linter.test.cpp.marks"
      />
    </section>

    <section>
      <CLionVsRadlerMemoryUsageChart
        label="After Check Test Config (Repl.test.cpp) (luau)"
        measure="afterCheckLocalTestConfig"
        project="luau/checkLocalTestConfig/Repl.test.cpp.marks"
      />
    </section>

    <section>
      <CLionVsRadlerMemoryUsageChart
        label="After Check Test Config (TypeInfer.unionTypes.test.cpp) (luau)"
        measure="afterCheckLocalTestConfig"
        project="luau/checkLocalTestConfig/TypeInfer.unionTypes.test.cpp.marks"
      />
    </section>

    <!-- END -->
  </DashboardPage>
</template>

<script setup lang="ts">
import DashboardPage from "../common/DashboardPage.vue"
import Divider from "../common/Divider.vue"
import CLionVsRadlerMemoryUsageChart from "./CLionVsRadlerMemoryUsageChart.vue"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
</script>
