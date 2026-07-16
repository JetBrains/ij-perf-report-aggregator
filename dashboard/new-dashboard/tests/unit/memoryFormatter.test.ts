import { describe, expect, it } from "vitest"
import { formatMeasureValue, resolveMeasureUnit } from "../../src/components/common/formatter"

describe("measure unit resolution", () => {
  it.each([
    ["rd.memory.workingSetMb", "mebibytes"],
    ["cacheKb", "kibibytes"],
    ["MEM.avgRamMegabytes", "mebibytes"],
    ["JVM.maxHeapMegabytes", "mebibytes"],
    ["Memory | IDE | RESIDENT SIZE (MB) 95th pctl", "mebibytes"],
    ["Lucene Files - Index Size (KB) - Cold", "kibibytes"],
  ])("detects %s as a binary memory size (%s)", (measureName, unit) => {
    expect(resolveMeasureUnit(measureName)).toBe(unit)
  })

  it.each([
    ["freedMemoryByGC", "mebibytes"],
    ["totalHeapUsedMax", "mebibytes"],
    ["freedMemory", "mebibytes"],
    ["indexSize", "kilobytes"],
    ["processingSpeedAvg#PHP", "kilobytesPerSecond"],
    ["lexingSpeed#Kotlin", "kilobytesPerSecond"],
    ["parsingSize#Java", "bytes"],
    ["numberOfIndexedFiles", "counter"],
    ["gcPauseCount", "counter"],
  ])("resolves %s to its declared unit (%s)", (measureName, unit) => {
    expect(resolveMeasureUnit(measureName)).toBe(unit)
  })

  it.each([
    ["fileCount", "counter"],
    ["undeclaredAction", "milliseconds"],
  ])("falls back to a name-based unit for the undeclared metric %s (%s)", (measureName, unit) => {
    expect(resolveMeasureUnit(measureName)).toBe(unit)
  })

  it("lets a declared unit win over the chart value-unit", () => {
    expect(resolveMeasureUnit("processingSpeedAvg#PHP", { valueUnit: "counter" })).toBe("kilobytesPerSecond")
  })

  it("lets an explicit value-unit override an undeclared metric", () => {
    expect(resolveMeasureUnit("undeclaredAction", { valueUnit: "counter" })).toBe("counter")
    expect(resolveMeasureUnit("fileCount", { valueUnit: "ms" })).toBe("milliseconds")
    expect(resolveMeasureUnit("anything", { valueUnit: "ns" })).toBe("nanoseconds")
  })

  it("renders a physical quantity as a counter while scaling", () => {
    expect(resolveMeasureUnit("MEM.avgRamMegabytes", { scaling: true })).toBe("counter")
    expect(resolveMeasureUnit("processingSpeedAvg#PHP", { scaling: true })).toBe("counter")
  })

  it("lets a clear count name win over a mis-typed stored duration", () => {
    // The perf pipeline stores this count as "d" by mistake; the "Count" name must still win.
    expect(resolveMeasureUnit("fileCount", { storedType: "d" })).toBe("counter")
    expect(resolveMeasureUnit("classLoadingLoadedCount", { storedType: "d" })).toBe("counter")
  })

  it("honours the stored type for names that are not obviously a count or size", () => {
    expect(resolveMeasureUnit("undeclaredAction", { storedType: "c" })).toBe("counter")
    expect(resolveMeasureUnit("undeclaredAction", { storedType: "d" })).toBe("milliseconds")
  })

  // The perf pipeline stores sizes and plain counts alike as "c", so a memory name must still be
  // detected as a size even when the stored type is "c" (the condition every production chart passes).
  it.each([
    ["rd.memory.workingSetMb", "mebibytes"],
    ["classLoadingMetrics/totalSizeKb", "kibibytes"],
    ["Memory | IDE | RESIDENT SIZE (MB) 95th pctl", "mebibytes"],
  ])("detects the counter-stored memory metric %s as a size (%s)", (measureName, unit) => {
    expect(resolveMeasureUnit(measureName, { storedType: "c" })).toBe(unit)
  })

  // A series' measureName is a composite of producer parts joined by " – " (machine/branch/dimension
  // tokens precede the metric), so the declared unit must resolve from the metric token, not the whole
  // string. These are the production shapes that previously fell through to "counter".
  it.each([
    ["master – processingSpeedAvg#Go – processingSpeedAvg#Go", "kilobytesPerSecond"],
    ["master – freedMemoryByGC – freedMemoryByGC", "mebibytes"],
    ["master – indexSize – indexSize", "kilobytes"],
  ])("resolves the declared unit from a composite measure name %s (%s)", (measureName, unit) => {
    expect(resolveMeasureUnit(measureName, { storedType: "c" })).toBe(unit)
  })
})

describe("measure value formatting", () => {
  it.each([
    // binary memory in IEC units
    [512, "mebibytes", "512 MiB"],
    [2048, "mebibytes", "2 GiB"],
    [0.5, "mebibytes", "512 KiB"],
    [0, "mebibytes", "0 MiB"],
    [2048, "kibibytes", "2 MiB"],
    [2048, "bytes", "2 KiB"],
    // decimal sizes and throughput in SI units
    [512, "kilobytes", "512 kB"],
    [2000, "kilobytes", "2 MB"],
    [800, "kilobytesPerSecond", "800 kB/s"],
    [2_000_000, "kilobytesPerSecond", "2 GB/s"],
    [0, "counter", "0"],
  ] as const)("formats %d %s as %s", (value, unit, expected) => {
    expect(formatMeasureValue(value, unit)).toBe(expected)
  })

  // The decimal separator follows the runtime locale, so match instead of comparing exactly.
  it.each([
    [1536, "mebibytes", /^1[.,]5 GiB$/],
    [1500, "kilobytesPerSecond", /^1[.,]5 MB\/s$/],
  ] as const)("formats %d %s with a locale decimal separator", (value, unit, pattern) => {
    expect(formatMeasureValue(value, unit)).toMatch(pattern)
  })
})
