package progress

import (
	"time"

	"github.com/vbauerster/mpb/v4"
	"github.com/vbauerster/mpb/v4/decor"
)

func NewProgress() *mpb.Progress {
	return mpb.New(
		mpb.WithRefreshRate(time.Second),
	)
}

func AddBar(name string, p *mpb.Progress, total int64) *mpb.Bar {
	return p.AddBar(
		total,
		mpb.PrependDecorators(
			decor.Name(name),
			decor.CountersNoUnit("% d / % d"),
		),
		mpb.AppendDecorators(
			decor.Percentage(decor.WCSyncSpace),
			decor.Elapsed(decor.ET_STYLE_HHMMSS),
			decor.AverageETA(decor.ET_STYLE_HHMMSS),
		),
	)
}
