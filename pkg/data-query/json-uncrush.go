package data_query

import (
  "strings"
)

var replacers = []*strings.Replacer{
  strings.NewReplacer(
    "\"", "'",
    "'", "\"",
  ),
  strings.NewReplacer(
    "':", "!",
    "!", "':",
  ),
  strings.NewReplacer(
    ",'", "~",
    "~", ",'",
  ),
  strings.NewReplacer(
    "}", ")",
    ")", "}",
  ),
  strings.NewReplacer(
    "{", "(",
    "(", "{",
  ),
}

func swap(s string) string {
  for _, replacer := range replacers {
    s = replacer.Replace(s)
  }
  return s
}

func unswap(s string) string {
  for i := len(replacers) - 1; i >= 0; i-- {
    s = replacers[i].Replace(s)
  }
  return s
}

func uncrush(s string) string {
  //s = unswap(s)

  // unsplit the string using the delimiter
  stringParts := strings.Split(s, "\u0001")
  result := stringParts[0]
  if len(stringParts) > 1 {
    splitString := stringParts[1]
    for _, character := range splitString {
      // split the string using the current splitCharacter
      c := string(character)
      splitArray := strings.Split(result, c)
      // rejoin the string with the last element from the split
      lastIndex := len(splitArray) - 1
      result = strings.Join(splitArray[0:lastIndex], splitArray[lastIndex])
    }
  }
  return result
}
