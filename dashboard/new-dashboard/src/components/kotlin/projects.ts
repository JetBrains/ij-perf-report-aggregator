/**
 * A project category is a project name prefix such as "intellij_commit/" and "kotlin_lang/" with an associated, human-readable label.
 */
interface ProjectCategory {
  label: string
  prefix: string
}

export const projectCategories: ProjectCategory[] = [
  buildCategory("Hello World", "kotlin_empty/"),
  buildCategory("IntelliJ", "intellij_commit/"),
  buildCategory("Kotlin Compiler", "kotlin_lang/"),
  buildCategory("Kotlin Language Server", "kotlin_language_server/"),
  buildCategory("Toolbox Enterprise", "toolbox_enterprise/"),
  buildCategory("Spring Framework", "spring-framework/"),
  buildCategory("Rust Plugin", "rust_commit/"),
  buildCategory("Kotlin PetClinic", "kotlin_petclinic/"),

  // The `ktor_samples` category is open by design (hence no closing `/`).
  buildCategory("Ktor Samples", "ktor_samples"),

  buildCategory("LeakCanary", "leak-canary-android/"),
  buildCategory("Arrow", "arrow/"),
  buildCategory("Empty Script (.kts)", "kotlin_empty_kts/"),
]

function buildCategory(label: string, prefix: string) {
  return { label, prefix }
}

/**
 * Encapsulates all data which is needed to render a chart on one of the Kotlin dev dashboards. For each chart definition, a chart for K1
 * and a chart for K2 will be rendered.
 *
 * The project data is also used by the K1 vs. K2 comparison dashboard to find out all performance test names.
 */
export interface ProjectsChartDefinition {
  /**
   * The label of the chart. A "K1" or "K2" qualifier will be appended to the label.
   */
  label: string

  measure: string

  /**
   * All project names in the chart *without* their `_k1` and `_k2` suffixes.
   */
  projects: string[]
}

export const completionProjects = {
  kotlinEmpty: ["kotlin_empty/completion/empty_place_with_library_cache", "kotlin_empty/completion/empty_place_typing_with_library_cache"],
  intelliJ: [
    "intellij_commit/completion/empty_place_with_library_cache",
    "intellij_commit/completion/after_parameter_with_library_cache",
    "intellij_commit/completion/empty_place_typing_with_library_cache",
  ],
  intelliJ2: [
    "intellij_commit/completion/KotlinHighLevelFunctionParameterInfoHandler_emptyPlace_updateUIOrFail_typing_with_library_cache",
    "intellij_commit/completion/KotlinHighLevelFunctionParameterInfoHandler_emptyPlace_updateUIOrFail_with_library_cache",
    "intellij_commit/completion/KtOCSwiftChangeSignatureTest_emptyPlace_changeReturnType_typing_with_library_cache",
    "intellij_commit/completion/KtOCSwiftChangeSignatureTest_emptyPlace_changeReturnType_with_library_cache",
  ],
  intelliJTyping2: [
    "intellij_commit/completion/IdeMenuBar_emptyPlace_sout_typing_with_library_cache",
    "intellij_commit/completion/TestModelParser_emptyPlace_if_typing_with_library_cache",
    "intellij_commit/completion/AndroidModuleSystem_emptyPlace_get_typing_with_library_cache",
  ],
  kotlinLang: [
    "kotlin_lang/completion/after_parameter_with_library_cache",
    "kotlin_lang/completion/empty_place_with_library_cache",
    "kotlin_lang/completion/empty_place_typing_with_library_cache",
  ],
  kotlinLanguageServer: [
    "kotlin_language_server/completion/Completions_emptyPlace_completions_typing_with_library_cache",
    "kotlin_language_server/completion/QuickFixesTest_emptyPlace_completions_typing_with_library_cache",
    "kotlin_language_server/completion/Completions_emptyPlace_completions_with_library_cache",
    "kotlin_language_server/completion/QuickFixesTest_emptyPlace_completions_with_library_cache",
  ],
  tbe: [
    "toolbox_enterprise/completion/ProfileController_constructor_typing_with_library_cache",
    "toolbox_enterprise/completion/ProfileController_dataclass_typing_with_library_cache",
    "toolbox_enterprise/completion/ProfileController_fun_typing_with_library_cache",
    "toolbox_enterprise/completion/ProfileServiceTest_body_with_library_cache",
    "toolbox_enterprise/completion/ProfileServiceTest_constructor_typing_with_library_cache",
    "toolbox_enterprise/completion/ProfileServiceTest_constructor_with_library_cache",
    "toolbox_enterprise/completion/ProfileServiceTest_emptyPlace_FileEnd_typing_with_library_cache",
    "toolbox_enterprise/completion/ProfileServiceTest_emptyPlace_FileEnd_with_library_cache",
  ],
  tbeCaseWithAssert: [
    "toolbox_enterprise/completion/IntelliJPluginSettingTest_assert_0_with_library_cache",
    "toolbox_enterprise/completion/IntelliJPluginSettingTest_assert_2_typing_with_library_cache",
    "toolbox_enterprise/completion/IntelliJPluginSettingTest_assert_4_typing_with_library_cache",
    "toolbox_enterprise/completion/IntelliJPluginSettingTest_assert_7_typing_with_library_cache",
    "toolbox_enterprise/completion/IntelliJPluginSettingTest_assert_10_typing_with_library_cache",
  ],
  kotlinScript: [
    "arrow/completion/completion_kts_with_library_cache",
    "kotlin_empty_kts/completion/completion_kts_with_library_cache",
    "kotlin_lang/completion/completion_kts_with_library_cache",
  ],
  kotlinLanguageServerEvaluateExpression: [
    "kotlin_language_server/completion/Completions_emptyPlace_completions_typing_with_library_cache",
    "kotlin_language_server/completion/Completions_emptyPlace_completions_with_library_cache",
    "kotlin_language_server/completion/QuickFixesTest_emptyPlace_completions_typing_with_library_cache",
    "kotlin_language_server/completion/QuickFixesTest_emptyPlace_completions_with_library_cache",
  ],
}

