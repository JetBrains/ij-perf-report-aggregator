import { ElMessage } from "element-plus"
import { Ref } from "vue"
import { IMessageHandle } from "element-plus/es/el-message/src/types"
import { TaskHandle } from "./debounce"

const serverNotAvailableErrorMessage = "Server is not available. Please check that server is running and VPN connection is established."

let errorMessageHandle: IMessageHandle | null = null

export function loadJson<T>(url: string,
                            loading: Ref<boolean> | null,
                            taskHandle: TaskHandle,
                            dataConsumer: (data: T) => void): Promise<unknown> {
  if (taskHandle.isCancelled) {
    return Promise.resolve()
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

  return fetch(url, {credentials: "omit", signal: controller.signal})
    .then(response => {
      if (taskHandle.isCancelled) {
        clearTimeout(timeoutId)
        return Promise.resolve()
      }

      if (response.ok) {
        if (errorMessageHandle != null) {
          errorMessageHandle.close()
          errorMessageHandle = null
        }
        return response.json()
      }
      else {
        return Promise.reject(new Error(`cannot load data (url=${url}, status=${response.status}`))
      }
    })
    .then(data => {
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
      return null
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
          showError(`Cannot load data from ${url}: ${e.message}`)
        }
      }
      return null
    })
}

function showError(message: string) {
  if (errorMessageHandle != null) {
    return
  }

  errorMessageHandle = ElMessage({
    message,
    type: "error",
    duration: 0,
    showClose: true,
    onClose: () => {
      errorMessageHandle = null
    }
  })
}