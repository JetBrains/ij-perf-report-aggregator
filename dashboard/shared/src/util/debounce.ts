interface DebouncedFunction {
  (): void

  clear(): void
}

export function debounceSync(func: () => void, wait: number = 100): DebouncedFunction {
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

  const debounced: DebouncedFunction = function () {
    timestamp = Date.now()
    if (!isScheduled) {
      isScheduled = true
      timeout = window.setTimeout(later, wait)
    }
  }
  debounced.clear = function () {
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

  timestamp = 0
  timeoutHandle = 0

  constructor(private readonly task: (taskHandle: TaskHandle) => Promise<unknown>, private readonly wait: number = 100) {
  }

  private later(): void {
    const last = Date.now() - this.timestamp
    if (last < this.wait && last >= 0) {
      this.timeoutHandle = window.setTimeout(this.laterFunctionReference, this.wait - last)
      return
    }

    this.timeoutHandle = 0

    let taskHandle = this.taskHandle
    if (taskHandle != null) {
      console.error("taskHandle must be null")
    }

    taskHandle = new TaskHandleImpl()
    this.taskHandle = taskHandle
    this.task(taskHandle)
      .then(() => {
        this.taskHandle = null
      })
      .catch(error => {
        const taskHandle = this.taskHandle
        if (taskHandle != null) {
          this.taskHandle = null
          taskHandle.cancel()
        }
        throw error
      })
  }

  /**
   * Do not pass function by reference - use executeFunctionReference instead.
   * @param immediately
   */
  execute(immediately: boolean = false): void {
    if (immediately) {
      if (this.timeoutHandle !== 0) {
        clearTimeout(this.timeoutHandle)
        this.timeoutHandle = 0
      }

      if (this.taskHandle !== null) {
        this.taskHandle.cancel()
        this.taskHandle = null
      }

      this.timestamp = 0
      this.later()
    }
    else {
      if (this.taskHandle != null) {
        this.taskHandle.cancel()
        this.taskHandle = null
      }

      if (this.timeoutHandle === 0) {
        this.timeoutHandle = window.setTimeout(this.laterFunctionReference, this.wait)
      }
      this.timestamp = Date.now()
    }
  }
}