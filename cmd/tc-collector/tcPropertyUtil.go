package main

import (
  "github.com/magiconair/properties"
  "strings"
)

var propertyParserLoader = &properties.Loader{Encoding: properties.UTF8}

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

func readProperties(data []byte) ([]byte, error) {
  p, err := propertyParserLoader.LoadBytes(data)
  if err != nil {
    return nil, err
  }

  json := PropertiesToJson(p)
  return []byte(json), nil
}

// cat '/Volumes/data/Downloads/build.finish.properties' | python -m json.tool > f.json
func isExcludedProperty(key string) bool {
  if excludedTcProperties[key] ||
    // ignore dep
    strings.HasPrefix(key, "dep.") ||
    strings.HasPrefix(key, "teamcity.nuget.") ||
    strings.HasPrefix(key, "teamcity.torrent.") ||
    strings.HasPrefix(key, "secure:teamcity.") ||
    strings.HasPrefix(key, "intellij.plugins.pluginrobot.") ||
    strings.HasPrefix(key, "intellij.influx.startup.") ||
    strings.HasPrefix(key, "env.JAVA_MAIN_CLASS_") ||
    strings.HasPrefix(key, "tools.xcode.arch.appletvos.") ||
    strings.HasPrefix(key, "tools.xcode.arch.iphoneos.") ||
    strings.HasPrefix(key, "tools.xcode.arch.watchos.") ||
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
    strings.Contains(key, ".DotNetFramework1.") ||
    strings.Contains(key, ".npmjs.com.auth.") ||
    strings.Contains(key, ".npm.auth.") {
    return true
  }
  return false
}
