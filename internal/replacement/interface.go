package replacement

import (
    "fmt"

    "github.com/pliffdax/spz-lab3/internal/config"
    "github.com/pliffdax/spz-lab3/internal/model"
)

type Replacer interface {
    Name() string
    PickVictim(now int64, pm *model.PhysMem, procs []*model.Process) (ppn int, err error)
    OnAccess(now int64, pm *model.PhysMem, procs []*model.Process, ppn int)
}

func New(cfg config.Config) (Replacer, error) {
    switch cfg.Algorithm {
    case "random":
        return NewRandom(cfg.Seed), nil
    case "wsclock":
        return NewWSClock(cfg.Delta), nil
    default:
        return nil, fmt.Errorf("unknown algorithm: %s (use random|wsclock)", cfg.Algorithm)
    }
}
