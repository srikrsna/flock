package server

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"github.com/srikrsna/flock/pkg"
	pb "github.com/srikrsna/flock/protos"
)

func handleBatch(ctx context.Context, db DB, tables map[string]flock.Table, req *pb.BatchInsertRequest) (*pb.BatchInsertResponse, error) {
	var rows []map[string]interface{}
	if err := gob.NewDecoder(bytes.NewReader(req.GetData())).Decode(&rows); err != nil {
		return nil, err
	}

	table, ok := tables[req.GetTable()]
	if !ok {
		return nil, errors.New("table not configured")
	}

	if err := flock.InsertBulk(ctx, db, rows, table, req.GetTableName()); err != nil {
		return nil, err
	}

	return &pb.BatchInsertResponse{Success: true}, nil
}
