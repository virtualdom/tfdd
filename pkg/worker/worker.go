package worker

import (
	"fmt"
	"os"
	"strings"

	"github.com/virtualdom/tfdd/pkg/terraform"
)

const NO_CHANGE = `No changes. Infrastructure is up-to-date.

This means that Terraform did not detect any differences between your
configuration and real physical resources that exist. As a result, no
actions need to be performed.`

func deposit(ch chan struct{}) {
	ch <- struct{}{}
}

// Process takes a path to a Terraform service and a channel through which to
// send a completion message. It runs `terraform init` and `terraform plan` and
// outputs a message if the given Terraform service doesn't contain the typical
// "No changes..." string.
func Process(path string, ch chan struct{}) {
	defer deposit(ch)

	tf := terraform.New(path)

	err := tf.Init()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to init %s: %+v", path, err)
		return
	}

	out, errlog, err := tf.Plan()
	if err != nil || errlog != "" {
		fmt.Fprintf(os.Stderr, "failed to plan %s: %+v\n%s", path, err, errlog)
		return
	}

	// todo: update this to parse TF output and list TF resource names that are
	// being changed
	// todo: delegate this to a `reporting` class that can optionally save these
	// to a database to avoid frequent reporting (if this is being run in a cron)
	if !strings.Contains(out, NO_CHANGE) {
		fmt.Printf("drift detected in %s\n", path)
	}
}
