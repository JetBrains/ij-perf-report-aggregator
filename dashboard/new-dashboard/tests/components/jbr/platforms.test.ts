import { describe, expect, it } from "vitest"
import { compareConfigurations, configurationOf, platformOf } from "../../../src/components/jbr/platforms"

describe("JBR platforms", () => {
  it("takes the configuration after the benchmark name, keeping its own underscores", () => {
    expect(configurationOf("Render_Ubuntu2404x86_64MM_Wayland")).toBe("Ubuntu2404x86_64MM_Wayland")
    expect(configurationOf("SPECjbb2015_macOS15aarch64")).toBe("macOS15aarch64")
  })

  it("tells the platform from the start of the configuration", () => {
    expect(platformOf("Render_macOS15aarch64mini_Metal")).toBe("macOS")
    expect(platformOf("Render_Ubuntu2404aarch64RB_Wayland")).toBe("Linux")
    expect(platformOf("Render_Windows11x86_64")).toBe("Windows")
    expect(platformOf("Render_FreeBSD14")).toBe("Other")
  })

  it("orders macOS, Linux, Windows, then the rest, each by configuration", () => {
    const projects = [
      "Render_Windows11x86_64",
      "Render_FreeBSD14",
      "Render_Ubuntu2404x86_64MM_Wayland",
      "Render_macOS15x86_64OGL",
      "Render_Windows10x86_64",
      "Render_macOS15aarch64Metal",
    ]
    expect(projects.toSorted(compareConfigurations)).toStrictEqual([
      "Render_macOS15aarch64Metal",
      "Render_macOS15x86_64OGL",
      "Render_Ubuntu2404x86_64MM_Wayland",
      "Render_Windows10x86_64",
      "Render_Windows11x86_64",
      "Render_FreeBSD14",
    ])
  })
})
