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

  machines: string[]
}

const machines = {
  linux: "linux-blade-hetzner",
  mac: "Mac Mini M2 Pro (10 vCPU, 32 GB)",
}

const labelByProject: { [key: string]: string } = {
  kotlinEmpty: "'Kotlin empty'",
  intelliJ: "'IntelliJ'",
  intelliJ2: "'IntelliJ suite 2'",
  intelliJTyping2: "'IntelliJ with typing suite 2'",
  kotlinLang: "'Kotlin lang'",
  kotlinLanguageServer: "'Kotlin language server'",
  tbe: "'Toolbox Enterprise (TBE)'",
  ktorSamples: "'Ktor samples'",
  androidCanaryLeak: "'Android canary leak'",
  anki: "'Android anki project'",
  kotlinScript: "'Kotlin script'",
  removedImports: "'Files with removed imports'",
  springFramework: "'Spring framework'",
  rustPlugin: "'Rust plugin'",
  kotlinCoroutines: "'Kotlin Coroutines'",
  space: "'Space'",
  sqliter: "'SQLiter'",
  ktor: "'Ktor'",
  kotlinCoroutinesQg: "'Kotlin Coroutines QG'",
  tbeCaseWithAssert: "'Toolbox Enterprise (TBE) different length'",
  mppNativeAcceptance: "'Native-acceptance'",
}

export const completionLinuxProjects = {
  kotlinEmpty: ["kotlin_empty/completion/Main_empty_place_with_library_cache", "kotlin_empty/completion/Main_empty_place_typing_with_library_cache"],
  intelliJ: [
    "intellij_commit/completion/DexInlineCallStackComparisonTest_empty_place_with_library_cache",
    "intellij_commit/completion/DexInlineCallStackComparisonTest_after_parameter_with_library_cache",
    "intellij_commit/completion/DexInlineCallStackComparisonTest_empty_place_typing_with_library_cache",
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
    "kotlin_lang/completion/CommonParser_after_parameter_with_library_cache",
    "kotlin_lang/completion/CommonParser_empty_place_with_library_cache",
    "kotlin_lang/completion/CommonParser_empty_place_typing_with_library_cache",
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
    "arrow/completion/build.gradle_completion_kts_with_library_cache",
    "kotlin_empty_kts/completion/build.gradle_completion_kts_with_library_cache",
    "kotlin_lang/completion/build.gradle_completion_kts_typing_with_library_cache",
    "kotlin_lang/completion/build.gradle_completion_kts_with_library_cache",
    "kotlin_lang/completion/build.gradle_top_level_typing_with_library_cache",
  ],
  kotlinCoroutines: [
    "kotlin_coroutines_commit/completion/CoroutineScheduler_in_constructor_typing_with_library_cache",
    "kotlin_coroutines_commit/completion/CoroutineScheduler_in_function_typing_with_library_cache",
    "kotlin_coroutines_commit/completion/CoroutineScheduler_in_function_with_library_cache",
  ],
  mppNativeAcceptance: [
    "kotlin_kmp_native_acceptance/completion/Sample.jvm_with_library_cache",
    "kotlin_kmp_native_acceptance/completion/Sample.linux_with_library_cache",
    "kotlin_kmp_native_acceptance/completion/Sample.macOS_with_library_cache",
    "kotlin_kmp_native_acceptance/completion/Sample.mingw_with_library_cache",
  ],
  space: ["space_specific/completion/Dimensions_typealias_with_library_cache"],
}

export const completionMacProjects = {
  kotlinCoroutinesQG: [
    "kotlin_coroutines_qg/completion/AwaitStressTest_test_AwaitStressTest_typing_with_library_cache",
    "kotlin_coroutines_qg/completion/CoroutineContext_enqueue_CoroutineContext_typing_with_library_cache",
    "kotlin_coroutines_qg/completion/CoroutineContext_isJsdom_typing_with_library_cache",
    "kotlin_coroutines_qg/completion/Dispatchers_dispatch_async_Dispatchers_typing_with_library_cache",
    "kotlin_coroutines_qg/completion/Future_JobNode_Future_typing_with_library_cache",
    "kotlin_coroutines_qg/completion/Launcher_testLauncherEntryPoint_Launcher_typing_with_library_cache",
    "kotlin_coroutines_qg/completion/SchedulerTask_afterTask_SchedulerTask_typing_with_library_cache",
    "kotlin_coroutines_qg/completion/TestBase_finish_TestBase_typing_with_library_cache",
    "kotlin_coroutines_qg/completion/WorkerTest_test_WorkerTest_typing_with_library_cache",
  ],
  sqliter: [
    "SQLiter/completion/SqliteDatabase_SqliteDatabase_typing_with_library_cache",
    "SQLiter/completion/TracingSqliteStatement_TracingSqliteStatement_typing_with_library_cache",
    "SQLiter/completion/Types_Types_typing_with_library_cache",
  ],
  ktor: [
    "ktor_before_add_wasm_client/completion/CharsetLinux_encodeImpl_typing_with_library_cache",
    "ktor_before_add_wasm_client/completion/CharsetMingw_findCharset_with_library_cache",
    "ktor_before_add_wasm_client/completion/URLBuilderJs_URLBuilderJs_typing_with_library_cache",
    "ktor_before_add_wasm_client/completion/DarwinUtils_DarwinUtils_with_library_cache",
  ],
}

