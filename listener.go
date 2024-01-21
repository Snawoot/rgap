package rgap

import (
	"errors"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type GroupConfig struct {
	ID             uint64
	PSK            *PSK
	Expire         time.Duration
	ClockSkew      time.Duration `yaml:"clock_skew"`
	ReadinessDelay time.Duration `yaml:"readiness_delay"`
}

type OutputConfig struct {
	Kind string
	Spec yaml.Node
}

type ListenerConfig struct {
	Listen  []string
	Groups  []GroupConfig
	Outputs []OutputConfig
}

type Listener struct {
	cfg *ListenerConfig
}

func NewListener(cfg *ListenerConfig) (*Listener, error) {
	enc := yaml.NewEncoder(os.Stdout)
	if err := enc.Encode(cfg); err != nil {
		panic(err)
	}
	return nil, errors.New("not implemented")
}