package main

import (
  "embed"
  "github.com/develar/errors"
  "github.com/zeebo/xxh3"
  "io/fs"
  "log"
  "net/http"
  "path/filepath"
  "strconv"
  "strings"
)

var (
  //go:embed resources
  assetFs embed.FS
)

func main() {
  err := run()
  if err != nil {
    log.Fatal(err)
  }
}

type assetInfo struct {
  data            []byte
  dataBr          []byte
  eTag            string
  contentLength   string
  contentLengthBr string

  contentType string
}

func run() error {
  var pathToAsset = map[string]*assetInfo{}

  err := fs.WalkDir(assetFs, "resources", func(path string, d fs.DirEntry, err error) error {
    if err != nil {
      return err
    }

    if !strings.ContainsRune(path, '.') || strings.HasPrefix(path, "/.") {
      // is a directory
      return nil
    }

    data, err := assetFs.ReadFile(path)
    if err != nil {
      return err
    }

    if path == "resources" {
      return nil
    }

    key := strings.TrimPrefix(path, "resources")
    isBr := strings.HasSuffix(key, ".br")
    if isBr {
      key = strings.TrimSuffix(key, ".br")
    }

    info := pathToAsset[key]
    if info == nil {
      info = &assetInfo{}
      pathToAsset[key] = info
    }

    lengthAsString := strconv.Itoa(len(data))
    if isBr {
      info.dataBr = data
      info.contentLengthBr = lengthAsString
    } else {
      info.data = data
      info.contentLength = lengthAsString

      // no mime type on distroless image
      if strings.HasSuffix(key, ".woff2") {
        info.contentType = "font/woff2"
      } else if strings.HasSuffix(key, ".svg") {
        info.contentType = "image/svg+xml"
      } else if strings.HasSuffix(key, ".js") {
        info.contentType = "text/javascript"
      } else if strings.HasSuffix(key, ".json") {
        info.contentType = "application/json"
      } else if strings.HasSuffix(key, ".html") {
        info.contentType = "text/html"
      } else if strings.HasSuffix(key, ".css") {
        info.contentType = "text/css"
      } else if strings.HasSuffix(key, ".wasm") {
        info.contentType = "application/wasm"
      } else if strings.HasSuffix(key, ".dictionary") {
        info.contentType = "application/octet-stream"
      } else {
        return errors.New("cannot determinate content-type by file extension: " + filepath.Ext(path))
      }

      // assets are immutable
      if !strings.HasPrefix(key, "/assets/") {
        hash := xxh3.Hash128(data)
        info.eTag = strconv.FormatUint(hash.Hi, 36) + "-" + strconv.FormatUint(hash.Lo, 36)
      }
    }

    return nil
  })
  if err != nil {
    return err
  }

  // no need to keep it anymore
  assetFs = embed.FS{}

  indexHtml := pathToAsset["/index.html"]
  pathToAsset["/"] = indexHtml

  http.Handle("/index.html", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
    newPath := "./"
    q := request.URL.RawQuery
    if q != "" {
      newPath += "?" + q
    }
    writer.Header().Set("Location", newPath)
    writer.WriteHeader(http.StatusMovedPermanently)
  }))
  http.Handle("/", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
    header := writer.Header()

    path := request.URL.Path
    asset := pathToAsset[path]
    if asset == nil {
      if strings.ContainsRune(path, '.') {
        http.NotFound(writer, request)
        return
      }

      // vue router HTML5 mode (https://next.router.vuejs.org/guide/essentials/history-mode.html#html5-mode)
      asset = indexHtml
    }

    header.Set("Vary", "Accept-Encoding")
    header.Set("Content-Type", asset.contentType)

    // https://medium.com/adobetech/an-http-caching-strategy-for-static-assets-configuring-the-server-1192452ce06a
    if len(asset.eTag) == 0 {
      header.Set("Cache-Control", "public,max-age=31536000,immutable")
    } else {
      header.Set("Cache-Control", "no-cache")

      if request.Header.Get("If-None-Match") == asset.eTag {
        writer.WriteHeader(http.StatusNotModified)
        return
      }

      header.Set("ETag", asset.eTag)
    }

    if len(asset.dataBr) != 0 && strings.Contains(request.Header.Get("Accept-Encoding"), "br") {
      header.Set("Content-Length", asset.contentLengthBr)
      header.Set("Content-Encoding", "br")
      _, _ = writer.Write(asset.dataBr)
    } else {
      header.Set("Content-Length", asset.contentLength)
      _, _ = writer.Write(asset.data)
    }
  }))
  return http.ListenAndServe(":8080", nil)
}
