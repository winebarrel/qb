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
qb -d root@/ -s 100
01:00 | 1 agents / run 761354 queries (1216 tps)

{
  "DSN": "root@/",
  "StartedAt": "2022-05-08T15:49:08.856942+09:00",
  "FinishedAt": "2022-05-08T15:50:08.855813+09:00",
  "ElapsedTime": 60,
  "NAgents": 1,
  "Rate": 0,
  "TransactionType": "tpcb-like",
  "Scale": 100,
  "Engine": "",
  "Token": "8814ce70-14c4-4f44-8953-9e091a0c1306",
  "GOMAXPROCS": 16,
  "QueryCount": 761354,
  "AvgTPS": 1409.90335270003,
  "MaxTPS": 1567,
  "MinTPS": 1177.6666666666667,
  "MedianTPS": 1409.5555555555557,
  "ExpectedTPS": 0,
  "Response": {
    "Time": {
      "Cumulative": "58.477221629s",
      "HMean": "64.625µs",
      "Avg": "76.806µs",
      "P50": "64.161µs",
      "P75": "80.607µs",
      "P95": "166.448µs",
      "P99": "191.518µs",
      "P999": "263.202µs",
      "Long5p": "197.654µs",
      "Short5p": "39.828µs",
      "Max": "39.216943ms",
      "Min": "31.043µs",
      "Range": "39.1859ms",
      "StdDev": "114.756µs"
    },
    "Rate": {
      "Second": 13019.667808951266
    },
    "Samples": 761354,
    "Count": 761354,
    "Histogram": [
      {
        "31µs - 3.949ms": 761305
      },
      {
        "3.949ms - 7.868ms": 36
      },
      {
        "7.868ms - 11.786ms": 3
      },
      {
        "11.786ms - 15.705ms": 2
      },
      {
        "15.705ms - 19.623ms": 1
      },
      {
        "19.623ms - 23.542ms": 1
      },
      {
        "23.542ms - 27.461ms": 3
      },
      {
        "27.461ms - 31.379ms": 1
      },
      {
        "31.379ms - 35.298ms": 1
      },
      {
        "35.298ms - 39.216ms": 1
      }
    ]
  }
}
```
