package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/showbiz-io/showbiz/services/api/internal/config"
	"github.com/showbiz-io/showbiz/services/api/internal/database"
	"github.com/showbiz-io/showbiz/services/api/internal/handler"
	"github.com/showbiz-io/showbiz/services/api/internal/middleware"
	"github.com/showbiz-io/showbiz/services/api/internal/provider"
	"github.com/showbiz-io/showbiz/services/api/internal/repository"
	"github.com/showbiz-io/showbiz/services/api/internal/service"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))

	cfg := config.Load()

	db, err := database.Connect(cfg.DSN())
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	userRepo := repository.NewUserRepo(db)
	tokenRepo := repository.NewTokenRepo(db)
	orgRepo := repository.NewOrgRepo(db)
	memberRepo := repository.NewMemberRepo(db)
	connectionRepo := repository.NewConnectionRepo(db)
	resourceRepo := repository.NewResourceRepo(db)

	authService := service.NewAuthService(userRepo, tokenRepo, cfg.JWTSecret)
	authHandler := handler.NewAuthHandler(authService)

	projectRepo := repository.NewProjectRepo(db)

	orgService := service.NewOrgService(orgRepo, memberRepo)
	orgHandler := handler.NewOrgHandler(orgService)

	projectService := service.NewProjectService(projectRepo, orgRepo)
	projectHandler := handler.NewProjectHandler(projectService)

	providerRegistry := provider.NewRegistry()
	providerRegistry.Register("stub", provider.NewStubProvider())
	providerRegistry.Register("fakeprovider", provider.NewFakeProvider(cfg.FakeProviderURL))
	providerHandler := handler.NewProviderHandler(providerRegistry)

	connectionService := service.NewConnectionService(connectionRepo, providerRegistry)
	connectionHandler := handler.NewConnectionHandler(connectionService)

	policyRepo := repository.NewPolicyRepo(db)
	attachmentRepo := repository.NewPolicyAttachmentRepo(db)
	iamService := service.NewIAMService(policyRepo, attachmentRepo, projectRepo, userRepo)
	iamHandler := handler.NewIAMHandler(iamService)

	resourceService := service.NewResourceService(resourceRepo, connectionRepo, providerRegistry)
	resourceHandler := handler.NewResourceHandler(resourceService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.CORS)
	r.Use(middleware.ContentTypeJSON)

	r.Get("/health", handler.Health)

	r.Route("/v1", func(r chi.Router) {
		// Auth routes (public)
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authHandler.Register)
			r.Post("/login", authHandler.Login)
			r.Post("/refresh", authHandler.Refresh)

			r.Group(func(r chi.Router) {
				r.Use(middleware.AuthMiddleware(cfg.JWTSecret))
				r.Get("/me", authHandler.Me)
			})
		})

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(cfg.JWTSecret))

			// Provider routes (read-only)
			r.Route("/providers", func(r chi.Router) {
				r.Get("/", providerHandler.List)
				r.Get("/{id}", providerHandler.Get)
			})

			// Organization routes
			r.Route("/organizations", func(r chi.Router) {
				r.Post("/", orgHandler.Create)
				r.Get("/", orgHandler.ListOrganizations)
				r.Get("/{id}", orgHandler.Get)
				r.Put("/{id}", orgHandler.Update)
				r.Post("/{id}/deactivate", orgHandler.Deactivate)
				r.Post("/{id}/activate", orgHandler.Activate)
				r.Get("/{id}/members", orgHandler.ListMembers)
				r.Post("/{id}/members", orgHandler.AddMember)
				r.Delete("/{id}/members/{email}", orgHandler.RemoveMember)

				// Project routes (nested under org)
				r.Route("/{orgId}/projects", func(r chi.Router) {
					r.Post("/", projectHandler.Create)
					r.Get("/", projectHandler.List)
					r.Get("/{projectId}", projectHandler.Get)
					r.Put("/{projectId}", projectHandler.Update)
					r.Delete("/{projectId}", projectHandler.Delete)
				})

				// Org-scoped IAM policies
				r.Route("/{orgId}/policies", func(r chi.Router) {
					r.Get("/", iamHandler.ListOrgPolicies)
					r.Post("/", iamHandler.CreateOrgPolicy)
					r.Delete("/{policyId}", iamHandler.DeleteOrgPolicy)
				})

				// Project attachments (nested under org/project)
				r.Route("/{orgId}/projects/{projectId}/attachments", func(r chi.Router) {
					r.Get("/", iamHandler.ListProjectAttachments)
					r.Post("/", iamHandler.AttachPolicy)
					r.Delete("/", iamHandler.DetachPolicy)
				})
			})

			// User routes
			r.Route("/users", func(r chi.Router) {
				// TODO: GET /{email}, PUT /{email}, POST /{email}/deactivate, POST /{email}/activate
			})

			// Project routes
			r.Route("/projects", func(r chi.Router) {
				r.Route("/{projectId}/connections", func(r chi.Router) {
					r.Post("/", connectionHandler.Create)
					r.Get("/", connectionHandler.List)
					r.Get("/{connectionId}", connectionHandler.Get)
					r.Put("/{connectionId}", connectionHandler.Update)
					r.Delete("/{connectionId}", connectionHandler.Delete)
				})

				r.Route("/{projectId}/resources", func(r chi.Router) {
					r.Post("/", resourceHandler.Create)
					r.Get("/", resourceHandler.List)
					r.Get("/{resourceId}", resourceHandler.Get)
					r.Put("/{resourceId}", resourceHandler.Update)
					r.Delete("/{resourceId}", resourceHandler.Delete)
				})
			})

			// Global IAM policies (read-only)
			r.Route("/iam/policies", func(r chi.Router) {
				r.Get("/", iamHandler.ListGlobalPolicies)
				r.Get("/{policyId}", iamHandler.GetPolicy)
			})
		})
	})

	addr := ":" + cfg.APIPort
	slog.Info("starting server", "addr", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}
