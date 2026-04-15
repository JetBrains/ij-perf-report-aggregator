import { createPinia, setActivePinia } from "pinia"
import { expect, it, describe, beforeEach } from "vitest"
import { timeFormatWithoutSeconds } from "../../src/components/common/formatter"
import { getInfoDataFrom } from "../../src/components/common/sideBar/InfoSidebarPerformance"
import { dbTypeStore } from "../../src/shared/dbTypes"
import type { DefaultLabelFormatterCallbackParams as CallbackDataParams } from "echarts"

describe("InfoSideBar Test", () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it("parses JBR data", () => {
    dbTypeStore().setDbType("jbr", "report")
    const result = getInfoDataFrom(
      {
        seriesName: "test",
        value: [
          1702766751000,
          2987,
          "avrora",
          "c",
          "intellij-macos-hw-munit-713",
          419534345,
          "DaCapo_MacOS12x86_64OGL",
          "1136.1",
          "241",
          {
            prev: 3049,
            next: 3010,
          },
        ],
      } as CallbackDataParams,
      "ms",
      null,
      ""
    )
    expect(result).toMatchObject({
      branch: "241",
      seriesName: "test",
      machineName: "intellij-macos-hw-munit-713",
      build: "1136.1",
      buildId: 419534345,
      deltaNext: "+23 (+0.8%)",
      deltaPrevious: "+62 (+2.1%)",
      date: timeFormatWithoutSeconds.format(1702766751000),
      projectName: "DaCapo_MacOS12x86_64OGL",
    })
  })

  it("parses IntelliJ data", () => {
    dbTypeStore().setDbType("perfint", "ruby")
    const result = getInfoDataFrom(
      {
        seriesName: "test",
        value: [
          1706228613000,
          6560,
          "vfs_initial_refresh",
          "d",
          "intellij-linux-performance-aws-i-0bfabf7ea619cf696",
          437084149,
          "intellij_sources/vfsRefresh/with-1-thread(s)",
          437052912,
          233,
          14283,
          0,
          "233",
          {
            prev: 6091,
            next: 5824,
          },
        ],
      } as CallbackDataParams,
      "ms",
      null,
      ""
    )
    expect(result).toMatchObject({
      branch: "233",
      seriesName: "test",
      machineName: "intellij-linux-performance-aws-i-0bfabf7ea619cf696",
      build: "233.14283",
      buildId: 437084149,
      installerId: 437052912,
      deltaNext: "-736 ms (-11.2%)",
      deltaPrevious: "-469 ms (-7.1%)",
      date: timeFormatWithoutSeconds.format(1706228613000),
      projectName: "intellij_sources/vfsRefresh/with-1-thread(s)",
    })
  })

  it("parses IntelliJ Dev data", () => {
    dbTypeStore().setDbType("perfintDev", "ruby")
    const result = getInfoDataFrom(
      {
        seriesName: "test",
        value: [
          1710347658000,
          61234,
          "vfs_initial_refresh",
          "d",
          "intellij-linux-performance-aws-i-0ed95da6a22b126e5",
          465279364,
          "intellij_commit/vfsRefresh/git-status",
          "master",
          {
            prev: 10533,
            next: 1812,
          },
        ],
      } as CallbackDataParams,
      "ms",
      null,
      ""
    )
    expect(result).toMatchObject({
      branch: "master",
      seriesName: "test",
      machineName: "intellij-linux-performance-aws-i-0ed95da6a22b126e5",
      buildId: 465279364,
      deltaNext: "-59 s, 422 ms (-97.0%)",
      deltaPrevious: "-50 s, 701 ms (-82.8%)",
      date: timeFormatWithoutSeconds.format(1710347658000),
      projectName: "intellij_commit/vfsRefresh/git-status",
    })
  })

  it("parses IntelliJ Startup data", () => {
    dbTypeStore().setDbType("ij", "report")
    const result = getInfoDataFrom(
      {
        seriesName: "test",
        value: [
          1706864572000,
          6685,
          "reopenProjectPerformance/fusCodeVisibleInEditorDurationMs",
          "intellij-linux-hw-munit-095",
          440736226,
          "simple for IJ",
          440729452,
          241,
          11368,
          0,
          "master",
          {
            prev: 6745,
            next: 6524,
          },
        ],
      } as CallbackDataParams,
      "ms",
      null,
      ""
    )
    expect(result).toMatchObject({
      branch: "master",
      seriesName: "test",
      machineName: "intellij-linux-hw-munit-095",
      buildId: 440736226,
      installerId: 440729452,
      build: "241.11368",
      deltaNext: "-161 ms (-2.4%)",
      deltaPrevious: "+60 ms (+0.9%)",
      date: timeFormatWithoutSeconds.format(1706864572000),
      projectName: "simple for IJ",
    })
  })

  it("parses IntelliJ Startup Dev data", () => {
    dbTypeStore().setDbType("ijDev", "report")
    const result = getInfoDataFrom(
      {
        seriesName: "test",
        value: [
          1707405609000,
          2870,
          "editorRestoring",
          "intellij-linux-hw-munit-095",
          443891895,
          "simple for alfio",
          "master",
          {
            prev: 1933,
            next: 1843,
          },
        ],
      } as CallbackDataParams,
      "ms",
      null,
      ""
    )
    expect(result).toMatchObject({
      branch: "master",
      seriesName: "test",
      machineName: "intellij-linux-hw-munit-095",
      buildId: 443891895,
      deltaNext: "-1 s, 27 ms (-35.8%)",
      deltaPrevious: "-937 ms (-32.6%)",
      date: timeFormatWithoutSeconds.format(1707405609000),
      projectName: "simple for alfio",
    })
  })

  it("parses Bazel data", () => {
    dbTypeStore().setDbType("bazel", "report")
    const result = getInfoDataFrom(
      {
        seriesName: "test",
        value: [
          1705108226000,
          216,
          "apply.project.tree.view.fix.memory.mb",
          "c",
          "default-linux-aws-large-disk-A-i-007485cabc9cfacef",
          429902567,
          "Synthetic 20000 project",
          "master",
          {
            prev: 230,
            next: 216,
          },
        ],
      } as CallbackDataParams,
      "ms",
      null,
      ""
    )
    expect(result).toMatchObject({
      branch: "master",
      seriesName: "test",
      machineName: "default-linux-aws-large-disk-A-i-007485cabc9cfacef",
      buildId: 429902567,
      deltaNext: "0 (0.0%)",
      deltaPrevious: "+14 (+6.5%)",
      date: timeFormatWithoutSeconds.format(1705108226000),
      projectName: "Synthetic 20000 project",
    })
  })

  it("parses Qodana data", () => {
    dbTypeStore().setDbType("qodana", "report")
    const result = getInfoDataFrom(
      {
        seriesName: "test",
        value: [
          1710484278000,
          650,
          "metric_name",
          "d",
          "qodana-linux-amd64-xl1-A-i-0147baff0a791a4ea",
          466545404,
          "Byte_Buddy",
          "qodana-jvm:2023.2-nightly",
          {
            prev: 668,
            next: 554,
          },
        ],
      } as CallbackDataParams,
      "ms",
      null,
      ""
    )
    expect(result).toMatchObject({
      branch: "qodana-jvm:2023.2-nightly",
      seriesName: "test",
      machineName: "qodana-linux-amd64-xl1-A-i-0147baff0a791a4ea",
      buildId: 466545404,
      deltaNext: "-96 ms (-14.8%)",
      deltaPrevious: "+18 ms (+2.8%)",
      date: timeFormatWithoutSeconds.format(1710484278000),
      projectName: "Byte_Buddy",
    })
  })

  it("parses unit test data", () => {
    dbTypeStore().setDbType("perfUnitTests", "report")
    const result = getInfoDataFrom(
      {
        seriesName: "test",
        value: [
          1709480472000,
          3,
          "attempt.count",
          "c",
          "intellij-linux-hw-hetzner-agent-17",
          458553579,
          "com.intellij.codeInsight.JavaCommentByLineTest.testUncommentLargeFilePerformance - Uncommenting large file",
          "master",
          {
            prev: 3,
            next: 3,
          },
        ],
      } as CallbackDataParams,
      "ms",
      null,
      ""
    )
    expect(result).toMatchObject({
      branch: "master",
      seriesName: "test",
      machineName: "intellij-linux-hw-hetzner-agent-17",
      buildId: 458553579,
      deltaNext: "0 (0.0%)",
      deltaPrevious: "0 (0.0%)",
      date: timeFormatWithoutSeconds.format(1709480472000),
      projectName: "com.intellij.codeInsight.JavaCommentByLineTest.testUncommentLargeFilePerformance - Uncommenting large file",
    })
  })

  it("parses Fleet Startup data", () => {
    dbTypeStore().setDbType("fleet", "report")
    const result = getInfoDataFrom(
      {
        seriesName: "test",
        value: [
          1706075619000,
          5669,
          "editor appeared.end",
          "intellij-linux-hw-munit-095",
          435764822,
          "fleet",
          435755382,
          1,
          31,
          4,
          "master",
          {
            prev: 5669,
            next: 4389,
          },
        ],
      } as CallbackDataParams,
      "ms",
      null,
      ""
    )
    expect(result).toMatchObject({
      branch: "master",
      seriesName: "test",
      machineName: "intellij-linux-hw-munit-095",
      buildId: 435764822,
      installerId: 435755382,
      build: "1.31.4",
      deltaNext: "-1 s, 280 ms (-22.6%)",
      deltaPrevious: "0 (0.0%)",
      date: timeFormatWithoutSeconds.format(1706075619000),
      projectName: "fleet",
    })
  })
})
