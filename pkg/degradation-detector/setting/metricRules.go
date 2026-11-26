package setting

import "strings"

// The rules that pick the metrics of a test by its name, for a block of metric-degradation that lists no metrics.
// The key is the value of `metricRules` in the spec of the notification.
var metricRules = map[string]func(test string) []string{
	"common":      getMetricFromTestName,
	"clion":       getClionMetricFromTestName,
	"kotlin-idea": withoutIndexingAndScanning(getKotlinMetricsFromTestName),
	"java":        withoutIndexingAndScanning(getJavaMetricsFromTestsNames),
	"ultimate":    getUltimateMetricsFromTestsNames,
	"cloud":       getCloudMetricsFromTestsNames,
	"rust":        getRustMetricFromTestName,
	"ruby":        getRubyMetricFromTestName,
}

// The Kotlin and Java rules never watched the indexing and scanning tests of their projects.
func withoutIndexingAndScanning(rule func(test string) []string) func(test string) []string {
	return func(test string) []string {
		if strings.Contains(test, "indexing") || strings.Contains(test, "scanning") {
			return nil
		}
		return rule(test)
	}
}

func getMetricFromTestName(test string) []string {
	if strings.Contains(test, "/vfsRefresh") {
		return []string{"vfs_initial_refresh"}
	}
	if strings.Contains(test, "/rebuild") {
		return []string{"build_compilation_duration"}
	}
	if strings.Contains(test, "/inspection") {
		return []string{"globalInspections"}
	}
	if strings.Contains(test, "/localInspection") {
		return []string{"localInspections", "firstCodeAnalysis"}
	}
	if strings.Contains(test, "/completion") {
		return []string{"completion"}
	}
	if strings.Contains(test, "/debug") {
		return []string{"debugRunConfiguration", "debugStep_into"}
	}
	if strings.Contains(test, "/showIntentions") {
		return []string{"Test#average_awt_delay", "showQuickFixes"}
	}
	if strings.Contains(test, "/showFileHistory") {
		return []string{"showFileHistory"}
	}
	if strings.Contains(test, "/expandProjectMenu") {
		return []string{"%expandProjectMenu"}
	}
	if strings.Contains(test, "/expandMainMenu") {
		return []string{"%expandMainMenu"}
	}
	if strings.Contains(test, "/expandEditorMenu") {
		return []string{"%expandEditorMenu"}
	}
	if strings.Contains(test, "/highlight") {
		return []string{"highlighting"}
	}
	if strings.Contains(test, "/FileStructureDialog") {
		return []string{"FileStructurePopup"}
	}
	if strings.Contains(test, "/createJavaClass") {
		return []string{"createJavaFile"}
	}
	if strings.Contains(test, "/createKotlinClass") {
		return []string{"createKotlinFile"}
	}
	if strings.Contains(test, "/indexing") {
		return []string{"scanningTimeWithoutPauses", "indexingTimeWithoutPauses"}
	}
	if strings.Contains(test, "/inlineRename") {
		return []string{"startInlineRename"}
	}
	if strings.Contains(test, "/-scanning") {
		return []string{"scanningTimeWithoutPauses"}
	}
	if strings.Contains(test, "/findUsages") {
		return []string{"findUsages", "fus_find_usages_all", "fus_find_usages_first"}
	}
	if strings.Contains(test, "/go-to-") {
		return []string{"searchEverywhere"}
	}
	if strings.Contains(test, "/ultimate") {
		return []string{"localInspections", "firstCodeAnalysis", "typingCodeAnalyzing", "completion"}
	}
	if strings.Contains(test, "/typing") {
		return []string{"typingCodeAnalyzing", "typing"}
	}
	if strings.Contains(test, "/scrollEditor") {
		return []string{"scrollEditor#average_awt_delay", "scrollEditor#max_awt_delay", "scrollEditor#average_cpu_load", "scrollEditor#max_cpu_load"}
	}
	return []string{}
}

