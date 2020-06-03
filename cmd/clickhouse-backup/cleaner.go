package main

import (
  "github.com/develar/errors"
  "go.uber.org/zap"
  "io/ioutil"
  "os"
  "sort"
  "time"
)

type LocalBackupFile struct {
  name     string
  modified time.Time
}

func (t *LocalBackupFile) String() string {
  return t.name
}

func (t *BackupManager) collectLocalBackups(backupParentDir string) ([]string, error) {
  unfilteredFiles, err := ioutil.ReadDir(backupParentDir)
  if err != nil {
    if os.IsNotExist(err) {
      return nil, nil
    }
    return nil, errors.WithStack(err)
  }

  localBackups := make([]LocalBackupFile, 0, len(unfilteredFiles))
  for _, f := range unfilteredFiles {
    if !f.IsDir() {
      continue
    }

    entry := LocalBackupFile{
      name: f.Name(),
    }

    // "In the absence of a time zone indicator, Parse returns a time in UTC."
    modTimeFromName, err := time.Parse(timeFormat, entry.name)
    if err != nil {
      t.Logger.Warn("cannot infer modification time from directory name", zap.Error(err))
      entry.modified = f.ModTime()
    } else {
      entry.modified = modTimeFromName
    }

    localBackups = append(localBackups, entry)
  }

  if len(localBackups) == 0 {
    return nil, nil
  }

  sort.SliceStable(localBackups, func(i, j int) bool {
    return localBackups[i].modified.Before(localBackups[j].modified)
  })

  names := make([]string, len(localBackups))
  for i, backup := range localBackups {
    names[i] = backup.name
  }

  return names, nil
}
