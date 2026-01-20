package routes

import (
	"spsc-loaneasy/internal/adapters/http/handlers"
	"spsc-loaneasy/internal/adapters/http/middleware"
	"spsc-loaneasy/internal/adapters/persistence/repositories"
	"spsc-loaneasy/internal/config"
	"spsc-loaneasy/internal/core/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"gorm.io/gorm"
)

// Setup configures all routes for the application
func Setup(app *fiber.App, db *gorm.DB, cfg *config.Config) {
	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	refreshTokenRepo := repositories.NewRefreshTokenRepository(db)
	memberRepo := repositories.NewMemberRepository(db)

	// Phase 4: Master repositories
	loanTypeRepo := repositories.NewLoanTypeRepository(db)
	loanStepRepo := repositories.NewLoanStepRepository(db)
	loanDocRepo := repositories.NewLoanDocRepository(db)
	loanApptRepo := repositories.NewLoanApptRepository(db)

	// Phase 4: Mortgage repositories
	mortgageRepo := repositories.NewMortgageRepository(db)
	transactionRepo := repositories.NewTransactionRepository(db)
	loanDocCurrRepo := repositories.NewLoanDocCurrentRepository(db)
	loanApptCurrRepo := repositories.NewLoanApptCurrentRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo, refreshTokenRepo, memberRepo, cfg)
	userService := services.NewUserService(userRepo, memberRepo)

	// Phase 4: Notification service
	notifyService := services.NewNotificationService()

	// Phase 4: Mortgage service
	mortgageService := services.NewMortgageService(
		mortgageRepo,
		transactionRepo,
		loanTypeRepo,
		loanStepRepo,
		loanDocRepo,
		loanApptRepo,
		loanDocCurrRepo,
		loanApptCurrRepo,
		memberRepo,
		userRepo,
		notifyService,
	)

	// Phase 5: Dashboard service
	dashboardService := services.NewDashboardService(db)

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler()
	authHandler := handlers.NewAuthHandler(authService, cfg)
	userHandler := handlers.NewUserHandler(userService)

	// Phase 4: Handlers
	mortgageHandler := handlers.NewMortgageHandler(mortgageService)
	masterHandler := handlers.NewMasterHandler(loanTypeRepo, loanStepRepo, loanDocRepo, loanApptRepo)

	// Phase 5: Dashboard handler
	dashboardHandler := handlers.NewDashboardHandler(dashboardService)

	// Health check & root routes
	app.Get("/", healthHandler.Root)
	app.Get("/health", healthHandler.HealthCheck)

	// Swagger documentation
	app.Get("/swagger/*", swagger.HandlerDefault)

	// API v1 group
	apiV1 := app.Group("/api/v1")
	setupAPIV1Routes(apiV1, healthHandler, authHandler, userHandler, mortgageHandler, masterHandler, dashboardHandler, cfg)
}

// setupAPIV1Routes configures API v1 routes
func setupAPIV1Routes(
	router fiber.Router,
	healthHandler *handlers.HealthHandler,
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	mortgageHandler *handlers.MortgageHandler,
	masterHandler *handlers.MasterHandler,
	dashboardHandler *handlers.DashboardHandler,
	cfg *config.Config,
) {
	// API Info
	router.Get("/", healthHandler.APIInfo)

	// Auth routes (public)
	authRoutes := router.Group("/auth")
	setupAuthRoutes(authRoutes, authHandler, cfg)

	// User management routes (Admin only)
	userRoutes := router.Group("/users")
	userRoutes.Use(middleware.AuthMiddleware(cfg))
	userRoutes.Use(middleware.AdminOnly())
	setupUserRoutes(userRoutes, userHandler)

	// Profile routes (Authenticated users)
	profileRoutes := router.Group("/profile")
	profileRoutes.Use(middleware.AuthMiddleware(cfg))
	setupProfileRoutes(profileRoutes, userHandler)

	// Phase 4: Mortgage routes (Officer/Admin)
	mortgageRoutes := router.Group("/mortgages")
	mortgageRoutes.Use(middleware.AuthMiddleware(cfg))
	setupMortgageRoutes(mortgageRoutes, mortgageHandler, cfg)

	// Phase 4: Master routes (Admin only)
	masterRoutes := router.Group("/master")
	masterRoutes.Use(middleware.AuthMiddleware(cfg))
	masterRoutes.Use(middleware.AdminOnly())
	setupMasterRoutes(masterRoutes, masterHandler)

	// Phase 5: Dashboard routes
	dashboardRoutes := router.Group("/dashboard")
	dashboardRoutes.Use(middleware.AuthMiddleware(cfg))
	setupDashboardRoutes(dashboardRoutes, dashboardHandler)
}

