# qb

qb is a MySQL benchmarking tool using TPC-B(pgbench).

## Usage

```
qb - MySQL benchmarking tool using TCP-B(pgbench).

  Flags:
       --version       Displays the program version string.
    -h --help          Displays help with available flag, subcommand, and positional value parameters.
    -d --dsn           Data Source Name, see https://github.com/go-sql-driver/mysql#examples.
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
$ qb -d root@/
01:00 | 1 agents / run 5371 queries (64 qps)

{
  "DSN": "root@/",
  "StartedAt": "2022-05-08T15:18:47.989408+09:00",
  "FinishedAt": "2022-05-08T15:19:47.992289+09:00",
  "ElapsedTime": 60,
  "NAgents": 1,
  "Rate": 0,
  "TransactionType": "tpcb-like",
  "Scale": 1,
  "Engine": "",
  "Token": "c71e9cfa-4348-4c89-bd6a-166fb990cc5d",
  "GOMAXPROCS": 16,
  "QueryCount": 5371,
  "AvgQPS": 89.51505524797236,
  "MaxQPS": 125,
  "MinQPS": 1,
  "MedianQPS": 90,
  "ExpectedQPS": 0,
  "Response": {
    "Time": {
      "Cumulative": "59.034815794s",
      "HMean": "154.71µs",
      "Avg": "10.991401ms",
      "P50": "153.853µs",
      "P75": "824.916µs",
      "P95": "58.429983ms",
      "P99": "82.055708ms",
      "P999": "88.942401ms",
      "Long5p": "71.506152ms",
      "Short5p": "60.321µs",
      "Max": "134.260474ms",
      "Min": "41.666µs",
      "Range": "134.218808ms",
      "StdDev": "21.344865ms"
    },
    "Rate": {
      "Second": 90.98021104600248
    },
    "Samples": 5371,
    "Count": 5371,
    "Histogram": [
      {
        "41µs - 13.463ms": 4167
      },
      {
        "13.463ms - 26.885ms": 5
      },
      {
        "26.885ms - 40.307ms": 374
      },
      {
        "40.307ms - 53.729ms": 470
      },
      {
        "53.729ms - 67.151ms": 193
      },
      {
        "67.151ms - 80.572ms": 84
      },
      {
        "80.572ms - 93.994ms": 73
      },
      {
        "93.994ms - 107.416ms": 2
      },
      {
        "107.416ms - 120.838ms": 1
      },
      {
        "120.838ms - 134.26ms": 2
      }
    ]
  }
}
```
