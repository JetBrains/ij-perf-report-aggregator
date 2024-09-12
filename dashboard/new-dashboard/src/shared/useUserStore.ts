import { defineStore } from "pinia"
import { ref } from "vue"
import { ServerWithCompressConfigurator } from "../configurators/ServerWithCompressConfigurator"

export interface User {
  email: string
  family_name: string
  given_name: string
  hd: string
  id: string
  locale: string
  name: string
  picture: string
  verified_email: boolean
}

export const useUserStore = defineStore("userStore", () => {
  // State
  const user = ref<User | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  // Fetch User Data from API
  const fetchUser = async (): Promise<User> => {
    const response = await fetch(`${ServerWithCompressConfigurator.DEFAULT_SERVER_URL}/api/auth/userinfo`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
    })

    if (!response.ok) {
      throw new Error("Failed to fetch user data")
    }

    return (await response.json()) as User
  }

  // Initialize user data
  const initializeUser = () => {
    loading.value = true
    error.value = null

    fetchUser()
      .then((data) => {
        user.value = data
      })
      .catch((error_: unknown) => {
        error.value = error_ instanceof Error ? error_.message : "Unknown error occurred";
      })
      .finally(() => {
        loading.value = false
      })
  }

  // Fetch user data when the store is initialized
  initializeUser()

  // Return the state and actions from the store
  return { user, loading, error }
})
