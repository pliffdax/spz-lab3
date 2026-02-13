package replacement

import (
    "fmt"

    "github.com/pliffdax/spz-lab3/internal/model"
)

type WSClock struct {
    Hand  int
    Delta int64
}

func NewWSClock(delta int64) *WSClock {
    if delta <= 0 {
        delta = 2000
    }
    return &WSClock{Hand: 0, Delta: delta}
}

func (w *WSClock) Name() string { return "wsclock" }

func (w *WSClock) PickVictim(now int64, pm *model.PhysMem, procs []*model.Process) (int, error) {
    n := len(pm.Pages)
    if n == 0 {
        return -1, fmt.Errorf("no physical pages")
    }

    start := w.Hand
    firstDirtyOld := -1

    for i := 0; i < n; i++ {
        ppn := (start + i) % n
        page := &pm.Pages[ppn]
        if !page.Occupied {
            continue
        }

        owner := page.OwnerPID
        vpn := page.VPN
        if owner < 0 || owner >= len(procs) || vpn < 0 || vpn >= len(procs[owner].PageTable) {
            continue
        }

        pte := &procs[owner].PageTable[vpn]

        if pte.R {
            pte.R = false
            page.LastAccess = now
            continue
        }

        age := now - page.LastAccess
        if age <= w.Delta {
            continue
        }

        if !pte.M {
            w.Hand = (ppn + 1) % n
            return ppn, nil
        }

        if firstDirtyOld == -1 {
            firstDirtyOld = ppn
        }
    }

    if firstDirtyOld != -1 {
        w.Hand = (firstDirtyOld + 1) % n
        return firstDirtyOld, nil
    }

    oldest := -1
    oldestAge := int64(-1)
    for i := 0; i < n; i++ {
        ppn := (start + i) % n
        page := &pm.Pages[ppn]
        if !page.Occupied {
            continue
        }
        age := now - page.LastAccess
        if age > oldestAge {
            oldestAge = age
            oldest = ppn
        }
    }

    if oldest != -1 {
        w.Hand = (oldest + 1) % n
        return oldest, nil
    }

    return -1, fmt.Errorf("no victim found")
}

func (w *WSClock) OnAccess(_ int64, _ *model.PhysMem, _ []*model.Process, _ int) {}