func getClionMetricFromTestName(test string) []string {
	if strings.Contains(test, "/gotoDeclaration") {
		return []string{
			"clionGotoDeclaration",
		}
	}
	if strings.Contains(test, "checkLocalTestConfig") {
		return []string{"waitFirstTestGutter"}
	}
	if strings.Contains(test, "/indexing") {
		return []string{"ocSymbolBuildingTimeMs", "backendIndexingTimeMs", "rd.memory.allocatedManagedMemoryMb/afterIndexing"}
	}
	if strings.Contains(test, "/completion") {
		return []string{"fus_time_to_show_90p"}
	}
	if strings.Contains(test, "/findUsages") {
		return []string{"findUsagesInToolWindow", "findUsagesInToolWindow#number"}
	}
	if strings.Contains(test, "/gotoDeclaration") {
		return []string{"clionGotoDeclaration"}
	}
	if strings.Contains(test, "/go-to-all-with-warmup") {
		return []string{"searchEverywhere", "searchEverywhere_first_elements_added"}
	}
	if strings.Contains(test, "/checkLocalTestConfig") {
		return []string{"waitFirstTestGutter"}
	}
	if strings.Contains(test, "/measureResolve") {
		return []string{"nova_total_memory_mb", "nova_resolving_references_ms"}
	}
	if strings.Contains(test, "/typing") {
		return []string{"typing#latency#mean_value"}
	}
	return getMetricFromTestName(test)
}

func getKotlinMetricsFromTestName(test string) []string {
	if strings.Contains(test, "/editor-kotlin-highlighting") {
		return []string{"typing_EditorBackSpace_duration", "typing_EditorBackSpace_warmup_duration", "typing_}_duration", "typing_}_warmup_duration"}
	}
	if strings.Contains(test, "/red-code-kotlin") {
		return []string{"replaceTextCodeAnalysis"}
	}
	if strings.Contains(test, "/createKotlinClass") {
		return []string{"createKotlinFile"}
	}
	if strings.Contains(test, "kotlin/localInspection") || strings.Contains(test, "kotlin_coroutines/localInspection") {
		return []string{"firstCodeAnalysis", "localInspections", "fus_file_types_usage_duration_ms", "fus_file_types_usage_time_to_show_ms"}
	}
	if strings.Contains(test, "kotlin/inspection") || strings.Contains(test, "kotlin_coroutines/inspection") {
		return []string{"globalInspections", "fullGCPause", "JVM.GC.collectionTimesMs"}
	}
	return []string{}
}

func getJavaMetricsFromTestsNames(test string) []string {
	if strings.Contains(test, "/inspection") {
		return []string{"globalInspections"}
	}
	if strings.Contains(test, "/showIntentions") {
		return []string{"Test#average_awt_delay", "showQuickFixes"}
	}
	if strings.Contains(test, "/localInspection/java_file") {
		return []string{"localInspections", "firstCodeAnalysis", "fus_file_types_usage_duration_ms"}
	}
	if strings.Contains(test, "/completion") {
		return []string{"completion"}
	}
	if strings.Contains(test, "/showIntentions") {
		return []string{"Test#average_awt_delay", "showQuickFixes"}
	}
	if strings.Contains(test, "/createJavaClass") {
		return []string{"createJavaFile"}
	}
	if strings.Contains(test, "/rename-method") {
		return []string{"performInlineRename"}
	}
	if strings.Contains(test, "/rename-class") {
		return []string{"performInlineRename"}
	}
	if strings.Contains(test, "/change-signature") {
		return []string{"changeJavaSignature: add parameter"}
	}
	if strings.Contains(test, "/move-class") {
		return []string{"moveClassToPackage"}
	}
	if strings.Contains(test, "/inline-method") {
		return []string{"inlineJavaMethod"}
	}
	if strings.Contains(test, "/rename-package") {
		return []string{"renameDirectoryAsPackage"}
	}
	if strings.Contains(test, "/editor-highlighting") {
		return []string{"typing_EditorBackSpace_duration", "typing_EditorBackSpace_warmup_duration", "typing_}_duration", "typing_}_warmup_duration"}
	}

	// for the future cases
	if strings.Contains(test, "/findUsages") {
		return []string{"findUsages", "fus_find_usages_all", "fus_find_usages_first"}
	}
	return []string{}
}

