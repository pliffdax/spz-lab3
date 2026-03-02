package config

type Config struct {
	Algorithm string

	Frames   int
	Procs    int
	PagesMin int
	PagesMax int

	WS      int
	WSShift int
	Quantum int

	Steps int64
	Seed  int64

	WriteProb float64

	Delta int64
}

func Default() Config {
	return Config{
		Algorithm: "wsclock",
		Frames:    64,
		Procs:     6,
		PagesMin:  64,
		PagesMax:  256,
		WS:        8,
		WSShift:   2000,
		Quantum:   50,
		Steps:     200000,
		Seed:      42,
		WriteProb: 0.30,
		Delta:     50,
	}
}
