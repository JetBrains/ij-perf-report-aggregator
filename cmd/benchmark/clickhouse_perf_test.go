package benchmark

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
	"testing"
	"time"
)

func TestClickhouse(t *testing.T) {
	t.Parallel()
	baseURL := "http://localhost:8123/"

	// List of branches to iterate over
	branches := []string{"master"}
	oses := []string{"intellij-linux-performance-aws-%"}

	projects := []string{
		"project-import-maven-quarkus/measureStartup",
		"project-reimport-maven-quarkus/measureStartup",
		"project-import-from-cache-maven-quarkus/measureStartup",
		"project-import-maven-1000-modules/measureStartup",
		"project-import-maven-5000-modules/measureStartup",
		"project-import-maven-keycloak/measureStartup",
		"project-import-maven-javaee7/measureStartup",
		"project-import-maven-javaee8/measureStartup",
		"project-import-maven-jersey/measureStartup",
		"project-import-maven-flink/measureStartup",
		"project-import-maven-drill/measureStartup",
		"project-import-maven-azure-sdk-java/measureStartup",
		"project-import-maven-hive/measureStartup",
		"project-import-maven-quarkus-to-legacy-model/measureStartup",
		"project-import-maven-1000-modules-to-legacy-model/measureStartup",
	}

	metrics := []string{
		"maven.sync.duration",
		"maven.import.after.import.configuration",
		"maven.import.stats.applying.model.task",
		"maven.import.stats.importing.task",
		"maven.import.stats.importing.task.old",
		"maven.project.importer.base.refreshing.files.task",
		"maven.projects.processor.plugin.resolving.task",
		"maven.projects.processor.reading.task",
		"maven.projects.processor.resolving.task",
		"maven.projects.processor.wait.for.completion.task",
		"quarkus.maven.importer.base.task",
		"quarkus.maven.post.processor.task",
		"jpa.facet.importer.lambdas",
		"jpa.facet.importer.reimport.facet",
		"importer_run.com.intellij.jpa.importer.maven.JpaFacetImporter.total_duration_ms",
		"importer_run.com.intellij.quarkus.run.maven.QsMavenImporter.total_duration_ms",
		"importer_run.com.intellij.spring.facet.importer.maven.SpringFacetImporter.total_duration_ms",
		"importer_run.com.intellij.spring.mvc.importer.boot.SpringBootWebFacetImporter.total_duration_ms",
		"importer_run.org.jetbrains.idea.maven.ext.javaee.ear.EarFacetImporter.total_duration_ms",
		"importer_run.org.jetbrains.idea.maven.ext.javaee.web.WebFacetImporter.total_duration_ms",
		"importer_run.org.jetbrains.idea.maven.importing.MavenAnnotationProcessorConfigurator.total_duration_ms",
		"importer_run.org.jetbrains.idea.maven.importing.MavenCompilerConfigurator.total_duration_ms",
		"importer_run.org.jetbrains.idea.maven.importing.MavenEncodingConfigurator.total_duration_ms",
		"importer_run.org.jetbrains.idea.maven.importing.MavenExternalAnnotationsConfigurator.total_duration_ms",
		"importer_run.org.jetbrains.idea.maven.importing.MavenRemoteRepositoriesConfigurator.total_duration_ms",
		"importer_run.org.jetbrains.kotlin.idea.maven.KotlinMavenImporter.total_duration_ms",
		"workspace_commit.attempts",
		"workspace_commit.duration_in_background_ms",
		"workspace_commit.duration_in_write_action_ms",
		"workspace_commit.duration_of_workspace_update_call_ms",
		"workspace_import.commit.duration_ms",
		"workspace_import.configurator_run.com.intellij.spring.facet.importer.maven.SpringFacetImporter.after_apply_duration_ms",
		"workspace_import.configurator_run.com.intellij.spring.facet.importer.maven.SpringFacetImporter.before_apply_duration_ms",
		"workspace_import.configurator_run.com.intellij.spring.facet.importer.maven.SpringFacetImporter.collect_folders_duration_ms",
		"workspace_import.configurator_run.com.intellij.spring.facet.importer.maven.SpringFacetImporter.config_modules_duration_ms",
		"workspace_import.configurator_run.com.intellij.spring.facet.importer.maven.SpringFacetImporter.total_duration_ms",
		"workspace_import.configurator_run.com.intellij.spring.mvc.importer.boot.SpringBootWebFacetImporter.after_apply_duration_ms",
		"workspace_import.configurator_run.com.intellij.spring.mvc.importer.boot.SpringBootWebFacetImporter.before_apply_duration_ms",
		"workspace_import.configurator_run.com.intellij.spring.mvc.importer.boot.SpringBootWebFacetImporter.collect_folders_duration_ms",
		"workspace_import.configurator_run.com.intellij.spring.mvc.importer.boot.SpringBootWebFacetImporter.config_modules_duration_ms",
		"workspace_import.configurator_run.com.intellij.spring.mvc.importer.boot.SpringBootWebFacetImporter.total_duration_ms",
		"workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenAnnotationProcessorConfigurator.after_apply_duration_ms",
		"workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenAnnotationProcessorConfigurator.before_apply_duration_ms",
		"workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenAnnotationProcessorConfigurator.collect_folders_duration_ms",
		"workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenAnnotationProcessorConfigurator.config_modules_duration_ms",
		"workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenAnnotationProcessorConfigurator.total_duration_ms",
		"workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenCompilerConfigurator.after_apply_duration_ms",
		"workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenCompilerConfigurator.before_apply_duration_ms",
		"workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenCompilerConfigurator.collect_folders_duration_ms",
		"workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenCompilerConfigurator.config_modules_duration_ms",
		"workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenCompilerConfigurator.total_duration_ms",
		"workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenEncodingConfigurator.after_apply_duration_ms",
		"workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenEncodingConfigurator.before_apply_duration_ms",
		"workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenEncodingConfigurator.collect_folders_duration_ms",
		"workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenEncodingConfigurator.config_modules_duration_ms",
		"workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenEncodingConfigurator.total_duration_ms",
		"workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenExternalAnnotationsConfigurator.after_apply_duration_ms",
		"workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenExternalAnnotationsConfigurator.before_apply_duration_ms",
		"workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenExternalAnnotationsConfigurator.collect_folders_duration_ms",
		"workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenExternalAnnotationsConfigurator.config_modules_duration_ms",
		"workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenExternalAnnotationsConfigurator.total_duration_ms",
		"workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenRemoteRepositoriesConfigurator.after_apply_duration_ms",
		"workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenRemoteRepositoriesConfigurator.before_apply_duration_ms",
		"workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenRemoteRepositoriesConfigurator.collect_folders_duration_ms",
		"workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenRemoteRepositoriesConfigurator.config_modules_duration_ms",
		"workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenRemoteRepositoriesConfigurator.total_duration_ms",
		"workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenWslTargetConfigurator.after_apply_duration_ms",
		"workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenWslTargetConfigurator.before_apply_duration_ms",
		"workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenWslTargetConfigurator.collect_folders_duration_ms",
		"workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenWslTargetConfigurator.config_modules_duration_ms",
		"workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenWslTargetConfigurator.total_duration_ms",
		"workspace_import.configurator_run.org.jetbrains.idea.maven.plugins.groovy.GroovyPluginConfigurator.after_apply_duration_ms",
		"workspace_import.configurator_run.org.jetbrains.idea.maven.plugins.groovy.GroovyPluginConfigurator.before_apply_duration_ms",
		"workspace_import.configurator_run.org.jetbrains.idea.maven.plugins.groovy.GroovyPluginConfigurator.collect_folders_duration_ms",
		"workspace_import.configurator_run.org.jetbrains.idea.maven.plugins.groovy.GroovyPluginConfigurator.config_modules_duration_ms",
		"workspace_import.configurator_run.org.jetbrains.idea.maven.plugins.groovy.GroovyPluginConfigurator.total_duration_ms",
		"workspace_import.configurator_run.org.jetbrains.idea.maven.ext.javaee.web.WebFacetImporterafter_apply_duration_ms",
		"workspace_import.configurator_run.org.jetbrains.idea.maven.ext.javaee.web.WebFacetImporterbefore_apply_duration_ms",
		"workspace_import.configurator_run.org.jetbrains.idea.maven.ext.javaee.web.WebFacetImportercollect_folders_duration_ms",
		"workspace_import.configurator_run.org.jetbrains.idea.maven.ext.javaee.web.WebFacetImporterconfig_modules_duration_ms",
		"workspace_import.configurator_run.org.jetbrains.idea.maven.ext.javaee.web.WebFacetImportertotal_duration_ms",
		"workspace_import.configurator_run.org.jetbrains.idea.maven.ext.javaee.ear.EarFacetImporter.after_apply_duration_ms",
		"workspace_import.configurator_run.org.jetbrains.idea.maven.ext.javaee.ear.EarFacetImporter.before_apply_duration_ms",
		"workspace_import.configurator_run.org.jetbrains.idea.maven.ext.javaee.ear.EarFacetImporter.collect_folders_duration_ms",
		"workspace_import.configurator_run.org.jetbrains.idea.maven.ext.javaee.ear.EarFacetImporter.config_modules_duration_ms",
		"workspace_import.configurator_run.org.jetbrains.idea.maven.ext.javaee.ear.EarFacetImporter.total_duration_ms",
		"workspace_import.configurator_run.org.jetbrains.idea.maven.ext.javaee.web.WebFacetImporterEx.after_apply_duration_ms",
		"workspace_import.configurator_run.org.jetbrains.idea.maven.ext.javaee.web.WebFacetImporterEx.before_apply_duration_ms",
		"workspace_import.configurator_run.org.jetbrains.idea.maven.ext.javaee.web.WebFacetImporterEx.collect_folders_duration_ms",
		"workspace_import.configurator_run.org.jetbrains.idea.maven.ext.javaee.web.WebFacetImporterEx.config_modules_duration_ms",
		"workspace_import.configurator_run.org.jetbrains.idea.maven.ext.javaee.web.WebFacetImporterEx.total_duration_ms",
		"workspace_import.configurator_run.org.jetbrains.idea.maven.ext.javaee.ear.EarFacetImporterEx.after_apply_duration_ms",
		"workspace_import.configurator_run.org.jetbrains.idea.maven.ext.javaee.ear.EarFacetImporterEx.before_apply_duration_ms",
		"workspace_import.configurator_run.org.jetbrains.idea.maven.ext.javaee.ear.EarFacetImporterEx.collect_folders_duration_ms",
		"workspace_import.configurator_run.org.jetbrains.idea.maven.ext.javaee.ear.EarFacetImporterEx.config_modules_duration_ms",
		"workspace_import.configurator_run.org.jetbrains.idea.maven.ext.javaee.ear.EarFacetImporterEx.total_duration_ms",
		"workspace_import.duration_ms",
		"workspace_import.legacy_importers.duration_ms",
		"workspace_import.legacy_importers.stats.duration_of_bridges_creation_ms",
		"workspace_import.legacy_importers.stats.duration_of_bridges_commit_ms",
		"workspace_import.populate.duration_ms",
		"maven.project.importer.post.importing.task.marker",
		"post_import_tasks_run.total_duration_ms",
		"AWTEventQueue.dispatchTimeTotal",
		"CPU | Load |Total % 95th pctl",
		"Memory | IDE | RESIDENT SIZE (MB) 95th pctl",
		"Memory | IDE | VIRTUAL SIZE (MB) 95th pctl",
		"gcPause",
		"gcPauseCount",
		"fullGCPause",
		"freedMemoryByGC",
		"totalHeapUsedMax",
	}

	const numWorkers = 1

	requests := make(chan string, numWorkers)
	var wg sync.WaitGroup

	worker := func() {
		for query := range requests {
			encodedQuery := url.QueryEscape(query)
			resp, err := http.Get(fmt.Sprintf("%s?query=%s", baseURL, encodedQuery))
			if err != nil {
				log.Printf("Error: %v", err)
				errorsCounterCH.Inc()
			} else {
				resp.Body.Close()
			}
			requestCounterCH.Inc()
			wg.Done()
		}
	}
	// Start the workers.
	for range numWorkers {
		go worker()
	}

	start := time.Now()

	for _, metric := range metrics {
		for _, project := range projects {
			for _, os := range oses {
				for _, branch := range branches {
					query := fmt.Sprintf("select toUnixTimestamp(generated_time)*1000 as `t`, measures.value, measures.name, machine, tc_build_id, project, tc_installer_build_id, build_c1, build_c2, build_c3 from perfint.idea array join measures where branch = '%s' and generated_time >subtractMonths(now(),12) and triggeredBy = '' and machine like '%s' and build_c3=0 and project = '%s' and measures.name = '%s' order by t", branch, os, project, metric)
					// query := fmt.Sprintf("select toUnixTimestamp(generated_time)*1000 as `t`, measures.value, measures.name, machine, tc_build_id, project, tc_installer_build_id, build_c1, build_c2, build_c3 from perfint.idea2 where branch = '%s' and generated_time >subtractMonths(now(),12) and measures.name = '%s' and triggeredBy = '' and machine like '%s' and build_c3=0 and project = '%s' order by t", branch, metric, os, project)
					wg.Add(1)
					requests <- query
				}
			}
		}
	}

	wg.Wait()
	close(requests)

	fmt.Println("Total requests: ", requestCounterCH.Value())
	fmt.Println("Total errors: ", errorsCounterCH.Value())
	elapsed := time.Since(start) // calculate the elapsed time
	fmt.Printf("The code executed in %s\n", elapsed)
}

var (
	requestCounterCH CounterCH
	errorsCounterCH  CounterCH
)

type CounterCH struct {
	mu sync.Mutex
	n  int
}

func (c *CounterCH) Inc() {
	c.mu.Lock()
	c.n++
	c.mu.Unlock()
}

func (c *CounterCH) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.n
}
