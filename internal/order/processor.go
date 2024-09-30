package order

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
)

type Processor struct {
	dbConn *sql.DB
}

func NewProcessor(conn *sql.DB) *Processor {
	return &Processor{conn}
}

func (p *Processor) Process(ctx context.Context, o Order) error {
	orderJSON, err := o.ToJSON()
	if err != nil {
		return fmt.Errorf("error while marshalling order: %w", err)
	}

	p.dbConn.ExecContext(ctx, `insert into orders (id, payload) values ($1, $2)`, o.ID, orderJSON)
	slog.Info("order processed successfully", "order_id", o.ID.String())

	return nil
}
