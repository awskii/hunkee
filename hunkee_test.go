package hunkee

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func TestNewParser(t *testing.T) {
	var s struct {
		Int int `hunk:"int"`
	}

	_, err := NewParser(":ab: :c", &s)
	if err == nil {
		t.Error("expected init error, got nil")
	}
}

func TestSetTimeLayout(t *testing.T) {
	var s struct {
		T  time.Time `hunk:"t"`
		Tr string    `hunk:"t_raw"`
	}

	p, err := NewParser(":t", &s)
	if err != nil {
		t.Error("init error " + err.Error())
	}

	p.SetTimeLayout("t", time.Kitchen)
	str := "5:43PM"
	if err := p.ParseLine(str, &s); err != nil {
		t.Error(err)
	}

	if s.T.Minute() != 43 {
		t.Errorf("unexpected value after parsing:\nhave: %d\nwant %d", s.T.Minute(), 43)
	}

	if s.Tr != "5:43PM" {
		t.Errorf("expected another raw value:\nhave: %q\nwant: %q", s.Tr, str)
	}
}

func TestSetMultiplyTimeLayouts(t *testing.T) {
	var s struct {
		A  time.Time `hunk:"a"`
		B  time.Time `hunk:"b"`
		C  time.Time `hunk:"c"`
		D  time.Time `hunk:"d"`
		E  time.Time `hunk:"e"`
		Cr string    `hunk:"c_raw"`
	}

	p, err := NewParser(":a :b :d :e :c", &s)
	if err != nil {
		t.Error("unexpected init error " + err.Error())
	}
	if p == nil {
		t.Fatal("returned nil parser")
	}

	layouts := map[string]string{
		"a": time.RFC3339,
		"b": time.RFC1123,
		"c": time.RFC1123,
		"d": time.Kitchen,
		"e": "2006-01-02 15:04:05",
	}

	v := map[string]string{
		"a": "2018-07-28T21:10:45+10:00",
		"b": "Tue, 10 Apr 2018 19:17:21 UTC",
		"c": "Tue, 10 Apr 2018 19:17:33 UTC",
		"d": "5:43PM",
		"e": "2006-01-02 03:04:05",
	}

	p.SetMultiplyTimeLayout(layouts)
	p.SetTokenSeparator('"')
	str := fmt.Sprintf("\"%s\" \"%s\" \"%s\" \"%s\" \"%s\"", v["a"], v["b"], v["d"], v["e"], v["c"])
	if err := p.ParseLine(str, &s); err != nil {
		t.Error(err)
	}

	if _, o := s.A.Zone(); o != 36000 || s.A.Year() != 2018 || s.A.Hour() != 21 {
		t.Errorf("wrong parsed time with options:\nhave: %q\nwant: %q", s.A.String(), v["a"])
	}

	if _, o := s.C.Zone(); o != 0 || s.C.Month() != 4 || s.C.Second() != 33 {
		t.Errorf("wrong parsed time with options:\nhave: %q\nwant: %q", s.C.String(), v["c"])
	}

	if s.D.Hour() != 17 || s.D.Minute() != 43 {
		t.Errorf("wrong parsed time with options:\nhave: %q\nwant: %q", s.D.String(), v["d"])
	}
}

func TestParseLine(t *testing.T) {
	var s struct {
		ID   int    `hunk:"id"`
		Name string `hunk:"name"`
	}

	p, err := NewParser(":id :name", &s)
	if err != nil {
		t.Error("unexpected error: " + err.Error())
	}

	if err := p.parseLine("998 Gordon", &s); err != nil {
		t.Error(err)
	}

	if s.ID != 998 {
		t.Errorf("unexpected result of parsing commented string:\nhave: %d\nwant: %d", s.ID, 998)
	}
	if s.Name != "Gordon" {
		t.Errorf("unexpected result of parsing commented string:\nhave: %s\nwant: %s", s.Name, "Gordon")
	}
}

func TestParseCommentedLine(t *testing.T) {
	var s struct {
		ID   int    `hunk:"id"`
		Name string `hunk:"name"`
	}

	p, err := NewParser(":id :name", &s)
	if err != nil {
		t.Error("unexpected error: " + err.Error())
	}

	if err := p.parseLine("#17 your_name_here\n", &s); err != nil {
		t.Error(err)
	}
	if s.ID != 0 {
		t.Errorf("unexpected result of parsing commented string:\nhave: %d\nwant: %d", s.ID, 0)
	}
}

