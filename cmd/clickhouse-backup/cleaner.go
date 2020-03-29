package main

import (
  "github.com/develar/errors"
  "go.uber.org/zap"
  "io/ioutil"
  "os"
  "path/filepath"
  "sort"
  "time"
)

type LocalBackupFile struct {
  name     string
  modified time.Time
}

func (t *BackupManager) removeOldLocalBackups(backupParentDir string, maxNumberOfFiles int) (string, error) {
  unfilteredFiles, err := ioutil.ReadDir(backupParentDir)
  if err != nil {
    if os.IsNotExist(err) {
      return "", nil
    }
    return "", errors.WithStack(err)
  }

  localBackups := make([]LocalBackupFile, 0, len(unfilteredFiles))
  for _, f := range unfilteredFiles {
    if !f.IsDir() {
      continue
    }

    e := LocalBackupFile{
      name: f.Name(),
    }

    // "In the absence of a time zone indicator, Parse returns a time in UTC."
    modTimeFromName, err := time.Parse(timeFormat, e.name)
    if err != nil {
      t.Logger.Warn("cannot infer modification time from directory name", zap.Error(err))
      e.modified = f.ModTime()
    } else {
      e.modified = modTimeFromName
    }

    localBackups = append(localBackups, e)
  }

  if len(localBackups) == 0 {
    return "", nil
  }

    sort.SliceStable(localBackups, func(i, j int) bool {
    return localBackups[i].modified.Before(localBackups[j].modified)
  })

  end := len(localBackups) - maxNumberOfFiles
  if end <= 0 {
    // incremental backup only if limit is not exceed - if exceeded it means that full backup should be created
    return filepath.Join(t.backupParentDir, localBackups[len(localBackups)-1].name), nil
  }

  for _, item := range localBackups[0:end] {
    t.Logger.Info("remove old local backup", zap.String("backup", item.name))
    err = os.RemoveAll(filepath.Join(backupParentDir, item.name))
    if err != nil {
      return "", errors.WithStack(err)
    }
  }

  return "", nil
}
