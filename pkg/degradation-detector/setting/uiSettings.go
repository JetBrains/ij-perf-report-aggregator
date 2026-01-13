package setting

import detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"

func GenerateUISettings() []detector.PerformanceSettings {
	testMetrics := []testMetricDef{
		{test: []string{"intellij_commit/expandProjectMenu"}, metric: []string{"%expandProjectMenu"}},
		{test: []string{"intellij_commit/expandMainMenu"}, metric: []string{"%expandMainMenu"}},
		{test: []string{"intellij_commit/expandEditorMenu"}, metric: []string{"%expandEditorMenu"}},

		{test: []string{"intellij_commit/scrollEditor/java_file_ContentManagerImpl"}, metric: []string{
			"scrollEditor#max_awt_delay",
			"scrollEditor#average_awt_delay",
			"scrollEditor#max_cpu_load",
			"scrollEditor#average_cpu_load",
		}},

		{test: []string{"intellij_commit/projectView", "kotlin_petclinic/projectView"}, metric: []string{"projectViewInit", "projectViewInit#cachedNodesLoaded"}},

		{test: []string{"intellij_commit/find-in-files", "intellij_commit/find-in-files-old"}, metric: []string{
			"findInFiles#openDialog",
			"findInFiles#search: newInstance",
			"findInFiles#search: intellij-ide-starter",
		}},

		{test: []string{"popups-performance-test/test-popups"}, metric: []string{
			"popupShown#EditorContextMenu",
			"popupShown#ProjectViewContextMenu",
			"popupShown#ProjectWidget",
			"popupShown#RunConfigurations",
			"popupShown#VcsLogBranchFilter",
			"popupShown#VcsLogDateFilter",
			"popupShown#VcsLogPathFilter",
			"popupShown#VcsLogUserFilter",
			"popupShown#VcsWidget",
			"afterShow#GitBranchesTreePopup",
			"afterShow#GitDefaultBranchesPopup",
		}},

		{test: []string{"intellij_commit/FileStructureDialog/java_file", "intellij_commit/FileStructureDialog/kotlin_file"}, metric: []string{"FileStructurePopup"}},

		{test: []string{"grails/showIntentions/Find cause", "kotlin/showIntention/Import", "spring_boot/showIntentions"}, metric: []string{"showQuickFixes"}},

		{test: []string{"intellij_commit/findUsages/PsiManager_getInstance_firstCall", "intellij_commit/findUsages/Library_getName"}, metric: []string{"findUsage_popup"}},

		{test: []string{
			"community/go-to-all/Runtime/typingLetterByLetter",
			"community/go-to-all-with-warmup/Runtime/typingLetterByLetter",
			"intellij_commit/go-to-all/Runtime/typingLetterByLetter",
			"intellij_commit/go-to-all-with-warmup/Runtime/typingLetterByLetter",
		}, metric: []string{"searchEverywhere_dialog_shown"}},

		{test: []string{
			"community/go-to-action/Kotlin/typingLetterByLetter",
			"community/go-to-action/Editor/typingLetterByLetter",
			"community/go-to-action/Runtime/typingLetterByLetter",
			"community/go-to-action-finished-embeddings/Runtime/typingLetterByLetter",
			"java/go-to-action/Runtime/typingLetterByLetter",
			"intellij_commit/go-to-action/Kotlin/typingLetterByLetter",
			"intellij_commit/go-to-action/Editor/typingLetterByLetter",
			"intellij_commit/go-to-action/Runtime/typingLetterByLetter",
			"intellij_commit/go-to-action-finished-embeddings/Runtime/typingLetterByLetter",
			"intellij_commit/new-se-go-to-action/Kotlin/typingLetterByLetter",
			"intellij_commit/new-se-go-to-action/Editor/typingLetterByLetter",
			"intellij_commit/new-se-go-to-action/Runtime/typingLetterByLetter",
		}, metric: []string{"searchEverywhere", "searchEverywhere_first_elements_added"}},

		{test: []string{
			"community/go-to-action-with-warmup/Kotlin/typingLetterByLetter",
			"community/go-to-action-with-warmup/Editor/typingLetterByLetter",
			"community/go-to-action-with-warmup/Runtime/typingLetterByLetter",
			"community/go-to-action-with-warmup-finished-embeddings/Runtime/typingLetterByLetter",
			"java/go-to-action-with-warmup/Runtime/typingLetterByLetter",
			"intellij_commit/go-to-action-with-warmup/Kotlin/typingLetterByLetter",
			"intellij_commit/go-to-action-with-warmup/Editor/typingLetterByLetter",
			"intellij_commit/go-to-action-with-warmup/Runtime/typingLetterByLetter",
			"intellij_commit/go-to-action-with-warmup-finished-embeddings/Runtime/typingLetterByLetter",
			"intellij_commit/new-se-go-to-action-with-warmup/Kotlin/typingLetterByLetter",
			"intellij_commit/new-se-go-to-action-with-warmup/Editor/typingLetterByLetter",
			"intellij_commit/new-se-go-to-action-with-warmup/Runtime/typingLetterByLetter",
		}, metric: []string{"searchEverywhere", "searchEverywhere_first_elements_added"}},

		{test: []string{
			"community/go-to-class/Kotlin/typingLetterByLetter",
			"community/go-to-class/Editor/typingLetterByLetter",
			"community/go-to-class/Runtime/typingLetterByLetter",
			"community/go-to-class-finished-embeddings/Runtime/typingLetterByLetter",
			"java/go-to-class/Runtime/typingLetterByLetter",
			"intellij_commit/go-to-class/Kotlin/typingLetterByLetter",
			"intellij_commit/go-to-class/Editor/typingLetterByLetter",
			"intellij_commit/go-to-class/Runtime/typingLetterByLetter",
			"intellij_commit/go-to-class-finished-embeddings/Runtime/typingLetterByLetter",
			"intellij_commit/new-se-go-to-class/Kotlin/typingLetterByLetter",
			"intellij_commit/new-se-go-to-class/Editor/typingLetterByLetter",
			"intellij_commit/new-se-go-to-class/Runtime/typingLetterByLetter",
		}, metric: []string{"searchEverywhere", "searchEverywhere_first_elements_added"}},

		{test: []string{
			"community/go-to-class-with-warmup/Kotlin/typingLetterByLetter",
			"community/go-to-class-with-warmup/Editor/typingLetterByLetter",
			"community/go-to-class-with-warmup/Runtime/typingLetterByLetter",
			"community/go-to-class-with-warmup-finished-embeddings/Runtime/typingLetterByLetter",
			"java/go-to-class-with-warmup/Runtime/typingLetterByLetter",
			"intellij_commit/go-to-class-with-warmup/Kotlin/typingLetterByLetter",
			"intellij_commit/go-to-class-with-warmup/Editor/typingLetterByLetter",
			"intellij_commit/go-to-class-with-warmup/Runtime/typingLetterByLetter",
			"intellij_commit/go-to-class-with-warmup-finished-embeddings/Runtime/typingLetterByLetter",
			"intellij_commit/new-se-go-to-class-with-warmup/Kotlin/typingLetterByLetter",
			"intellij_commit/new-se-go-to-class-with-warmup/Editor/typingLetterByLetter",
			"intellij_commit/new-se-go-to-class-with-warmup/Runtime/typingLetterByLetter",
		}, metric: []string{"searchEverywhere", "searchEverywhere_first_elements_added"}},

		{test: []string{
			"community/go-to-file/Editor/typingLetterByLetter",
			"community/go-to-file/Kotlin/typingLetterByLetter",
			"community/go-to-file/Runtime/typingLetterByLetter",
			"community/go-to-file-finished-embeddings/Runtime/typingLetterByLetter",
			"java/go-to-file/Runtime/typingLetterByLetter",
			"intellij_commit/go-to-file/Kotlin/typingLetterByLetter",
			"intellij_commit/go-to-file/Editor/typingLetterByLetter",
			"intellij_commit/go-to-file/Runtime/typingLetterByLetter",
			"intellij_commit/go-to-file-finished-embeddings/Runtime/typingLetterByLetter",
			"intellij_commit/new-se-go-to-file/Kotlin/typingLetterByLetter",
			"intellij_commit/new-se-go-to-file/Editor/typingLetterByLetter",
			"intellij_commit/new-se-go-to-file/Runtime/typingLetterByLetter",
		}, metric: []string{"searchEverywhere", "searchEverywhere_first_elements_added"}},

		{test: []string{
			"community/go-to-file-with-warmup/Editor/typingLetterByLetter",
			"community/go-to-file-with-warmup/Kotlin/typingLetterByLetter",
			"community/go-to-file-with-warmup/Runtime/typingLetterByLetter",
			"community/go-to-file-with-warmup-finished-embeddings/Runtime/typingLetterByLetter",
			"java/go-to-file-with-warmup/Runtime/typingLetterByLetter",
			"intellij_commit/go-to-file-with-warmup/Kotlin/typingLetterByLetter",
			"intellij_commit/go-to-file-with-warmup/Editor/typingLetterByLetter",
			"intellij_commit/go-to-file-with-warmup/Runtime/typingLetterByLetter",
			"intellij_commit/go-to-file-with-warmup-finished-embeddings/Runtime/typingLetterByLetter",
			"intellij_commit/new-se-go-to-file-with-warmup/Kotlin/typingLetterByLetter",
			"intellij_commit/new-se-go-to-file-with-warmup/Editor/typingLetterByLetter",
			"intellij_commit/new-se-go-to-file-with-warmup/Runtime/typingLetterByLetter",
		}, metric: []string{"searchEverywhere", "searchEverywhere_first_elements_added"}},

		{test: []string{
			"community/go-to-all/Editor/typingLetterByLetter",
			"community/go-to-all/Kotlin/typingLetterByLetter",
			"community/go-to-all/Runtime/typingLetterByLetter",
			"community/go-to-all-finished-embeddings/Runtime/typingLetterByLetter",
			"java/go-to-all/Runtime/typingLetterByLetter",
			"intellij_commit/go-to-all/Kotlin/typingLetterByLetter",
			"intellij_commit/go-to-all/Editor/typingLetterByLetter",
			"intellij_commit/go-to-all/Runtime/typingLetterByLetter",
			"intellij_commit/go-to-all-finished-embeddings/Runtime/typingLetterByLetter",
			"intellij_commit/new-se-go-to-all/Kotlin/typingLetterByLetter",
			"intellij_commit/new-se-go-to-all/Editor/typingLetterByLetter",
			"intellij_commit/new-se-go-to-all/Runtime/typingLetterByLetter",
		}, metric: []string{"searchEverywhere", "searchEverywhere_first_elements_added"}},

		{test: []string{
			"community/go-to-all-with-warmup/Editor/typingLetterByLetter",
			"community/go-to-all-with-warmup/Kotlin/typingLetterByLetter",
			"community/go-to-all-with-warmup/Runtime/typingLetterByLetter",
			"community/go-to-all-with-warmup-finished-embeddings/Runtime/typingLetterByLetter",
			"java/go-to-all-with-warmup/Runtime/typingLetterByLetter",
			"intellij_commit/go-to-all-with-warmup/Kotlin/typingLetterByLetter",
			"intellij_commit/go-to-all-with-warmup/Editor/typingLetterByLetter",
			"intellij_commit/go-to-all-with-warmup/Runtime/typingLetterByLetter",
			"intellij_commit/go-to-all-with-warmup-finished-embeddings/Runtime/typingLetterByLetter",
			"intellij_commit/new-se-go-to-all-with-warmup/Kotlin/typingLetterByLetter",
			"intellij_commit/new-se-go-to-all-with-warmup/Editor/typingLetterByLetter",
			"intellij_commit/new-se-go-to-all-with-warmup/Runtime/typingLetterByLetter",
		}, metric: []string{"searchEverywhere", "searchEverywhere_first_elements_added"}},

		{test: []string{
			"community/go-to-symbol/Editor/typingLetterByLetter",
			"community/go-to-symbol/Kotlin/typingLetterByLetter",
			"community/go-to-symbol/Runtime/typingLetterByLetter",
			"community/go-to-symbol-finished-embeddings/Runtime/typingLetterByLetter",
			"java/go-to-symbol/Runtime/typingLetterByLetter",
			"intellij_commit/go-to-symbol/Kotlin/typingLetterByLetter",
			"intellij_commit/go-to-symbol/Editor/typingLetterByLetter",
			"intellij_commit/go-to-symbol/Runtime/typingLetterByLetter",
			"intellij_commit/go-to-symbol-finished-embeddings/Runtime/typingLetterByLetter",
			"intellij_commit/new-se-go-to-symbol/Kotlin/typingLetterByLetter",
			"intellij_commit/new-se-go-to-symbol/Editor/typingLetterByLetter",
			"intellij_commit/new-se-go-to-symbol/Runtime/typingLetterByLetter",
		}, metric: []string{"searchEverywhere", "searchEverywhere_first_elements_added"}},

		{test: []string{
			"community/go-to-symbol-with-warmup/Editor/typingLetterByLetter",
			"community/go-to-symbol-with-warmup/Kotlin/typingLetterByLetter",
			"community/go-to-symbol-with-warmup/Runtime/typingLetterByLetter",
			"community/go-to-symbol-with-warmup-finished-embeddings/Runtime/typingLetterByLetter",
			"java/go-to-symbol-with-warmup/Runtime/typingLetterByLetter",
			"intellij_commit/go-to-symbol-with-warmup/Kotlin/typingLetterByLetter",
			"intellij_commit/go-to-symbol-with-warmup/Editor/typingLetterByLetter",
			"intellij_commit/go-to-symbol-with-warmup/Runtime/typingLetterByLetter",
			"intellij_commit/go-to-symbol-with-warmup-finished-embeddings/Runtime/typingLetterByLetter",
			"intellij_commit/new-se-go-to-symbol-with-warmup/Kotlin/typingLetterByLetter",
			"intellij_commit/new-se-go-to-symbol-with-warmup/Editor/typingLetterByLetter",
			"intellij_commit/new-se-go-to-symbol-with-warmup/Runtime/typingLetterByLetter",
		}, metric: []string{"searchEverywhere", "searchEverywhere_first_elements_added"}},

		{test: []string{
			"community/go-to-text/Editor/typingLetterByLetter",
			"community/go-to-text/Kotlin/typingLetterByLetter",
			"community/go-to-text/Runtime/typingLetterByLetter",
			"community/go-to-text-finished-embeddings/Runtime/typingLetterByLetter",
			"java/go-to-text/Runtime/typingLetterByLetter",
			"intellij_commit/go-to-text/Kotlin/typingLetterByLetter",
			"intellij_commit/go-to-text/Editor/typingLetterByLetter",
			"intellij_commit/go-to-text/Runtime/typingLetterByLetter",
			"intellij_commit/go-to-text-finished-embeddings/Runtime/typingLetterByLetter",
			"intellij_commit/new-se-go-to-text/Kotlin/typingLetterByLetter",
			"intellij_commit/new-se-go-to-text/Editor/typingLetterByLetter",
			"intellij_commit/new-se-go-to-text/Runtime/typingLetterByLetter",
		}, metric: []string{"searchEverywhere", "searchEverywhere_first_elements_added"}},

		{test: []string{
			"community/go-to-text-with-warmup/Editor/typingLetterByLetter",
			"community/go-to-text-with-warmup/Kotlin/typingLetterByLetter",
			"community/go-to-text-with-warmup/Runtime/typingLetterByLetter",
			"community/go-to-text-with-warmup-finished-embeddings/Runtime/typingLetterByLetter",
			"java/go-to-text-with-warmup/Runtime/typingLetterByLetter",
			"intellij_commit/go-to-text-with-warmup/Kotlin/typingLetterByLetter",
			"intellij_commit/go-to-text-with-warmup/Editor/typingLetterByLetter",
			"intellij_commit/go-to-text-with-warmup/Runtime/typingLetterByLetter",
			"intellij_commit/go-to-text-with-warmup-finished-embeddings/Runtime/typingLetterByLetter",
			"intellij_commit/new-se-go-to-text-with-warmup/Kotlin/typingLetterByLetter",
			"intellij_commit/new-se-go-to-text-with-warmup/Editor/typingLetterByLetter",
			"intellij_commit/new-se-go-to-text-with-warmup/Runtime/typingLetterByLetter",
		}, metric: []string{"searchEverywhere", "searchEverywhere_first_elements_added"}},
	}

	machines := []string{"intellij-linux-performance-aws-%", "intellij-windows-performance-%"}

	settings := make([]detector.PerformanceSettings, 0, 100)
	for _, testMetric := range testMetrics {
		for _, test := range testMetric.test {
			for _, metric := range testMetric.metric {
				for _, machine := range machines {
					settings = append(settings, detector.PerformanceSettings{
						Db:      "perfintDev",
						Table:   "idea",
						Project: test,
						BaseSettings: detector.BaseSettings{
							Machine: machine,
							Metric:  metric,
							Branch:  "master",
							SlackSettings: detector.SlackSettings{
								Channel:     "ij-ui-performance-alerts",
								ProductLink: "intellij",
							},
							AnalysisSettings: detector.AnalysisSettings{
								MinimumSegmentLength: 8,
							},
						},
					})
				}
			}
		}
	}
	return settings
}
