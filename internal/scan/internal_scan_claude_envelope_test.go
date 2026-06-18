package scan

import "testing"

// Array/streaming shape from Claude Code CLI v2.1.x (issue #17).
const cliArray = `[{"type":"system","subtype":"init","session_id":"x"},` +
	`{"type":"system","subtype":"api_retry","error_status":401,"error":"authentication_failed"},` +
	`{"type":"assistant","message":{"role":"assistant"}},` +
	`{"type":"result","subtype":"success","is_error":false,` +
	`"result":"{\"verdict\":\"OK\",\"confidence\":88,\"summary\":\"fine\",\"findings\":[]}",` +
	`"total_cost_usd":0.05,"usage":{"input_tokens":2,"output_tokens":430,"cache_read_input_tokens":12629}}]`

const cliObject = `{"type":"result","result":"{\"verdict\":\"MALICIOUS\",\"confidence\":90,\"summary\":\"bad\",\"findings\":[]}","total_cost_usd":0.01,"usage":{"input_tokens":10,"output_tokens":20}}`

func TestParseClaudeArrayShape(t *testing.T) {
	text, u, ok := parseClaudeEnvelope(cliArray, 100)
	if !ok {
		t.Fatal("array envelope not parsed (issue #17 regression)")
	}
	v := parseVerdict(text)
	if v.Verdict != "OK" || v.Confidence != 88 {
		t.Fatalf("got verdict=%q conf=%v", v.Verdict, v.Confidence)
	}
	if !u.HaveCost || u.In != 2+12629 {
		t.Fatalf("usage wrong: %+v", u)
	}
}

func TestParseClaudeObjectShape(t *testing.T) {
	text, _, ok := parseClaudeEnvelope(cliObject, 100)
	if !ok {
		t.Fatal("object envelope not parsed")
	}
	if v := parseVerdict(text); v.Verdict != "MALICIOUS" {
		t.Fatalf("got %q", v.Verdict)
	}
}

func TestParseClaudeGarbage(t *testing.T) {
	if _, _, ok := parseClaudeEnvelope("not json at all", 100); ok {
		t.Fatal("garbage should not parse as envelope")
	}
}
