import { ElMessage } from "element-plus"
import { Ref } from "vue"

const serverNotAvailableErrorMessage = "Server is not available. Please check that server is running and VPN connection is established."

let isErrorMessageDisplayed = false

export function loadJson<T>(url: string, loading: Ref<boolean> | null, dataConsumer: (data: T) => void): void {
  // console.debug("load", url)
  if (loading != null) {
    loading.value = true
  }

  function showError(message: string) {
    if (isErrorMessageDisplayed) {
      return
    }

    isErrorMessageDisplayed = true
    ElMessage({
      message,
      type: "error",
      duration: 0,
      showClose: true,
      onClose: () => {
        isErrorMessageDisplayed = false
      }
    })
  }

  let isCancelledByTimeout = false
  const controller = new AbortController()
  const timeoutId = setTimeout(() => {
    isCancelledByTimeout = true
    controller.abort()
    showError(serverNotAvailableErrorMessage)
  }, 8000)

  fetch(url, {credentials: "omit", signal: controller.signal})
    .then(response => {
      if (response.ok) {
        return response.json()
      }
      else {
        return Promise.reject(new Error(`cannot load data (url=${url}, status=${response.status}`))
      }
    })
    .then(data => {
      clearTimeout(timeoutId)

      if (loading != null) {
        loading.value = false
      }

      if (data == null) {
        console.error("empty result", url)
        showError("Server returned empty result")
        return null
      }
      dataConsumer(data)
      return null
    })
    .catch((e) => {
      clearTimeout(timeoutId)

      if (loading != null) {
        loading.value = false
      }

      console.error("cannot load data", url, e)
      if (!isCancelledByTimeout) {
        if (e instanceof TypeError) {
          showError(serverNotAvailableErrorMessage)
        }
        else {
          showError(`Cannot load data from ${url}: ${e.message}`)
        }
      }
    })
}
