package qb

import (
	"runtime"
	"sort"
	"sync"
	"time"

	"github.com/winebarrel/tachymeter"
)

type RecorderReport struct {
	DSN         string
	StartedAt   time.Time
	FinishedAt  time.Time
	ElapsedTime time.Duration
	TaskOpts
	Token       string
	GOMAXPROCS  int
	QueryCount  int
	AvgTPS      float64
	MaxTPS      float64
	MinTPS      float64
	MedianTPS   float64
	ExpectedTPS int
	Response    *tachymeter.Metrics
}

type RecorderOpts struct {
	DSN       string
	HInterval time.Duration
}

type Recorder struct {
	sync.Mutex
	RecorderOpts
	TaskOpts
	startedAt  time.Time
	finishedAt time.Time
	token      string
	channel    chan []recorderDataPoint
	dataPoints []recorderDataPoint
	stmtSize   int
}

func newRecorder(recOpts *RecorderOpts, taskOpts *TaskOpts, token string, stmtSize int) *Recorder {
	return &Recorder{
		RecorderOpts: *recOpts,
		TaskOpts:     *taskOpts,
		token:        token,
		stmtSize:     stmtSize,
	}
}

func (rec *Recorder) start(bufsize int) {
	rec.dataPoints = []recorderDataPoint{}
	ch := make(chan []recorderDataPoint, bufsize)
	rec.channel = ch

	go func() {
		for redDps := range ch {
			rec.appendDataPoints(redDps)
		}
	}()

	rec.startedAt = time.Now()
}

func (rec *Recorder) appendDataPoints(recDps []recorderDataPoint) {
	rec.Lock()
	defer rec.Unlock()
	rec.dataPoints = append(rec.dataPoints, recDps...)
}

func (rec *Recorder) close() {
	close(rec.channel)
	rec.finishedAt = time.Now()
}

func (rec *Recorder) tpsHist() []float64 {
	recDps := rec.dataPoints

	if len(recDps) == 0 {
		return []float64{}
	}

	sort.Slice(recDps, func(i, j int) bool {
		return recDps[i].timestamp.Before(recDps[j].timestamp)
	})

	minTm := recDps[0].timestamp
	hist := []int{0}

	for _, v := range recDps {
		if minTm.Add(1 * time.Second).Before(v.timestamp) {
			minTm = minTm.Add(1 * time.Second)
			hist = append(hist, 0)
		}

		hist[len(hist)-1]++
	}

	f64Hist := make([]float64, len(hist))

	for i, v := range hist {
		f64Hist[i] = float64(v) / float64(rec.stmtSize)
	}

	return f64Hist
}

func (rec *Recorder) qps() (minQPS float64, maxQPS float64, medianQPS float64) {
	tpsHist := rec.tpsHist()

	if len(tpsHist) == 0 {
		return
	}

	sort.Slice(tpsHist, func(i, j int) bool {
		return tpsHist[i] < tpsHist[j]
	})

	minQPS = tpsHist[0]
	maxQPS = tpsHist[len(tpsHist)-1]

	median := len(tpsHist) / 2
	medianNext := median + 1

	if len(tpsHist) == 1 {
		medianQPS = tpsHist[0]
	} else if len(tpsHist) == 2 {
		medianQPS = (tpsHist[0] + tpsHist[1]) / 2
	} else if len(tpsHist)%2 == 0 {
		medianQPS = (tpsHist[median] + tpsHist[medianNext]) / 2
	} else {
		medianQPS = tpsHist[medianNext]
	}

	return
}

type recorderDataPoint struct {
	timestamp time.Time
	resTime   time.Duration
}

func (rec *Recorder) add(recDps []recorderDataPoint) {
	rec.channel <- recDps
}

func (rec *Recorder) Report() (rr *RecorderReport) {
	nanoElapsed := rec.finishedAt.Sub(rec.startedAt)
	queryCnt := rec.Count()

	rr = &RecorderReport{
		DSN:         rec.DSN,
		StartedAt:   rec.startedAt,
		FinishedAt:  rec.finishedAt,
		ElapsedTime: nanoElapsed / time.Second,
		TaskOpts:    rec.TaskOpts,
		Token:       rec.token,
		GOMAXPROCS:  runtime.GOMAXPROCS(0),
		QueryCount:  queryCnt,
		AvgTPS:      float64(queryCnt) * float64(time.Second) / float64(nanoElapsed) / float64(rec.stmtSize),
		ExpectedTPS: rec.NAgents * rec.Rate,
	}

	t := tachymeter.New(&tachymeter.Config{
		Size:      len(rec.dataPoints),
		HBins:     10,
		HInterval: rec.HInterval,
	})

	for _, v := range rec.dataPoints {
		t.AddTime(v.resTime)
	}

	rr.Response = t.Calc()
	rr.MinTPS, rr.MaxTPS, rr.MedianTPS = rec.qps()

	return
}

func (rec *Recorder) Count() int {
	rec.Lock()
	defer rec.Unlock()
	return len(rec.dataPoints)
}
