import { ToastSeverity } from "primevue/api"
import { ToastMessageOptions } from "primevue/toast"
import ToastEventBus from "primevue/toasteventbus"
import { Ref } from "vue"
import { TaskHandle } from "./debounce"

const serverNotAvailableErrorMessage = "Server is not available. Please check that server is running."

let message: ToastMessageOptions | null = null

function removeOldMessage() {
  if (message != null) {
    ToastEventBus.emit("remove", message)
    message = null
  }
}

function showError(detail: string) {
  removeOldMessage()
  message = {severity: ToastSeverity.ERROR, detail}
  ToastEventBus.emit("add", message)
}

export function loadJson<T>(url: string,
                            loading: Ref<boolean> | null,
                            taskHandle: TaskHandle,
                            dataConsumer: (data: T) => void): Promise<T | null> {
  if (taskHandle.isCancelled) {
    return Promise.resolve(null)
  }

  // console.debug("load", url)
  if (loading != null) {
    loading.value = true
  }

  let isCancelled = false
  const controller = new AbortController()
  const timeoutId = setTimeout(() => {
    isCancelled = true
    controller.abort()
    showError(serverNotAvailableErrorMessage)
  }, 8000)

  taskHandle.onCancel(() => {
    isCancelled = true
    controller.abort()
  })

  function loaded() {
    clearTimeout(timeoutId)
    if (loading != null) {
      loading.value = false
    }
  }

  return fetch(url, {signal: controller.signal})
    .then(response => {
      if (taskHandle.isCancelled) {
        clearTimeout(timeoutId)
        return Promise.resolve()
      }

      if (response.ok) {
        removeOldMessage()
        return response.json()
      }
      else {
        return Promise.reject(new Error(`cannot load data (url=${url}, status=${response.status}`))
      }
    })
    .then((data: T) => {
      loaded()
      if (taskHandle.isCancelled) {
        return null
      }

      if (data == null) {
        console.error("empty result", url)
        showError("Server returned empty result")
      }
      else {
        dataConsumer(data)
      }
      return data
    })
    .catch(e => {
      loaded()

      if (taskHandle.isCancelled) {
        return null
      }

      console.error("cannot load data", url, e)
      if (!isCancelled) {
        if (e instanceof TypeError) {
          showError(serverNotAvailableErrorMessage)
        }
        else {
          showError(`Cannot load data from ${url}: ${(e as Error).message}`)
        }
      }
      return null
    })
}