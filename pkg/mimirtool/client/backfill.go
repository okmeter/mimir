package client

import (
	"context"
	"net/http"
)

func (c *MimirClient) Backfill(ctx context.Context, source string) error {
	res, err := c.doRequest("/api/v1/backfill", http.MethodPost, nil)
	if err != nil {
		return err
	}
	res.Body.Close()

	return nil
}
