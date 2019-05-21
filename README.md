# tist

A simple project for parallel running sql against multiple TiDB clusters

## Prepare

1. Start multiple TiDB clusters or MySQL servers. Due to my dev machine is a 2017 MBP 13 inch, I started 3 MySQL servers by docker.
2. Create some sql files with same file prefix and postfix. e.g. `sql-0.sql sql-1.sql`

## Usage

```text
A test tool for parallel running sql against multiple TiDB clusters.

Usage:
  tist [flags]

Flags:
      --client-number int         Number of client (default 3)
  -h, --help                      help for tist
      --sql-file-postfix string   Postfix for SQL files (default "sql")
      --sql-file-prefix string    Prefix for SQL files (default "./config/sql")
      --tidb-config string        TiDB clusters JSON file (default "./config/tidb-clusters.json")
```
