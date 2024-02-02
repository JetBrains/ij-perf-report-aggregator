import { expect, test } from "vitest"
import { removeCommonSegments } from "../../src/util/removeCommonPrefixes"
import { replaceToLink } from "../../src/util/linkReplacer"

test("make link from YT issue", () => {
  expect(replaceToLink("IDEA-123")).toStrictEqual('<a class="underline decoration-dotted hover:no-underline" href="https://youtrack.jetbrains.com/issue/IDEA-123">IDEA-123</a>')
})

test("remove and replace YT link", () => {
  expect(replaceToLink("https://youtrack.jetbrains.com/issue/IDEA-317733")).toStrictEqual(
    '<a class="underline decoration-dotted hover:no-underline" href="https://youtrack.jetbrains.com/issue/IDEA-317733">IDEA-317733</a>'
  )
})

test("transform YT link with description", () => {
  expect(replaceToLink("https://youtrack.jetbrains.com/issue/AT-626/Migrate-Workspace-Model-Benchmarks-tests-to-use-common-approach-of-unit-perf-tests")).toStrictEqual(
    '<a class="underline decoration-dotted hover:no-underline" href="https://youtrack.jetbrains.com/issue/AT-626">AT-626</a> Migrate Workspace Model Benchmarks tests to use common approach of unit perf tests'
  )
})

test("transform YT link with description with prefix", () => {
  expect(
    replaceToLink("New SDK are enabled: https://youtrack.jetbrains.com/issue/IJPL-165/Migrate-implementation-of-Sdk-and-ProjectJdkTable-to-workspace-model-storage")
  ).toStrictEqual(
    '<a class="underline decoration-dotted hover:no-underline" href="https://youtrack.jetbrains.com/issue/IJPL-165">IJPL-165</a> Migrate implementation of Sdk and ProjectJdkTable to workspace model storage'
  )
})

test("transform slack", () => {
  expect(replaceToLink("Kotlin workspace model was enabled by default - https://jetbrains.slack.com/archives/C02JWL8P48K/p1705309556008919")).toStrictEqual(
    'Kotlin workspace model was enabled by default - <a class="underline decoration-dotted hover:no-underline" href="https://jetbrains.slack.com/archives/C02JWL8P48K/p1705309556008919">slack</a>'
  )
})

test("transform commit", () => {
  expect(replaceToLink("Speedup JPS Sync test https://jetbrains.team/p/ij/repositories/intellij/revision/69f4102715f053745592433a771385f28cde8e3d")).toStrictEqual(
    'Speedup JPS Sync test <a class="underline decoration-dotted hover:no-underline" href="https://jetbrains.team/p/ij/repositories/intellij/revision/69f4102715f053745592433a771385f28cde8e3d">69f4102715f053745592433a771385f28cde8e3d</a>'
  )
})