func getUltimateMetricsFromTestsNames(test string) []string {
	if strings.Contains(test, "/localInspection") {
		return []string{"localInspections", "firstCodeAnalysis", "fus_file_types_usage_duration_ms", "fus_file_types_usage_time_to_show_ms"}
	}
	if strings.Contains(test, "/ultimateCase") {
		return []string{"localInspections", "firstCodeAnalysis", "completion", "fus_file_types_usage_duration_ms", "fus_file_types_usage_time_to_show_ms"}
	}
	if strings.Contains(test, "/typing") {
		return []string{"typingCodeAnalyzing", "test#average_awt_delay", "test#max_awt_delay", "firstCodeAnalysis", "fus_file_types_usage_duration_ms", "fus_file_types_usage_time_to_show_ms"}
	}
	if strings.Contains(test, "/completion") {
		return []string{"completion", "firstCodeAnalysis", "fus_file_types_usage_duration_ms", "fus_file_types_usage_time_to_show_ms"}
	}
	if strings.Contains(test, "/indexing") {
		return []string{"indexingTimeWithoutPauses", "scanningTimeWithoutPauses"}
	}

	return []string{}
}

func getCloudMetricsFromTestsNames(test string) []string {
	if strings.Contains(test, "/localInspection") {
		return []string{"localInspections", "firstCodeAnalysis", "fus_file_types_usage_duration_ms", "fus_file_types_usage_time_to_show_ms"}
	}
	if strings.Contains(test, "/typing") {
		return []string{"typingCodeAnalyzing", "test#average_awt_delay", "test#max_awt_delay", "firstCodeAnalysis", "fus_file_types_usage_duration_ms", "fus_file_types_usage_time_to_show_ms"}
	}
	if strings.Contains(test, "/completion") {
		return []string{"completion", "firstCodeAnalysis", "fus_file_types_usage_duration_ms", "fus_file_types_usage_time_to_show_ms"}
	}
	if strings.Contains(test, "/indexing") {
		return []string{"indexingTimeWithoutPauses", "scanningTimeWithoutPauses"}
	}

	return []string{}
}

func getRustMetricFromTestName(test string) []string {
	if strings.Contains(test, "/indexing") {
		return []string{
			"rust_duration_from_start_to_work",
			"rust_def_maps_execution_time", "rust_macro_expansion_execution_time", "rust_class_instances_tree_size_mb", "indexingTimeWithoutPauses",
		}
	}
	if strings.Contains(test, "/local-inspection") {
		return []string{"typingCodeAnalyzing#mean_value", "firstCodeAnalysis"}
	}
	if strings.Contains(test, "/global-inspection") {
		return []string{"globalInspections"}
	}
	if strings.Contains(test, "/completion") {
		return []string{"completion#mean_value"}
	}
	if strings.Contains(test, "/findUsages") {
		return []string{"findUsages"}
	}
	if strings.Contains(test, "/typing") {
		return []string{"typing#latency#mean_value"}
	}
	return []string{}
}

func getRubyMetricFromTestName(test string) []string {
	if strings.Contains(test, "/indexing") {
		return []string{
			"indexingTimeWithoutPauses",
			"numberOfIndexedFilesWritingIndexValue", "indexSize", "scanningTimeWithoutPauses",
			"gcPause", "freedMemoryByGC",
		}
	}
	if strings.Contains(test, "/firstCodeAnalysis") {
		return []string{"firstCodeAnalysis#mean_value"}
	}
	if strings.Contains(test, "inspections-test/") {
		return []string{"globalInspections", "gcPause", "freedMemoryByGC"}
	}
	if strings.Contains(test, "/typing") {
		return []string{"test#average_awt_delay", "typing#median_value"}
	}
	if strings.Contains(test, "/completion") {
		return []string{"completion#mean_value", "completion#firstElementShown#mean_value", "completion#number#mean_value"}
	}
	if strings.Contains(test, "/getSymbolMembers") {
		return []string{"getSymbolMembers#mean_value", "getSymbolMembers#number#mean_value"}
	}
	// very likely find usages
	return []string{"findUsages", "findUsages#number", "gcPause", "freedMemoryByGC"}
}
