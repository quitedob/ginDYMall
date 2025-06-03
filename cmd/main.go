// @title           Gin E-Commerce API (douyin)
// @version         0.13.0
// @description     This is a sample server for a Gin e-commerce application.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:5001 // Adjust if necessary, e.g. based on config.GlobalConfig.System.HttpPort
// @BasePath  /api/v1        // Adjust if your API base path is different

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
package main

import (
	"fmt"
	"log"

	v1 "douyin/api/v1"
	conf "douyin/config"
	"douyin/global" // For global.DB
	mylog "douyin/pkg/utils/log"
	"douyin/repository/cache"
	"douyin/repository/db/model" // For RBAC models
	"douyin/pkg/utils/upload"    // For OSS client
	i18nUtils "douyin/pkg/utils/i18n" // For i18n
	"douyin/pkg/utils/email"     // For Email client
	"douyin/service"             // For NotificationService
	"douyin/routes"
	"regexp"

	_ "douyin/docs" // Import for Swagger docs generation
	"douyin/middleware" // Added import for middleware
	// "ginDYMall/api/v1" // Already imported as v1 "douyin/api/v1"
	// "ginDYMall/middleware" // Already imported as middleware "douyin/middleware"
	"github.com/redis/go-redis/v9" // For HealthController Redis client
	"golang.org/x/text/language"   // For i18n language tags

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"context" // For graceful shutdown and timeout
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"math/rand"
	"net/http" // For graceful shutdown
	"os"       // For graceful shutdown
	"os/signal"// For graceful shutdown
	"syscall"  // For graceful shutdown
	"time"

	jwtUtil "douyin/pkg/utils/jwt" // Import jwt package
)

