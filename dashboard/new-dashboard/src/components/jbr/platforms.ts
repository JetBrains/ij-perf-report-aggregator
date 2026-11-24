/**
 * JBR benchmark projects are named "<benchmark>_<configuration>", the configuration starting with the OS it runs on
 * ("Render_macOS15aarch64Metal", "Render_Ubuntu2404x86_64MM_Wayland", "Render_Windows11x86_64"). Benchmark names hold
 * no underscore.
 */
const PLATFORMS = ["macOS", "Linux", "Windows"] as const

/** The configuration part of a project name. */
export function configurationOf(project: string): string {
  const separator = project.indexOf("_")
  return separator === -1 ? project : project.slice(separator + 1)
}

/** The platform a project runs on, from the start of its configuration; "Other" for an unknown OS. */
export function platformOf(project: string): string {
  const configuration = configurationOf(project)
  if (configuration.startsWith("macOS")) return "macOS"
  if (configuration.startsWith("Windows")) return "Windows"
  if (/^(Ubuntu|Fedora|Debian|Linux|Alpine|Arch)/.test(configuration)) return "Linux"
  return "Other"
}

function platformRank(project: string): number {
  const index = (PLATFORMS as readonly string[]).indexOf(platformOf(project))
  return index === -1 ? PLATFORMS.length : index
}

/** Orders projects by platform (macOS, Linux, Windows, then the rest), then by configuration. */
export function compareConfigurations(a: string, b: string): number {
  return platformRank(a) - platformRank(b) || configurationOf(a).localeCompare(configurationOf(b))
}
