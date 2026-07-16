import { describe, expect, it } from "vitest"
import { getMachineGroupName } from "../../src/configurators/MachineConfigurator"

describe("machine group mapping", () => {
  it("maps known agent classes", () => {
    expect(getMachineGroupName("intellij-linux-performance-aws-i-08aec6c8ee5a71bba")).toBe("Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)")
    expect(getMachineGroupName("intellij-linux-performance-aws-lt-a-i-045667485579a157a")).toBe("Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)")
    expect(getMachineGroupName("intellij-linux-performance-tiny-aws-on-demand-i-0abc12345")).toBe("Linux EC2 C6id.xlarge (4 vCPU Xeon, 8 GB)")
    expect(getMachineGroupName("intellij-windows-performance-mem-aws-i-0deadbeef")).toBe("Windows EC2 C6id.4xlarge or i4i.4xlarge (16 vCPU Xeon, 32 or 128 GB)")
    expect(getMachineGroupName("intellij-macos-unit-2200-large-10298")).toBe("mac large")
  })

  it("supports the regex rule", () => {
    expect(getMachineGroupName("ij-w11u-azr7")).toBe("windows-azure")
  })

  it("returns Unknown for unmapped names", () => {
    expect(getMachineGroupName("something-totally-unknown")).toBe("Unknown")
  })

  it("keeps powerful and weak linux classes distinct", () => {
    expect(getMachineGroupName("intellij-linux-performance-aws-i-08aec6c8ee5a71bba")).not.toBe(getMachineGroupName("intellij-linux-performance-tiny-aws-on-demand-i-0abc12345"))
  })
})
