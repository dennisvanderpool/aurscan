package ui

import (
	"strings"
	"testing"

	"github.com/manticore-projects/aurscan/internal/scan"
)

func mal() []scan.Result {
	return []scan.Result{{Pkg: "evil", V: scan.Verdict{Verdict: "MALICIOUS", Confidence: 95, Summary: "bad"}}}
}

func TestGateViaAbortDefault(t *testing.T) {
	var out strings.Builder
	if GateVia(mal(), strings.NewReader("\n"), &out, false) {
		t.Fatal("empty input should abort (return false)")
	}
}

func TestGateViaContinueRequiresINSTALL(t *testing.T) {
	var out strings.Builder
	if GateVia(mal(), strings.NewReader("c\nnope\n"), &out, false) {
		t.Fatal("wrong confirm word should not proceed")
	}
	out.Reset()
	if !GateVia(mal(), strings.NewReader("c\nINSTALL\n"), &out, false) {
		t.Fatal("c then INSTALL should proceed")
	}
}

func TestGateViaOKProceeds(t *testing.T) {
	ok := []scan.Result{{Pkg: "good", V: scan.Verdict{Verdict: "OK", Confidence: 90}}}
	var out strings.Builder
	if !GateVia(ok, strings.NewReader(""), &out, false) {
		t.Fatal("OK should proceed without prompting")
	}
}
