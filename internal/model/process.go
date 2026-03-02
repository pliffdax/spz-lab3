package model

import "math/rand"

type Process struct {
	PID   int
	Pages int
	WS    int

	PageTable []PTE

	wsBase int
}

func NewProcess(pid int, pages int, ws int) *Process {
	pt := make([]PTE, pages)
	for i := range pages {
		pt[i] = PTE{Present: false, R: false, M: false, PPN: -1}
	}
	if ws > pages {
		ws = pages
	}
	return &Process{
		PID:       pid,
		Pages:     pages,
		WS:        ws,
		PageTable: pt,
		wsBase:    0,
	}
}

func (p *Process) ShiftWorkingSet() {
	if p.Pages == 0 || p.WS == 0 {
		return
	}
	p.wsBase++
	if p.wsBase >= p.Pages {
		p.wsBase = 0
	}
}

func (p *Process) NextVPN(r *rand.Rand) int {
	if p.Pages == 0 {
		return 0
	}

	if r.Float64() < 0.90 && p.WS > 0 {
		off := r.Intn(p.WS)
		return (p.wsBase + off) % p.Pages
	}

	return r.Intn(p.Pages)
}
