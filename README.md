# hunkee

[![GoDoc](https://godoc.org/github.com/awskii/hunkee?status.svg)](https://godoc.org/github.com/awskii/hunkee)
[![Go Report Card](https://goreportcard.com/badge/github.com/awskii/hunkee)](https://goreportcard.com/report/github.com/awskii/hunkee)
[![License: MIT](https://img.shields.io/github/license/mashape/apistatus.svg)](https://opensource.org/licenses/MIT)

Convenient way to parse logs

All you need to parse log file - add "format stirng" and provide line to parse and structure to parse into.

You can specify raw field, which will be filled with raw value (string) of token.
Simple:
```go
type s struct {
  ID     int64  `hunk:"id"`
  IDRaw string `hunk:"id_raw"`
}
```
For that example format string might be as simple as `":id"`.

You can use raw values to parse not supported types.

Note that dots in tags are not supported. Embedded structs are not supported too.

## Supported types
* int, int8, int16, int32, int64
* uint, uint8, uint16, uint32, uint64
* bool
* string
* time.Time (with layout and timezone parsing)
* time.Duration
* net.IP
* url.URL

## Usage
Take a glance on that example (same at example/main.go):
```go
// define your structure with needed field. Notice that you can set *_raw fields, which will be
// filled by raw value
var s struct {
	ParsedAtTimestamp          time.Time `hunk:"date"`
	RemoteUser                 string    `hunk:"remote_user"`
	RemoteIP                   net.IP    `hunk:"remote_addr"`
	RemoteIPRaw                string    `hunk:"remote_addr_raw"`
	GeoIPCountryCode           string    `hunk:"geoip_country_code"`
	GeoIPCity                  string    `hunk:"geoip_city"`
	GcdnTimetamp               time.Time `hunk:"time_local"`
	GcdnResponderName          string    `hunk:"responder_name"`
	GcdnAPIClientID            uint64    `hunk:"gcdn_api_client_id"`
	GcdnResourceID             int64     `hunk:"gcdn_api_resource_id"`
	RQHeaderHost               string    `hunk:"host"`
	RQHeaderUserAgent          string    `hunk:"http_user_agent"`
	RQLength                   int64     `hunk:"request_length"`
	Size                       uint64    `hunk:"body_bytes_sent"`
	HTTPRequestRaw             string    `hunk:"request"`
	HTTPStatus                 int       `hunk:"status"`
	HTTPReferer                string    `hunk:"http_referer"`
	HTTPScheme                 string    `hunk:"scheme"`
	HTTPRangeRaw               string    `hunk:"http_range"`
	ServerToClientBytesSent    uint64    `hunk:"bytes_sent"`
	ServerToClientBytesSentRaw string    `hunk:"bytes_sent_raw"`
	SentHTTPContentSize        uint64    `hunk:"sent_http_content_size"`
	SentHTTPContentSizeRaw     string    `hunk:"sent_http_content_size_raw"`
	UpstreamResponseTimeRaw    string    `hunk:"upstream_response_time"`
	UpstreamResponseLengthRaw  string    `hunk:"upstream_response_length"`
	CacheStatus                string    `hunk:"upstream_cache_status"`
	ProcessingTime             float64   `hunk:"request_time"`
	UpstreamIPRaw              string    `hunk:"upstream_addr"`
	UIDCookieGot               string    `hunk:"uid_got"`
	UIDCookieSet               string    `hunk:"uid_set"`
	ShieldUsedRaw              string    `hunk:"shield_type"`
}
// define format string
f := `:remote_addr - :remote_user :time_local :request :status ` +
	`:body_bytes_sent :http_referer :http_user_agent :bytes_sent :sent_http_content_size ` +
	`:scheme :host :request_time :upstream_response_time :request_length :http_range ` +
	`:responder_name :upstream_cache_status :upstream_response_length :upstream_addr ` +
	`:gcdn_api_client_id :gcdn_api_resource_id :uid_got :uid_set :geoip_country_code ` +
	`:geoip_city :shield_type`

// log line
l := `"62.149.10.131" "-" "-" "[04/Jan/2018:19:15:39 +0000]" "GET /dino.jpg HTTP/1.1" "200" "207402" "" "saelmon" "207957" "-" "https" "di.gcdn.co" "0.000" "-" "88" "-" "[gn]" "HIT" "-" "-" "777" "1337" "-" "-" "UA" "-" "shield_no"`

// initialize parser
p, err := NewParser(f, &s)
if err != nil {
	fmt.Println(err)
}

// set time parsing format for time_local field and say that all tokens separated by '"'
p.SetTokenSeparator('"')
p.SetTimeLayout("time_local", "[02/Jan/2006:15:04:05 -0700]")

// parse line into s
if err = p.ParseLine(l, &s); err != nil {
	fmt.Println(err)
}
fmt.Printf("%#v\n", s)
```

Note that all concurrency dispatch is lying on your shoulders.

## Benchmarks
```
goos: linux
goarch: amd64
pkg: github.com/awskii/hunkee
BenchmarkParse-4                	 1000000	      1061 ns/op	      32 B/op	       1 allocs/op
BenchmarkParseWithoutTime-4     	 5000000	       406 ns/op	       0 B/op	       0 allocs/op
BenchmarkParseRE-4              	 1000000	      2192 ns/op	     448 B/op	       6 allocs/op
BenchmarkParseREWithoutTime-4   	 1000000	      1246 ns/op	     256 B/op	       4 allocs/op
```

## Don't be an enemy of yourself
If you passing an unsupported interface or structure, dont't start an issue about something goes wrong.
If you create structure with raw field of any other type than string, don't be confused.
