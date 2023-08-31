const KOTLIN_TEST_REGEXP = new RegExp("_with_library_cache_k.+")
export function replaceKotlinName(name: string) {
  if (name.indexOf("/") > 0) {
    const project = name.slice(0, name.indexOf("/"))
    const test = name.slice(name.lastIndexOf("/") + 1, name.length).replace(KOTLIN_TEST_REGEXP, "")
    return projectByName(project) + "/" + test
  }
  return name
}

function projectByName(name: string) {
  switch (name) {
    case "intellij_commit": {
      return "IJ"
    }
    case "kotlin_language_server": {
      return "KLS"
    }
    case "kotlin_lang": {
      return "KL"
    }
    case "toolbox_enterprise": {
      return "TBE"
    }
    default:
      return name
  }
}
