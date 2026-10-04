// jevsql runs a SQL statement against a SQLite file with jev_classify() available.
//
//	export TYPESAFE_API_KEY=...
//	go run ./cmd/jevsql tickets.db "SELECT id, body FROM tickets WHERE
//	  jev_classify(body, 'How urgent is this?',
//	    json_object('urgent','needs attention now','normal','can wait')) = 'urgent'"
package main

import (
	"database/sql"
	"fmt"
	"jev-sqlite-extension/jevsqlite"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "use like so: jevsql <db file> <sql>")
		os.Exit(2)
	}

	sql.Register("sqlite3_jev", jevsqlite.Driver(os.Getenv("TYPESAFE_API_KEY")))

	db, err := sql.Open("sqlite3_jev", os.Args[1])
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	defer db.Close()

	rows, err := db.Query(os.Args[2])
	if err != nil {
		fmt.Print(err.Error())
		return
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
			fmt.Print(err.Error())
			return
		}
		out := make([]string, len(cols))
		for i, v := range vals {
			out[i] = v.String
		}
		fmt.Println(strings.Join(out, "\t"))
	}
	if err := rows.Err(); err != nil {
		fmt.Print(err.Error())
		return
	}
}
