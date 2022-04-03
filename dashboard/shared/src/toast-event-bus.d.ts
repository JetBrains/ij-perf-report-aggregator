declare module "primevue/toasteventbus" {
  export function emit(type: string, obj: unknown): void
  export function on(type: string, handler: unknown): void
  export function off(type: string, handler: unknown): void
}