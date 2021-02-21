export function debounceSync(func: () => void, wait: number = 100): () => void {
  let timeout = 0
  let isScheduled = false
  let timestamp = 0
  function later() {
    const last = Date.now() - timestamp
    if (last < wait && last >= 0) {
      timeout = window.setTimeout(later, wait - last)
    }
    else {
      isScheduled = false
      func()
    }
  }

  const debounced = function () {
    timestamp = Date.now()
    if (!isScheduled) {
      isScheduled = true
      timeout = window.setTimeout(later, wait)
    }
  }
  debounced.clear = () => {
    if (isScheduled) {
      isScheduled = false
      clearTimeout(timeout)
    }
  }
  return debounced
}

export interface TaskHandle {
  readonly isCancelled: boolean

  onCancel(listener: () => void): void
}

class TaskHandleImpl implements TaskHandle {
  private readonly cancelListeners: Array<() => void> = []

  timeoutHandle = 0
  timestamp = 0

  isCancelled = false

  onCancel(listener: () => void): void {
    this.cancelListeners.push(listener)
  }

  cancel(): void {
    this.isCancelled = true
    for (const cancelListener of this.cancelListeners) {
      try {
        cancelListener()
      }
      catch (e) {
        console.error("cannot execute cancel listener", e)
      }
    }
    this.cancelListeners.length = 0
  }
}

export class DebouncedTask {
  private readonly laterFunctionReference = (): void => this.later()
  readonly executeFunctionReference = (): void => this.execute()
  private taskHandle: TaskHandleImpl | null = null

  constructor(private readonly task: (taskHandle: TaskHandle) => Promise<unknown>, private readonly wait: number = 100) {
  }

  private later(): void {
    const taskHandle = this.taskHandle
    if (taskHandle == null) {
      console.error("taskHandle must be not null")
      return
    }

    const last = Date.now() - taskHandle.timestamp
    if (last < this.wait && last >= 0) {
      taskHandle.timeoutHandle = window.setTimeout(this.laterFunctionReference, this.wait - last)
      return
    }

    const done = () => {
      if (taskHandle != null) {
        taskHandle.cancel()
        this.taskHandle = null
      }
    }
    this.task(taskHandle)
      .then(done)
      .catch(error => {
        done()
        throw error
      })
  }

  /**
   * Do not pass function by reference - use executeFunctionReference instead.
   * @param immediately
   */
  execute(immediately: boolean = false): void {
    if (immediately) {
      const taskHandle = this.taskHandle
      if (taskHandle !== null) {
        const timeoutHandle = taskHandle.timeoutHandle
        if (timeoutHandle !== 0) {
          taskHandle.timeoutHandle = 0
          clearTimeout(timeoutHandle)
        }
        taskHandle.cancel()
      }

      this.taskHandle = new TaskHandleImpl()
      this.later()
    }
    else {
      if (this.taskHandle === null) {
        this.taskHandle = new TaskHandleImpl()
        this.taskHandle.timeoutHandle = window.setTimeout(this.laterFunctionReference, this.wait)
      }
      this.taskHandle.timestamp = Date.now()
    }
  }
}