export const completionCharts: ProjectsChartDefinition[] = [
  ...generateCompletionDefinitions("'Hello-world'", completionProjects.kotlinEmpty),
  ...generateCompletionDefinitions("'IntelliJ'", completionProjects.intelliJ),
  ...generateCompletionDefinitions("'IntelliJ suite 2'", completionProjects.intelliJ2),
  ...generateCompletionDefinitions("'IntelliJ with typing suite 2'", completionProjects.intelliJTyping2),
  ...generateCompletionDefinitions("'Kotlin lang'", completionProjects.kotlinLang),
  ...generateCompletionDefinitions("'Kotlin language server'", completionProjects.kotlinLanguageServer),
  ...generateCompletionDefinitions("'Toolbox Enterprise (TBE)'", completionProjects.tbe),
  ...generateCompletionDefinitions("'Toolbox Enterprise (TBE) different length'", completionProjects.tbeCaseWithAssert),
  ...generateCompletionDefinitions("'Kotlin script'", completionProjects.kotlinScript),
]

function generateCompletionDefinitions(labelPrefix: string, projects: string[]): ProjectsChartDefinition[] {
  return [
    {
      label: `${labelPrefix} completion mean value`,
      measure: "completion#mean_value",
      projects,
    },
    {
      label: `${labelPrefix} firstElementShown mean value`,
      measure: "completion#firstElementShown#mean_value",
      projects,
    },
  ]
}

/**
 * Highlighting projects are also the projects for local inspections.
 */