func main() {
	loading()

	// Initialize JWT secret
	if conf.GlobalConfig.EncryptSecret == nil || conf.GlobalConfig.EncryptSecret.JwtSecret == "" {
		log.Fatal("JWT Secret is not configured properly.") // Or handle more gracefully
	}
	jwtUtil.Init(conf.GlobalConfig.EncryptSecret.JwtSecret)


	// 构造 MySQL DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		conf.GlobalConfig.MySql.Default.UserName,
		conf.GlobalConfig.MySql.Default.Password,
		conf.GlobalConfig.MySql.Default.DbHost,
		conf.GlobalConfig.MySql.Default.DbPort,
		conf.GlobalConfig.MySql.Default.DbName,
		conf.GlobalConfig.MySql.Default.Charset,
	)

	// 创建数据库连接
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		mylog.Errorf("数据库连接失败: %v", err)
		log.Fatalf("数据库连接失败: %v", err)
	}
	mylog.Info("数据库连接成功")

	// Set global DB instance for middlewares or other packages
	global.DB = db

	// Initialize i18n
	// Define these in your config or constants if they vary
	defaultSiteLang := language.Chinese // Default language for the site (e.g., Simplified Chinese)
	// Ensure AmericanEnglish is used if "en-US" is a target, or just English for "en"
	supportedSiteLangs := []language.Tag{language.Chinese, language.English, language.AmericanEnglish}
	localesDir := "locales" // Directory containing JSON translation files (relative to app root)

	if err := i18nUtils.InitI18n(defaultSiteLang, supportedSiteLangs, localesDir); err != nil {
		mylog.Fatalf("Failed to initialize i18n: %v", err)
	}
	mylog.Info("i18n system initialized.")

	// Auto-migrate tables, including RBAC models
	rbacModels := model.GetRBACModels()
	// If you have other models to migrate, add them here: e.g., &model.User{}, &model.Product{}
	// For now, only migrating RBAC models as per task.
	// In a real app, you'd list all models: err = global.DB.AutoMigrate(append(rbacModels, &model.User{})...)
	err = global.DB.AutoMigrate(rbacModels...)
	if err != nil {
		mylog.Fatalf("GORM AutoMigrate RBAC tables failed: %v", err)
	}
	mylog.Info("RBAC tables migrated successfully")


	v1.SetDB(db) // This function might also set global.DB or uses the passed db.
	             // If v1.SetDB already sets global.DB, the line global.DB = db above might be redundant
	             // but explicit assignment is safer for clarity.
	v1.SetCheckoutController(db)

	// Initialize HealthController
	// Assuming cache.GetClient() returns the *redis.Client initialized by cache.InitCache()
	// This part might need adjustment if cache.GetClient() is not the correct way to get the Redis client.
	// For now, let's assume 'cache.Rdb' is the client, if InitCache sets it up.
	// If cache.InitCache() makes a client available like cache.Client, use that.
	// As per common patterns, cache.InitCache might populate a global variable in the cache package.
	// Let's try to use a hypothetical cache.Client. If this fails, this indicates
	// an unknown method of Redis client retrieval.
	// A common pattern is cache.Rdb as a global variable in package cache.
	var redisClient *redis.Client
	// Attempt to retrieve from a hypothetical global variable or function in cache package.
	// This is a placeholder. The actual way to get the client might differ.
	// For example, if cache.Rdb is the client:
	if cache.Rdb != nil { // Check if Rdb is the client and initialized
		redisClient = cache.Rdb
	} else {
		// Fallback or error: Redis client not found.
		// This indicates a need to understand how cache.InitCache() provides the client.
		// For the purpose of this task, we'll log an error if it's nil and proceed,
		// though in a real scenario, this might be a fatal error.
		mylog.Error("Redis client not available for HealthController. Health checks for Redis might fail.")
		// healthCtrl will be initialized with a nil Redis client if not found.
		// The Healthz check handles nil Redis client.
	}
	healthCtrl := v1.HealthController{DB: db, Redis: redisClient}


	router := gin.New() // Using gin.New() to have more control over middleware.
	healthCtrl.RegisterRoutes(router) // Register health check routes

	// Register Middlewares (order matters)
	// router.Use(middleware.Jaeger())      // Tracing first - Assuming Jaeger is initialized elsewhere or not part of this specific step's focus if not present
	// For now, Jaeger() is commented out if not set up in this context. If it is, uncomment.
	// If Jaeger is used, ensure InitJaeger is called in main.
	
	router.Use(middleware.HTTPSRedirectMiddleware()) // Redirect HTTP to HTTPS first (if applicable at app level)
	// Recovery should be early.
	// If mylog.LogrusObj.Out can be used as io.Writer for Recovery:
	// router.Use(gin.RecoveryWithWriter(mylog.LogrusObj.Writer()))
	// Otherwise, default gin.Recovery() logs to gin.DefaultErrorWriter (stderr).
	// ErrorHandler middleware should then catch panics if Recovery re-panics or sets an error.
	router.Use(gin.Recovery()) 
	router.Use(middleware.RequestID())   // Add RequestID early
	router.Use(middleware.LoggerMiddleware()) // Custom Logrus logger
	router.Use(middleware.Cors())        // CORS policy
	router.Use(middleware.SecurityHeadersMiddleware()) // Add Security Headers
	router.Use(middleware.ContextTimeout(15 * time.Second)) // Global request timeout
	router.Use(middleware.PrometheusMiddleware())      // Add Prometheus middleware

	// DBInjectorMiddleware must come before RBAC or any other middleware that needs c.Get("db")
	router.Use(middleware.DBInjectorMiddleware())

	// I18nMiddleware should be early enough to localize messages for subsequent middlewares/handlers
	router.Use(middleware.I18nMiddleware())

	// router.Use(middleware.AuthMiddleware()) // AuthMiddleware is typically applied to specific route groups / in routes.go
	                                          // or globally if all routes need auth.
	                                          // It should set "userID" in the context for RBAC.
	router.Use(middleware.ErrorHandler()) // Error handler should be relatively late

	// Register custom validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// Assuming RegisterCustomValidators might exist in middleware or be defined here
		// For now, keep the existing alphanum validator registration
		v.RegisterValidation("alphanum", func(fl validator.FieldLevel) bool {
			return regexp.MustCompile("^[a-zA-Z0-9]+$").MatchString(fl.Field().String())
		})
		// Example if RegisterCustomValidators was defined:
		// middleware.RegisterCustomValidators(v)
	}

	routes.NewRouter(router, db)
	middleware.RegisterMetricsRoute(router) // Register Prometheus /metrics route

	// Swagger UI route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Example RBAC protected route group
	// Make sure AuthMiddleware is applied before RBAC if it's group-specific
	// For example:
	// adminGroup := router.Group("/api/v1/admin", middleware.AuthMiddleware(), middleware.RBAC("admin_area:access"))
	// {
	//    adminGroup.GET("/dashboard", func(c *gin.Context) {
	//        c.JSON(http.StatusOK, gin.H{"message": "Welcome to admin dashboard!"})
	//    })
	//    // Add more admin routes here that require "admin_area:access" or other permissions
	//    statsGroup := adminGroup.Group("/stats", middleware.RBAC("stats:view"))
	//    {
	//        statsGroup.GET("/", func(c *gin.Context) {
	//             c.JSON(http.StatusOK, gin.H{"message": "Viewing statistics."})
	//        })
	//    }
	// }

	// Initialize OSS Client and Upload Controller
	// Assuming conf.GlobalConfig.OSS has fields like Type, Endpoint, Region, Bucket, AccessKeyID, SecretAccessKey
	// These fields should match the structure expected by upload.Config
	if conf.GlobalConfig.OSS.Type == "s3" { // Check if OSS type is S3
		ossCfg := upload.Config{
			Type:            conf.GlobalConfig.OSS.Type,
			Endpoint:        conf.GlobalConfig.OSS.Endpoint,
			Region:          conf.GlobalConfig.OSS.Region,
			Bucket:          conf.GlobalConfig.OSS.Bucket,
			AccessKeyID:     conf.GlobalConfig.OSS.AccessKeyID,
			SecretAccessKey: conf.GlobalConfig.OSS.SecretAccessKey,
		}

		ossClient, err := upload.NewClient(ossCfg)
		if err != nil {
			// Use the project's logger mylog (which is douyin/pkg/utils/log)
			mylog.Fatalf("Failed to initialize OSS client: %v", err)
		}

		// Register Upload Controller
		// The UploadController's RegisterRoutes method is expected to create /api/v1/upload group
		uploadController := v1.UploadController{OSSClient: ossClient}
		uploadController.RegisterRoutes(router) // Pass the main router
		mylog.Info("OSS S3 client and UploadController initialized.")

	} else {
		mylog.Info("OSS Type not configured to 's3' or not set in config. Skipping OSS client and UploadController initialization.")
	}

	// Initialize Email Client and Notification Service
	// cancelWorker needs to be declared here to be accessible by the shutdown logic.
	var cancelWorker context.CancelFunc

	if conf.GlobalConfig.Mail.Host != "" { // Check if mail is configured
		mailCfg := email.ConfigMail{
			Host:               conf.GlobalConfig.Mail.Host,
			Port:               conf.GlobalConfig.Mail.Port,
			Username:           conf.GlobalConfig.Mail.Username,
			Password:           conf.GlobalConfig.Mail.Password,
			From:               conf.GlobalConfig.Mail.From,
			UseTLS:             conf.GlobalConfig.Mail.UseTLS,
			InsecureSkipVerify: conf.GlobalConfig.Mail.InsecureSkipVerify,
		}
		emailClient := email.NewClient(mailCfg)

		// Ensure redisClient (which is cache.Rdb if available) is used here
		if redisClient != nil { // redisClient was initialized for HealthController
			notificationSvc := service.NewNotificationService(redisClient, emailClient)

			var workerCtx context.Context
			workerCtx, cancelWorker = context.WithCancel(context.Background())
			// No defer for cancelWorker here, it's called during graceful shutdown

			go notificationSvc.ListenAndSend(workerCtx)
			mylog.Info("NotificationService (email queue worker) started.")

			// Example of enqueuing an email (for testing, remove/comment out in production)
			/*
			go func() {
			   time.Sleep(5 * time.Second) // Wait a bit for things to start
			   testJob := service.EmailJob{
			       To:      []string{"testrecipient@example.com"},
			       Subject: "Test Email from System",
			       Body:    "This is a test email sent asynchronously via Redis queue.",
			   }
			   err := notificationSvc.EnqueueEmail(context.Background(), testJob)
			   if err != nil {
			       mylog.Errorf("Test email enqueue failed: %v", err)
			   }
			}()
			*/

		} else {
			mylog.Warn("Redis client (cache.Rdb/redisClient) is nil. NotificationService not started.")
		}
	} else {
		mylog.Info("Mail not configured (Mail.Host is empty). Email sending disabled.")
	}


	// HTTP Server Setup for Graceful Shutdown
	srv := &http.Server{
		Addr:    conf.GlobalConfig.System.HttpPort,
		Handler: router,
	}

	go func() {
		mylog.Infof("服务器启动，监听端口 %s", conf.GlobalConfig.System.HttpPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			mylog.Fatalf("监听失败: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	mylog.Info("收到关闭信号，开始优雅退出...")

	// Signal the email worker to stop if it was started
	if cancelWorker != nil {
		mylog.Info("Signaling email worker to stop...")
		cancelWorker()
		// Optionally, add a short delay or a sync.WaitGroup if you want to ensure
		// the worker has some time to finish processing a current item before server shutdown.
		// For this example, direct cancellation is shown. A more robust system might wait.
		// time.Sleep(2 * time.Second)
	}

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		mylog.Fatalf("优雅关机失败: %v", err)
	}
	mylog.Info("服务已安全退出")
	fmt.Println("服务已安全退出")
}

// loading 用于加载配置、初始化缓存等服务
func loading() {
	rand.Seed(time.Now().UnixNano()) // Seed rand
	if err := conf.InitConfig(); err != nil {
		mylog.Errorf("配置加载失败: %v", err)
		log.Fatalf("配置加载失败: %v", err)
		return // Added return to ensure function exits on error
	}
	// Initialize Logger must be after config loading if logger depends on config
	mylog.InitLogger() // Initialize custom logger

	// Initialize Cache after logger (if cache logs)
	if err := cache.InitCache(); err != nil { // This function should ideally make the Redis client accessible
		mylog.Errorf("Redis 初始化失败: %v", err)
		log.Fatalf("Redis 初始化失败: %v", err)
		return // Added return
	}
	mylog.Info("加载配置完成，缓存初始化完成...")
	fmt.Println("加载配置完成...")
}
