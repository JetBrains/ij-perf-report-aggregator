export const aiaModels = ["NativeLLaMA", "JBAICloudWithoutControl", "JBAICloudControl"]
export const aiaLanguages = ["css", "csharp", "cpp", "go", "html", "javascript", "java", "kotlin", "php", "python", "ruby", "rust", "scala", "terraform"]

export function getAllProjects(prefix: string) {
  return aiaLanguages.map((project) => aiaModels.map((model) => prefix + "_" + project + "_" + model)).flat()
}
