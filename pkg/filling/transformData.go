package filling

import (
  "database/sql"
  "github.com/alecthomas/kingpin"
  "github.com/develar/errors"
  "github.com/json-iterator/go"
  "github.com/magiconair/properties"
  "go.uber.org/zap"
  "io"
  "os"
  "path/filepath"
  "report-aggregator/pkg/teamcity"
  "report-aggregator/pkg/util"
  "sort"
  "strconv"
  "strings"
)

func ConfigureTransformCommand(app *kingpin.Application, logger *zap.Logger) {
  command := app.Command("transform", "")
  dbPath := command.Flag("db", "").Required().String()
  command.Action(func(context *kingpin.ParseContext) error {
    return transform(*dbPath, logger)
  })
}

func copyDatabaseFile(filePath string, logger *zap.Logger) (string, error) {
  dirPath := filepath.Dir(filePath)
  fileBaseName := filepath.Base(filePath)
  newFilePath := filepath.Join(dirPath, strings.TrimSuffix(fileBaseName, filepath.Ext(fileBaseName))+"-new.sqlite")
  newFile, err := os.Create(newFilePath)
  if err != nil {
    return "", errors.WithStack(err)
  }

  defer util.Close(newFile, logger)

  file, err := os.Open(filePath)
  if err != nil {
    return "", errors.WithStack(err)
  }

  defer util.Close(file, logger)

  _, err = io.Copy(newFile, file)
  if err != nil {
    return "", errors.WithStack(err)
  }

  return newFilePath, nil
}

func transform(dbPath string, logger *zap.Logger) error {
  readDb, err := sql.Open("sqlite3", "file:"+dbPath+"?mode=ro")
  if err != nil {
    return errors.WithStack(err)
  }

  // https://github.com/mattn/go-sqlite3/issues/274#issuecomment-232897193
  newDbFile, err := copyDatabaseFile(dbPath, logger)
  if err != nil {
    return errors.WithStack(err)
  }

  writeDb, err := sql.Open("sqlite3", "file:"+newDbFile+"?cache=shared&_busy_timeout=5000&_journal_mode=OFF&mode=rw")
  if err != nil {
    return errors.WithStack(err)
  }

  writeDb.SetMaxOpenConns(1)

  //noinspection SqlResolve
  rows, err := readDb.Query("select id, tc_build_properties from report where tc_build_properties is not null")
  if err != nil {
    return errors.WithStack(err)
  }

  defer util.Close(rows, logger)

  for rows.Next() {
    var id string
    var data sql.RawBytes
    err = rows.Scan(&id, &data)
    if err != nil {
      return errors.WithStack(err)
    }

    var json map[string]interface{}
    err = jsoniter.ConfigFastest.Unmarshal(data, &json)
    if err != nil {
      return errors.WithStack(err)
    }

    keys := make([]string, 0, len(json))
    for k := range json {
      keys = append(keys, k)
    }

    sort.Strings(keys)
    p := properties.NewProperties()

    for _, key := range keys {
      i := json[key]
      switch v := i.(type) {
      case bool:
        if v {
          _, _, _ = p.Set(key, "true")
        } else {
          _, _, _ = p.Set(key, "false")
        }
      case string:
        _, _, _ = p.Set(key, v)
      case int:
        _, _, _ = p.Set(key, strconv.Itoa(v))
      case float64:
        vInt := int(v)
        if float64(vInt) != v {
          return errors.Errorf("really float?")
        }
        _, _, _ = p.Set(key, strconv.Itoa(vInt))
      default:
        return errors.Errorf("unknown type: %v", v)
      }
    }

    newData := teamcity.PropertiesToJson(p)

    // check that json is valid
    var validate map[string]interface{}
    err = jsoniter.UnmarshalFromString(newData, &validate)
    if err != nil {
      return errors.WithStack(err)
    }

    if string(data) == newData {
      continue
    }

    logger.Info("update row", zap.String("id", id))
    _, err = writeDb.Exec("update report set tc_build_properties = ? where id = ?", newData, id)
    if err != nil {
      return errors.WithStack(err)
    }
  }

  err = rows.Err()
  if err != nil {
    return errors.WithStack(err)
  }
  return nil
}
