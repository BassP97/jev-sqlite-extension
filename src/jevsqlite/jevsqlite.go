package jevsqlite

import (
	"context"
	"encoding/json"
	"fmt"
	"jev-sqlite-extension/jev"
	"time"

	"github.com/mattn/go-sqlite3"
)

const Timeout = 10 * time.Second

func Register(conn *sqlite3.SQLiteConn, apiKey string) error {
	return conn.RegisterFunc("jev_classify", func(text, question, options string) (string, error) {
		var opts map[string]string
		if err := json.Unmarshal([]byte(options), &opts); err != nil {
			return "", fmt.Errorf("jev_classify: options must be a JSON object of label -> description: %w", err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), Timeout)
		defer cancel()
		label, _, err := jev.Classify(ctx, apiKey, text, question, opts)
		return label, err
	}, false)
}

func Driver(apiKey string) *sqlite3.SQLiteDriver {
	return &sqlite3.SQLiteDriver{
		ConnectHook: func(conn *sqlite3.SQLiteConn) error { return Register(conn, apiKey) },
	}
}
