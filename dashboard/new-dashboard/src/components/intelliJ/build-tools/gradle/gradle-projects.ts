const PROJECT_NAMES = [
  "grazie-platform-project-import-gradle/",
  "project-import-gradle-monolith-51-modules-4000-dependencies-2000000-files/",
  "project-import-gradle-micronaut/",
  "project-import-gradle-hibernate-orm/",
  "project-import-gradle-cas/",
  "project-reimport-gradle-cas/",
  "project-import-from-cache-gradle-cas/",
  "project-import-gradle-1000-modules/",
  "project-import-gradle-1000-modules-limited-ram/",
  "project-import-gradle-5000-modules/",
  "project-import-gradle-android-extra-large/",
  "project-import-android-500-modules/",
  "project-reimport-space/",
  "project-import-space/",
  "project-import-open-telemetry/",
]

export const GRADLE_PROJECTS = PROJECT_NAMES.map((name) => name + "measureStartup")

export const GRADLE_PROJECTS_FAST_INSTALLERS = PROJECT_NAMES.map((name) => name + "fastInstaller")
