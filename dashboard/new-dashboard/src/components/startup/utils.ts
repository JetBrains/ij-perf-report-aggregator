import { ServerConfigurator } from "shared/src/configurators/ServerConfigurator"
import { Ref } from "vue"

export function fetchHighlightingPasses(highlightingPasses: Ref<Array<string>|undefined>){
  fetch(ServerConfigurator.DEFAULT_SERVER_URL + "/api/highlightingPasses")
    .then(response => response.json())
    .then((data: Array<string>) => {
      highlightingPasses.value = data.map(it => "metrics."  + it)
    })
    .catch(error => console.error(error))
}