import { defineStore } from "pinia"

export const useMessageStore = defineStore("Message", {
  state() {
    return {isError: false, message: ""}
  },
  actions: {
    showMessage(msg: string) {
      this.isError = true
      this.message = msg
    },
  },
})