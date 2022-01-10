package routes

import (
	//"context"
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/syedomair/ex-paygate-approve/routes/approve"
	"github.com/syedomair/ex-paygate-lib/lib/container"
	log "github.com/syedomair/ex-paygate-lib/lib/tools/logger"
	"github.com/syedomair/ex-paygate-lib/lib/tools/request"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

// NewRouter comment
func NewRouter(c container.Container) *chi.Mux {
	router := chi.NewRouter()
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"POST"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	router.Use(cors.Handler)

	router.Route("/v1", func(r chi.Router) {
		r.Mount("/", routerSetup(
			c.Db(),
			c.Logger(),
			c.SigningKey(),
		))
	})

	return router
}

// Route Public
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
	Access      string
}

// Routes List Public
type Routes []Route

func routerSetup(db *gorm.DB, logger log.Logger, signingKey string) *chi.Mux {

	repoApprove := approve.NewPostgresRepository(db, logger)
	payApprove := approve.NewPaymentService(logger)

	router := chi.NewRouter()
	approveController := &approve.Controller{
		Logger: logger,
		Repo:   repoApprove,
		Pay:    payApprove,
	}

	var routes = Routes{
		Route{
			"Ping",
			"GET",
			"/ping",
			approveController.Ping,
			"public",
		},
		Route{
			"Approve",
			"POST",
			"/authorize",
			approveController.ApproveAction,
			"public",
		},
	}

	for _, route := range routes {
		var handler http.Handler
		strKey := uuid.NewString()
		handler = route.HandlerFunc
		handler = loggerMiddleware(strKey, logger, handler, route.Method, route.Pattern)
		router.Method(route.Method, route.Pattern, handler)
	}
	return router
}

func loggerMiddleware(strKey string, logger log.Logger, handle http.Handler, method string, pattern string) http.Handler {
	methodName := "loggerMiddleware"
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testNum := r.Header.Get("Testnum")
		ctx := r.Context()
		ctx = context.WithValue(ctx, request.ContextKeyRequestID, strKey)
		r = r.WithContext(ctx)
		logger.Debug(strKey, "******* %v ****%v:%v**********", testNum, method, pattern)
		logger.Debug(strKey, "M:%v ", methodName)
		handle.ServeHTTP(w, r)
	})
}
