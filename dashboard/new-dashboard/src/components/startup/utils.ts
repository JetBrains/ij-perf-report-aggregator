import { computedAsync } from "@vueuse/core"
import { Ref } from "vue"
import { ServerWithCompressConfigurator } from "../../configurators/ServerWithCompressConfigurator"

export function fetchHighlightingPasses(): Ref<string[] | null> {
  return computedAsync(
    () =>
      fetch(ServerWithCompressConfigurator.DEFAULT_SERVER_URL + "/api/highlightingPasses")
        .then((response) => response.json())
        .then((data: string[]) => {
          return data.map((it) => "metrics." + it)
        })
        .catch((error) => {
          console.error(error)
          return null
        }),
    null
  )
}
