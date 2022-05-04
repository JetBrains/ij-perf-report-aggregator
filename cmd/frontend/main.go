package main

import (
  "embed"
  "io/fs"
  "log"
  "net/http"
  "strings"
)

var (
  //go:embed resources
  res embed.FS
)

func main() {
  err := run()
  if err != nil {
    log.Fatal(err)
  }
}

func run() error {
  sub, _ := fs.Sub(res, "resources")
  fileServer := http.FileServer(http.FS(sub))
  return http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
    urlPath := request.URL.Path
    if strings.ContainsRune(urlPath, '.') || urlPath == "/" {
      fileServer.ServeHTTP(w, request)
    } else {
      // vue router HTML5 mode (https://next.router.vuejs.org/guide/essentials/history-mode.html#html5-mode)
      http.ServeFile(w, request, "/index.html")
    }
  }))
}
