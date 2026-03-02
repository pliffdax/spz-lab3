package replacement

import (
	"fmt"
	"math/rand"

	"github.com/pliffdax/spz-lab3/internal/model"
)

type Random struct {
	r *rand.Rand
}

func NewRandom(seed int64) *Random {
	return &Random{r: rand.New(rand.NewSource(seed))}
}

func (a *Random) Name() string { return "random" }

func (a *Random) PickVictim(_ int64, pm *model.PhysMem, _ []*model.Process) (int, error) {
	if len(pm.Pages) == 0 {
		return -1, fmt.Errorf("no physical pages")
	}
	for range 10_000 {
		ppn := a.r.Intn(len(pm.Pages))
		if pm.Pages[ppn].Occupied {
			return ppn, nil
		}
	}
	for ppn := range pm.Pages {
		if pm.Pages[ppn].Occupied {
			return ppn, nil
		}
	}
	return -1, fmt.Errorf("no victim found")
}

func (a *Random) OnAccess(now int64, pm *model.PhysMem, _ []*model.Process, ppn int) {
	pm.Pages[ppn].LastAccess = now
}
