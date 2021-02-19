declare module "rison-node" {
  export function encode(obj: unknown): string

  export function decode(data: string): unknown
}