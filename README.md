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

$ qb -d root@/
01:00 | 1 agents / run 618 queries (1 tps)

{
  "DSN": "root@/",
  "StartedAt": "2022-05-08T16:38:08.659458+09:00",
  "FinishedAt": "2022-05-08T16:39:08.658546+09:00",
  "ElapsedTime": 60,
  "NAgents": 1,
  "Rate": 0,
  "TransactionType": "tpcb-like",
  "Scale": 1,
  "Engine": "",
  "Token": "12cd77c8-b5a2-4882-9e99-5e17d1a39177",
  "GOMAXPROCS": 16,
  "QueryCount": 618,
  "AvgTPS": 1.14443101181896,
  "MaxTPS": 1.8888888888888888,
  "MinTPS": 0.1111111111111111,
  "MedianTPS": 1,
  "ExpectedTPS": 0,
  "Response": {
    "Time": {
      "Cumulative": "59.438755342s",
      "HMean": "154.035µs",
      "Avg": "96.179215ms",
      "P50": "131.713µs",
      "P75": "1.102764ms",
      "P95": "489.049785ms",
      "P99": "536.563846ms",
      "P999": "829.05208ms",
      "Long5p": "541.774974ms",
      "Short5p": "66.104µs",
      "Max": "897.154843ms",
      "Min": "45.857µs",
      "Range": "897.108986ms",
      "StdDev": "185.490362ms"
    },
    "Rate": {
      "Second": 10.397256746783109
    },
    "Samples": 618,
    "Count": 618,
    "Histogram": [
      {
        "45µs - 89.756ms": 481
      },
      {
        "89.756ms - 179.467ms": 1
      },
      {
        "179.467ms - 269.178ms": 1
      },
      {
        "269.178ms - 358.889ms": 59
      },
      {
        "358.889ms - 448.6ms": 4
      },
      {
        "448.6ms - 538.311ms": 66
      },
      {
        "538.311ms - 628.022ms": 3
      },
      {
        "628.022ms - 717.733ms": 1
      },
      {
        "717.733ms - 807.443ms": 1
      },
      {
        "807.443ms - 897.154ms": 1
      }
    ]
  }
}
```
