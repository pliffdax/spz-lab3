package sim

import (
    "fmt"

    "github.com/pliffdax/spz-lab3/internal/config"
)

func FormatReport(cfg config.Config, algName string, r Result) string {
    faultRate := 0.0
    if r.TotalAccesses > 0 {
        faultRate = float64(r.PageFaults) / float64(r.TotalAccesses)
    }

    out := ""
    out += "Lab3 paging simulation\n"
    out += fmt.Sprintf("alg=%s frames=%d pmax=%d ws=%d wsShift=%d quantum=%d steps=%d seed=%d delta=%d writeProb=%.2f\n",
        algName, cfg.Frames, cfg.PagesMax, cfg.WS, cfg.WSShift, cfg.Quantum, cfg.Steps, cfg.Seed, cfg.Delta, cfg.WriteProb)

    out += "\nSystem stats:\n"
    out += fmt.Sprintf("refs=%d faults=%d fault_rate=%.4f evictions=%d writebacks=%d\n",
        r.TotalAccesses, r.PageFaults, faultRate, r.Evictions, r.Writebacks)

    out += "\nPer-process summary:\n"
    for _, p := range r.PerProc {
        pr := 0.0
        if p.Refs > 0 {
            pr = float64(p.Faults) / float64(p.Refs)
        }
        out += fmt.Sprintf("pid=%d pages=%d refs=%d faults=%d fault_rate=%.4f evictions=%d writebacks=%d\n",
            p.PID, p.Pages, p.Refs, p.Faults, pr, p.Evictions, p.Writebacks)
    }

    return out
}
