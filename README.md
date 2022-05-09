# qb

qb is a MySQL benchmarking tool using TPC-B(pgbench).

## Usage

```
qb - MySQL benchmarking tool using TCP-B(same as pgbench).

  Flags:
       --version       Displays the program version string.
    -h --help          Displays help with available flag, subcommand, and positional value parameters.
    -d --dsn           Data Source Name, see https://github.com/go-sql-driver/mysql#examples.
    -i --initialize    Invokes initialization mode.
    -n --nagents       Number of agents. (default: 1)
    -t --time          Test run time (sec). Zero is infinity. (default: 60)
    -r --rate          Rate limit for each agent (qps). Zero is unlimited. (default: 0)
    -T --type          Transaction type (tpcb-like,simple-update,select-only). (default: tpcb-like)
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
01:00 | 1 agents / run 742111 queries (1015 tps)

{
  "DSN": "root@/",
  "StartedAt": "2022-05-09T17:47:45.535164+09:00",
  "FinishedAt": "2022-05-09T17:48:45.536835+09:00",
  "ElapsedTime": 60,
  "NAgents": 1,
  "Rate": 0,
  "TransactionType": "tpcb-like",
  "Scale": 1,
  "Engine": "",
  "Token": "8c5c8e33-0356-42f0-b95b-5e47e74a7e85",
  "GOMAXPROCS": 16,
  "QueryCount": 742111,
  "AvgTPS": 1124.395525948911,
  "MaxTPS": 1219,
  "MinTPS": 1012.6363636363636,
  "MedianTPS": 1152.2727272727273,
  "ExpectedTPS": 0,
  "Response": {
    "Time": {
      "Cumulative": "58.485125316s",
      "HMean": "64.226µs",
      "Avg": "78.809µs",
      "P50": "69.669µs",
      "P75": "83.322µs",
      "P95": "162.796µs",
      "P99": "215.801µs",
      "P999": "317.144µs",
      "Long5p": "213.344µs",
      "Short5p": "36.897µs",
      "Max": "36.677174ms",
      "Min": "30.205µs",
      "Range": "36.646969ms",
      "StdDev": "133.311µs"
    },
    "Rate": {
      "Second": 12688.884498243142
    },
    "Samples": 742111,
    "Count": 742111,
    "Histogram": [
      {
        "30µs - 3.694ms": 742047
      },
      {
        "3.694ms - 7.359ms": 49
      },
      {
        "7.359ms - 11.024ms": 1
      },
      {
        "11.024ms - 14.688ms": 1
      },
      {
        "14.688ms - 18.353ms": 1
      },
      {
        "18.353ms - 22.018ms": 2
      },
      {
        "22.018ms - 25.683ms": 1
      },
      {
        "25.683ms - 29.347ms": 4
      },
      {
        "29.347ms - 33.012ms": 3
      },
      {
        "33.012ms - 36.677ms": 2
      }
    ]
  }
}
```
