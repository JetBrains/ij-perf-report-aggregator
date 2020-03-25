package main

import (
  "github.com/develar/errors"
  "go.uber.org/zap"
  "io/ioutil"
  "os"
  "path/filepath"
  "sort"
)

func (t *BackupManager) removeOldLocalBackups(backupParentDir string, backupsToKeepLocal int) ([]string, error) {
  unfilteredFiles, err := ioutil.ReadDir(backupParentDir)
  if err != nil {
    if os.IsNotExist(err) {
      return nil, nil
    }
    return nil, errors.WithStack(err)
  }

  sort.SliceStable(unfilteredFiles, func(i, j int) bool {
    return unfilteredFiles[i].ModTime().Before(unfilteredFiles[j].ModTime())
  })

  localBackups := make([]string, 0, len(unfilteredFiles))
  for _, f := range unfilteredFiles {
    if !f.IsDir() {
      continue
    }

    localBackups = append(localBackups, f.Name())
  }

  end := len(localBackups) - backupsToKeepLocal
  if end <= 0 {
    return localBackups, nil
  }

  for _, name := range localBackups[0:end] {
    t.Logger.Info("remove old local backup", zap.String("backup", name))
    err = os.RemoveAll(filepath.Join(backupParentDir, name))
    if err != nil {
      return nil, errors.WithStack(err)
    }
  }

  return localBackups[end:], nil
}
