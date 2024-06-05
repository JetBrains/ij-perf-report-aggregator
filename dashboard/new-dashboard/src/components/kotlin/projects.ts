import { ChartDefinition } from "../charts/DashboardCharts"
import { SimpleMeasureConfigurator } from "../../configurators/SimpleMeasureConfigurator"
import { computed, ComputedRef } from "vue"

export const KOTLIN_MAIN_METRICS = [
  "completion#mean_value",
  "completion#firstElementShown#mean_value",
  "localInspections#mean_value",
  "semanticHighlighting#mean_value",
  "findUsages#mean_value",
]

const MEASURES = {
  completionMeasures: [
    { name: "completion#mean_value", label: "completion mean value" },
    { name: "completion#firstElementShown#mean_value", label: "firstElementShown mean value" },
  ],
  highlightingMeasures: [{ name: "semanticHighlighting#mean_value", label: "semantic highlighting mean value" }],
  codeAnalysisMeasures: [{ name: "localInspections#mean_value", label: "Code Analysis mean value" }],
  refactoringMeasures: [
    { name: "performInlineRename#mean_value", label: "PerformInlineRename" },
    { name: "startInlineRename#mean_value", label: "StartInlineRename" },
    { name: "prepareForRename#mean_value", label: "PrepareForRename" },
  ],
  optimizeImportsMeasures: [{ name: "execute_editor_optimizeimports", label: "Optimize imports" }],
  insertCodeMeasures: [{ name: "execute_editor_paste", label: "Insert code" }],
  findUsagesMeasures: [{ name: "findUsages#mean_value", label: "findUsages mean value" }],
  evaluateExpressionMeasures: [{ name: "evaluateExpression#mean_value", label: "evaluate expression mean value" }],
  convertJavaToKotlinProjectsMeasures: [{ name: "convertJavaToKotlin", label: "convert java to kotlin" }],
}

export const MACHINES = {
  linux: "linux-blade-hetzner",
  mac: "Mac Mini M2 Pro (10 vCPU, 32 GB)",
}

export const PROJECT_CATEGORIES: Record<string, ProjectCategory> = {
  kotlinEmpty: buildCategory("Kotlin empty", "kotlin_empty/"),
  intelliJ: buildCategory("IntelliJ", "intellij_commit/"),
  // Same intelliJ. Need to avoid lot of lines on chart
  intelliJ2: buildCategory("IntelliJ suite 2", ""),

  intelliJSources: buildCategory("Intellij Sources", "intellij_sources/"),
  intelliJTyping: buildCategory("IntelliJ with typing", ""),
  kotlinLang: buildCategory("Kotlin lang", "kotlin_lang/"),
  // Same kotlinLang. Need to avoid lot of lines on chart
  kotlinScript: buildCategory("Kotlin script", ""),

  kotlinLanguageServer: buildCategory("Kotlin language server", "kotlin_language_server/"),
  tbe: buildCategory("Toolbox Enterprise (TBE)", "toolbox_enterprise/"),
  tbeCaseWithAssert: buildCategory("Toolbox Enterprise (TBE) different length", ""),

  ktorSamples: buildCategory("Ktor samples", "ktor_samples"),
  androidCanaryLeak: buildCategory("Android canary leak", "leak-canary-android/"),
  anki: buildCategory("Android anki project", "anki-android/"),
  removedImports: buildCategory("Files with removed imports", ""),
  springFramework: buildCategory("Spring framework", "spring-framework/"),
  rustPlugin: buildCategory("Rust plugin", "rust_commit/"),
  kotlinCoroutines: buildCategory("Kotlin Coroutines", "kotlin_coroutines_commit/"),
  sqliter: buildCategory("SQLiter", "SQLiter/"),
  ktor: buildCategory("Ktor", "ktor_before_add_wasm_client/"),
  kotlinCoroutinesQG: buildCategory("Kotlin Coroutines QG", "kotlin_coroutines_qg/"),

  mppNativeAcceptance: buildCategory("Native-acceptance", "kotlin_kmp_native_acceptance/"),
  petClinic: buildCategory("Pet Clinic", "kotlin_petclinic/"),
  arrow: buildCategory("Arrow", "arrow/"),
  kotlinEmptyScript: buildCategory("Empty Script (.kts)", "kotlin_empty_kts/"),
}