export const highlightingProjects = {
  kotlinEmpty: ["kotlin_empty/highlight/Main_with_library_cache"],
  intelliJ: [
    "intellij_commit/highlight/KtOCSwiftChangeSignatureTest_with_library_cache",
    "intellij_commit/highlight/KotlinHighLevelFunctionParameterInfoHandler_with_library_cache",
    "intellij_commit/highlight/JdkList_with_library_cache",
    "intellij_commit/highlight/ComposeCompletionContributorTest_with_library_cache",
    "intellij_commit/highlight/AgpUpgradeRefactoringProcessor_with_library_cache",
  ],
  intelliJ2: [
    "intellij_commit/highlight/AndroidModelTest_with_library_cache",
    "intellij_commit/highlight/SecureWireOverStreamTransport_with_library_cache",
    "intellij_commit/highlight/DexInlineCallStackComparisonTest_with_library_cache",
    "intellij_commit/highlight/DexLocalVariableTableBreakpointTest_with_library_cache",
    "intellij_commit/highlight/OraIntrospector_with_library_cache",
    "intellij_commit/highlight/SolutionModel.Generated_with_library_cache",
    "intellij_commit/highlight/UIAutomationInteractionModel.Generated_with_library_cache",
    "kotlin_lang/highlight/LazyJVM_with_library_cache",
  ],
  kotlinLang: [
    "kotlin_lang/highlight/CommonParser_with_library_cache",
    "kotlin_lang/highlight/FirErrors_with_library_cache",
    "kotlin_lang/highlight/Flag_with_library_cache",
    "kotlin_lang/highlight/KtFirDataClassConverters_with_library_cache",
    "kotlin_lang/highlight/DefaultArgumentStubGenerator_with_library_cache",
    "kotlin_lang/highlight/RawFirBuilder_with_library_cache",
  ],
  kotlinLanguageServer: [
    "kotlin_language_server/highlight/Compiler_with_library_cache",
    "kotlin_language_server/highlight/Completions_with_library_cache",
    "kotlin_language_server/highlight/CompletionsTest_with_library_cache",
    "kotlin_language_server/highlight/JavaElementConverter_with_library_cache",
    "kotlin_language_server/highlight/KotlinTextDocumentService_with_library_cache",
    "kotlin_language_server/highlight/QuickFixesTest_with_library_cache",
    "kotlin_language_server/highlight/SourcePath_with_library_cache",
  ],
  tbe: [
    "toolbox_enterprise/highlight/IdeSettingControllerTest_with_library_cache",
    "toolbox_enterprise/highlight/IntelliJPluginSettingTest_with_library_cache",
    "toolbox_enterprise/highlight/LoginTests_with_library_cache",
    "toolbox_enterprise/highlight/PluginAuditLogService_with_library_cache",
    "toolbox_enterprise/highlight/PluginControllerTest_with_library_cache",
    "toolbox_enterprise/highlight/ProfileController_with_library_cache",
    "toolbox_enterprise/highlight/ProfileService_with_library_cache",
    "toolbox_enterprise/highlight/ProfileServiceTest_with_library_cache",
    "toolbox_enterprise/highlight/SecurityTests_with_library_cache",
    "toolbox_enterprise/highlight/UsageDataFlowTests_with_library_cache",
    "toolbox_enterprise/highlight/VmOptionSettingTest_with_library_cache",
  ],
  ktorSamples: [
    "ktor_samples_mongodb/highlight/ApplicationTest_with_library_cache",
    "ktor_samples_httpbin/highlight/HttpBinApplication_with_library_cache",
    "ktor_samples_youkube/highlight/Upload_with_library_cache",
    "ktor_samples_youkube/highlight/Videos_with_library_cache",
    "ktor_samples_location-header/highlight/LocationHeaderApplication_with_library_cache",
    "ktor_samples_reverse-proxy/highlight/ReverseProxyApplication_with_library_cache",
    "ktor_samples_sse/highlight/SseApplication_with_library_cache",
  ],
  androidCanaryLeak: [
    "leak-canary-android/highlight/ByteArrayTimSort_with_library_cache",
    "leak-canary-android/highlight/HeapObject_with_library_cache",
    "leak-canary-android/highlight/HprofInMemoryIndex_with_library_cache",
    "leak-canary-android/highlight/HprofRecordReader_with_library_cache",
    "leak-canary-android/highlight/HprofWriter_with_library_cache",
    "leak-canary-android/highlight/LeakStatusTest_with_library_cache",
    "leak-canary-android/highlight/Neo4JCommand_with_library_cache",
  ],
  kotlinScript: [
    "arrow/highlight/build.gradle_with_library_cache",
    "kotlin_empty_kts/highlight/build.gradle_with_library_cache",
    "kotlin_lang/highlight/build.gradle_with_library_cache",
  ],
  removedImports: [
    "toolbox_enterprise/highlight/removedImports/IdeSettingControllerTest_with_library_cache",
    "intellij_commit/highlight/removedImports/UIAutomationInteractionModel.Generated_with_library_cache",
    "kotlin_language_server/highlight/removedImports/CompletionsTest_with_library_cache",
    "kotlin_lang/highlight/removedImports/DefaultArgumentStubGenerator_with_library_cache",
  ],
  springFramework: [
    "spring-framework/highlight/BeanDefinitionDsl_with_library_cache",
    "spring-framework/highlight/CoRouterFunctionDsl_with_library_cache",
    "spring-framework/highlight/reactive/RouterFunctionDsl_with_library_cache",
    "spring-framework/highlight/servlet/RouterFunctionDsl_with_library_cache",
    "spring-framework/highlight/StatusResultMatchersDsl_with_library_cache",
  ],
  rustPlugin: [
    "rust_commit/highlight/CargoBuildManagerTest_with_library_cache",
    "rust_commit/highlight/CargoWorkspace_with_library_cache",
    "rust_commit/highlight/ClippyLints_with_library_cache",
    "rust_commit/highlight/MacroExpansionManager_with_library_cache",
    "rust_commit/highlight/MemoryCategorization_with_library_cache",
    "rust_commit/highlight/Processors_with_library_cache",
    "rust_commit/highlight/RsErrorAnnotatorTest_with_library_cache",
    "rust_commit/highlight/RsReferenceImporter_with_library_cache",
    "rust_commit/highlight/StubImplementations_with_library_cache",
    "rust_commit/highlight/TypeInferenceWalker_with_library_cache",
    "rust_commit/highlight/TypeInference_with_library_cache",
  ],
}

