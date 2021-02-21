package main

import (
  "log"
  "net/http"
  "os"
  "path"
  "strings"
)

func main() {
  err := run()
  if err != nil {
    log.Fatal(err)
  }
}

func run() error {
  publicDir := http.Dir(os.Getenv("KO_DATA_PATH"))
  indexHtml := path.Join(string(publicDir), "/index.html")
  fileServer := http.FileServer(publicDir)
  return http.ListenAndServe(":8081", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
    urlPath := request.URL.Path
    if strings.ContainsRune(urlPath, '.') || urlPath == "/" {
      fileServer.ServeHTTP(writer, request)
    } else {
      // vue router HTML5 mode (https://next.router.vuejs.org/guide/essentials/history-mode.html#html5-mode)
      http.ServeFile(writer, request, indexHtml)
    }
  }))
}
