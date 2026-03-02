package sim

import (
	"math/rand"

	"github.com/pliffdax/spz-lab3/internal/config"
	"github.com/pliffdax/spz-lab3/internal/model"
	"github.com/pliffdax/spz-lab3/internal/replacement"
)

type Result struct {
	TotalAccesses int64
	PageFaults    int64
	Evictions     int64
	Writebacks    int64

	PerProc []ProcStats
}

type ProcStats struct {
	PID        int
	Pages      int
	Refs       int64
	Faults     int64
	Evictions  int64
	Writebacks int64
}

type Simulator struct {
	cfg config.Config
	rep replacement.Replacer

	rng *rand.Rand

	pm    *model.PhysMem
	procs []*model.Process

	now int64
}

func New(cfg config.Config, rep replacement.Replacer) *Simulator {
	r := rand.New(rand.NewSource(cfg.Seed))
	pm := model.NewPhysMem(cfg.Frames)

	procs := make([]*model.Process, 0, cfg.Procs)
	for pid := 0; pid < cfg.Procs; pid++ {
		pages := cfg.PagesMin
		if cfg.PagesMax > cfg.PagesMin {
			pages = cfg.PagesMin + r.Intn(cfg.PagesMax-cfg.PagesMin+1)
		}
		procs = append(procs, model.NewProcess(pid, pages, cfg.WS))
	}

	return &Simulator{
		cfg:   cfg,
		rep:   rep,
		rng:   r,
		pm:    pm,
		procs: procs,
		now:   0,
	}
}

func (s *Simulator) Run() Result {
	res := Result{
		TotalAccesses: 0,
		PageFaults:    0,
		Evictions:     0,
		Writebacks:    0,
		PerProc:       make([]ProcStats, len(s.procs)),
	}
	for i, p := range s.procs {
		res.PerProc[i] = ProcStats{PID: p.PID, Pages: p.Pages}
	}

	pid := 0
	qLeft := s.cfg.Quantum

	for step := int64(0); step < s.cfg.Steps; step++ {
		s.now++

		if s.cfg.WSShift > 0 && (s.now%int64(s.cfg.WSShift) == 0) {
			for _, p := range s.procs {
				p.ShiftWorkingSet()
			}
		}

		if qLeft <= 0 {
			pid = (pid + 1) % len(s.procs)
			qLeft = s.cfg.Quantum
		}
		qLeft--

		p := s.procs[pid]
		vpn := p.NextVPN(s.rng)
		isWrite := s.rng.Float64() < s.cfg.WriteProb

		res.TotalAccesses++
		res.PerProc[pid].Refs++

		s.handleAccess(&res, p, vpn, isWrite)
	}

	return res
}

func (s *Simulator) handleAccess(res *Result, p *model.Process, vpn int, isWrite bool) {
	pte := &p.PageTable[vpn]

	if pte.Present {
		pte.R = true
		if isWrite {
			pte.M = true
		}
		if pte.PPN >= 0 && pte.PPN < len(s.pm.Pages) {
			s.rep.OnAccess(s.now, s.pm, s.procs, pte.PPN)
		}
		return
	}

	res.PageFaults++
	res.PerProc[p.PID].Faults++

	var ppn int
	if s.pm.HasFree() {
		ppn = s.pm.AllocFree()
	} else {
		victim, err := s.rep.PickVictim(s.now, s.pm, s.procs)
		if err != nil {
			return
		}
		s.evictVictim(res, victim)
		ppn = victim
	}

	pte.Present = true
	pte.R = true
	pte.M = isWrite
	pte.PPN = ppn

	s.pm.Pages[ppn].Occupied = true
	s.pm.Pages[ppn].OwnerPID = p.PID
	s.pm.Pages[ppn].VPN = vpn
	s.pm.Pages[ppn].LastAccess = s.now
}

func (s *Simulator) evictVictim(res *Result, victimPPN int) {
	v := &s.pm.Pages[victimPPN]
	if !v.Occupied {
		return
	}

	owner := v.OwnerPID
	vpn := v.VPN
	if owner < 0 || owner >= len(s.procs) || vpn < 0 || vpn >= len(s.procs[owner].PageTable) {
		s.pm.FreeFrame(victimPPN)
		return
	}

	proc := s.procs[owner]
	pte := &proc.PageTable[vpn]

	if pte.M {
		res.Writebacks++
		res.PerProc[owner].Writebacks++
		pte.M = false
	}

	pte.Present = false
	pte.R = false
	pte.PPN = -1

	res.Evictions++
	res.PerProc[owner].Evictions++

	v.Occupied = false
	v.OwnerPID = model.NoOwner
	v.VPN = -1
	v.LastAccess = 0
}
