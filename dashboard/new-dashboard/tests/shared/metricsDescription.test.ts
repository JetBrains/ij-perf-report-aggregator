import { describe, expect, it } from "vitest"
import { getMetricDescription } from "../../src/shared/metricsDescription"
import { formatMeasureValue, resolveMeasureUnit } from "../../src/components/common/formatter"

describe("metricsDescription after-# sub-metric resolution", () => {
  it("resolves a #jvm.alloc.mb sub-metric to mebibytes / lower", () => {
    const info = getMetricDescription("coldStartHighlighting_fast_importGraph#jvm.alloc.mb")
    expect(info?.unit).toBe("mebibytes")
    expect(info?.betterDirection).toBe("lower")
  })

  it("resolves the GC count sub-metric to a counter", () => {
    expect(getMetricDescription("warmStartHighlighting_slow_httpType#jvm.gc.count")?.unit).toBe("counter")
  })

  it("resolves the CPU-time sub-metric to milliseconds", () => {
    expect(getMetricDescription("typingHighlighting_medium_userStore#jvm.cpu.time.ms")?.unit).toBe("milliseconds")
  })

  it("still resolves an existing before-# wildcard (no regression)", () => {
    // lexingSize#100 must keep hitting the "lexingSize#*" wildcard, not the new after-# fallback.
    expect(getMetricDescription("lexingSize#100")?.unit).toBe("bytes")
  })

  it("surfaces the declared sub-metric unit through resolveMeasureUnit", () => {
    expect(resolveMeasureUnit("coldStartHighlighting_slow_generatedPb#jvm.alloc.mb")).toBe("mebibytes")
    expect(resolveMeasureUnit("coldStartHighlighting_slow_generatedPb#jvm.gc.count")).toBe("counter")
    expect(resolveMeasureUnit("coldStartHighlighting_slow_generatedPb#jvm.cpu.time.ms")).toBe("milliseconds")
  })

  it("formats a large MiB allocation as GiB-scaled binary size", () => {
    // kubernetes slow allocation churn (~26 970 MiB) must render in GiB, not milliseconds.
    expect(formatMeasureValue(26970, resolveMeasureUnit("coldStartHighlighting_slow_generatedPb#jvm.alloc.mb"))).toMatch(/GiB$/)
  })
})
