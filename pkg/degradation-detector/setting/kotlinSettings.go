package setting

import detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"

func GenerateKotlinSettings() []detector.PerformanceSettings {
	testNames := []string{"intellij_commit/completion/DexInlineCallStackComparisonTest_empty_place_with_library_cache",
		"intellij_commit/completion/DexInlineCallStackComparisonTest_after_parameter_with_library_cache",
		"intellij_commit/completion/DexInlineCallStackComparisonTest_empty_place_typing_with_library_cache",
		"intellij_commit/completion/KotlinHighLevelFunctionParameterInfoHandler_emptyPlace_updateUIOrFail_typing_with_library_cache",
		"intellij_commit/completion/KotlinHighLevelFunctionParameterInfoHandler_emptyPlace_updateUIOrFail_with_library_cache",
		"intellij_commit/completion/KtOCSwiftChangeSignatureTest_emptyPlace_changeReturnType_typing_with_library_cache",
		"intellij_commit/completion/KtOCSwiftChangeSignatureTest_emptyPlace_changeReturnType_with_library_cache",
		"intellij_commit/completion/IdeMenuBar_emptyPlace_sout_typing_with_library_cache",
		"intellij_commit/completion/TestModelParser_emptyPlace_if_typing_with_library_cache",
		"intellij_commit/completion/AndroidModuleSystem_emptyPlace_get_typing_with_library_cache",
		"kotlin_lang/completion/CommonParser_after_parameter_with_library_cache",
		"kotlin_lang/completion/CommonParser_empty_place_with_library_cache",
		"kotlin_lang/completion/CommonParser_empty_place_typing_with_library_cache",
		"kotlin_language_server/completion/Completions_emptyPlace_completions_typing_with_library_cache",
		"kotlin_language_server/completion/QuickFixesTest_emptyPlace_completions_typing_with_library_cache",
		"kotlin_language_server/completion/Completions_emptyPlace_completions_with_library_cache",
		"kotlin_language_server/completion/QuickFixesTest_emptyPlace_completions_with_library_cache",
		"toolbox_enterprise/completion/ProfileController_constructor_typing_with_library_cache",
		"toolbox_enterprise/completion/ProfileController_dataclass_typing_with_library_cache",
		"toolbox_enterprise/completion/ProfileController_fun_typing_with_library_cache",
		"toolbox_enterprise/completion/ProfileServiceTest_body_with_library_cache",
		"toolbox_enterprise/completion/ProfileServiceTest_constructor_typing_with_library_cache",
		"toolbox_enterprise/completion/ProfileServiceTest_constructor_with_library_cache",
		"toolbox_enterprise/completion/ProfileServiceTest_emptyPlace_FileEnd_typing_with_library_cache",
		"toolbox_enterprise/completion/ProfileServiceTest_emptyPlace_FileEnd_with_library_cache",
		"toolbox_enterprise/completion/IntelliJPluginSettingTest_assert_0_with_library_cache",
		"toolbox_enterprise/completion/IntelliJPluginSettingTest_assert_2_typing_with_library_cache",
		"toolbox_enterprise/completion/IntelliJPluginSettingTest_assert_4_typing_with_library_cache",
		"toolbox_enterprise/completion/IntelliJPluginSettingTest_assert_7_typing_with_library_cache",
		"toolbox_enterprise/completion/IntelliJPluginSettingTest_assert_10_typing_with_library_cache",
		"arrow/completion/build.gradle_completion_kts_with_library_cache",
		"kotlin_empty_kts/completion/build.gradle_completion_kts_with_library_cache",
		"kotlin_lang/completion/build.gradle_completion_kts_with_library_cache",
		"kotlin_language_server/completion/Completions_emptyPlace_completions_typing_with_library_cache",
		"kotlin_language_server/completion/Completions_emptyPlace_completions_with_library_cache",
		"kotlin_language_server/completion/QuickFixesTest_emptyPlace_completions_typing_with_library_cache",
		"kotlin_language_server/completion/QuickFixesTest_emptyPlace_completions_with_library_cache",
		"kotlin_coroutines_commit/completion/CoroutineScheduler_in_constructor_typing_with_library_cache",
		"kotlin_coroutines_commit/completion/CoroutineScheduler_in_function_typing_with_library_cache",
		"kotlin_coroutines_commit/completion/CoroutineScheduler_in_function_with_library_cache",
		"kotlin_kmp_native_acceptance/completion/Sample.jvm_with_library_cache",
		"kotlin_kmp_native_acceptance/completion/Sample.linux_with_library_cache",
		"kotlin_kmp_native_acceptance/completion/Sample.macOS_with_library_cache",
		"kotlin_kmp_native_acceptance/completion/Sample.mingw_with_library_cache",
		"ktor_commit/completion/ContentNegotiationTest_fun_typing_with_library_cache",
		"ktor_commit/completion/DarwinClientEngineConfig_fun_with_library_cache",
		"ktor_commit/completion/HighLoadHttpGenerator_end_constructor_typing_with_library_cache",
		"ktor_commit/completion/HighLoadHttpGenerator_mid_constructor_typing_with_library_cache",
		"ktor_commit/completion/LockFreeLinkedList_getter_typing_with_library_cache",
		"ktor_commit/completion/LockFreeLinkedList_typealias_typing_with_library_cache",
		"ktor_commit/completion/RequestResponseBuilderJs_fun_typing_with_library_cache",
		"ktor_commit/completion/RequestResponseBuilderNative_fun_typing_with_library_cache",
		"space_specific/completion/Dimensions_typealias_with_library_cache", "intellij_commit/highlight/KtOCSwiftChangeSignatureTest_with_library_cache",
		"intellij_commit/highlight/KotlinHighLevelFunctionParameterInfoHandler_with_library_cache",
		"intellij_commit/highlight/JdkList_with_library_cache",
		"intellij_commit/highlight/ComposeCompletionContributorTest_with_library_cache",
		"intellij_commit/highlight/AgpUpgradeRefactoringProcessor_with_library_cache",
		"intellij_commit/highlight/AndroidModelTest_with_library_cache",
		"intellij_commit/highlight/SecureWireOverStreamTransport_with_library_cache",
		"intellij_commit/highlight/DexInlineCallStackComparisonTest_with_library_cache",
		"intellij_commit/highlight/DexLocalVariableTableBreakpointTest_with_library_cache",
		"intellij_commit/highlight/OraIntrospector_with_library_cache",
		"intellij_commit/highlight/SolutionModel.Generated_with_library_cache",
		"intellij_commit/highlight/UIAutomationInteractionModel.Generated_with_library_cache",
		"kotlin_lang/highlight/LazyJVM_with_library_cache",
		"kotlin_lang/highlight/CommonParser_with_library_cache",
		"kotlin_lang/highlight/FirErrors_with_library_cache",
		"kotlin_lang/highlight/Flag_with_library_cache",
		"kotlin_lang/highlight/KtFirDataClassConverters_with_library_cache",
		"kotlin_lang/highlight/DefaultArgumentStubGenerator_with_library_cache",
		"kotlin_lang/highlight/RawFirBuilder_with_library_cache",
		"kotlin_language_server/highlight/Compiler_with_library_cache",
		"kotlin_language_server/highlight/Completions_with_library_cache",
		"kotlin_language_server/highlight/CompletionsTest_with_library_cache",
		"kotlin_language_server/highlight/JavaElementConverter_with_library_cache",
		"kotlin_language_server/highlight/KotlinTextDocumentService_with_library_cache",
		"kotlin_language_server/highlight/QuickFixesTest_with_library_cache",
		"kotlin_language_server/highlight/SourcePath_with_library_cache",
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
		"ktor_samples_mongodb/highlight/ApplicationTest_with_library_cache",
		"ktor_samples_httpbin/highlight/HttpBinApplication_with_library_cache",
		"ktor_samples_youkube/highlight/Upload_with_library_cache",
		"ktor_samples_youkube/highlight/Videos_with_library_cache",
		"ktor_samples_location-header/highlight/LocationHeaderApplication_with_library_cache",
		"ktor_samples_reverse-proxy/highlight/ReverseProxyApplication_with_library_cache",
		"ktor_samples_sse/highlight/SseApplication_with_library_cache",
		"leak-canary-android/highlight/ByteArrayTimSort_with_library_cache",
		"leak-canary-android/highlight/HeapObject_with_library_cache",
		"leak-canary-android/highlight/HprofInMemoryIndex_with_library_cache",
		"leak-canary-android/highlight/HprofRecordReader_with_library_cache",
		"leak-canary-android/highlight/HprofWriter_with_library_cache",
		"leak-canary-android/highlight/LeakStatusTest_with_library_cache",
		"leak-canary-android/highlight/Neo4JCommand_with_library_cache",
		"anki-android/highlight/AbstractSchedTest_with_library_cache",
		"anki-android/highlight/ACRATest_with_library_cache",
		"anki-android/highlight/Finder_with_library_cache",
		"anki-android/highlight/FlashCardsContract_with_library_cache",
		"anki-android/highlight/MetaDB_with_library_cache",
		"anki-android/highlight/TagsUtilTest_with_library_cache",
		"anki-android/highlight/TokenizerTest_with_library_cache",
		"anki-android/highlight/UniqueArrayListTest_with_library_cache",
		"arrow/highlight/build.gradle_with_library_cache",
		"kotlin_empty_kts/highlight/build.gradle_with_library_cache",
		"kotlin_lang/highlight/build.gradle_with_library_cache",
		"toolbox_enterprise/highlight/removedImports/IdeSettingControllerTest_with_library_cache",
		"intellij_commit/highlight/removedImports/UIAutomationInteractionModel.Generated_with_library_cache",
		"kotlin_language_server/highlight/removedImports/CompletionsTest_with_library_cache",
		"kotlin_lang/highlight/removedImports/DefaultArgumentStubGenerator_with_library_cache",
		"spring-framework/highlight/BeanDefinitionDsl_with_library_cache",
		"spring-framework/highlight/CoRouterFunctionDsl_with_library_cache",
		"spring-framework/highlight/reactive/RouterFunctionDsl_with_library_cache",
		"spring-framework/highlight/servlet/RouterFunctionDsl_with_library_cache",
		"spring-framework/highlight/StatusResultMatchersDsl_with_library_cache",
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
		"kotlin_coroutines_commit/highlight/BufferedChannel_with_library_cache",
		"kotlin_coroutines_commit/highlight/CoroutineScheduler_with_library_cache",
		"kotlin_coroutines_commit/highlight/JobSupport_with_library_cache",
		"ktor_commit/highlight/AuthBuildersTest_with_library_cache",
		"ktor_commit/highlight/BufferPrimitives_with_library_cache",
		"ktor_commit/highlight/CacheControlMergeTest_with_library_cache",
		"ktor_commit/highlight/ChunkBufferNativeTest_with_library_cache",
		"ktor_commit/highlight/ConcurrentMapJs_with_library_cache",
		"ktor_commit/highlight/ContentTestSuite_with_library_cache",
		"ktor_commit/highlight/CryptoMingw_with_library_cache",
		"ktor_commit/highlight/HighLoadHttpGenerator_with_library_cache",
		"ktor_commit/highlight/LockFreeLinkedList_with_library_cache",
		"ktor_commit/highlight/OAuth2_with_library_cache",
		"ktor_commit/highlight/ThreadInfo_with_library_cache",
		"space_specific/highlight/DocumentsStatsTest_with_library_cache",
		"space_specific/highlight/EventCounters_with_library_cache",
		"space_specific/highlight/HttpApiClientTest_with_library_cache",
		"space_specific/highlight/M2ChannelMessageListVmV2_with_library_cache",
		"space_specific/highlight/M2ChannelVm_with_library_cache",
		"space_specific/highlight/MagicBarComponent_with_library_cache",
		"space_specific/highlight/ProfilesImpl_with_library_cache",
		"space_specific/highlight/SchemaMigration_with_library_cache",
		"space_specific/highlight/SecretsTests_with_library_cache",
		"space_specific/highlight/XScrollable_with_library_cache", "intellij_commit/rename/SqlBlock_SqlBlockRenamed",
		"kotlin_language_server/insertCode/Rename_renameSymbol",
		"kotlin_language_server/insertCode/SpecialJavaFileForTest_j2k",
		"intellij_commit/otimizeImports/AbstractKotlinMavenImporterTest",
		"intellij_commit/otimizeImports/FSD",
		"intellij_commit/otimizeImports/DiagramsModel.Generated",
		"intellij_commit/otimizeImports/OraIntrospector",
		"intellij_commit/otimizeImports/QuickFixRegistrar",
		"intellij_commit/otimizeImports/SwiftTypeAssignabilityTest",
		"intellij_commit/otimizeImports/TerraformConfigCompletionContributor", "intellij_commit/findUsages/loadModuleEntity_with_library_cache",
		"intellij_commit/findUsages/setUp_with_library_cache",
		"intellij_commit/findUsages/SolutionModel_with_library_cache",
		"intellij_commit/findUsages/SqlBlock_with_library_cache",
		"intellij_commit/findUsages/UIAutomationInteractionModel_with_library_cache",
		"kotlin_lang/findUsages/CommonParser_with_library_cache",
		"kotlin_lang/findUsages/DefaultArgumentStubGenerator_with_library_cache",
		"kotlin_lang/findUsages/FirErrors_with_library_cache",
		"kotlin_lang/findUsages/Flag_with_library_cache",
		"kotlin_lang/findUsages/ReferenceSymbolTable_with_library_cache",
		"toolbox_enterprise/findUsages/ErrorReport_with_library_cache",
		"toolbox_enterprise/findUsages/genUuid_with_library_cache",
		"toolbox_enterprise/findUsages/getAll_with_library_cache",
		"toolbox_enterprise/findUsages/getTempDirectory_with_library_cache",
		"toolbox_enterprise/findUsages/PrincipalContext_with_library_cache",
		"toolbox_enterprise/findUsages/RequestMapping_with_library_cache",
		"toolbox_enterprise/findUsages/RestController_with_library_cache",
		"toolbox_enterprise/findUsages/ROLE_ADMIN_with_library_cache",
		"leak-canary-android/findUsages/BOOLEAN_with_library_cache",
		"leak-canary-android/findUsages/HeapGraph_with_library_cache",
		"leak-canary-android/findUsages/HeapObject_with_library_cache",
		"leak-canary-android/findUsages/HprofRecordTag_with_library_cache",
		"leak-canary-android/findUsages/INT_with_library_cache",
		"leak-canary-android/findUsages/PrimitiveType_with_library_cache",
		"anki-android/findUsages/Card_with_library_cache",
		"anki-android/findUsages/CompatHelper_with_library_cache",
		"anki-android/findUsages/Decks_with_library_cache",
		"anki-android/findUsages/load_with_library_cache",
		"kotlin_coroutines_commit/findUsages/assert_with_library_cache",
		"kotlin_coroutines_commit/findUsages/emit_with_library_cache",
		"kotlin_coroutines_commit/findUsages/Flow_with_library_cache",
		"kotlin_coroutines_commit/findUsages/FlowCollector_with_library_cache",
		"kotlin_coroutines_commit/findUsages/runBlocking_with_library_cache",
		"ktor_commit/findUsages/ByteReadChannel_with_library_cache",
		"ktor_commit/findUsages/HttpClient_jvm_with_library_cache",
		"ktor_commit/findUsages/HttpClient_with_library_cache",
		"ktor_commit/findUsages/toHttpDateString_with_library_cache",
		"space_specific/findUsages/ApiFlag_with_library_cache",
		"space_specific/findUsages/Http_with_library_cache",
		"space_specific/findUsages/IntSizePx_with_library_cache",
		"space_specific/findUsages/UniqueConstraint_with_library_cache", "kotlin_language_server/evaluate-expression/ClassPathTest_with_library_cache",
		"kotlin_language_server/evaluate-expression/Debouncer_with_library_cache",
		"kotlin_language_server/evaluate-expression/KotlinTextDocumentService_with_library_cache",
		"kotlin_petclinic/evaluate-expression/CacheConfig/sleep-1000_with_library_cache",
		"intellij_commit/evaluate-expression/DumbServiceStartupActivity_with_library_cache",
		"intellij_commit/completion/evaluate-expression_with_library_cache",
		"kotlin_petclinic/completion/evaluate-expression/typing-it_with_library_cache",
		"kotlin_petclinic/completion/evaluate-expression/typing-system_with_library_cache"}
	tests := generateKotlinTests(testNames)
	metrics := []string{
		"completion#mean_value", "findUsages#mean_value",
		"semanticHighlighting#mean_value", "localInspections#mean_value",
		"completion#firstElementShown#mean_value", "evaluateExpression#mean_value",
		"performInlineRename#mean_value", "startInlineRename#mean_value",
		"prepareForRename#mean_value", "execute_editor_optimizeimports"}
	aliases := map[string]string{
		"completion#mean_value":                   "completion",
		"completion#firstElementShown#mean_value": "completion",
		"findUsages#mean_value":                   "findUsages",
		"semanticHighlighting#mean_value":         "highlighting",
		"localInspections#mean_value":             "highlighting",
		"performInlineRename#mean_value":          "rename",
		"prepareForRename#mean_value":             "rename",
		"startInlineRename#mean_value":            "rename",
		"execute_editor_optimizeimports":          "optimizeimports",
		"evaluateExpression#mean_value":           "debugger",
	}
	settings := make([]detector.PerformanceSettings, 0, len(testNames)*len(metrics)*2)

	for _, test := range tests {
		for _, metric := range metrics {
			alias := getAlias(metric, aliases)
			settings = append(settings, detector.PerformanceSettings{
				Db:          "perfint",
				Table:       "kotlin",
				Machine:     "intellij-linux-hw-hetzner%",
				Project:     test,
				Metric:      metric,
				Branch:      "master",
				MetricAlias: alias,
				SlackSettings: detector.SlackSettings{
					Channel:     "kotlin-plugin-perf-tests",
					ProductLink: "kotlin",
				},
				AnalysisSettings: detector.AnalysisSettings{
					ReportType: detector.DegradationEvent,
				},
			})
		}
	}
	for _, test := range tests {
		for _, metric := range metrics {
			settings = append(settings, detector.PerformanceSettings{
				Db:          "perfint",
				Table:       "kotlin",
				Machine:     "intellij-linux-hw-hetzner%",
				Project:     test,
				Metric:      metric,
				Branch:      "master",
				MetricAlias: "perf",
				SlackSettings: detector.SlackSettings{
					Channel:     "kotlin-plugin-perf-test-merged",
					ProductLink: "kotlin",
				},
				AnalysisSettings: detector.AnalysisSettings{
					ReportType: detector.DegradationEvent,
				},
			})
		}
	}
	for _, test := range tests {
		for _, metric := range metrics {
			alias := getAlias(metric, aliases)
			settings = append(settings, detector.PerformanceSettings{
				Db:          "perfint",
				Table:       "kotlin",
				Machine:     "intellij-linux-hw-hetzner%",
				Project:     test,
				Metric:      metric,
				Branch:      "master",
				MetricAlias: alias,
				SlackSettings: detector.SlackSettings{
					Channel:     "kotlin-plugin-perf-tests-optimization",
					ProductLink: "kotlin",
				},
				AnalysisSettings: detector.AnalysisSettings{
					ReportType: detector.ImprovementEvent,
				},
			})
		}
	}
	for _, test := range tests {
		for _, metric := range metrics {
			alias := getAlias(metric, aliases)
			settings = append(settings, detector.PerformanceSettings{
				Db:          "perfintDev",
				Table:       "kotlin",
				Machine:     "intellij-linux-hw-hetzner%",
				Project:     test,
				Metric:      metric,
				Branch:      "kt-master",
				MetricAlias: alias,
				SlackSettings: detector.SlackSettings{
					Channel:     "kotlin-plugin-perf-tests-kt-master",
					ProductLink: "kotlin",
				},
				AnalysisSettings: detector.AnalysisSettings{
					ReportType: detector.DegradationEvent,
				},
			})
		}
	}
	return settings
}

func generateKotlinTests(tests []string) []string {
	k1K2tests := make([]string, 0, len(tests)*2)
	for _, test := range tests {
		k1K2tests = append(k1K2tests, test+"_k1")
		k1K2tests = append(k1K2tests, test+"_k2")
	}
	return k1K2tests
}

func getAlias(metric string, aliases map[string]string) string {
	alias, ok := aliases[metric]
	if !ok {
		alias = metric
	}
	return alias
}
