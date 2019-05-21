# tist

A simple project for parallel running sql against multiple TiDB clusters

## Prepare

1. Start multiple TiDB clusters or MySQL servers. I started several MySQL servers on my 2017 MBP 13 inch. They are totally enough for testing.
2. Create TiDB configruation file. Refer to [config/tidb-clusters.json](config/tidb-clusters.json)
3. Create some sql files with same file prefix and postfix. e.g. `sql-0.sql sql-1.sql`. Leave a blank line at the end of the SQL files.

### TiDB Configguration File

```json
[
    {
        "host": "127.0.0.1",
        "port": 3306,
        "user": "root",
        "password": "root",
        "database": "interview"
    },
    {
        "host": "127.0.0.1",
        "port": 3307,
        "user": "root",
        "password": "root",
        "database": "interview"
    },
    {
        "host": "127.0.0.1",
        "port": 3308,
        "user": "root",
        "password": "root",
        "database": "interview"
    }
]
```

## Usage

```text
A test tool for parallel running sql against multiple TiDB clusters.

Usage:
  tist [flags]

Flags:
  -c, --client-number int        number of client (default 3)
  -h, --help                     help for tist
  -p, --sql-file-prefix string   prefix for SQL files (default "./config/sql")
  -s, --sql-file-suffix string   suffix for SQL files (default "sql")
  -t, --tidb-config string       TiDB clusters JSON file (default "./config/tidb-clusters.json")
  -v, --verbose                  verbose output
```

## Example

```bash
$ ./tist -c 3 -p ./config/sql -s sql -t ./config/tidb-clusters.json
INFO[2019-05-21T16:52:50+08:00] ClientID: 0, SQL Num: 10
INFO[2019-05-21T16:52:50+08:00] ClientID: 1, SQL Num: 10
INFO[2019-05-21T16:52:50+08:00] ClientID: 2, SQL Num: 10
INFO[2019-05-21T16:52:50+08:00] Total SQL Num: 30
INFO[2019-05-21T16:54:15+08:00] Working on No.1000 permutation
INFO[2019-05-21T16:55:53+08:00] Working on No.2000 permutation
INFO[2019-05-21T16:56:24+08:00] SQL Num: 30, Permutation Num: 2286, Time Duration: 3m34.237561496s
```