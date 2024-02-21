import { Mock, SpyInstance, vitest } from "vitest"

type Test = () => boolean
export const awaitCallbackTrue = async (callback: Test, timeoutMs: number = 1000) => {
  await vitest.waitFor(
    () => {
      if (!callback()) {
        throw new Error(`The callback did not return the true value during the timeout ${timeoutMs}`)
      }
    },
    {
      timeout: timeoutMs,
    }
  )
}

export const awaitMockCallsCount = async (mock: Mock | SpyInstance, count: number) => {
  await awaitCallbackTrue(() => {
    return mock.mock.calls.length >= count
  })
}
