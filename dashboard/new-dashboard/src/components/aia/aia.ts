export const aiaModels = ["NativeLLaMA", "JBAICloudWithoutControl", "JBAICloudControl", ""]
export const aiaLanguages = ["css", "csharp", "cpp", "html", "javascript", "java", "kotlin", "php", "python", "ruby", "rust", "scala", "terraform", "typescript"]

export function getAllProjects(prefix: string) {
  return aiaLanguages.flatMap((project) => aiaModels.map((model) => prefix + "_" + project + "_" + model))
}
