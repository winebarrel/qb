# qb

qb is a MySQL benchmarking tool using TPC-B(pgbench).

## Usage

```
qb - MySQL benchmarking tool using TCP-B(pgbench).

  Flags:
       --version       Displays the program version string.
    -h --help          Displays help with available flag, subcommand, and positional value parameters.
    -d --dsn           Data Source Name, see https://github.com/go-sql-driver/mysql#examples.
    -i --initialize    Invokes initialization mode.
    -n --nagents       Number of agents. (default: 1)
    -t --time          Test run time (sec). Zero is infinity. (default: 60)
    -r --rate          Rate limit for each agent (qps). Zero is unlimited. (default: 0)
    -T --type          Transaction type (tpcb-like). (default: tpcb-like)
    -s --scale         Scaling factor. (default: 1)
    -e --engine        Engine of the table to be created.
       --hinterval     Histogram interval, e.g. '100ms'. (default: 0)
       --only-print    Just print SQL without connecting to DB.
       --no-progress   Do not show progress.
```

```
$ qb -d root@/ -i -s 10
dropping old database...
creating database...
creating tables...
generating data...
analyzing tables...

$ qb -d root@/ -n 10
01:00 | 10 agents / run 2603958 queries (2559 tps)

{
  "DSN": "root@/",
  "StartedAt": "2022-05-08T17:36:46.173678+09:00",
  "FinishedAt": "2022-05-08T17:37:46.177411+09:00",
  "ElapsedTime": 60,
  "NAgents": 10,
  "Rate": 0,
  "TransactionType": "tpcb-like",
  "Scale": 1,
  "Engine": "",
  "Token": "4e5626d2-0e42-4159-a98b-61b4d6b9e7da",
  "GOMAXPROCS": 16,
  "QueryCount": 2613797,
  "AvgTPS": 4840.1553833887765,
  "MaxTPS": 6650.777777777777,
  "MinTPS": 0.2222222222222222,
  "MedianTPS": 5125.222222222223,
  "ExpectedTPS": 0,
  "Response": {
    "Time": {
      "Cumulative": "9m52.304185844s",
      "HMean": "158.97µs",
      "Avg": "226.606µs",
      "P50": "163.068µs",
      "P75": "232.565µs",
      "P95": "436.563µs",
      "P99": "782.877µs",
      "P999": "6.11045ms",
      "Long5p": "1.101047ms",
      "Short5p": "82.425µs",
      "Max": "72.844002ms",
      "Min": "846ns",
      "Range": "72.843156ms",
      "StdDev": "645.026µs"
    },
    "Rate": {
      "Second": 4412.930150536564
    },
    "Samples": 2613797,
    "Count": 2613797,
    "Histogram": [
      {
        "0s - 7.285ms": 2612000
      },
      {
        "7.285ms - 14.569ms": 1167
      },
      {
        "14.569ms - 21.853ms": 142
      },
      {
        "21.853ms - 29.138ms": 193
      },
      {
        "29.138ms - 36.422ms": 88
      },
      {
        "36.422ms - 43.706ms": 43
      },
      {
        "43.706ms - 50.991ms": 71
      },
      {
        "50.991ms - 58.275ms": 74
      },
      {
        "58.275ms - 65.559ms": 18
      },
      {
        "65.559ms - 72.844ms": 1
      }
    ]
  }
}
```
