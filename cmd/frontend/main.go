package main

import (
  "log"
  "net/http"
  "os"
)

func main() {
  err := run()
  if err != nil {
    log.Fatal(err)
  }
}

func run() error {
  return http.ListenAndServe(":8080", http.FileServer(http.Dir(os.Getenv("KO_DATA_PATH"))))
}