// setupAuthRoutes configures authentication routes
func setupAuthRoutes(router fiber.Router, handler *handlers.AuthHandler, cfg *config.Config) {
	// Public routes
	router.Post("/register", handler.Register)
	router.Post("/login", handler.Login)
	router.Post("/refresh", handler.RefreshToken)
	router.Post("/logout", handler.Logout)

	// Protected routes
	router.Get("/me", middleware.AuthMiddleware(cfg), handler.Me)
	router.Post("/logout-all", middleware.AuthMiddleware(cfg), handler.LogoutAll)
}

// setupUserRoutes configures user management routes (Admin only)
func setupUserRoutes(router fiber.Router, handler *handlers.UserHandler) {
	router.Get("/", handler.ListUsers)
	router.Get("/:id", handler.GetUser)
	router.Put("/:id", handler.UpdateUser)
	router.Delete("/:id", handler.DeleteUser)
	router.Put("/:id/role", handler.SetUserRole)
}

// setupProfileRoutes configures profile routes (Authenticated)
func setupProfileRoutes(router fiber.Router, handler *handlers.UserHandler) {
	router.Get("/", handler.GetProfile)
	router.Put("/", handler.UpdateProfile)
	router.Put("/password", handler.ChangePassword)
}

// setupMortgageRoutes configures mortgage routes (Phase 4)
func setupMortgageRoutes(router fiber.Router, handler *handlers.MortgageHandler, cfg *config.Config) {
	// Member can view their own mortgages
	router.Get("/my", handler.GetMyMortgages)

	// Officer/Admin routes
	officerRoutes := router.Group("")
	officerRoutes.Use(middleware.OfficerOrAdmin())

	officerRoutes.Post("/", handler.Create)
	officerRoutes.Get("/", handler.List)
	officerRoutes.Get("/:id", handler.GetByID)
	officerRoutes.Get("/:id/history", handler.GetHistory)
	officerRoutes.Get("/:id/docs", handler.GetDocs)
	officerRoutes.Put("/:id/docs", handler.UpdateDoc)
	officerRoutes.Get("/:id/appts", handler.GetAppts)
	officerRoutes.Post("/:id/appts", handler.CreateAppt)
	officerRoutes.Put("/:id/appts/:appt_id/complete", handler.CompleteAppt)
	officerRoutes.Put("/:id/step", handler.ChangeStep)
	officerRoutes.Put("/:id/approve", handler.Approve)
	officerRoutes.Put("/:id/reject", handler.Reject)

	// Admin only
	adminRoutes := router.Group("")
	adminRoutes.Use(middleware.AdminOnly())
	adminRoutes.Put("/:id/officer", handler.ChangeOfficer)
}

// setupMasterRoutes configures master data routes (Admin only) (Phase 4)
func setupMasterRoutes(router fiber.Router, handler *handlers.MasterHandler) {
	// Loan Types
	router.Get("/loan-types", handler.ListLoanTypes)
	router.Get("/loan-types/:id", handler.GetLoanType)
	router.Post("/loan-types", handler.CreateLoanType)
	router.Put("/loan-types/:id", handler.UpdateLoanType)
	router.Delete("/loan-types/:id", handler.DeleteLoanType)

	// Loan Steps
	router.Get("/loan-steps", handler.ListLoanSteps)
	router.Get("/loan-steps/:id", handler.GetLoanStep)
	router.Post("/loan-steps", handler.CreateLoanStep)
	router.Put("/loan-steps/:id", handler.UpdateLoanStep)
	router.Delete("/loan-steps/:id", handler.DeleteLoanStep)

	// Loan Docs
	router.Get("/loan-docs", handler.ListLoanDocs)
	router.Get("/loan-docs/:id", handler.GetLoanDoc)
	router.Post("/loan-docs", handler.CreateLoanDoc)
	router.Put("/loan-docs/:id", handler.UpdateLoanDoc)
	router.Delete("/loan-docs/:id", handler.DeleteLoanDoc)

	// Loan Appts
	router.Get("/loan-appts", handler.ListLoanAppts)
	router.Get("/loan-appts/:id", handler.GetLoanAppt)
	router.Post("/loan-appts", handler.CreateLoanAppt)
	router.Put("/loan-appts/:id", handler.UpdateLoanAppt)
	router.Delete("/loan-appts/:id", handler.DeleteLoanAppt)
}

// setupDashboardRoutes configures dashboard routes (Phase 5)
func setupDashboardRoutes(router fiber.Router, handler *handlers.DashboardHandler) {
	// Auto-detect role dashboard (All authenticated users)
	router.Get("/", handler.GetMyDashboard)

	// User dashboard (All authenticated users)
	router.Get("/user", handler.GetUserDashboard)

	// Officer dashboard (Officer/Admin only)
	router.Get("/officer", middleware.OfficerOrAdmin(), handler.GetOfficerDashboard)

	// Admin dashboard (Admin only)
	router.Get("/admin", middleware.AdminOnly(), handler.GetAdminDashboard)
}
