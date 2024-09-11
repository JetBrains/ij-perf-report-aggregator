package data_query

import (
	"context"
	"errors"
	"github.com/ClickHouse/ch-go"
	"github.com/ClickHouse/ch-go/proto"
	"github.com/jackc/puddle/v2"
	"io"
	"net"
	"time"
)

type DatabaseConnectionSupplier interface {
	AcquireDatabase(ctx context.Context, name string) (*puddle.Resource[*ch.Client], error)
}

func executeQuery(ctx context.Context, sqlQuery string, query Query, dbSupplier DatabaseConnectionSupplier, resultHandler func(ctx context.Context, block proto.Block, result *proto.Results) error) error {
	for range 120 {
		dbResource, err := dbSupplier.AcquireDatabase(ctx, query.Database)
		if err != nil {
			return err
		}

		done, err := doExecution(ctx, sqlQuery, dbResource, resultHandler)
		if err != nil {
			return err
		}

		if done {
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}

	return errors.New("cannot acquire database")
}

func doExecution(ctx context.Context, sqlQuery string, client *puddle.Resource[*ch.Client], resultHandler func(ctx context.Context, block proto.Block, result *proto.Results) error) (bool, error) {
	isDestroyed := false
	defer func() {
		if !isDestroyed {
			client.Release()
		}
	}()

	var result proto.Results
	err := client.Value().Do(ctx, ch.Query{
		Body:   sqlQuery,
		Result: result.Auto(),
		OnResult: func(ctx context.Context, block proto.Block) error {
			return resultHandler(ctx, block, &result)
		},
	})

	if err == nil {
		return true, nil
	}

	// if net error or io error - connection was closed due to inactivity, destroy it and acquire a new one
	if !isNetError(err) && !client.Value().IsClosed() {
		return true, err
	}

	isDestroyed = true
	client.Destroy()
	return false, nil
}

func isNetError(err error) bool {
	if err != nil {
		if errors.Is(err, io.ErrUnexpectedEOF) {
			return true
		}
		var netError net.Error
		if errors.As(err, &netError) {
			return true
		}
	}
	return false
}
