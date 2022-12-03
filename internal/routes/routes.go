package routes

import (
	ProductHandler "github.com/SapanPatibandha/golangWithDynamoDB/internal/handlers/product"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-ghi/chi/middleware"
)

type Router struct {
	config *Config
	router *chi.Mux
}

func NewRouter() *Router {
	return &Router{
		config: NewRouter().SetTimeout(serviceConfig.GetConfig().Timeout),
		router: chi.NewRouter(),
	}
}

func (r *Router) SetRouters(repository adapter.Interface) *chi.Mux {
	r.setConfigsRouters()
	r.RouterHealth(repository)
	r.RouterProduct(repository)

	return r.router
}

func (r *Router) setConfigsRouters() {
	r.EnableCORS()
	r.EnableLogger()
	r.EnableTimeout()
	r.EnableRecover()
	r.EnableRequestID()
	r.EnableRealIP()
}

func (r *Router) RouterHealth(repository adapter.Interface) {
	handler := HealthHandler.newHandler(repository)

	r.router.Route("/health", func(route chi.Router) {
		route.Post("/", handler.Post)
		route.Get("/", handler.Get)
		route.Put("/", handler.Put)
		route.Delete("/", handler.Delete)
		route.Options("/", handler.Options)
	})
}

func (r *Router) RouterProduct(repository adaptor.Interface) {

	handler := ProductHandler.NewHandler(repository)

	r.router.route("/product", func(route chi.Router) {
		route.Post("/", handler.Post)
		route.Get("/", handler.Get)
		route.Put("/{ID}", handler.Put)
		route.Delete("/{ID}", handler.Delete)
		route.Options("/", handler.Options)
	})
}

func (r *Router) EnableLogger() *Router {
	r.router.Use(middleware.Logger)
	return r
}

func (r *Router) EnableTimeout() *Router {
	r.router.Use(middleware.Timeout(r.config.GetTimeout()))
	return r

}

func (r *Router) EnableCORS() *Router {
	r.router.Use(r.config.Corse)
	return r
}

func (r *Router) EnableRecover() *Router {
	r.router.Use(middleware.Recoverer)
	return r
}

func (r *Router) EnableRequestID() *Router {
	r.router.Use(middleware.RequestID)
	return r
}

func (r *Router) EnableRealIP() *Router {
	r.router.Use(middleware.RealIP)
	return r
}
