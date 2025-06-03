// config.go
package config

import (
	"fmt"
	"os"

	"fmt"
	"os"
	"strings" // Added for strings.NewReplacer

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

// GlobalConfig 为全局配置对象
var GlobalConfig *Conf

// Conf 定义整个项目的配置结构体
type Conf struct {
	System        *System                 `yaml:"system"`        // 系统相关配置
	Oss           *Oss                    `yaml:"oss"`           // 对象存储配置（可选）
	MySql         *MySql                  `yaml:"mysql"`         // MySQL 数据库配置
	Email         *Email                  `yaml:"email"`         // 邮件配置
	Redis         *Redis                  `yaml:"redis"`         // Redis 配置
	EncryptSecret *EncryptSecret          `yaml:"encryptSecret"` // 加密密钥配置
	Cache         *Cache                  `yaml:"cache"`         // 缓存配置
	KafKa         map[string]*KafkaConfig `yaml:"kafKa"`         // Kafka 配置
	RabbitMq      *RabbitMq               `yaml:"rabbitMq"`      // RabbitMQ 配置
	Es            *Es                     `yaml:"es"`            // ElasticSearch 配置
	PhotoPath     *LocalPhotoPath         `yaml:"photoPath"`     // 本地图片存储路径配置
}

// 以下为各部分配置结构体定义（部分可根据实际需求扩展）
type System struct {
	AppEnv      string `yaml:"env"`         // 运行环境
	Domain      string `yaml:"domain"`      // 项目域名
	Version     string `yaml:"version"`     // 版本号
	HttpPort    string `yaml:"HttpPort"`    // HTTP 端口
	Host        string `yaml:"Host"`        // 主机地址
	UploadModel string `yaml:"UploadModel"` // 上传模式
}

type Oss struct {
	BucketName      string `yaml:"bucketName"`
	AccessKeyId     string `yaml:"accessKeyId"`
	AccessKeySecret string `yaml:"accessKeySecret"`
	Endpoint        string `yaml:"endPoint"`
	EndpointOut     string `yaml:"endpointOut"`
	QiNiuServer     string `yaml:"qiNiuServer"`
}

type MySql struct {
	Default *MySqlConfig `yaml:"default"`
}

type MySqlConfig struct {
	Dialect  string `yaml:"dialect"`
	DbHost   string `yaml:"dbHost"`
	DbPort   string `yaml:"dbPort"`
	DbName   string `yaml:"dbName"`
	UserName string `yaml:"userName"`
	Password string `yaml:"password"`
	Charset  string `yaml:"charset"`
}

type Email struct {
	ValidEmail string `yaml:"validEmail"`
	SmtpHost   string `yaml:"smtpHost"`
	SmtpEmail  string `yaml:"smtpEmail"`
	SmtpPass   string `yaml:"smtpPass"`
}

type Redis struct {
	RedisHost     string `yaml:"redisHost"`     // Redis 服务器地址
	RedisPort     string `yaml:"redisPort"`     // Redis 端口
	RedisUsername string `yaml:"redisUsername"` // Redis 用户名
	RedisPassword string `yaml:"redisPwd"`      // Redis 密码
	RedisDbName   int    `yaml:"redisDbName"`   // Redis 数据库编号
	RedisNetwork  string `yaml:"redisNetwork"`  // 网络协议（tcp）
}

type EncryptSecret struct {
	JwtSecret   string `yaml:"jwtSecret"`
	EmailSecret string `yaml:"emailSecret"`
	PhoneSecret string `yaml:"phoneSecret"`
	MoneySecret string `yaml:"moneySecret"`
}

type LocalPhotoPath struct {
	PhotoHost   string `yaml:"photoHost"`
	ProductPath string `yaml:"productPath"`
	AvatarPath  string `yaml:"avatarPath"`
}

type Cache struct {
	CacheType    string `yaml:"cacheType"`
	CacheExpires int64  `yaml:"cacheEmpires"`
	CacheWarmUp  bool   `yaml:"cacheWarmUp"`
	CacheServer  string `yaml:"cacheServer"`
}

type Es struct {
	EsHost  string `yaml:"esHost"`
	EsPort  string `yaml:"esPort"`
	EsIndex string `yaml:"esIndex"`
}

type RabbitMq struct {
	RabbitMQ         string `yaml:"rabbitMq"`
	RabbitMQUser     string `yaml:"rabbitMqUser"`
	RabbitMQPassWord string `yaml:"rabbitMqPassWord"`
	RabbitMQHost     string `yaml:"rabbitMqHost"`
	RabbitMQPort     string `yaml:"rabbitMqPort"`
}

type KafkaConfig struct {
	DisableConsumer bool   `yaml:"disableConsumer"`
	Debug           bool   `yaml:"debug"`
	Address         string `yaml:"address"`
	RequiredAck     int    `yaml:"requiredAck"`
	ReadTimeout     int64  `yaml:"readTimeout"`
	WriteTimeout    int64  `yaml:"writeTimeout"`
	MaxOpenRequests int    `yaml:"maxOpenRequests"`
	Partition       int    `yaml:"partition"`
}

// LoadConfig 从指定路径加载配置文件，并反序列化到 Conf 对象中
// This function might become obsolete if InitConfig handles all loading.
// Or it can be kept for specific use cases like loading a config file not based on APP_ENV.
func LoadConfig(path string) (*Conf, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("无法打开配置文件: %v", err)
	}
	defer file.Close()

	var config Conf
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}
	return &config, nil
}

