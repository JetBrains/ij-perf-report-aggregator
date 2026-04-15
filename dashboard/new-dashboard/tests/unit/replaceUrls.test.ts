import { describe, expect, it } from "vitest"
import { replaceToLink } from "../../src/util/linkReplacer"

describe("link replacement", () => {
  it("make link from YT issue", () => {
    expect(replaceToLink("IDEA-123")).toBe('<a class="underline decoration-dotted hover:no-underline" href="https://youtrack.jetbrains.com/issue/IDEA-123">IDEA-123</a>')
  })

  it("remove and replace YT link", () => {
    expect(replaceToLink("https://youtrack.jetbrains.com/issue/IDEA-317733")).toBe(
      '<a class="underline decoration-dotted hover:no-underline" href="https://youtrack.jetbrains.com/issue/IDEA-317733">IDEA-317733</a>'
    )
  })

  it("transform YT link with description", () => {
    expect(replaceToLink("https://youtrack.jetbrains.com/issue/AT-626/Migrate-Workspace-Model-Benchmarks-tests-to-use-common-approach-of-unit-perf-tests")).toBe(
      '<a class="underline decoration-dotted hover:no-underline" href="https://youtrack.jetbrains.com/issue/AT-626">AT-626</a> Migrate Workspace Model Benchmarks tests to use common approach of unit perf tests'
    )
  })

  it("transform YT link with description with prefix", () => {
    expect(replaceToLink("New SDK are enabled: https://youtrack.jetbrains.com/issue/IJPL-165/Migrate-implementation-of-Sdk-and-ProjectJdkTable-to-workspace-model-storage")).toBe(
      '<a class="underline decoration-dotted hover:no-underline" href="https://youtrack.jetbrains.com/issue/IJPL-165">IJPL-165</a> Migrate implementation of Sdk and ProjectJdkTable to workspace model storage'
    )
  })

  it("transform slack", () => {
    expect(replaceToLink("Kotlin workspace model was enabled by default - https://jetbrains.slack.com/archives/C02JWL8P48K/p1705309556008919")).toBe(
      'Kotlin workspace model was enabled by default - <a class="underline decoration-dotted hover:no-underline" href="https://jetbrains.slack.com/archives/C02JWL8P48K/p1705309556008919">slack</a>'
    )
  })

  it("transform commit", () => {
    expect(replaceToLink("Speedup JPS Sync test https://jetbrains.team/p/ij/repositories/intellij/revision/69f4102715f053745592433a771385f28cde8e3d")).toBe(
      'Speedup JPS Sync test <a class="underline decoration-dotted hover:no-underline" href="https://jetbrains.team/p/ij/repositories/intellij/revision/69f4102715f053745592433a771385f28cde8e3d">69f4102715f053745592433a771385f28cde8e3d</a>'
    )
  })

  it("transform commit hash", () => {
    expect(replaceToLink("Project structure was changed: [python, ds, jupyter]: Migrate Python support to V2 Ilya Kazakevich 129 files 98f418c52d90")).toBe(
      'Project structure was changed: [python, ds, jupyter]: Migrate Python support to V2 Ilya Kazakevich 129 files <a class="underline decoration-dotted hover:no-underline" href="https://jetbrains.team/p/ij/repositories/ultimate/revision/98f418c52d90">98f418c52d90</a>'
    )
  })

  it("transform review", () => {
    expect(replaceToLink("New metrics were added https://jetbrains.team/p/ij/reviews/120177/timeline https://jetbrains.team/p/ij/reviews/120151/timeline")).toBe(
      'New metrics were added <a class="underline decoration-dotted hover:no-underline" href="https://jetbrains.team/p/ij/reviews/120177/timeline">review</a> <a class="underline decoration-dotted hover:no-underline" href="https://jetbrains.team/p/ij/reviews/120177/timeline">review</a>'
    )
  })
})
