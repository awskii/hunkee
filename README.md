# hunkee
Convenient way to parse logs

Currently unstable. No autotests, no benchmarks. But example works fine (with debug output for now).
## Usage
`go get github.com/awskii/hunkee`

You can take a glance of concurrent usage of lib in `example/main.go`.

Note that all concurrency dispatch is lying on your shoulders.

`$ go build && ./example test.log`
## Benchmarks
```
BenchmarkParse-4     	 2000000	       953 ns/op	      32 B/op	       1 allocs/op
BenchmarkParseRE-4   	  500000	      2482 ns/op	     448 B/op	       6 allocs/op
```

## Don't be an enemy of yourself
If you passing an unsupported interface or structure, dont't start an issue about something goes wrong.
If you create structure with raw field of any other type than string, don't be confused.
