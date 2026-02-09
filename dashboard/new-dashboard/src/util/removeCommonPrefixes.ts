function collectCommonPart(paths: string[]): string | null {
  const commonSegments = new Map<string, number>()
  for (const path of paths) {
    const segments = path.split("/").slice(1)
    for (const segment of segments) {
      commonSegments.set(segment, (commonSegments.get(segment) ?? 0) + 1)
    }
  }
  let commonSegment = null
  for (const [key, value] of commonSegments.entries()) {
    if (value == paths.length) {
      commonSegment = key
    }
  }
  return commonSegment
}

export function removeCommonSegments(paths: string[]): string[] {
  if (paths.length < 2) return paths

  const commonSegment = collectCommonPart(paths)
  if (commonSegment !== "" && commonSegment != null) {
    paths = paths.map((path) => path.replaceAll("/" + commonSegment, ""))
    paths = paths.map((path) => {
      path = path.replace("//", "/")
      while (path.startsWith("/")) {
        path = path.slice(1)
      }
      while (path.endsWith("/")) {
        path = path.slice(0, Math.max(0, path.length - 1))
      }

      return path
    })

    return paths
  }
  return paths
}
