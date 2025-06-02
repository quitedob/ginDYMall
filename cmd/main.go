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
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"math/rand"
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

	router := gin.Default()
	// Standard Gin middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	// Custom error handler
	router.Use(middleware.ErrorHandler())

	// Register custom validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("alphanum", func(fl validator.FieldLevel) bool {
			return regexp.MustCompile("^[a-zA-Z0-9]+$").MatchString(fl.Field().String())
		})
	}

	routes.NewRouter(router, db)

	// 启动 HTTP 服务
	if err := router.Run(conf.GlobalConfig.System.HttpPort); err != nil {
		mylog.Errorf("启动 HTTP 服务失败: %v", err)
		log.Fatalf("启动 HTTP 服务失败: %v", err)
	}
	mylog.Infof("服务启动成功，监听端口: %s", conf.GlobalConfig.System.HttpPort)
	fmt.Println("服务启动成功...")
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
}
