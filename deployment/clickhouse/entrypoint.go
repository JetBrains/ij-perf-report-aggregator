package main

import (
  "bytes"
  _ "embed"
  "errors"
  "github.com/nats-io/nats.go"
  "github.com/valyala/fastjson"
  "log"
  "os"
  "os/exec"
  "path/filepath"
  "strings"
  "syscall"
  "time"
)

//go:embed config.xml
var clickhouseConfig []byte

func main() {
  clickhouseExecutable := "/usr/bin/clickhouse"

  isLocalRun := len(os.Getenv("KUBERNETES_SERVICE_HOST")) == 0
  if isLocalRun {
    clickhouseExecutable = "/usr/local/bin/clickhouse"
    setS3EnvForLocalRun()
  }

  bucket := getEnvOrFile("S3_BUCKET", "/etc/s3/bucket")
  s3AccessKey := getEnvOrFile("S3_ACCESS_KEY", "/etc/s3/accessKey")
  s3SecretKey := getEnvOrFile("S3_SECRET_KEY", "/etc/s3/secretKey")

  restoreData := os.Getenv("RESTORE_DB") == "true"

  configFile := "/var/lib/clickhouse/config.xml"
  if isLocalRun {
    configFile = "/Volumes/data/Documents/report-aggregator/deployment/ch-local-config.xml"
  }

  if restoreData {
    err := prepareConfigAndDir(isLocalRun, bucket, s3AccessKey, s3SecretKey, configFile)
    if err != nil {
      log.Fatal(err)
    }
  }

  cmd := exec.Command(clickhouseExecutable, "server", "--config-file="+configFile)
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr

  err := cmd.Start()
  if err != nil {
    log.Fatal(err)
  }

  if restoreData {
    defer func() {
      err = cmd.Process.Signal(syscall.SIGTERM)
      if err != nil {
        log.Println(err)
        _ = cmd.Process.Kill()
      }

      err = cmd.Wait()
      if err != nil {
        log.Println(err)
      }
    }()

    chBackupExecutable := "/usr/bin/clickhouse-backup"
    if isLocalRun {
      chBackupExecutable = "/usr/local/bin/clickhouse-backup"
    }
    err = restoreDb(chBackupExecutable, s3AccessKey, s3SecretKey, bucket)
    if err != nil {
      log.Fatal(err)
    }
    return
  }

  go func() {
    // wait for clickhouse server start
    time.Sleep(10 * time.Second)
    requestClearCache()
  }()

  err = cmd.Wait()
  if err != nil {
    log.Fatal(err)
  }
}

func prepareConfigAndDir(isLocalRun bool, bucket string, s3AccessKey string, s3SecretKey string, configFile string) error {
  chDir := "/var/lib/clickhouse"
  if isLocalRun {
    chDir = "/Volumes/data/ij-perf-db/clickhouse"
  }

  entries, err := os.ReadDir(chDir)
  if err != nil && !os.IsNotExist(err) {
    return err
  }

  if entries != nil {
    for _, entry := range entries {
      err = os.RemoveAll(filepath.Join(chDir, entry.Name()))
      if err != nil && !os.IsNotExist(err) {
        return err
      }
    }
  }

  if !isLocalRun {
    s3Url := "https://" + bucket + ".s3.eu-west-1.amazonaws.com/data/"
    log.Print("S3 URL: " + s3Url)

    s := strings.NewReplacer(
      "$S3_URL", s3Url,
      "$S3_ACCESS_KEY", s3AccessKey,
      "$S3_SECRET_KEY", s3SecretKey,
    ).Replace(string(clickhouseConfig))

    // /etc is not writeable
    err := os.WriteFile(configFile, []byte(s), 0666)
    if err != nil {
      return err
    }
  }
  return nil
}

