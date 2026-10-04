package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"jev-sqlite-extension/jev"
	"os"
	"strings"
	"time"

	"github.com/mattn/go-sqlite3"
)

func classify(text, question, options string) (string, error) {
	var opts map[string]string
	if err := json.Unmarshal([]byte(options), &opts); err != nil {
		return "", fmt.Errorf("jev_classify: options must be a JSON object of label -> description: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	label, _, err := jev.Classify(ctx, os.Getenv("TYPESAFE_API_KEY"), text, question, opts)
	return label, err
}

func main() {
	if len(os.Args) != 3 {
		os.Exit(2)
	}

	sql.Register("sqlite3_jev", &sqlite3.SQLiteDriver{
		ConnectHook: func(c *sqlite3.SQLiteConn) error {
			return c.RegisterFunc("jev_classify", classify, false)
		},
	})

	db, err := sql.Open("sqlite3_jev", os.Args[1])
	if err != nil {
		fmt.Printf("uh oh")
		return
	}
	defer db.Close()

	rows, err := db.Query(os.Args[2])
	if err != nil {
		os.Exit(1)
	}
	defer rows.Close()

	cols, _ := rows.Columns()
	for rows.Next() {
		vals := make([]sql.NullString, len(cols))
		ptrs := make([]any, len(cols))
		for i := range vals {
			ptrs[i] = &vals[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			fmt.Printf("uh oh")
			return
		}
		out := make([]string, len(cols))
		for i, v := range vals {
			out[i] = v.String
		}
		fmt.Println(strings.Join(out, "\t"))
	}
	if err := rows.Err(); err != nil {
		fmt.Printf("uh oh")
		return
	}
}
