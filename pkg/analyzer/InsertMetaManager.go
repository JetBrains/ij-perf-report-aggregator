package analyzer

import (
  "context"
  "github.com/develar/errors"
  "github.com/jackc/pgx/v5/pgxpool"
  "os"
)

type InsertMetaManager struct {
  dbPool  *pgxpool.Pool
  context context.Context
}

func NewInsertMetaManager(insertContext context.Context) (*InsertMetaManager, error) {
  dbpool, err := pgxpool.New(insertContext, os.Getenv("DATABASE_URL"))
  if err != nil {
    return nil, errors.WithStack(err)
  }

  manager := &InsertMetaManager{
    dbPool:  dbpool,
    context: insertContext,
  }

  return manager, nil
}

func (t *InsertMetaManager) InsertProjectDescription(project string, branch string, url string, methodName string, description string) error {
  if project != "" && branch != "" && (methodName != "" || description != "" || url != "") {
    _, err := t.dbPool.Exec(t.context, "INSERT INTO project_description (project, branch, url, methodName, description) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (project, branch) DO UPDATE SET methodName = excluded.methodName, url = excluded.url, description = excluded.description;", project, branch, url, methodName, description)
    if err != nil {
      return err
    }
  }
  return nil
}
