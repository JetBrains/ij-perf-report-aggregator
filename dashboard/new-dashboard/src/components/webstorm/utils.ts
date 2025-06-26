export function groupBy3<T>(array: T[]): T[][] {
  const result = []
  for (let i = 0; i < array.length; i += 3) {
    const component = [array[i]]
    if (i + 1 < array.length) component.push(array[i + 1])
    if (i + 2 < array.length) component.push(array[i + 2])
    result.push(component)
  }
  return result
}