export const KOTLIN_PROJECTS = {
  linux: {
    completion: {
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
      intelliJTyping: [
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
      arrow: ["arrow/completion/build.gradle_completion_kts_with_library_cache"],
      kotlinEmptyScript: ["kotlin_empty_kts/completion/build.gradle_completion_kts_with_library_cache"],
      kotlinScript: [
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
    },
    highlighting: {
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
      arrow: ["arrow/highlight/build.gradle_with_library_cache"],
      kotlinEmptyScript: ["kotlin_empty_kts/highlight/build.gradle_with_library_cache"],
      kotlinScript: [
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
    },
    refactoringRename: {
      intelliJ: ["intellij_commit/rename/SqlBlock_SqlBlockRenamed"],
      kotlinLanguageServer: ["kotlin_language_server/rename/Compiler_createKtFileTest"],
    },
    refactoringInsertCode: {
      kotlinLanguageServer: ["kotlin_language_server/insertCode/Rename_renameSymbol", "kotlin_language_server/insertCode/SpecialJavaFileForTest_j2k"],
    },
    optimizeImports: {
      intelliJ: [
        "intellij_commit/otimizeImports/AbstractKotlinMavenImporterTest",
        "intellij_commit/otimizeImports/FSD",
        "intellij_commit/otimizeImports/DiagramsModel.Generated",
        "intellij_commit/otimizeImports/OraIntrospector",
        "intellij_commit/otimizeImports/QuickFixRegistrar",
        "intellij_commit/otimizeImports/SwiftTypeAssignabilityTest",
        "intellij_commit/otimizeImports/TerraformConfigCompletionContributor",
      ],
    },
    findUsages: {
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
      kotlinScript: [
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
    },
    evaluateExpression: {
      kotlinLanguageServer: [
        "kotlin_language_server/evaluate-expression/ClassPathTest_with_library_cache",
        "kotlin_language_server/evaluate-expression/Debouncer_with_library_cache",
        "kotlin_language_server/evaluate-expression/KotlinTextDocumentService_with_library_cache",
      ],
      petClinic: ["kotlin_petclinic/evaluate-expression/CacheConfig/sleep-1000_with_library_cache"],
      intelliJ: ["intellij_commit/evaluate-expression/DumbServiceStartupActivity_with_library_cache"],
    },
    completionInEvaluateExpression: {
      intelliJ: ["intellij_commit/completion/evaluate-expression_with_library_cache"],
      petClinic: [
        "kotlin_petclinic/completion/evaluate-expression/typing-it_with_library_cache",
        "kotlin_petclinic/completion/evaluate-expression/typing-system_with_library_cache",
      ],
    },
    convertJavaToKotlin: {
      intelliJSources: [
        "intellij_sources/javaToKotlin/HighlightFixUtil",
        "intellij_sources/javaToKotlin/CompilerConfigurationImpl",
        "intellij_sources/javaToKotlin/DaemonCodeAnalyzerImpl",
        "intellij_sources/javaToKotlin/HighlightMethodUtil",
        "intellij_sources/javaToKotlin/EditorImpl",
      ],
    },
  },
  mac: {
    completion: {
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
    },
    highlighting: {
      sqliter: [
        "SQLiter/highlight/ActualSqliteStatement_with_library_cache",
        "SQLiter/highlight/Logger_with_library_cache",
        "SQLiter/highlight/SQLiteException_with_library_cache",
        "SQLiter/highlight/SqliteDatabase_with_library_cache",
        "SQLiter/highlight/SqliteStatement_with_library_cache",
        "SQLiter/highlight/TracingSqliteStatement_with_library_cache",
        "SQLiter/highlight/Types_with_library_cache",
      ],
      kotlinCoroutinesQG: [
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
    },
    findUsages: {
      ktor: [
        "ktor_before_add_wasm_client/findUsages/SelectorManager_with_library_cache",
        "ktor_before_add_wasm_client/findUsages/Encoding_with_library_cache",
        "ktor_before_add_wasm_client/findUsages/CharsetNative_with_library_cache",
        "ktor_before_add_wasm_client/findUsages/Memory_with_library_cache",
      ],
      sqliter: ["SQLiter/findUsages/sqlite3_column_type_with_library_cache", "SQLiter/findUsages/SQLITE_OK_with_library_cache"],
      kotlinCoroutinesQG: [
        "kotlin_coroutines_qg/findUsages/createMainDispatcher_with_library_cache",
        "kotlin_coroutines_qg/findUsages/hexAddress_with_library_cache",
        "kotlin_coroutines_qg/findUsages/Runnable_with_library_cache",
        "kotlin_coroutines_qg/findUsages/SynchronizedObject_with_library_cache",
      ],
    },
  },
}

export const KOTLIN_PROJECT_CONFIGURATOR = new SimpleMeasureConfigurator("project", null)
KOTLIN_PROJECT_CONFIGURATOR.initData(Object.values(PROJECT_CATEGORIES).flatMap((c) => c.label))

function buildCategory(label: string, prefix: string): ProjectCategory {
  return { label, prefix }
}

function projectsToDefinition(projectsByOS: ProjectsByOS[]): ComputedRef<ChartDefinition[]> {
  return computed(() =>
    projectsByOS
      .flatMap(({ projects, measures, machines }) =>
        measures.flatMap(({ name, label }) =>
          Object.entries(projects).flatMap(([key, value]) => {
            return {
              labels: [`'${PROJECT_CATEGORIES[key].label}' ${label}`],
              measures: [name],
              projects: value,
              machines,
            }
          })
        )
      )
      .filter((chart) => KOTLIN_PROJECT_CONFIGURATOR.selected.value?.some((selected) => chart.labels[0].startsWith(`'${selected}'`)))
  )
}

export const completionProjects = { ...KOTLIN_PROJECTS.linux.completion, ...KOTLIN_PROJECTS.mac.completion }

export const completionCharts = projectsToDefinition([
  {
    projects: KOTLIN_PROJECTS.linux.completion,
    measures: MEASURES.completionMeasures,
    machines: [MACHINES.linux],
  },
  {
    projects: KOTLIN_PROJECTS.mac.completion,
    measures: MEASURES.completionMeasures,
    machines: [MACHINES.mac],
  },
])

/**
 * Highlighting projects are also the projects for local inspections.
 */

export const highlightingProjects = { ...KOTLIN_PROJECTS.linux.highlighting, ...KOTLIN_PROJECTS.mac.highlighting }

export const highlightingCharts = projectsToDefinition([
  {
    projects: KOTLIN_PROJECTS.linux.highlighting,
    measures: MEASURES.highlightingMeasures,
    machines: [MACHINES.linux],
  },
  {
    projects: KOTLIN_PROJECTS.mac.highlighting,
    measures: MEASURES.highlightingMeasures,
    machines: [MACHINES.mac],
  },
])

export const codeAnalysisCharts = projectsToDefinition([
  {
    projects: KOTLIN_PROJECTS.linux.highlighting,
    measures: MEASURES.codeAnalysisMeasures,
    machines: [MACHINES.linux],
  },
  {
    projects: KOTLIN_PROJECTS.mac.highlighting,
    measures: MEASURES.codeAnalysisMeasures,
    machines: [MACHINES.mac],
  },
])

export const refactoringCharts = projectsToDefinition([
  {
    projects: KOTLIN_PROJECTS.linux.refactoringRename,
    measures: MEASURES.refactoringMeasures,
    machines: [MACHINES.linux],
  },
  {
    projects: KOTLIN_PROJECTS.linux.refactoringInsertCode,
    measures: MEASURES.insertCodeMeasures,
    machines: [MACHINES.linux],
  },
  {
    projects: KOTLIN_PROJECTS.linux.optimizeImports,
    measures: MEASURES.optimizeImportsMeasures,
    machines: [MACHINES.linux],
  },
])

export const findUsagesProjects = { ...KOTLIN_PROJECTS.linux.findUsages, ...KOTLIN_PROJECTS.mac.findUsages }

export const findUsagesCharts = projectsToDefinition([
  {
    projects: KOTLIN_PROJECTS.linux.findUsages,
    measures: MEASURES.findUsagesMeasures,
    machines: [MACHINES.linux],
  },
  {
    projects: KOTLIN_PROJECTS.mac.findUsages,
    measures: MEASURES.findUsagesMeasures,
    machines: [MACHINES.mac],
  },
])

export const evaluateExpressionChars = projectsToDefinition([
  {
    projects: KOTLIN_PROJECTS.linux.evaluateExpression,
    measures: MEASURES.evaluateExpressionMeasures,
    machines: [MACHINES.linux],
  },
  {
    projects: KOTLIN_PROJECTS.linux.completionInEvaluateExpression,
    measures: MEASURES.completionMeasures,
    machines: [MACHINES.linux],
  },
])

export const convertJavaToKotlinProjectsChars = projectsToDefinition([
  {
    projects: KOTLIN_PROJECTS.linux.convertJavaToKotlin,
    measures: MEASURES.convertJavaToKotlinProjectsMeasures,
    machines: [MACHINES.linux],
  },
])

const scriptHighlight = { kotlinScript: KOTLIN_PROJECTS.linux.highlighting.kotlinScript }
export const highlightingScriptCharts = projectsToDefinition([
  {
    projects: scriptHighlight,
    measures: MEASURES.highlightingMeasures,
    machines: [MACHINES.linux],
  },
])
export const codeAnalysisScriptCharts = projectsToDefinition([
  {
    projects: scriptHighlight,
    measures: MEASURES.codeAnalysisMeasures,
    machines: [MACHINES.linux],
  },
])

export const scriptCompletionCharts = projectsToDefinition([
  {
    projects: { kotlinScript: KOTLIN_PROJECTS.linux.completion.kotlinScript },
    measures: MEASURES.completionMeasures,
    machines: [MACHINES.linux],
  },
])

export const scriptFindUsagesCharts = projectsToDefinition([
  {
    projects: { kotlinScript: KOTLIN_PROJECTS.linux.findUsages.kotlinScript },
    measures: MEASURES.findUsagesMeasures,
    machines: [MACHINES.linux],
  },
])

type Projects = Record<string, string[]>

interface ProjectsByOS {
  projects: Projects
  measures: Measure[]
  machines: string[]
}

/**
 * A project category is a project name prefix such as "intellij_commit/" and "kotlin_lang/" with an associated, human-readable label.
 */
interface ProjectCategory {
  label: string
  prefix: string
}

interface Measure {
  name: string
  label: string
}
