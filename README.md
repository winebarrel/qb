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
01:00 | 10 agents / run 2907294 queries (4260 tps)

{
  "DSN": "root@/",
  "StartedAt": "2022-05-09T17:55:42.674904+09:00",
  "FinishedAt": "2022-05-09T17:56:42.677028+09:00",
  "ElapsedTime": 60,
  "NAgents": 10,
  "Rate": 0,
  "TransactionType": "tpcb-like",
  "Engine": "",
  "Token": "1665fb47-eb84-4ef4-9f12-f8ccc84b9248",
  "GOMAXPROCS": 16,
  "QueryCount": 2907294,
  "AvgTPS": 4404.898497183622,
  "MaxTPS": 5315.090909090909,
  "MinTPS": 0.2727272727272727,
  "MedianTPS": 4452.090909090909,
  "ExpectedTPS": 0,
  "Response": {
    "Time": {
      "Cumulative": "9m47.730783642s",
      "HMean": "150.102µs",
      "Avg": "202.157µs",
      "P50": "156.668µs",
      "P75": "208.138µs",
      "P95": "383.366µs",
      "P99": "669.418µs",
      "P999": "5.50043ms",
      "Long5p": "890.04µs",
      "Short5p": "80.223µs",
      "Max": "39.475841ms",
      "Min": "37.719µs",
      "Range": "39.438122ms",
      "StdDev": "393.538µs"
    },
    "Rate": {
      "Second": 4946.642375926489
    },
    "Samples": 2907294,
    "Count": 2907294,
    "Histogram": [
      {
        "37µs - 3.981ms": 2900756
      },
      {
        "3.981ms - 7.925ms": 5548
      },
      {
        "7.925ms - 11.869ms": 643
      },
      {
        "11.869ms - 15.812ms": 137
      },
      {
        "15.812ms - 19.756ms": 37
      },
      {
        "19.756ms - 23.7ms": 39
      },
      {
        "23.7ms - 27.644ms": 47
      },
      {
        "27.644ms - 31.588ms": 67
      },
      {
        "31.588ms - 35.532ms": 15
      },
      {
        "35.532ms - 39.475ms": 5
      }
    ]
  }
}
```
