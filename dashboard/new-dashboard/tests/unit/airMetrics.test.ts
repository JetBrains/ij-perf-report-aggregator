import { describe, expect, it } from "vitest"
import { formatMeasureValue, resolveMeasureUnit } from "../../src/components/common/formatter"

describe("Air metric units", () => {
  it("reads the request size of the ACP adapter as bytes", () => {
    expect(resolveMeasureUnit("continue.XXL.requestBytes")).toBe("bytes")
    expect(resolveMeasureUnit("live.second.backgroundBytes")).toBe("bytes")
    expect(formatMeasureValue(3 * 1024 * 1024, "bytes")).not.toContain("ms")
  })

  it("reads the memory of the ACP adapter and of the IDE as mebibytes", () => {
    expect(resolveMeasureUnit("open.M.agentRssPeakMiB")).toBe("mebibytes")
    expect(resolveMeasureUnit("live.agentRssPeakMiB")).toBe("mebibytes")
    expect(resolveMeasureUnit("heap.XL.heapLive.MiB")).toBe("mebibytes")
    expect(resolveMeasureUnit("sessionList.heapLive.MiB")).toBe("mebibytes")
  })

  it("keeps the times of the ACP adapter as milliseconds", () => {
    expect(resolveMeasureUnit("live.agentMs")).toBe("milliseconds")
    expect(resolveMeasureUnit("list.initialize.coldMs")).toBe("milliseconds")
  })
})
