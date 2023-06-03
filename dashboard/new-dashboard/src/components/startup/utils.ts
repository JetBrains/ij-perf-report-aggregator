import { Ref } from "vue"
import { ServerConfigurator } from "../../configurators/ServerConfigurator"

export function fetchHighlightingPasses(highlightingPasses: Ref<string[]|undefined>){
  fetch(ServerConfigurator.DEFAULT_SERVER_URL + "/api/highlightingPasses")
    .then(response => response.json())
    .then((data: string[]) => {
      highlightingPasses.value = data.map(it => "metrics."  + it)
    })
    .catch(error => console.error(error))
}