package browserworker

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"slices"
	"strings"
	"sync/atomic"
	"testing"
)

// The browser's every way out is the proxy: the proxy flag, no bypass for
// loopback, and WebRTC's UDP confined to it (agentgabo's review of
// kiff/planque#89: a STUN or TURN candidate would otherwise go around an
// HTTP proxy to any address).
func TestLaunchArgsSendEverythingThroughTheProxy(t *testing.T) {
	args := launchArgs("/tmp/profile", "127.0.0.1:9999")
	for _, want := range []string{
		"--proxy-server=http://127.0.0.1:9999",
		"--proxy-bypass-list=<-loopback>",
		"--force-webrtc-ip-handling-policy=disable_non_proxied_udp",
		"--site-per-process",
		"--remote-debugging-pipe",
		"--user-data-dir=/tmp/profile",
		"--headless=new",
	} {
		if !slices.Contains(args, want) {
			t.Errorf("the browser isn't launched with %s: %v", want, args)
		}
	}
}

// Inside the real browser: a page that tries WebRTC, with a STUN server
// at a private address and a remote candidate on loopback (a UDP trap),
// gets no peer connection at all, and the browser sends the trap nothing:
// no UDP leaves the browser. (With the policy flag alone, the browser
// sent the trap 57 STUN packets; hence the script in attach.)
func TestWorkerWebRTCProbesNoAddress(t *testing.T) {
	w, s := worker(t), shop(t)
	trap, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1)})
	if err != nil {
		t.Fatal(err)
	}
	defer trap.Close()
	var packets atomic.Int32
	go func() {
		buf := make([]byte, 1500)
		for {
			if _, _, err := trap.ReadFromUDP(buf); err != nil {
				return
			}
			packets.Add(1)
		}
	}()
	port := trap.LocalAddr().(*net.UDPAddr).Port
	ctx := context.Background()
	sess, err := w.Start(ctx, startFor(s))
	if err != nil {
		t.Fatal(err)
	}
	defer sess.Close(ctx)
	if _, err := sess.Open(ctx, s.URL+"/item/pen"); err != nil {
		t.Fatal(err)
	}
	probe := fmt.Sprintf(`new Promise(done => {
		const found = [];
		let pc;
		try { pc = new RTCPeerConnection({iceServers: [{urls: ['stun:10.0.0.7:3478', 'stun:169.254.169.254:3478']}]}); }
		catch (e) { done(JSON.stringify(['error: ' + e])); return; }
		pc.onicecandidate = e => { if (e.candidate) found.push(e.candidate.candidate); };
		pc.createDataChannel('probe');
		pc.createOffer().then(o => pc.setLocalDescription(o)).then(() => {
			// A remote peer on loopback: the browser would send STUN checks to it.
			const sdp = 'v=0\r\no=- 1 1 IN IP4 127.0.0.1\r\ns=-\r\nt=0 0\r\na=group:BUNDLE 0\r\nm=application 9 UDP/DTLS/SCTP webrtc-datachannel\r\nc=IN IP4 127.0.0.1\r\na=mid:0\r\na=ice-ufrag:trap\r\na=ice-pwd:trapxtrapxtrapxtrapxtrap\r\na=fingerprint:sha-256 ' + 'AB:'.repeat(31) + 'AB\r\na=setup:passive\r\na=sctp-port:5000\r\n';
			return pc.setRemoteDescription({type: 'answer', sdp});
		}).then(() => pc.addIceCandidate({candidate: 'candidate:1 1 udp 2130706431 127.0.0.1 %d typ host', sdpMid: '0'}))
		.catch(e => found.push('error: ' + e));
		setTimeout(() => done(JSON.stringify(found)), 3500);
	})`, port)
	var ev struct {
		Result struct {
			Value string `json:"value"`
		} `json:"result"`
	}
	in := sess.(*session)
	if err := in.c.call(ctx, in.sid, "Runtime.evaluate", map[string]any{"expression": probe, "returnByValue": true, "awaitPromise": true}, &ev); err != nil {
		t.Fatal(err)
	}
	var found []string
	if err := json.Unmarshal([]byte(ev.Result.Value), &found); err != nil {
		t.Fatalf("the probe's answer: %q", ev.Result.Value)
	}
	// The page gets no peer connection; nothing else is gathered.
	if len(found) != 1 || !strings.HasPrefix(found[0], "error:") || !strings.Contains(found[0], "RTCPeerConnection") {
		t.Errorf("WebRTC answered the page: %v", found)
	}
	if n := packets.Load(); n != 0 {
		t.Errorf("the browser sent %d UDP packets to the trap", n)
	}
}

