package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/storage/redis"
	"github.com/hibiken/asynq"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/mirusky-dev/challenge-18/core"
	"github.com/mirusky-dev/challenge-18/core/env"
	"github.com/mirusky-dev/challenge-18/core/mailer"
	"github.com/mirusky-dev/challenge-18/repositories"
	"github.com/mirusky-dev/challenge-18/router/middlewares"
	"github.com/mirusky-dev/challenge-18/services"
)

var (
	refreshTokenKey = "gobp.refresh-token"
)

func NewController(
	userService services.IUserService,
	taskService services.ITaskService,
	authService services.IAuthService,
	accountService services.IAccountService,
	tokenService services.ITokenService,

) Controller {
	return Controller{
		userService:    userService,
		taskService:    taskService,
		authService:    authService,
		accountService: accountService,
		tokenService:   tokenService,
	}
}

type Controller struct {
	// INFO: Services and dependencies goes here

	userService    services.IUserService
	taskService    services.ITaskService
	authService    services.IAuthService
	accountService services.IAccountService
	tokenService   services.ITokenService
}

func Setup(config env.Config, backgroundClient *asynq.Client) *fiber.App {

	db, err := gorm.Open(mysql.Open(config.DatabaseURL), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	mailService, _ := mailer.NewNoopMailer(config)

	// Dependencies Setup
	argonPasswordHasher := core.NewArgon2IDPasswordHasher()

	refreshTokenStorage := redis.New(redis.Config{URL: config.RedisURL + "/1"})
	tokenRevokationStorage := redis.New(redis.Config{URL: config.RedisURL + "/2"})
	emailVerificationStorage := redis.New(redis.Config{URL: config.RedisURL + "/3"})
	passwordResetStorage := redis.New(redis.Config{URL: config.RedisURL + "/4"})

	// Repositories Setup
	userRepository := repositories.NewUserRepository(db)
	taskRepository := repositories.NewTaskRepository(db)

	// Services Setup
	userService := services.NewUserService(userRepository, argonPasswordHasher)
	taskService := services.NewTaskService(taskRepository, backgroundClient)
	tokenService := services.NewTokenService(config, refreshTokenStorage, tokenRevokationStorage, userRepository)
	accountService := services.NewAccountService(config, mailService, argonPasswordHasher, emailVerificationStorage, userRepository)
	authService := services.NewAuthService(config, mailService, argonPasswordHasher, userRepository, passwordResetStorage, tokenService)

	ctrl := NewController(
		userService,
		taskService,
		authService,
		accountService,
		tokenService,
	)

	// New App with custom error handling
	app := fiber.New(fiber.Config{
		DisableStartupMessage: !config.EnableStartupMessage,
		EnablePrintRoutes:     config.EnablePrintRoutes,
		ErrorHandler:          core.ErrorHandler,
	})

	// Default middlewares
	app.Use(logger.New())
	app.Use(recover.New(recover.Config{EnableStackTrace: config.EnableStackTrace}))

	corsConfig := cors.ConfigDefault
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, Cache-Control"

	app.Use(cors.New(corsConfig))

	// App middlewares
	app.Use(middlewares.Context())

	// Routes
	api := app.Group("/api")

	v1 := api.Group("/v1")
	debug := api.Group("/debug")

	debug.Get("/ping", ctrl.pong("Pong!"))
	debug.Get("/empty", ctrl.empty)
	debug.Get("/friendly-error", ctrl.friendlyError)
	debug.Get("/context", ctrl.context)

	v1.Post("/auth/login", ctrl.login)
	v1.Post("/auth/refresh-token", ctrl.refreshToken)
	// v1.Post("/auth/register", ctrl.register)
	v1.Post("/auth/reset-password", ctrl.sendResetPasswordLink)
	v1.Post("/auth/reset-password/:id", ctrl.verifyResetPassword)
	v1.Post("/accounts/email-verification/:id", ctrl.verifyEmailVerification)

	{
		// INFO: Everything under this will need auth
		api.Use(middlewares.JWT(config))
		api.Use(middlewares.User(ctrl.tokenService))

		debug.Get("/authenticated-ping", ctrl.pong("Authenticated Pong!"))

		v1.Get("/auth/logout", ctrl.logout)
		v1.Post("/auth/logout", ctrl.logout)

		v1.Get("/accounts/me", ctrl.context)
		v1.Post("/accounts/change-password", ctrl.changePassword)
		v1.Post("/accounts/email-verification", ctrl.sendEmailVerificationLink)

		v1.Post("/users", middlewares.Authorize(core.HasRole("admin")), ctrl.createUser)
		v1.Get("/users", middlewares.Authorize(core.HasRole("admin")), ctrl.getAllUsers)
		v1.Get("/users/:id", middlewares.Authorize(core.HasRole("admin")), ctrl.getUserByID)
		v1.Put("/users/:id", middlewares.Authorize(core.HasRole("admin")), ctrl.updateUser)
		v1.Delete("/users/:id", middlewares.Authorize(core.HasRole("admin")), ctrl.deleteUser)

		v1.Post("/tasks", middlewares.Authorize(core.HasRole("tech")), ctrl.createTask)
		v1.Get("/tasks", middlewares.Authorize(core.HasRoles{Roles: []string{"tech", "manager"}}), ctrl.getAllTasks)
		v1.Get("/tasks/:id", middlewares.Authorize(core.HasRoles{Roles: []string{"tech", "manager"}}), ctrl.getTaskByID)
		v1.Put("/tasks/:id", middlewares.Authorize(core.HasRole("tech")), ctrl.updateTask)
		v1.Delete("/tasks/:id", middlewares.Authorize(core.HasRole("manager")), ctrl.deleteTask)
	}

	return app
}
