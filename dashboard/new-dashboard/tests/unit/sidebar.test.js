import { createPinia, setActivePinia } from "pinia"
import { expect, test, describe, beforeEach } from "vitest"
import { getInfoDataFrom } from "../../src/components/common/sideBar/InfoSidebarPerformance"
import { dbTypeStore } from "../../src/shared/dbTypes"

describe("InfoSideBar Test", () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  test("JBR", () => {
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
      },
      "ms",
      null
    )
    expect(result.branch).toEqual("241")
    expect(result.seriesName).toEqual("test")
    expect(result.machineName).toEqual("intellij-macos-hw-munit-713")
    expect(result.build).toEqual("1136.1")
    expect(result.buildId).toEqual(419534345)
    expect(result.deltaNext).toEqual("+23 (+0.8%)")
    expect(result.deltaPrevious).toEqual("+62 (+2.1%)")
    expect(result.date).toEqual("Dec 16, 2023, 11:45 PM")
    expect(result.projectName).toEqual("DaCapo_MacOS12x86_64OGL")
  })

  test("IntelliJ", () => {
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
      },
      "ms",
      null
    )
    expect(result.branch).toEqual("233")
    expect(result.seriesName).toEqual("test")
    expect(result.machineName).toEqual("intellij-linux-performance-aws-i-0bfabf7ea619cf696")
    expect(result.build).toEqual("233.14283")
    expect(result.buildId).toEqual(437084149)
    expect(result.installerId).toEqual(437052912)
    expect(result.deltaNext).toEqual("-736 ms (-11.2%)")
    expect(result.deltaPrevious).toEqual("-469 ms (-7.1%)")
    expect(result.date).toEqual("Jan 26, 2024, 1:23 AM")
    expect(result.projectName).toEqual("intellij_sources/vfsRefresh/with-1-thread(s)")
  })

  test("IntelliJ Dev", () => {
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
      },
      "ms",
      null
    )
    expect(result.branch).toEqual("master")
    expect(result.seriesName).toEqual("test")
    expect(result.machineName).toEqual("intellij-linux-performance-aws-i-0ed95da6a22b126e5")
    expect(result.buildId).toEqual(465279364)
    expect(result.deltaNext).toEqual("-59 s, 422 ms (-97.0%)")
    expect(result.deltaPrevious).toEqual("-50 s, 701 ms (-82.8%)")
    expect(result.date).toEqual("Mar 13, 2024, 5:34 PM")
    expect(result.projectName).toEqual("intellij_commit/vfsRefresh/git-status")
  })

  test("IntelliJ Startup", () => {
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
      },
      "ms",
      null
    )
    expect(result.branch).toEqual("master")
    expect(result.seriesName).toEqual("test")
    expect(result.machineName).toEqual("intellij-linux-hw-munit-095")
    expect(result.buildId).toEqual(440736226)
    expect(result.installerId).toEqual(440729452)
    expect(result.build).toEqual("241.11368")
    expect(result.deltaNext).toEqual("-161 ms (-2.4%)")
    expect(result.deltaPrevious).toEqual("+60 ms (+0.9%)")
    expect(result.date).toEqual("Feb 2, 2024, 10:02 AM")
    expect(result.projectName).toEqual("simple for IJ")
  })

  test("IntelliJ Startup Dev", () => {
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
      },
      "ms",
      null
    )
    expect(result.branch).toEqual("master")
    expect(result.seriesName).toEqual("test")
    expect(result.machineName).toEqual("intellij-linux-hw-munit-095")
    expect(result.buildId).toEqual(443891895)
    expect(result.deltaNext).toEqual("-1 s, 27 ms (-35.8%)")
    expect(result.deltaPrevious).toEqual("-937 ms (-32.6%)")
    expect(result.date).toEqual("Feb 8, 2024, 4:20 PM")
    expect(result.projectName).toEqual("simple for alfio")
  })

  test("Bazel", () => {
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
      },
      "ms",
      null
    )
    expect(result.branch).toEqual("master")
    expect(result.seriesName).toEqual("test")
    expect(result.machineName).toEqual("default-linux-aws-large-disk-A-i-007485cabc9cfacef")
    expect(result.buildId).toEqual(429902567)
    expect(result.deltaNext).toEqual("0 (0.0%)")
    expect(result.deltaPrevious).toEqual("+14 (+6.5%)")
    expect(result.date).toEqual("Jan 13, 2024, 2:10 AM")
    expect(result.projectName).toEqual("Synthetic 20000 project")
  })

  test("Qodana", () => {
    dbTypeStore().setDbType("qodana", "report")
    const result = getInfoDataFrom(
      {
        seriesName: "test",
        value: [
          1710484278000,
          650,
          "qodana-linux-amd64-xl1-A-i-0147baff0a791a4ea",
          466545404,
          "Byte_Buddy",
          "qodana-jvm:2023.2-nightly",
          {
            prev: 668,
            next: 554,
          },
        ],
      },
      "ms",
      null
    )
    expect(result.branch).toEqual("qodana-jvm:2023.2-nightly")
    expect(result.seriesName).toEqual("test")
    expect(result.machineName).toEqual("qodana-linux-amd64-xl1-A-i-0147baff0a791a4ea")
    expect(result.buildId).toEqual(466545404)
    expect(result.deltaNext).toEqual("-96 ms (-14.8%)")
    expect(result.deltaPrevious).toEqual("+18 ms (+2.8%)")
    expect(result.date).toEqual("Mar 15, 2024, 7:31 AM")
    expect(result.projectName).toEqual("Byte_Buddy")
  })

  test("Unit Test", () => {
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
      },
      "ms",
      null
    )
    expect(result.branch).toEqual("master")
    expect(result.seriesName).toEqual("test")
    expect(result.machineName).toEqual("intellij-linux-hw-hetzner-agent-17")
    expect(result.buildId).toEqual(458553579)
    expect(result.deltaNext).toEqual("0 (0.0%)")
    expect(result.deltaPrevious).toEqual("0 (0.0%)")
    expect(result.date).toEqual("Mar 3, 2024, 4:41 PM")
    expect(result.projectName).toEqual("com.intellij.codeInsight.JavaCommentByLineTest.testUncommentLargeFilePerformance - Uncommenting large file")
  })

  test("Fleet Startup", () => {
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
      },
      "ms",
      null
    )
    expect(result.branch).toEqual("master")
    expect(result.seriesName).toEqual("test")
    expect(result.machineName).toEqual("intellij-linux-hw-munit-095")
    expect(result.buildId).toEqual(435764822)
    expect(result.installerId).toEqual(435755382)
    expect(result.build).toEqual("1.31.4")
    expect(result.deltaNext).toEqual("-1 s, 280 ms (-22.6%)")
    expect(result.deltaPrevious).toEqual("0.00 ms (0.0%)")
    expect(result.date).toEqual("Jan 24, 2024, 6:53 AM")
    expect(result.projectName).toEqual("fleet")
  })

  test("Fleet Tests", () => {
    dbTypeStore().setDbType("fleet", "measure")
    const result = getInfoDataFrom(
      {
        seriesName: "test",
        value: [
          1706377569000,
          12702000000,
          "intellij-linux-hw-hetzner-agent-23",
          437946452,
          "multiCaretTyping",
          "master",
          {
            prev: 7404000000,
            next: 7905000000,
          },
        ],
      },
      "ns",
      null
    )
    expect(result.branch).toEqual("master")
    expect(result.seriesName).toEqual("test")
    expect(result.machineName).toEqual("intellij-linux-hw-hetzner-agent-23")
    expect(result.buildId).toEqual(437946452)
    expect(result.date).toEqual("Jan 27, 2024, 6:46 PM")
    expect(result.projectName).toEqual("multiCaretTyping")
    expect(result.deltaNext).toEqual("-4 s, 797 ms (-37.8%)")
    expect(result.deltaPrevious).toEqual("-5 s, 298 ms (-41.7%)")
  })
})
