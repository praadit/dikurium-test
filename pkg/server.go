package pkg

import (
	"context"
	"fmt"
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/praadit/dikurium-test/pkg/controller"
	generated "github.com/praadit/dikurium-test/pkg/graph/generated"
	loader "github.com/praadit/dikurium-test/pkg/graph/loaders"
	graph "github.com/praadit/dikurium-test/pkg/graph/resolvers"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Server struct {
	controller *controller.Controller
	loader     *loader.Loader
}

func NewServer(db *gorm.DB, logger *zap.Logger) *Server {
	loader := loader.NewLoader(db, logger)
	controller := controller.InitController(db, logger)
	return &Server{
		controller: controller,
		loader:     loader,
	}
}

// Defining the Playground handler
func (s *Server) PlaygroundHandler() http.Handler {
	h := playground.Handler("GraphQL", "/query")

	return h
}

func (s *Server) GraphqlHandler() http.Handler {
	c := s.GraphConfig()
	h := handler.NewDefaultServer(generated.NewExecutableSchema(c))
	return h
}

func (s *Server) GraphConfig() generated.Config {
	c := generated.Config{
		Resolvers: &graph.Resolver{
			Controller: s.controller,
			Loader:     s.loader,
		},
	}
	c.Directives.Auth = func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
		claims := ctx.Value("claims")
		if claims == nil {
			// block calling the next resolver
			return nil, fmt.Errorf("Access denied")
		}

		fmt.Println(claims)

		// or let it pass through
		return next(ctx)
	}

	return c
}
