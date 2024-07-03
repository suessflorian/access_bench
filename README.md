Running at 1M

```
 Documents/index_diff % go test -v -bench=. .
goos: darwin
goarch: arm64
pkg: github.com/suessflorian/index_diff
BenchmarkRandomSelectsFromNonUniqueIndex
BenchmarkRandomSelectsFromNonUniqueIndex-8   	   2319	   518879 ns/op
BenchmarkRandomSelectsFromUniqueIndex
BenchmarkRandomSelectsFromUniqueIndex-8      	   2725	   380747 ns/op
PASS
ok  	github.com/suessflorian/index_diff	21.846s
```
