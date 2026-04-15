package vocab

import (
	"context"
	"testing"
)

func TestValidateReportsValidatorPath(t *testing.T) {
	resp, err := Validate(context.Background(), []string{"graphs/demo.ttl"})
	if err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
	if got := resp["validator"]; got != "policy/tools/validate_all.py" {
		t.Fatalf("unexpected validator path: %v", got)
	}
	if got := resp["operation"]; got != "validate" {
		t.Fatalf("unexpected operation: %v", got)
	}
}