func setS3EnvForLocalRun() {
  cmd := exec.Command("doppler", "secrets", "download", "--project", "s3", "--config", "prd", "--no-file")
  stdout, err := cmd.Output()
  if err != nil {
    log.Println("failed to use doppler to retrieve credentials", err)
  } else {
    excludePrefix := []byte("DOPPLER_")
    fastjson.MustParseBytes(stdout).GetObject().Visit(func(key []byte, v *fastjson.Value) {
      if !bytes.HasPrefix(key, excludePrefix) {
        err = os.Setenv(string(key), string(v.GetStringBytes()))
        if err != nil {
          log.Println("cannot set env", err)
        }
      }
    })
  }
}

func restoreDb(chBackupExecutable string, s3AccessKey string, s3SecretKey string, bucket string) error {
  // wait a little bit for clickhouse start
  time.Sleep(1 * time.Second)

  var backupName string
  attemptCount := 3
  for i := 0; i < attemptCount; i++ {
    // just for debug - print all backups
    cmd := exec.Command(chBackupExecutable, "list", "remote")
    configureBackupToolEnv(cmd, s3AccessKey, s3SecretKey, bucket, false)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    err := cmd.Start()
    if err != nil {
      return err
    }
    err = cmd.Wait()

    var result []byte
    if err == nil {
      cmd = exec.Command(chBackupExecutable, "list", "remote", "latest")
      configureBackupToolEnv(cmd, s3AccessKey, s3SecretKey, bucket, true)

      cmd.Stderr = os.Stderr
      result, err = cmd.Output()
    }

    if err != nil {
      if i < attemptCount {
        time.Sleep(time.Duration((i+1)*2) * time.Second)
        continue
      } else {
        log.Println("cannot get latest backup name")
        return err
      }
    }

    backupName = strings.TrimSpace(string(result))
    break
  }

  if len(backupName) == 0 {
    return errors.New("no remote backup")
  }

  cmd := exec.Command(chBackupExecutable, "restore_remote", "--drop=false", backupName)
  configureBackupToolEnv(cmd, s3AccessKey, s3SecretKey, bucket, false)
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr
  err := cmd.Run()
  if err != nil {
    return err
  }

  log.Println("DB is restored (backup=" + backupName + ")")
  return nil
}

func configureBackupToolEnv(cmd *exec.Cmd, s3AccessKey string, s3SecretKey string, bucket string, logOnlyErrors bool) {
  cmd.Env = []string{}
  for _, s := range os.Environ() {
    if !strings.HasPrefix(s, "S3_") && !strings.HasPrefix(s, "CLICKHOUSE_") {
      cmd.Env = append(cmd.Env, s)
    }
  }
  if logOnlyErrors {
    cmd.Env = append(cmd.Env, "LOG_LEVEL=error")
  }
  cmd.Env = append(cmd.Env, "S3_ALLOW_MULTIPART_DOWNLOAD=true")
  cmd.Env = append(cmd.Env, "REMOTE_STORAGE=s3")
  cmd.Env = append(cmd.Env, "S3_ACCESS_KEY="+s3AccessKey)
  cmd.Env = append(cmd.Env, "S3_SECRET_KEY="+s3SecretKey)
  cmd.Env = append(cmd.Env, "S3_BUCKET="+bucket)
  cmd.Env = append(cmd.Env, "S3_REGION=eu-west-1")
}

func requestClearCache() {
  url := os.Getenv("NATS")
  if len(url) == 0 {
    url = "nats://nats:4222"
  }

  nc, err := nats.Connect(url, nats.Name("NATS Sample Publisher"))
  if err != nil {
    log.Fatal(err)
  }
  defer nc.Close()

  err = nc.Publish("server.clearCache", []byte("clickhouse"))
  if err != nil {
    log.Fatal(err)
  }

  err = nc.Flush()
  if err != nil {
    log.Fatal(err)
  }
}

func getEnvOrFile(envName string, file string) string {
  v := os.Getenv(envName)
  if len(v) == 0 {
    b, err := os.ReadFile(file)
    if err != nil {
      log.Fatal(err)
    }
    return string(b)
  }
  return v
}
