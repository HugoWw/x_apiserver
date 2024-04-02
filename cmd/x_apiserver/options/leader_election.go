package options

import (
	"github.com/google/uuid"
	"github.com/spf13/pflag"
	"time"
)

type LeaderElectionOptions struct {
	ID            string
	LockName      string
	RenewDeadline time.Duration
	LeaseDuration time.Duration
	RetryPeriod   time.Duration
}

func NewLeaderElectionOptions() *LeaderElectionOptions {
	return &LeaderElectionOptions{}
}

func (l *LeaderElectionOptions) Valid() []error {
	return nil
}

func (l *LeaderElectionOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&l.ID, "id", uuid.New().String(), "The holder identity name")
}
