import { expect, test } from "vitest"
import { removeCommonSegments } from "../../src/util/removeCommonPrefixes"

test("remove common suffix", () => {
  expect(removeCommonSegments(["foo/indexing", "bar/indexing"])).toStrictEqual(["foo", "bar"])
})

test("remove common part from the middle", () => {
  expect(removeCommonSegments(["foo/indexing/bar", "bar1/indexing/foo1"])).toStrictEqual(["foo/bar", "bar1/foo1"])
})

test("remove common part from the middle and end", () => {
  expect(removeCommonSegments(["foo/indexing", "bar/indexing/foo", "test/indexing"])).toStrictEqual(["foo", "bar/foo", "test"])
})

test("don't remove starting segment", () => {
  expect(removeCommonSegments(["grails/showIntentions/Find cause", "grails/foo/bar"])).toStrictEqual(["grails/showIntentions/Find cause", "grails/foo/bar"])
})

test("common segment is not replaced everywhere", () => {
  expect(removeCommonSegments(["typingInJavaFile_4Threads/typing", "typingInKotlinFile_16Threads/typing", "typingInKotlinFile_4Threads/typing"])).toStrictEqual([
    "typingInJavaFile_4Threads",
    "typingInKotlinFile_16Threads",
    "typingInKotlinFile_4Threads",
  ])
})
