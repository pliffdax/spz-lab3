package model

type PTE struct {
    Present bool
    R       bool
    M       bool
    PPN     int
}

type Access struct {
    PID     int
    VPN     int
    IsWrite bool
}