export const completionProjects = { ...completionMacProjects, ...completionLinuxProjects }

export const completionCharts = Object.entries(completionLinuxProjects)
  .map(([key, value]) => generateCompletionDefinitions(labelByProject[key], value, [machines.linux]))
  .concat(Object.entries(completionMacProjects).map(([key, value]) => generateCompletionDefinitions(labelByProject[key], value, [machines.mac])))

function generateCompletionDefinitions(labelPrefix: string, projects: string[], machines: string[]): ProjectsChartDefinition[] {
  return [
    {
      label: `${labelPrefix} completion mean value`,
      measure: "completion#mean_value",
      projects,
      machines,
    },
    {
      label: `${labelPrefix} firstElementShown mean value`,
      measure: "completion#firstElementShown#mean_value",
      projects,
      machines,
    },
  ]
}

/**
 * Highlighting projects are also the projects for local inspections.
 */
export const highlightingLinuxProjects = {
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
  anki: [
    "anki-android/highlight/AbstractSchedTest_with_library_cache",
    "anki-android/highlight/ACRATest_with_library_cache",
    "anki-android/highlight/Finder_with_library_cache",
    "anki-android/highlight/FlashCardsContract_with_library_cache",
    "anki-android/highlight/MetaDB_with_library_cache",
    "anki-android/highlight/TagsUtilTest_with_library_cache",
    "anki-android/highlight/TokenizerTest_with_library_cache",
    "anki-android/highlight/UniqueArrayListTest_with_library_cache",
  ],
  kotlinScript: [
    "arrow/highlight/build.gradle_with_library_cache",
    "kotlin_empty_kts/highlight/build.gradle_with_library_cache",
    "kotlin_lang/highlight/atomicfu_atomicfu-compiler/build.gradle_with_library_cache",
    "kotlin_lang/highlight/build.gradle_with_library_cache",
    "kotlin_lang/highlight/js_js.tests/build.gradle_with_library_cache",
    "kotlin_lang/highlight/kotlin-gradle-plugin/build.gradle_with_library_cache",
    "kotlin_lang/highlight/kotlin-gradle-plugin-integration-tests/build.gradle_with_library_cache",
    "kotlin_lang/highlight/libraries_kotlin.test/build.gradle_with_library_cache",
    "kotlin_lang/highlight/libraries_stdlib/build.gradle_with_library_cache",
    "kotlin_lang/highlight/prepare.compiler/build.gradle_with_library_cache",
    "kotlin_lang/highlight/wasm_wasm.tests/build.gradle_with_library_cache",
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
  kotlinCoroutines: [
    "kotlin_coroutines_commit/highlight/BufferedChannel_with_library_cache",
    "kotlin_coroutines_commit/highlight/CoroutineScheduler_with_library_cache",
    "kotlin_coroutines_commit/highlight/JobSupport_with_library_cache",
  ],
  space: [
    "space_specific/highlight/DocumentsStatsTest_with_library_cache",
    "space_specific/highlight/EventCounters_with_library_cache",
    "space_specific/highlight/HttpApiClientTest_with_library_cache",
    "space_specific/highlight/M2ChannelMessageListVmV2_with_library_cache",
    "space_specific/highlight/M2ChannelVm_with_library_cache",
    "space_specific/highlight/MagicBarComponent_with_library_cache",
    "space_specific/highlight/ProfilesImpl_with_library_cache",
    "space_specific/highlight/SchemaMigration_with_library_cache",
    "space_specific/highlight/SecretsTests_with_library_cache",
    "space_specific/highlight/XScrollable_with_library_cache",
  ],
}

export const highlightingMacProjects = {
  sqliter: [
    "SQLiter/highlight/ActualSqliteStatement_with_library_cache",
    "SQLiter/highlight/Logger_with_library_cache",
    "SQLiter/highlight/SQLiteException_with_library_cache",
    "SQLiter/highlight/SqliteDatabase_with_library_cache",
    "SQLiter/highlight/SqliteStatement_with_library_cache",
    "SQLiter/highlight/TracingSqliteStatement_with_library_cache",
    "SQLiter/highlight/Types_with_library_cache",
  ],
  kotlinCoroutinesQg: [
    "kotlin_coroutines_qg/highlight/BufferedChannel_with_library_cache",
    "kotlin_coroutines_qg/highlight/Builders_with_library_cache",
    "kotlin_coroutines_qg/highlight/CopyOnWriteList_with_library_cache",
    "kotlin_coroutines_qg/highlight/CoroutineExceptionHandlerImpl_with_library_cache",
    "kotlin_coroutines_qg/highlight/CoroutineScheduler_with_library_cache",
    "kotlin_coroutines_qg/highlight/Dispatchers_with_library_cache",
    "kotlin_coroutines_qg/highlight/JSDispatcher_with_library_cache",
    "kotlin_coroutines_qg/highlight/Launcher_with_library_cache",
    "kotlin_coroutines_qg/highlight/MainDispatcherTest_with_library_cache",
    "kotlin_coroutines_qg/highlight/MultithreadedDispatchers_with_library_cache",
    "kotlin_coroutines_qg/highlight/Select_with_library_cache",
    "kotlin_coroutines_qg/highlight/SharedFlow_with_library_cache",
    "kotlin_coroutines_qg/highlight/TestBase_with_library_cache",
    "kotlin_coroutines_qg/highlight/WorkerMain_with_library_cache",
  ],
  ktor: [
    "ktor_before_add_wasm_client/highlight/ByteBufferChannel_with_library_cache",
    "ktor_before_add_wasm_client/highlight/StaticContentTest_with_library_cache",
    "ktor_before_add_wasm_client/highlight/RoutingResolveTest_with_library_cache",
    "ktor_before_add_wasm_client/highlight/PrimitiveArraysNative_with_library_cache",
    "ktor_before_add_wasm_client/highlight/CertificatePinner_with_library_cache",
    "ktor_before_add_wasm_client/highlight/CharsetLinux_with_library_cache",
    "ktor_before_add_wasm_client/highlight/TestUtilsMacos_with_library_cache",
    "ktor_before_add_wasm_client/highlight/PrimitiveArraysJs_with_library_cache",
    "ktor_before_add_wasm_client/highlight/CharsetMingw_with_library_cache",
  ],
}

export const highlightingProjects = { ...highlightingLinuxProjects, ...highlightingMacProjects }

export const highlightingCharts = Object.entries(highlightingLinuxProjects)
  .map(([key, value]) => generateHighlightingDefinition(labelByProject[key], value, [machines.linux]))
  .concat(Object.entries(highlightingMacProjects).map(([key, value]) => generateHighlightingDefinition(labelByProject[key], value, [machines.mac])))

function generateHighlightingDefinition(labelPrefix: string, projects: string[], machines: string[]): ProjectsChartDefinition {
  return {
    label: `${labelPrefix} semantic highlighting mean value`,
    measure: "semanticHighlighting#mean_value",
    projects,
    machines,
  }
}

export const codeAnalysisCharts = Object.entries(highlightingLinuxProjects)
  .map(([key, value]) => generateCodeAnalysisChartsDefinition(labelByProject[key], value, [machines.linux]))
  .concat(Object.entries(highlightingMacProjects).map(([key, value]) => generateCodeAnalysisChartsDefinition(labelByProject[key], value, [machines.mac])))

function generateCodeAnalysisChartsDefinition(labelPrefix: string, projects: string[], machines: string[]): ProjectsChartDefinition {
  return {
    label: `${labelPrefix} Code Analysis mean value`,
    measure: "localInspections#mean_value",
    projects,
    machines,
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
    machines: [machines.linux],
  }
}

export const findUsagesLinuxProjects = {
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
  anki: [
    "anki-android/findUsages/Card_with_library_cache",
    "anki-android/findUsages/CompatHelper_with_library_cache",
    "anki-android/findUsages/Decks_with_library_cache",
    "anki-android/findUsages/load_with_library_cache",
  ],
  kotlinLangScript: [
    "kotlin_lang/findUsages/intellijVersionForIde_with_library_cache",
    "kotlin_lang/findUsages/JvmTestFramework_with_library_cache",
    "kotlin_lang/findUsages/SourceSet_with_library_cache",
  ],
  kotlinCoroutines: [
    "kotlin_coroutines_commit/findUsages/assert_with_library_cache",
    "kotlin_coroutines_commit/findUsages/emit_with_library_cache",
    "kotlin_coroutines_commit/findUsages/Flow_with_library_cache",
    "kotlin_coroutines_commit/findUsages/FlowCollector_with_library_cache",
    "kotlin_coroutines_commit/findUsages/runBlocking_with_library_cache",
  ],
  space: [
    "space_specific/findUsages/ApiFlag_with_library_cache",
    "space_specific/findUsages/Http_with_library_cache",
    "space_specific/findUsages/IntSizePx_with_library_cache",
    "space_specific/findUsages/UniqueConstraint_with_library_cache",
  ],
}

export const findUsagesMacProjects = {
  ktor: [
    "ktor_before_add_wasm_client/findUsages/SelectorManager_with_library_cache",
    "ktor_before_add_wasm_client/findUsages/Encoding_with_library_cache",
    "ktor_before_add_wasm_client/findUsages/CharsetNative_with_library_cache",
    "ktor_before_add_wasm_client/findUsages/Memory_with_library_cache",
  ],
  sqliter: ["SQLiter/findUsages/sqlite3_prepare16_v2_with_library_cache_with_library_cache", "SQLiter/findUsages/SQLITE_OK_with_library_cache"],
  kotlinCoroutinesQg: [
    "kotlin_coroutines_qg/findUsages/createMainDispatcher_with_library_cache",
    "kotlin_coroutines_qg/findUsages/hexAddress_with_library_cache",
    "kotlin_coroutines_qg/findUsages/Runnable_with_library_cache",
    "kotlin_coroutines_qg/findUsages/SynchronizedObject_with_library_cache",
  ],
}

export const findUsagesProjects = { ...findUsagesMacProjects, ...findUsagesLinuxProjects }

export const findUsagesCharts: ProjectsChartDefinition[] = Object.entries(findUsagesLinuxProjects)
  .map(([key, value]) => generateFindUsagesDefinition(labelByProject[key], value, [machines.linux]))
  .concat(Object.entries(findUsagesMacProjects).map(([key, value]) => generateFindUsagesDefinition(labelByProject[key], value, [machines.mac])))

function generateFindUsagesDefinition(labelPrefix: string, projects: string[], machines: string[]): ProjectsChartDefinition {
  return {
    label: `${labelPrefix} findUsages mean value`,
    measure: "findUsages#mean_value",
    projects,
    machines,
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
  ...generateCompletionDefinitions("'PetClinic completion in evaluate expression'", completionInEvaluateExpressionProjects.petClinic, [machines.linux]),
  ...generateCompletionDefinitions("'Intellij completion in evaluate expression'", completionInEvaluateExpressionProjects.intellij, [machines.linux]),
]

function generateEvaluateExpressionDefinition(labelPrefix: string, projects: string[]): ProjectsChartDefinition {
  return {
    label: `${labelPrefix} evaluate expression mean value`,
    measure: "evaluateExpression#mean_value",
    projects,
    machines: [machines.linux],
  }
}

const convertJavaToKotlinProjects = {
  intellijSources: [
    "intellij_sources/javaToKotlin/HighlightFixUtil",
    "intellij_sources/javaToKotlin/CompilerConfigurationImpl",
    "intellij_sources/javaToKotlin/DaemonCodeAnalyzerImpl",
    "intellij_sources/javaToKotlin/HighlightMethodUtil",
    "intellij_sources/javaToKotlin/EditorImpl",
  ],
}

export const evaluateConvertJavaToKotlinProjectsChars: ProjectsChartDefinition[] = [
  generateConvertJavaToKotlinDefinition("'Intellij Sources'", convertJavaToKotlinProjects.intellijSources, [machines.linux]),
]

function generateConvertJavaToKotlinDefinition(labelPrefix: string, projects: string[], machines: string[]): ProjectsChartDefinition {
  return {
    label: `${labelPrefix} convert java to kotlin`,
    measure: "convertJavaToKotlin",
    projects,
    machines,
  }
}

const scriptHighlight = [{ label: "'Kotlin script'", projects: highlightingLinuxProjects.kotlinScript }]
export const highlightingScriptCharts: ProjectsChartDefinition[] = scriptHighlight.map((v) => generateHighlightingDefinition(v.label, v.projects, [machines.linux]))
export const codeAnalysisScriptCharts: ProjectsChartDefinition[] = scriptHighlight.map((v) => generateCodeAnalysisChartsDefinition(v.label, v.projects, [machines.linux]))
export const scriptCompletionCharts: ProjectsChartDefinition[] = [...generateCompletionDefinitions("'Kotlin script'", completionLinuxProjects.kotlinScript, [machines.linux])]
export const scriptFindUsagesCharts: ProjectsChartDefinition[] = [generateFindUsagesDefinition("'Kotlin lang (kts)'", findUsagesLinuxProjects.kotlinLangScript, [machines.linux])]
export const KOTLIN_MAIN_METRICS = [
  "completion#mean_value",
  "completion#firstElementShown#mean_value",
  "localInspections#mean_value",
  "semanticHighlighting#mean_value",
  "findUsages#mean_value",
]
