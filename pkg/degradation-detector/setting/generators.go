package setting

import (
	"net/http"

	detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"
)

type Generator struct {
	Name     string
	Static   bool
	Generate func(backendUrl string, client *http.Client) []detector.Settings
}

func Generators() []Generator {
	return []Generator{
		static("workspace", GenerateWorkspaceSettings),
		static("kotlin", GenerateKotlinSettings),
		withBackend("kotlinIdea", GenerateKotlinIdeaSettings),
		static("maven", GenerateMavenSettings),
		static("gradle", GenerateGradleSettings),
		static("vcs", GenerateVCSSettings),
		withBackend("phpstorm", GeneratePhpStormSettings),
		withBackend("clion", GenerateClionSettings),
		withBackend("unitTests", GenerateAllUnitTestsSettings),
		withBackend("goland", GenerateGolandPerfSettings),
		withBackend("rust", GenerateRustPerfSettings),
		withBackend("fleetPerformance", GenerateFleetPerformanceSettings),
		withBackend("ruby", GenerateRubyPerfSettings),
		withBackend("java", GenerateJavaSettings),
		withBackend("ultimate", GenerateUltimateSettings),
		static("aia", GenerateAIASettings),
		static("aiaTestToken", GenerateAIATestTokenSettings),
		withBackend("kotlinBuildTools", GenerateKotlinBuildToolsSettings),
		withBackend("kotlinMultiplatformTooling", GenerateKotlinMultiplatformToolingSettings),
		static("ui", GenerateUISettings),
		static("editor", GenerateEditorSettings),
		withBackend("webstorm", GenerateWebStormSettings),
		withBackend("cloud", GenerateCloudSettings),
		withBackend("startupIdea", GenerateStartupSettingsForIDEA),
		withBackend("startupGoland", GenerateStartupSettingsForGoland),
		withBackend("startupPhpStorm", GenerateStartupSettingsForPhpStorm),
		static("fleetStartup", GenerateFleetStartupSettings),
	}
}

func toSettings[T detector.Settings](items []T) []detector.Settings {
	result := make([]detector.Settings, 0, len(items))
	for _, item := range items {
		result = append(result, item)
	}
	return result
}

func static[T detector.Settings](name string, f func() []T) Generator {
	return Generator{
		Name:     name,
		Static:   true,
		Generate: func(_ string, _ *http.Client) []detector.Settings { return toSettings(f()) },
	}
}

func withBackend[T detector.Settings](name string, f func(string, *http.Client) []T) Generator {
	return Generator{
		Name: name,
		Generate: func(backendUrl string, client *http.Client) []detector.Settings {
			return toSettings(f(backendUrl, client))
		},
	}
}
