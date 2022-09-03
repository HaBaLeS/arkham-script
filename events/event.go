package events

import (
	"arkham-script/actions"
	"fmt"
	"gopkg.in/yaml.v3"
)

type EngineEvent interface {
	Execute()
}

type Events struct {
	Event []ScriptEvent `yaml:"event"`
}

type ScriptEvent struct {
	Trigger   string           `yaml:"trigger"`
	Type      string           `yaml:"type"`
	Actions   []actions.Action `yaml:"actions"`
	eventImpl EngineEvent
}

func (se *ScriptEvent) UnmarshalYAML(value *yaml.Node) error {

	for k, v := range value.Content {
		fmt.Printf("%d, %s\n", k, v)
	}

	return nil
}
