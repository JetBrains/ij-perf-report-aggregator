export const aiaModels = ["NativeLLaMA", "JBAICloudWithoutControl", "JBAICloudControl", ""]
export const aiaLanguages = ["css", "csharp", "cpp", "html", "javascript", "java", "kotlin", "php", "python", "ruby", "rust", "scala", "terraform", "typescript"]

export function getAllProjects(prefix: string) {
  const projects = aiaLanguages.flatMap((project) => aiaModels.map((model) => {
    const projectId = prefix + "_" + project + "_" + model;
    if (projectId.includes("code-generation_java_")) {
      // Include both the original java project and its idea alias
      return [projectId, projectId.replace("code-generation_java_", "code-generation_idea_")];
    }
    return projectId;
  }));
  return projects.flat();
}
