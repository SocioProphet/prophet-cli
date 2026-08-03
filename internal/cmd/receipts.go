package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/socioprophet/prophet-cli/internal/receipt"
)

// emitReceipt writes a façade-local receipt for a delegated action when the
// operator passed --receipt. It never alters the command's own exit behavior:
// a receipt-write failure is surfaced on stderr but does not mask the delegated
// result. Receipts are convenience artifacts, not the estate ProofArtifact spine.
func emitReceipt(command, delegate, status string, args []string, started time.Time, execErr error) {
	if flags.Receipt == "" {
		return
	}
	r := receipt.New(command, delegate, status, args, started, time.Now(), execErr)
	path, err := receipt.Write(flags.Receipt, r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "prophet: receipt not written: %v\n", err)
		return
	}
	if !flags.Quiet {
		fmt.Fprintf(os.Stderr, "prophet: receipt %s\n", path)
	}
}