// InitConfig 利用 viper 加载配置文件并反序列化到全局变量 GlobalConfig 中
func InitConfig() error {
	workDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取工作目录失败: %v", err)
	}

	v := viper.New() // Use a local Viper instance

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev" // Default to 'dev' environment
	}
	fmt.Printf("当前运行环境 (APP_ENV): %s\n", env)

	v.SetConfigName("config." + env) // e.g., config.dev, config.prod
	v.SetConfigType("yaml")
	v.AddConfigPath(workDir + "/config") // Standard path for config files
	v.AddConfigPath(workDir)             // Fallback for simpler local setups

	// Allow environment variables to override config file settings
	v.AutomaticEnv()
	// Replace dots with underscores for environment variable mapping
	// e.g., EncryptSecret.JwtSecret becomes ENCRYPTSECRET_JWTSECRET
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件 '%s' 失败: %w", "config."+env+".yaml", err)
	}

	if err := v.Unmarshal(&GlobalConfig); err != nil {
		return fmt.Errorf("反序列化配置失败: %w", err)
	}

	// Explicitly override JWT secret if JWT_SECRET environment variable is set
	// This ensures the specific requirement from the issue is met,
	// even though AutomaticEnv + SetEnvKeyReplacer might also handle it
	// if the struct tags are `mapstructure:"jwtSecret"` and env var is `ENCRYPTSECRET_JWTSECRET`.
	// The current struct tag is `yaml:"jwtSecret"`.
	// For robust override, check explicitly.
	if envJwtSecret := os.Getenv("JWT_SECRET"); envJwtSecret != "" {
		if GlobalConfig.EncryptSecret == nil {
			GlobalConfig.EncryptSecret = &EncryptSecret{} // Ensure EncryptSecret struct is initialized
		}
		GlobalConfig.EncryptSecret.JwtSecret = envJwtSecret
		fmt.Println("JWT Secret overridden by JWT_SECRET environment variable.")
	} else {
		// Check if AutomaticEnv picked up a structured env var like ENCRYPTSECRET_JWTSECRET
		// This part is tricky because viper's AutomaticEnv is case-insensitive for keys but case-sensitive for values from env.
		// The SetEnvKeyReplacer helps map `EncryptSecret.JwtSecret` to `ENCRYPTSECRET_JWTSECRET`.
		// If `GlobalConfig.EncryptSecret.JwtSecret` is already correctly populated by viper.Unmarshal
		// (after considering env vars due to AutomaticEnv), this specific os.Getenv("JWT_SECRET") might be redundant
		// or could be a fallback if the structured one isn't set.
		// For now, the explicit os.Getenv("JWT_SECRET") takes precedence as per the issue.
	}


	fmt.Printf("配置文件 '%s' 加载成功！\n", "config."+env+".yaml")
	return nil
}

// GetExpiresTime 根据缓存配置获取过期时间
func GetExpiresTime() int64 {
	if GlobalConfig.Cache.CacheExpires == 0 {
		return 1800
	}
	if GlobalConfig.Cache.CacheExpires == -1 {
		return -1
	}
	return GlobalConfig.Cache.CacheExpires
}
