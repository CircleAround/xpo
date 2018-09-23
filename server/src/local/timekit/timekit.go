package timekit

import "time"

type TimeProvider interface {
	Now() time.Time
}

type RealTimeProvider struct {
}

func (t *RealTimeProvider) Now() time.Time {
	return time.Now()
}

type TestTimeProvider struct {
	now time.Time
}

func (p *TestTimeProvider) Now() time.Time {
	if p.now.IsZero() {
		return time.Now()
	} else {
		return p.now
	}
}

func (p *TestTimeProvider) StopAt(time time.Time) time.Time {
	p.now = time
	return time
}

func (p *TestTimeProvider) StopNow() time.Time {
	return p.StopAt(time.Now())
}

func (p *TestTimeProvider) Release() {
	p.now = time.Time{}
}
