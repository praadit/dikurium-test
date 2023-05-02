package graph

import (
	"github.com/praadit/dikurium-test/pkg/controller"
	loader "github.com/praadit/dikurium-test/pkg/graph/loaders"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Controller *controller.Controller
	Loader     *loader.Loader
}
