package main

import (
    "flag"
    "fmt"
    "os"

    "github.com/pliffdax/spz-lab3/internal/config"
    "github.com/pliffdax/spz-lab3/internal/replacement"
    "github.com/pliffdax/spz-lab3/internal/sim"
)

func main() {
    cfg := config.Default()

    flag.StringVar(&cfg.Algorithm, "alg", cfg.Algorithm, "replacement algorithm: random|wsclock")
    flag.IntVar(&cfg.Frames, "frames", cfg.Frames, "number of physical frames")
    flag.IntVar(&cfg.Procs, "procs", cfg.Procs, "number of processes")
    flag.IntVar(&cfg.PagesMin, "pagesMin", cfg.PagesMin, "min pages per process")
    flag.IntVar(&cfg.PagesMax, "pagesMax", cfg.PagesMax, "max pages per process")
    flag.IntVar(&cfg.WS, "ws", cfg.WS, "working set size")
    flag.IntVar(&cfg.WSShift, "wsShift", cfg.WSShift, "working set shift period (in accesses)")
    flag.IntVar(&cfg.Quantum, "quantum", cfg.Quantum, "scheduler quantum (in accesses)")
    flag.Int64Var(&cfg.Steps, "steps", cfg.Steps, "total memory accesses")
    flag.Int64Var(&cfg.Seed, "seed", cfg.Seed, "rng seed")
    flag.Int64Var(&cfg.Delta, "delta", cfg.Delta, "WSClock time window delta (in accesses)")
    flag.Float64Var(&cfg.WriteProb, "writeProb", cfg.WriteProb, "probability of write access [0..1]")
    flag.Parse()

    rep, err := replacement.New(cfg)
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }

    s := sim.New(cfg, rep)
    res := s.Run()
    fmt.Print(sim.FormatReport(cfg, rep.Name(), res))
}
