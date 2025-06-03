package main

import (
	"fmt"
	"log"

	v1 "douyin/api/v1"
	conf "douyin/config"
	mylog "douyin/pkg/utils/log"
	"douyin/repository/cache"
	"douyin/routes"
	"regexp"

	"douyin/middleware" // Added import for middleware

	"github.com/gin-gonic/gin"
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

	v1.SetDB(db)
	v1.SetCheckoutController(db)

	router := gin.New() // Using gin.New() to have more control over middleware.

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
	// router.Use(middleware.AuthMiddleware()) // AuthMiddleware is typically applied to specific route groups / in routes.go
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
		return
	}
	if err := cache.InitCache(); err != nil {
		mylog.Errorf("Redis 初始化失败: %v", err)
		log.Fatalf("Redis 初始化失败: %v", err)
		return
	}
	mylog.Info("加载配置完成，缓存初始化完成...")
	fmt.Println("加载配置完成...")

	// Initialize Logger
	mylog.InitLogger() // Initialize custom logger
}
