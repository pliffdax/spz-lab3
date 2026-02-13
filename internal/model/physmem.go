package model

const NoOwner = -1

type PhysPage struct {
    Occupied bool
    OwnerPID int
    VPN      int

    LastAccess int64
}

type PhysMem struct {
    Pages []PhysPage
    Free  []int
}

func NewPhysMem(frames int) *PhysMem {
    pm := &PhysMem{
        Pages: make([]PhysPage, frames),
        Free:  make([]int, 0, frames),
    }
    for i := 0; i < frames; i++ {
        pm.Pages[i] = PhysPage{
            Occupied:   false,
            OwnerPID:   NoOwner,
            VPN:        -1,
            LastAccess: 0,
        }
        pm.Free = append(pm.Free, i)
    }
    return pm
}

func (pm *PhysMem) HasFree() bool { return len(pm.Free) > 0 }

func (pm *PhysMem) AllocFree() int {
    n := len(pm.Free)
    ppn := pm.Free[n-1]
    pm.Free = pm.Free[:n-1]
    return ppn
}

func (pm *PhysMem) FreeFrame(ppn int) {
    p := &pm.Pages[ppn]
    p.Occupied = false
    p.OwnerPID = NoOwner
    p.VPN = -1
    p.LastAccess = 0
    pm.Free = append(pm.Free, ppn)
}
