import { computed, ComputedRef } from "vue"

export function usePlaceholder(
  props: { label: string },
  model: () => string | Array<unknown> | null,
  selected: () => string | Array<string> | null,
): ComputedRef<string> {
  return computed<string>(() => {
    const values = model()
    // PrimeVue doesn't show value if values are not yet loaded.
    // In our case model is stored on server, but selected value stored locally. So, selected value is resolved much faster.
    if ((values == null || values.length === 0)) {
      const value = selected()
      if (value != null && value.length != 0) {
        return Array.isArray(value) ? value.join(", ") : value
      }
    }
    return props.label
  })
}