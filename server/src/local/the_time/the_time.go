package the_time

import "time"

type Provider interface {
	Now() time.Time
}

type RealTimeProvider struct {
}

func Real() *RealTimeProvider {
	return new(RealTimeProvider)
}

func (t *RealTimeProvider) Now() time.Time {
	return time.Now()
}

type ControlableTimeProvider struct {
	now time.Time
}

func Machine() *ControlableTimeProvider {
	return new(ControlableTimeProvider)
}

func (p *ControlableTimeProvider) Now() time.Time {
	if p.now.IsZero() {
		return time.Now()
	}
	return p.now
}

func (p *ControlableTimeProvider) TravelTo(time time.Time) time.Time {
	p.now = time
	return time
}

func (p *ControlableTimeProvider) TravelToNow() time.Time {
	return p.TravelTo(time.Now())
}

func (p *ControlableTimeProvider) BackToReal() {
	p.now = time.Time{}
}
