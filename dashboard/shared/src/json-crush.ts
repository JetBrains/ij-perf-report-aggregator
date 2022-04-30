/////////////////////////////////////////////////////////////////////////////////
// JSONCrush v1.1.6 by Frank Force - https://github.com/KilledByAPixel/JSONCrush
/////////////////////////////////////////////////////////////////////////////////

"use strict"

const maxSubstringLength = 50

function byteLength(string: string) {
  return encodeURI(encodeURIComponent(string)).replace(/%../g, "i").length
}

function hasUnmatchedSurrogate(string: string) {
  // check ends of string for unmatched surrogate pairs
  const c1 = string.charCodeAt(0)
  const c2 = string.charCodeAt(string.length - 1)
  return (c1 >= 0xDC00 && c1 <= 0xDFFF) || (c2 >= 0xD800 && c2 <= 0xDBFF)
}

const unescapedCharacters = "-_.!~*'()"

// create a string of replacement characters
const replaceCharacters: Array<string> = []

// prefer replacing with characters that will not be escaped by encodeURIComponent
for (let i = 127; --i;) {
  if (
    (i >= 48 && i <= 57) || // 0-9
    (i >= 65 && i <= 90) || // A-Z
    (i >= 97 && i <= 122) || // a-z
    unescapedCharacters.includes(String.fromCharCode(i))
  ) {
    replaceCharacters.push(String.fromCharCode(i))
  }
}

// pick from extended set last
for (let i = 32; i < 255; ++i) {
  const c = String.fromCharCode(i)
  if (c !== "\\" && !replaceCharacters.includes(c)) {
    replaceCharacters.unshift(c)
  }
}

export function crush(s: string): string {
  // s = swap(s)

  // insert delimiter between JSCrush parts
  // JSCrush Algorithm (replace repeated substrings with single characters)
  let replaceCharacterPos = replaceCharacters.length
  let splitString = ""

  // count instances of substrings
  let substringCount = new Map<string, number>()
  for (let substringLength = 2; substringLength < maxSubstringLength; substringLength++) {
    for (let i = 0; i < s.length - substringLength; ++i) {
      const substring = s.substring(i, i + substringLength)

      // don't recount if already in list
      if (substringCount.has(substring)) {
        continue
      }

      // prevent breaking up unmatched surrogates
      if (hasUnmatchedSurrogate(substring)) {
        continue
      }

      // count how many times the substring appears
      let count = 1
      for (let substringPos = s.indexOf(substring, i + substringLength); substringPos >= 0; ++count) {
        substringPos = s.indexOf(substring, substringPos + substringLength)
      }

      // add to list if it appears multiple times
      if (count > 1) {
        substringCount.set(substring, count)
      }
    }
  }

  // loop while string can be crushed more
  // eslint-disable-next-line no-constant-condition
  while (true) {
    // get the next character that is not in the string
    // eslint-disable-next-line no-empty
    for (; replaceCharacterPos-- && s.includes(replaceCharacters[replaceCharacterPos]);) {
    }
    if (replaceCharacterPos < 0) {
      // ran out of replacement characters
      break
    }
    const replaceCharacter = replaceCharacters[replaceCharacterPos]

    // find the longest substring to replace
    let bestSubstring = ""
    let bestLengthDelta = 0
    const replaceByteLength = byteLength(replaceCharacter)
    // https://stackoverflow.com/questions/35940216/es6-is-it-dangerous-to-delete-elements-from-set-map-during-set-map-iteration
    for (const [substring, count] of substringCount) {
      // calculate change in length of string if it substring was replaced
      let lengthDelta = (count - 1) * byteLength(substring) - (count + 1) * replaceByteLength
      if (splitString.length === 0) {
        // include the delimiter length
        lengthDelta -= 3
      }
      if (lengthDelta <= 0) {
        substringCount.delete(substring)
      }
      else if (lengthDelta > bestLengthDelta) {
        bestSubstring = substring
        bestLengthDelta = lengthDelta
      }
    }
    if (bestSubstring.length === 0) {
      // string can't be compressed further
      break
    }

    // create new string with the split character
    s = s.split(bestSubstring).join(replaceCharacter) + replaceCharacter + bestSubstring
    splitString = replaceCharacter + splitString

    // update substring count list after the replacement
    const newSubstringCount = new Map<string, number>()
    for (const substring of substringCount.keys()) {
      // make a new substring with the replacement
      const newSubstring = substring.split(bestSubstring).join(replaceCharacter)

      // count how many times the new substring appears
      let count = 0
      for (let i = s.indexOf(newSubstring); i >= 0; ++count) {
        i = s.indexOf(newSubstring, i + newSubstring.length)
      }

      // add to list if it appears multiple times
      if (count > 1) {
        newSubstringCount.set(newSubstring, count)
      }
    }
    substringCount = newSubstringCount
  }

  s = escapePathSegment(s)
  if (splitString.length !== 0) {
    s += "%01" + escapePathSegment(splitString)
  }
  return s
}

// swap out characters for lesser used ones that wont get escaped
const swapGroups = [
  ["\"", "'"],
  ["':", "!"],
  [",'", "~"],
  ["}", ")", "\\", "\\"],
  ["{", "(", "\\", "\\"],
]

function swapInternal(string: string, g: Array<string>) {
  const regex = new RegExp(`${(g[2] ?? "")+g[0]}|${(g[3] ?? "")+g[1]}`, "g")
  return string.replace(regex, $1 => ($1 === g[0] ? g[1] : g[0]))
}

function swap(s: string) {
  for (const swapGroup of swapGroups) {
    s = swapInternal(s, swapGroup)
  }
  return s
}

export function escapePathSegment(x: string): string {
  const s = encodeURIComponent(x)
  if (s.length === x.length) {
    // nothing was escaped
    return s
  }

  return s
    .replace(/%2C/g, ",")
    .replace(/%3A/g, ":")
    .replace(/%40/g, "@")
    .replace(/%24/g, "$")
}