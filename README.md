# spz-lab3

## Build

```bash
go build -o ./bin/vm-sim ./cmd/vm-sim
```

## Run

```bash
./bin/vm-sim -alg=random
./bin/vm-sim -alg=wsclock
```

## Notes

Recommended WSClock time window:

```bash
./bin/vm-sim -alg=wsclock -delta=50
```
