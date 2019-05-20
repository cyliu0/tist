# tist

A simple project for parallel running sql against multiple TiDB clusters

## Usage

```
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