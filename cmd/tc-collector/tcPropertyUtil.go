package main

import "strings"

//noinspection SpellCheckingInspection
var excludedTcProperties = map[string]bool{
  "vcsroot.usernameStyle":                              true,
  "vcsroot.teamcitySshKey":                             true,
  "vcsroot.ijplatform_IntelliJMonorepo.usernameStyle":  true,
  "vcsroot.ijplatform_IntelliJMonorepo.teamcitySshKey": true,
  "vcsroot.ijplatform_IntelliJMonorepo.authMethod":     true,
  "env.ARTIFACTORY_API_KEY":                            true,
  "env.APPL_PASSWORD":                                  true,
  "jetbrains.sign.service.secret":                      true,
}

func isExcludedProperty(key string) bool {
  if excludedTcProperties[key] ||
    strings.HasPrefix(key, "teamcity.nuget.") ||
    strings.HasPrefix(key, "secure:teamcity.") ||
    strings.HasPrefix(key, "intellij.plugins.pluginrobot.") ||
    strings.HasPrefix(key, "intellij.influx.startup.") ||
    strings.HasPrefix(key, "env.JAVA_MAIN_CLASS_") ||
    strings.HasPrefix(key, "npmjs.com.auth.") ||
    strings.HasPrefix(key, "npm.auth.") {
    return true
  }

  if strings.HasSuffix(key, ".user.password") ||
    strings.HasSuffix(key, ".auth.password") ||
    strings.HasSuffix(key, ".user.password") ||
    strings.HasSuffix(key, ".user.name") {
    return true
  }

  // dep.ijplatform_master_Idea_Installers.
  if strings.Contains(key, ".teamcity.nuget.") ||
    strings.Contains(key, ".secure:teamcity.") ||
    strings.Contains(key, ".system.pin.builds.user.password") ||
    strings.Contains(key, ".system.pin.builds.user.name") ||
    strings.Contains(key, ".intellij.plugins.pluginrobot.") ||
    strings.Contains(key, ".intellij.influx.startup.") ||
    strings.Contains(key, ".env.JAVA_MAIN_CLASS_") ||
    strings.Contains(key, ".npmjs.com.auth.") ||
    strings.Contains(key, ".npm.auth.") {
    return true
  }
  return false
}
