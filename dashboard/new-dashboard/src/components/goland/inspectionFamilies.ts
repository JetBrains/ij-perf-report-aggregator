import type { MetricInfo } from "../../shared/metricsDescription"

export interface InspectionFamily {
  key: string
  label: string
  /** Inspection short names; at most 9, the size of the light palette, so that each keeps its color across charts. */
  inspections: string[]
}

/**
 * The inspections of GoLand Inspect Code, each in one family. The heaviest are the 9 that take 87–93% of the time on
 * every project of TeamCity build 1143364371 (2026-11-25); the others follow the groups of the inspection settings,
 * split where a group has more than 9, so that inspections that share a cost (reference search, type inference) share
 * a chart. An inspection missing here gets neither a chart nor a description.
 */
export const GO_INSPECTION_FAMILIES: InspectionFamily[] = [
  {
    key: "heaviest",
    label: "Heaviest",
    inspections: [
      "Annotator",
      "GoMaybeNil",
      "GoDfaInspectionRunner",
      "GoResourceLeak",
      "GoVetCopyLock",
      "GoBoolExpressions",
      "GoUnusedExportedFunction",
      "GoDeprecation",
      "GoVetIntToStringConversion",
    ],
  },
  {
    key: "unused",
    label: "Unused declarations",
    inspections: ["GoUnusedConst", "GoUnusedGlobalVariable", "GoUnusedFunction", "GoUnusedParameter", "GoUnusedType", "GoUnusedExportedType", "GoUnusedTypeParameter"],
  },
  {
    key: "redundant",
    label: "Redundant code",
    inspections: [
      "GoRedundantBlankArgInRange",
      "GoRedundantComma",
      "GoRedundantConversion",
      "GoRedundantImportAlias",
      "GoRedundantSecondIndexInSlices",
      "GoRedundantSemicolon",
      "GoRedundantTypeDeclInCompositeLit",
    ],
  },
  { key: "declarations", label: "Declaration redundancy", inspections: ["GoEmptyDeclaration", "GoPreferNilSlice", "GoSelfAssignment", "GoVarAndConstTypeMayBeOmitted"] },
  {
    key: "vet",
    label: "go vet",
    inspections: [
      "GoVetAtomic",
      "GoVetFailNowInNotTestGoroutine",
      "GoVetImpossibleInterfaceToInterfaceAssertion",
      "GoVetLostCancel",
      "GoVetStructTag",
      "GoVetUnmarshal",
      "GoVetUnsafePointer",
      "GoErrorsAs",
      "GoShift",
    ],
  },
  {
    key: "calls",
    label: "Probable bugs: calls and errors",
    inspections: [
      "GoDirectComparisonOfErrors",
      "GoTypeAssertionOnErrors",
      "GoUnhandledErrorResult",
      "GoUnusedCallResult",
      "GoPrintFunctions",
      "GoStringsReplaceCount",
      "GoStandardMethods",
      "GoDeferGo",
      "GoLoopClosure",
    ],
  },
  {
    key: "bugs",
    label: "Probable bugs: other",
    inspections: [
      "GoBuildTag",
      "GoDebugDirective",
      "GoLeadingWhitespaceInDirectiveComment",
      "GoImportUsedAsName",
      "GoReservedWordUsedAsName",
      "GoIrregularIota",
      "GoSwitchMissingCasesForIotaConsts",
      "GoMixedReceiverTypes",
      "GoDivisionByZero",
    ],
  },
  {
    key: "naming",
    label: "Code style: comments and names",
    inspections: [
      "GoCommentLeadingSpace",
      "GoCommentStart",
      "GoErrorStringFormat",
      "GoExportedOwnDeclaration",
      "GoNameStartsWithPackageName",
      "GoReceiverNames",
      "GoSnakeCaseUsage",
      "GoUnitSpecificDurationSuffix",
    ],
  },
  {
    key: "style",
    label: "Code style: other",
    inspections: ["GoRedundantElseInIf", "GoRedundantTrueInForCondition", "GoStructInitializationWithoutFieldNames", "GoUnsortedImport"],
  },
  { key: "controlFlow", label: "Control flow", inspections: ["GoAssignmentToReceiver", "GoDeferInLoop", "GoInfiniteFor", "GoUnreachableCode"] },
  {
    key: "other",
    label: "General",
    inspections: [
      "GoCyclicImports",
      "GoDisabledGopathIndexing",
      "GoExportedFuncWithUnexportedType",
      "GoMissingTrailingComma",
      "GoRedundantParens",
      "GoTestName",
      // Reported only on the projects whose Go SDK has runtime.AddCleanup.
      "GoDetectSetFinalizerUsages",
    ],
  },
  // Outside the Go section of the settings: Go Template, and the custom RegExp inspection of all languages.
  { key: "otherLanguages", label: "Templates and custom RegExp", inspections: ["GoTemplateDuplicateVariable", "GoTemplateUnknownVariable", "CustomRegExpInspection"] },
]

/** golangci-lint, run by the golangci-lint projects only; its time is the time of the external linter. */
export const GOLANGCI_LINT_FAMILY: InspectionFamily = { key: "golangciLint", label: "golangci-lint", inspections: ["GoLinter"] }

const TIME_DESCRIPTION =
  "Time of the inspection in Inspect Code, summed over files and threads. It includes the waits for the CPU and for the caches that the inspection fills for the others, so it finds the inspection behind a change of globalInspections, not its own cost. It counts whole milliseconds per file, so an inspection faster than 1 ms per file adds 0."

const PROBLEMS_DESCRIPTION = "Problems the inspection reported in Inspect Code. A time change that comes with a change here changed the results, not only the speed."

/** `globalInspections#<id>.ms` and `globalInspections#<id>.problems` of the inspections of all families. */
export const GO_INSPECTION_RUN_METRICS = new Map<string, MetricInfo>(
  [GOLANGCI_LINT_FAMILY, ...GO_INSPECTION_FAMILIES]
    .flatMap((family) => family.inspections)
    .flatMap((id): [string, MetricInfo][] => [
      [`globalInspections#${id}.ms`, { description: TIME_DESCRIPTION, unit: "milliseconds" }],
      [`globalInspections#${id}.problems`, { description: PROBLEMS_DESCRIPTION, unit: "counter", betterDirection: "stable" }],
    ])
)
