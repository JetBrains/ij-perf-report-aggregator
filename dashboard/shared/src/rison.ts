/* eslint-disable @typescript-eslint/no-explicit-any */
// https://github.com/w33ble/rison-node/blob/master/js/rison.js
// Uses CommonJS, AMD or browser globals to create a module.
//  the stringifier is based on
//    http://json.org/json.js as of 2006-04-28 from json.org
//  the parser is based on
//    http://osteele.com/sources/openlaszlo/json


/*
 * we divide the uri-safe glyphs into three sets
 *   <rison> - used by rison                         ' ! : ( ) ,
 *   <reserved> - not common in strings, reserved    * @ $ & ; =
 *
 * we define <identifier> as anything that's not forbidden
 */

/**
 * characters that are illegal inside ids.
 * <rison> and <reserved> classes are illegal in ids.
 *
 */
const notIdChar = " '!:(),*@$"

/**
 * characters that are illegal as the start of an id
 * this is so ids can't look like numbers.
 */
const notIdStart = "-0123456789"
const idOk = new RegExp(`^[^${notIdStart}${notIdChar}][^${notIdChar}]*$`)

/**
 * this is like encodeURIComponent() but quotes fewer characters.
 *
 * encodeURIComponent passes   ~!*()-_.'
 * rison.quote also passes   ,:@$
 */
export function makeUrlSafe(x: string): string {
  // if (/^[-A-Za-z\d~!*()_.',:@$]*$/.test(x)) {
  //   return x
  // }

  return encodeURIComponent(x)
    // .replace(/%2C/g, ",")
    // .replace(/%3A/g, ":")
    // .replace(/%40/g, "@")
    // .replace(/%24/g, "$")
}

function doEncode(value: any) {
  // typeof for array also object
  if (value instanceof Array) {
    return array(value)
  }
  else {
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore
    // eslint-disable-next-line @typescript-eslint/no-unsafe-call,@typescript-eslint/no-unsafe-assignment
    return encoders[typeof value](value) as string
  }
}

function object(x: any): string {
  let result = "("
  let isNonFirst = false
  // eslint-disable-next-line @typescript-eslint/no-unsafe-argument
  for (const [key, value] of Object.entries(x)) {
    if (value === undefined || value === null) {
      continue
    }

    if (isNonFirst) {
      result  += ","
    }
    else {
      isNonFirst = true
    }

    result += key + ":" + doEncode(value)
  }
  return result + ")"
}

function array(x: Array<unknown>): string {
  let result = "!("
  let isNonFirst, i
  const l = x.length
  for (i = 0; i < l; i++) {
    const value = x[i]
    if (value === undefined) {
      continue
    }

    if (isNonFirst) {
      result += ","
    }
    else {
      isNonFirst = true
    }

    if (value === null) {
      return "!n"
    }

    result += doEncode(value)
  }
  return result + ")"
}

const encoders = {
  object,
  "boolean"(x: boolean) {
    return x ? "!t" : "!f"
  },
  number(x: number): string {
    if (isFinite(x)) {
      // strip '+' out of exponent, '-' is ok though
      return x.toString().replace(/\+/, "")
    }
    else {
      return "!n"
    }
  },
  string(x: string): string {
    if (x.length === 0) {
      return "''"
    }
    else if (idOk.test(x)) {
      return x
    }

    x = x
      .replaceAll("!", "!!")
      .replaceAll("'", "!'")
    return `'${x}'`
  },
}

/**
 * rison-encode a javascript structure
 */
export function encodeRison(v: unknown): string {
  // eslint-disable-next-line @typescript-eslint/ban-ts-comment
  // @ts-ignore
  // eslint-disable-next-line @typescript-eslint/no-unsafe-call
  return v instanceof Array ? array(v) : object(v)
}