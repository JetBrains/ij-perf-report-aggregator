
export function toColor(str: string): string {
  let i
  let hash = 0
  for (i = 0; i < str.length; i++) {
    if (!is_numeric(str.charAt(i))) {
      const charCode = str.codePointAt(i)
      if (charCode == undefined) {
        return "black"
      }
      hash = charCode + ((hash << 5) - hash)
    }
  }
  let colour = "#"
  for (i = 0; i < 3; i++) {
    const value = (hash >> (i * 8)) & 0xFF
    colour += ("00" + value.toString(16)).slice(-2)
  }
  return colour
}

function is_numeric(str: string): boolean {
  return str >= "0" && str <= "9"
}