func TestParseLineWithEscape(t *testing.T) {
	var s struct {
		ID   int    `hunk:"id"`
		Name string `hunk:"name"`
	}

	p, err := NewParser(":id :name", &s)
	if err != nil {
		t.Error("unexpected error: " + err.Error())
	}

	p.SetTokenSeparator('"')
	if err := p.parseLine(`"998" "Gordon Freeman"`, &s); err != nil {
		t.Error(err)
	}

	if s.ID != 998 {
		t.Errorf("unexpected result of parsing commented string:\nhave: %d\nwant: %d", s.ID, 998)
	}
	if s.Name != "Gordon Freeman" {
		t.Errorf("unexpected result of parsing commented string:\nhave: %s\nwant: %s", s.Name, "Gordon Freeman")
	}
}

func TestParseLineReal(t *testing.T) {
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
	f := `:remote_addr - :remote_user :time_local :request :status ` +
		`:body_bytes_sent :http_referer :http_user_agent :bytes_sent :sent_http_content_size ` +
		`:scheme :host :request_time :upstream_response_time :request_length :http_range ` +
		`:responder_name :upstream_cache_status :upstream_response_length :upstream_addr ` +
		`:gcdn_api_client_id :gcdn_api_resource_id :uid_got :uid_set :geoip_country_code ` +
		`:geoip_city :shield_type`

	l := `"62.149.10.131" "-" "-" "[04/Jan/2018:19:15:39 +0000]" "GET /dino.jpg HTTP/1.1" "200" "207402" "" "saelmon" "207957" "-" "https" "di.gcdn.co" "0.000" "-" "88" "-" "[gn]" "HIT" "-" "-" "777" "1337" "-" "-" "UA" "-" "shield_no"`

	p, err := NewParser(f, &s)
	if err != nil {
		t.Errorf("unexpected error: :%s", err)
	}

	p.SetTokenSeparator('"')
	p.SetTimeLayout("time_local", "[02/Jan/2006:15:04:05 -0700]")

	if err = p.ParseLine(l, &s); err != nil {
		t.Error(err)
	}
	if !s.RemoteIP.Equal(net.IPv4(62, 149, 10, 131)) {
		t.Error("reomte_addr did not parsed properly")
	}
	if s.RemoteIPRaw != "62.149.10.131" {
		t.Error("reomte_addr_raw did not init with proper value")
	}
	if s.GcdnTimetamp.IsZero() {
		t.Error("time_local was not parsed properly")
	}
	if s.HTTPRequestRaw != "GET /dino.jpg HTTP/1.1" {
		t.Errorf("request was not parsed properly: %q != %q",
			s.HTTPRequestRaw, "GET /dino.jpg HTTP/1.1")
	}
	if s.HTTPStatus != 200 {
		t.Errorf("status was not parsed properly: %d != %d", s.HTTPStatus, 200)
	}
	if s.Size != 207402 {
		t.Errorf("bytes_sent was not parsed properly: %d != %d", s.Size, 207402)
	}
	if s.RQHeaderUserAgent != "saelmon" {
		t.Errorf("user_agent was not parsed properly: %q != %q",
			s.RQHeaderUserAgent, "salemon")
	}
	if s.ServerToClientBytesSent != 207957 {
		t.Errorf("bytes_sent was not parsed properly: %q != %q",
			s.ServerToClientBytesSent, 207957)
	}
	if s.ServerToClientBytesSentRaw != "207957" {
		t.Errorf("bytes_sent_raw was not set properly: %q != %q",
			s.ServerToClientBytesSentRaw, "207957")
	}
	if s.HTTPScheme != "https" {
		t.Errorf("scheme was not parsed properly: %q != %q", s.HTTPScheme, "https")
	}
	if s.RQHeaderHost != "di.gcdn.co" {
		t.Errorf("host was not parsed properly: %q != %q", s.HTTPScheme, "di.gcdn.co")
	}
	if s.RQLength != 88 {
		t.Errorf("request_length was not parsed properly: %d != %d", s.RQLength, 88)
	}
	if s.GcdnResponderName != "[gn]" {
		t.Errorf("responder_name was not parsed properly: %s != %s", s.GcdnResponderName, "[gn]")
	}
	if s.CacheStatus != "HIT" {
		t.Errorf("cache_status was not parsed properly: %s != %s", s.CacheStatus, "HIT")
	}
	if s.GcdnAPIClientID != 777 {
		t.Errorf("gcdn_api_client_id was not parsed properly: %d != %d", s.GcdnAPIClientID, 777)
	}
	if s.GcdnResourceID != 1337 {
		t.Errorf("gcdn_api_resource_id was not parsed properly: %d != %d", s.GcdnResourceID, 1337)
	}
	if s.GeoIPCountryCode != "UA" {
		t.Errorf("geoip_country_code was not parsed properly: %q != %q", s.GeoIPCountryCode, "UA")
	}
	if s.ShieldUsedRaw != "shield_no" {
		t.Errorf("shield_type was not parsed properly: %q != %q", s.ShieldUsedRaw, "shield_no")
	}

}
