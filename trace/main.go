package main

import (
	"crypto/tls"
	embedutil "git.randomchars.net/FreeNitori/EmbedUtil"
	multiplexer "git.randomchars.net/FreeNitori/Multiplexer"
	"net/http"
	"net/http/httptrace"
	"time"
)

// Setup sets up the plugin and returns route.
//goland:noinspection GoUnusedExportedFunction
func Setup() interface{} {
	return &multiplexer.Route{
		Pattern:       "trace",
		AliasPatterns: []string{},
		Description:   "Perform http trace to discord API.",
		Category:      multiplexer.SystemCategory,
		Handler: func(context *multiplexer.Context) {
			request, err := http.NewRequest("GET", "https://discord.com/api/v8/gateway", nil)
			if !context.HandleError(err) {
				return
			}

			embed := embedutil.New("HTTP Trace", "")
			var requestStart time.Time
			var dnsStart, tlsHandshakeStart, connectStart time.Time
			var dns, tlsHandshake, connect, ttfb time.Duration
			t := &httptrace.ClientTrace{
				DNSStart: func(info httptrace.DNSStartInfo) {
					dnsStart = time.Now()
				},
				DNSDone: func(info httptrace.DNSDoneInfo) {
					dns = time.Since(dnsStart)
				},
				TLSHandshakeStart: func() {
					tlsHandshakeStart = time.Now()
				},
				TLSHandshakeDone: func(state tls.ConnectionState, err error) {
					tlsHandshake = time.Since(tlsHandshakeStart)
				},
				ConnectStart: func(network, addr string) {
					connectStart = time.Now()
				},
				ConnectDone: func(network, addr string, err error) {
					connect = time.Since(connectStart)
				},
				GotFirstResponseByte: func() {
					ttfb = time.Since(requestStart)
				},
			}
			req := request.WithContext(httptrace.WithClientTrace(request.Context(), t))
			requestStart = time.Now()
			var resp *http.Response
			resp, err = http.DefaultTransport.RoundTrip(req)
			if !context.HandleError(err) {
				return
			}
			if !context.HandleError(resp.Body.Close()) {
				return
			}
			embed.AddField("DNS", dns.String(), false)
			embed.AddField("TLS Handshake", tlsHandshake.String(), false)
			embed.AddField("Connect", connect.String(), false)
			embed.AddField("Time to first byte", ttfb.String(), false)
			context.SendEmbed("", embed)
		},
	}
}