package configuration

import "github.com/maestre3d/lifetrack-sanbox/pkg/domain/model"

// NewConfiguration reads and returns a kernel configuration
func NewConfiguration() model.Configuration {
	return model.Configuration{
		Version: "0.1.0-alpha",
		Stage:   "prod",
	}
}