const highlightingLabelsAndProjects = [
  { label: "'Kotlin empty'", projects: highlightingProjects.kotlinEmpty },
  { label: "'IntelliJ'", projects: highlightingProjects.intelliJ },
  { label: "'IntelliJ suite 2'", projects: highlightingProjects.intelliJ2 },
  { label: "'Kotlin lang'", projects: highlightingProjects.kotlinLang },
  { label: "'Kotlin language server'", projects: highlightingProjects.kotlinLanguageServer },
  { label: "'Toolbox Enterprise (TBE)'", projects: highlightingProjects.tbe },
  { label: "'Ktor samples'", projects: highlightingProjects.ktorSamples },
  { label: "'Android canary leak'", projects: highlightingProjects.androidCanaryLeak },
  { label: "'Kotlin script'", projects: highlightingProjects.kotlinScript },
  { label: "'Spring framework'", projects: highlightingProjects.springFramework },
  { label: "'Rust plugin'", projects: highlightingProjects.rustPlugin },
  { label: "'Files with removed imports'", projects: highlightingProjects.removedImports },
]

export const highlightingCharts: ProjectsChartDefinition[] = highlightingLabelsAndProjects.map((v) => generateHighlightingDefinition(v.label, v.projects))

function generateHighlightingDefinition(labelPrefix: string, projects: string[]): ProjectsChartDefinition {
  return {
    label: `${labelPrefix} semantic highlighting mean value`,
    measure: "semanticHighlighting#mean_value",
    projects,
  }
}

export const codeAnalysisCharts: ProjectsChartDefinition[] = highlightingLabelsAndProjects.map((v) => generateCodeAnalysisChartsDefinition(v.label, v.projects))

function generateCodeAnalysisChartsDefinition(labelPrefix: string, projects: string[]): ProjectsChartDefinition {
  return {
    label: `${labelPrefix} Code Analysis mean value`,
    measure: "localInspections#mean_value",
    projects,
  }
}

export const refactoringProjects = [
  "intellij_commit/rename/SqlBlock_SqlBlockRenamed",
  "kotlin_language_server/insertCode/Rename_renameSymbol",
  "kotlin_language_server/insertCode/SpecialJavaFileForTest_j2k",
]
export const optimizeImportsProjects = [
  "intellij_commit/otimizeImports/AbstractKotlinMavenImporterTest",
  "intellij_commit/otimizeImports/FSD",
  "intellij_commit/otimizeImports/DiagramsModel.Generated",
  "intellij_commit/otimizeImports/OraIntrospector",
  "intellij_commit/otimizeImports/QuickFixRegistrar",
  "intellij_commit/otimizeImports/SwiftTypeAssignabilityTest",
  "intellij_commit/otimizeImports/TerraformConfigCompletionContributor",
]

export const refactoringCharts: ProjectsChartDefinition[] = [
  generateRefactoringDefinition("PerformInlineRename", "performInlineRename#mean_value", refactoringProjects),
  generateRefactoringDefinition("StartInlineRename", "startInlineRename#mean_value", refactoringProjects),
  generateRefactoringDefinition("PrepareForRename", "prepareForRename#mean_value", refactoringProjects),
  generateRefactoringDefinition("Optimize imports", "execute_editor_optimizeimports", optimizeImportsProjects),
]

function generateRefactoringDefinition(labelPrefix: string, measure: string, projectsData: string[]): ProjectsChartDefinition {
  return {
    label: `${labelPrefix} mean value`,
    measure,
    projects: projectsData,
  }
}

