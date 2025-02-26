export const aiaModels = ["NativeLLaMA", "JBAICloudWithoutControl", "JBAICloudControl", ""]
export const aiaLanguages = ["css", "csharp", "cpp", "html", "javascript", "idea", "kotlin", "php", "python", "ruby", "rust", "scala", "terraform"]

export function getAllProjects(prefix: string) {
  return aiaLanguages.flatMap((project) => aiaModels.map((model) => prefix + "_" + project + "_" + model))
}
