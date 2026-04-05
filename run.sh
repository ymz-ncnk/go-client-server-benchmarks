# #!/bin/bash
GEN_SIZE=400000 go test -bench BenchmarkQPS -count=10 -timeout=30m > ./results/qps/benchmarks.txt
go test -bench BenchmarkFixed -benchtime=100000x -benchmem -count=10 >./results/fixed/benchmarks.txt

go generate