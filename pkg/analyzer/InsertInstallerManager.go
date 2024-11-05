package analyzer

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/JetBrains/ij-perf-report-aggregator/pkg/sql-util"
	"golang.org/x/tools/container/intsets"
)

type InsertInstallerManager struct {
	sql_util.InsertDataManager

	maxId       uint32
	insertedIds intsets.Sparse
}

func NewInstallerInsertManager(insertContext context.Context, db driver.Conn) (*InsertInstallerManager, error) {
	//noinspection GrazieInspection
	insertManager := sql_util.NewBatchInsertManager(insertContext, db, "insert into installer", 1, slog.With("type", "installer"))

	manager := &InsertInstallerManager{
		InsertDataManager: sql_util.InsertDataManager{
			InsertManager: insertManager,
		},

		insertedIds: intsets.Sparse{},
	}

	//noinspection SqlResolve
	err := db.QueryRow(insertContext, "select max(id) from installer").Scan(&manager.maxId)
	if err != nil {
		return nil, fmt.Errorf("cannot query max id: %w", err)
	}

	return manager, nil
}

func (t *InsertInstallerManager) Insert(id int, changes []string) error {
	if t.insertedIds.Has(id) {
		return nil
	}

	if id <= int(t.maxId) {
		exists, err := t.CheckExists(t.InsertManager.Db.QueryRow(t.InsertManager.InsertContext, "select 1 from installer where id = "+strconv.Itoa(id)+" limit 1"))
		if err != nil {
			return err
		}

		if exists {
			return nil
		}
	}

	batch, err := t.InsertManager.PrepareForAppend()
	if err != nil {
		return err
	}

	err = batch.Append(uint32(id), changes)
	if err != nil {
		return fmt.Errorf("cannot append: %w", err)
	}

	t.insertedIds.Insert(id)
	return nil
}
