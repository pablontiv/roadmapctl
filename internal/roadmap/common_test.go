package roadmap

import (
	"github.com/pablontiv/roadmapctl/internal/config"
)

func testDefaultConfig() *config.Config {
	return &config.Config{
		Fields: config.FieldsConfig{
			Lifecycle:      "estado",
			RecordType:     "tipo",
			TaskValue:      "task",
			OutcomeValue:   "outcome",
			DisplayName:    "titulo",
			DependencyLink: "blocked_by",
		},
	}
}