export const findUsagesProjects = {
  intelliJ: [
    "intellij_commit/findUsages/loadModuleEntity_with_library_cache",
    "intellij_commit/findUsages/setUp_with_library_cache",
    "intellij_commit/findUsages/SolutionModel_with_library_cache",
    "intellij_commit/findUsages/SqlBlock_with_library_cache",
    "intellij_commit/findUsages/UIAutomationInteractionModel_with_library_cache",
  ],
  kotlinLang: [
    "kotlin_lang/findUsages/CommonParser_with_library_cache",
    "kotlin_lang/findUsages/DefaultArgumentStubGenerator_with_library_cache",
    "kotlin_lang/findUsages/FirErrors_with_library_cache",
    "kotlin_lang/findUsages/Flag_with_library_cache",
    "kotlin_lang/findUsages/ReferenceSymbolTable_with_library_cache",
  ],
  tbe: [
    "toolbox_enterprise/findUsages/ErrorReport_with_library_cache",
    "toolbox_enterprise/findUsages/genUuid_with_library_cache",
    "toolbox_enterprise/findUsages/getAll_with_library_cache",
    "toolbox_enterprise/findUsages/getTempDirectory_with_library_cache",
    "toolbox_enterprise/findUsages/PrincipalContext_with_library_cache",
    "toolbox_enterprise/findUsages/RequestMapping_with_library_cache",
    "toolbox_enterprise/findUsages/RestController_with_library_cache",
    "toolbox_enterprise/findUsages/ROLE_ADMIN_with_library_cache",
  ],
  androidCanaryLeak: [
    "leak-canary-android/findUsages/BOOLEAN_with_library_cache",
    "leak-canary-android/findUsages/HeapGraph_with_library_cache",
    "leak-canary-android/findUsages/HeapObject_with_library_cache",
    "leak-canary-android/findUsages/HprofRecordTag_with_library_cache",
    "leak-canary-android/findUsages/INT_with_library_cache",
    "leak-canary-android/findUsages/PrimitiveType_with_library_cache",
  ],
}

export const findUsagesCharts: ProjectsChartDefinition[] = [
  generateFindUsagesDefinition("'IntelliJ'", findUsagesProjects.intelliJ),
  generateFindUsagesDefinition("'Kotlin lang'", findUsagesProjects.kotlinLang),
  generateFindUsagesDefinition("'Toolbox Enterprise (TBE)'", findUsagesProjects.tbe),
  generateFindUsagesDefinition("'Android canary leak'", findUsagesProjects.androidCanaryLeak),
]

function generateFindUsagesDefinition(labelPrefix: string, projects: string[]): ProjectsChartDefinition {
  return {
    label: `${labelPrefix} findUsages mean value`,
    measure: "findUsages#mean_value",
    projects,
  }
}

export const evaluateExpressionProjects = {
  kotlinLanguageServer: [
    "kotlin_language_server/evaluate-expression/ClassPathTest_with_library_cache",
    "kotlin_language_server/evaluate-expression/Debouncer_with_library_cache",
    "kotlin_language_server/evaluate-expression/KotlinTextDocumentService_with_library_cache",
  ],
  petClinic: ["kotlin_petclinic/evaluate-expression/CacheConfig/sleep-1000_with_library_cache"],
  intellij: ["intellij_commit/evaluate-expression/DumbServiceStartupActivity_with_library_cache"],
}
export const completionInEvaluateExpressionProjects = {
  intellij: ["intellij_commit/completion/evaluate-expression_with_library_cache"],
  petClinic: ["kotlin_petclinic/completion/evaluate-expression/typing-it_with_library_cache", "kotlin_petclinic/completion/evaluate-expression/typing-system_with_library_cache"],
}
export const evaluateExpressionChars: ProjectsChartDefinition[] = [
  generateEvaluateExpressionDefinition("'Kotlin language server'", evaluateExpressionProjects.kotlinLanguageServer),
  generateEvaluateExpressionDefinition("'PetClinic'", evaluateExpressionProjects.petClinic),
  generateEvaluateExpressionDefinition("'Intellij'", evaluateExpressionProjects.intellij),
  ...generateCompletionDefinitions("'PetClinic completion in evaluate expression'", completionInEvaluateExpressionProjects.petClinic),
  ...generateCompletionDefinitions("'Intellij completion in evaluate expression'", completionInEvaluateExpressionProjects.intellij),
]
function generateEvaluateExpressionDefinition(labelPrefix: string, projects: string[]): ProjectsChartDefinition {
  return {
    label: `${labelPrefix} evaluate expression mean value`,
    measure: "evaluateExpression#mean_value",
    projects,
  }
}

export const KOTLIN_MAIN_METRICS = [
  "completion#mean_value",
  "completion#firstElementShown#mean_value",
  "localInspections#mean_value",
  "semanticHighlighting#mean_value",
  "findUsages#mean_value",
]
