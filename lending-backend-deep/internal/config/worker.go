package config

import "time"

// EnvWorker holds environment variables that configure the background workers.
type EnvWorker struct {
	// BusinessTimezone is an IANA name. It has no default: a wrong zone would
	// silently shift age and ID expiry results around midnight.
	BusinessTimezone             string `env:"BUSINESS_TIMEZONE,required,notEmpty"`
	ValidatorPollIntervalSeconds int    `env:"VALIDATOR_POLL_INTERVAL_SECONDS" envDefault:"1"`
	NotifierPollIntervalSeconds  int    `env:"NOTIFIER_POLL_INTERVAL_SECONDS"  envDefault:"5"`
}

// ValidatorPollInterval is the wait between polls when no application is due.
func (e *EnvWorker) ValidatorPollInterval() time.Duration {
	return time.Duration(e.ValidatorPollIntervalSeconds) * time.Second
}

// NotifierPollInterval is the wait between polls when no notification is due.
func (e *EnvWorker) NotifierPollInterval() time.Duration {
	return time.Duration(e.NotifierPollIntervalSeconds) * time.Second
}
