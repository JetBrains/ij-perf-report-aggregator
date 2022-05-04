package main

import (
  "embed"
  "github.com/zeebo/xxh3"
  "io/fs"
  "io/ioutil"
  "log"
  "net/http"
  "strconv"
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
  f, _ := res.Open("resources/index.html")
  indexHtml, _ := ioutil.ReadAll(f)
  indexHtmlEtag := strconv.FormatUint(xxh3.Hash(indexHtml), 36)
  indexHtmlLength := strconv.Itoa(len(indexHtml))
  return http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
    urlPath := request.URL.Path
    if strings.ContainsRune(urlPath, '.') || urlPath == "/" {
      fileServer.ServeHTTP(w, request)
    } else {
      // vue router HTML5 mode (https://next.router.vuejs.org/guide/essentials/history-mode.html#html5-mode)

      header := w.Header()
      header.Set("Content-Type", "text/html")
      header.Set("Content-Length", indexHtmlLength)
      header.Set("Vary", "Accept-Encoding")
      header.Set("Cache-Control", "public, must-revalidate, max-age=30")
      header.Set("ETag", indexHtmlEtag)

      prevEtag := request.Header.Get("If-None-Match")
      if prevEtag == indexHtmlEtag {
        w.WriteHeader(http.StatusNotModified)
        return
      }

      w.WriteHeader(http.StatusOK)
      _, _ = w.Write(indexHtml)
    }
  }))
}