// A cross-site iframe is its own target under site isolation: WebRTC is
// gone there too, before its script runs, and the trap gets nothing. The
// page on 127.0.0.1 embeds a frame from localhost, another site to the
// browser. (agentgabo's review of kiff/planque#89, round 2.)
func TestWorkerWebRTCIsGoneInACrossSiteFrame(t *testing.T) {
	w, s := worker(t), shop(t)
	trap, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1)})
	if err != nil {
		t.Fatal(err)
	}
	defer trap.Close()
	var packets atomic.Int32
	go func() {
		buf := make([]byte, 1500)
		for {
			if _, _, err := trap.ReadFromUDP(buf); err != nil {
				return
			}
			packets.Add(1)
		}
	}()
	port := trap.LocalAddr().(*net.UDPAddr).Port
	ctx := context.Background()
	sess, err := w.Start(ctx, startFor(s))
	if err != nil {
		t.Fatal(err)
	}
	defer sess.Close(ctx)
	inner := strings.Replace(s.URL, "127.0.0.1", "localhost", 1) + "/inner"
	if _, err := sess.Open(ctx, s.URL+"/frame?src="+inner); err != nil {
		t.Fatal(err)
	}
	in := sess.(*session)
	// The frame's own target, attached by the browser: run the probe in it.
	var targets struct {
		TargetInfos []struct {
			TargetID string `json:"targetId"`
			Type     string `json:"type"`
			URL      string `json:"url"`
		} `json:"targetInfos"`
	}
	if err := in.c.call(ctx, "", "Target.getTargets", nil, &targets); err != nil {
		t.Fatal(err)
	}
	frameTarget := ""
	for _, ti := range targets.TargetInfos {
		if ti.Type == "iframe" && strings.Contains(ti.URL, "localhost") {
			frameTarget = ti.TargetID
		}
	}
	if frameTarget == "" {
		t.Fatalf("no cross-site frame target: %+v (site isolation off?)", targets.TargetInfos)
	}
	var attached struct {
		SessionID string `json:"sessionId"`
	}
	if err := in.c.call(ctx, "", "Target.attachToTarget", map[string]any{"targetId": frameTarget, "flatten": true}, &attached); err != nil {
		t.Fatal(err)
	}
	probe := fmt.Sprintf(`new Promise(done => {
		const found = [];
		let pc;
		try { pc = new RTCPeerConnection({iceServers: [{urls: ['stun:10.0.0.7:3478']}]}); }
		catch (e) { done(JSON.stringify(['error: ' + e])); return; }
		pc.onicecandidate = e => { if (e.candidate) found.push(e.candidate.candidate); };
		pc.createDataChannel('probe');
		pc.createOffer().then(o => pc.setLocalDescription(o)).then(() => {
			const sdp = 'v=0\r\no=- 1 1 IN IP4 127.0.0.1\r\ns=-\r\nt=0 0\r\na=group:BUNDLE 0\r\nm=application 9 UDP/DTLS/SCTP webrtc-datachannel\r\nc=IN IP4 127.0.0.1\r\na=mid:0\r\na=ice-ufrag:trap\r\na=ice-pwd:trapxtrapxtrapxtrapxtrap\r\na=fingerprint:sha-256 ' + 'AB:'.repeat(31) + 'AB\r\na=setup:passive\r\na=sctp-port:5000\r\n';
			return pc.setRemoteDescription({type: 'answer', sdp});
		}).then(() => pc.addIceCandidate({candidate: 'candidate:1 1 udp 2130706431 127.0.0.1 %d typ host', sdpMid: '0'}))
		.catch(e => found.push('error: ' + e));
		setTimeout(() => done(JSON.stringify(found)), 3500);
	})`, port)
	var ev struct {
		Result struct {
			Value string `json:"value"`
		} `json:"result"`
	}
	if err := in.c.call(ctx, attached.SessionID, "Runtime.evaluate", map[string]any{"expression": probe, "returnByValue": true, "awaitPromise": true}, &ev); err != nil {
		t.Fatal(err)
	}
	var found []string
	if err := json.Unmarshal([]byte(ev.Result.Value), &found); err != nil {
		t.Fatalf("the frame's answer: %q", ev.Result.Value)
	}
	if len(found) != 1 || !strings.HasPrefix(found[0], "error:") || !strings.Contains(found[0], "RTCPeerConnection") {
		t.Errorf("WebRTC answered the frame: %v", found)
	}
	if n := packets.Load(); n != 0 {
		t.Errorf("the frame sent %d UDP packets to the trap", n)
	}
}
