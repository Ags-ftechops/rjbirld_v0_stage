```go
package config

import (
    "github.com/spf13/viper"
)

type Config struct {
    Server   ServerConfig
    DB       DatabaseConfig
    Redis    RedisConfig
    AWS      AWSConfig
    Temporal TemporalConfig
    JWT      JWTConfig
}

type ServerConfig struct {
    Port         string `mapstructure:"PORT"`
    Environment  string `mapstructure:"ENVIRONMENT"`
    AllowOrigins string `mapstructure:"ALLOW_ORIGINS"`
}

type DatabaseConfig struct {
    Host     string `mapstructure:"DB_HOST"`
    Port     string `mapstructure:"DB_PORT"`
    User     string `mapstructure:"DB_USER"`
    Password string `mapstructure:"DB_PASSWORD"`
    Name     string `mapstructure:"DB_NAME"`
    SSLMode  string `mapstructure:"DB_SSLMODE"`
}

type RedisConfig struct {
    Host     string `mapstructure:"REDIS_HOST"`
    Port     string `mapstructure:"REDIS_PORT"`
    Password string `mapstructure:"REDIS_PASSWORD"`
    DB       int    `mapstructure:"REDIS_DB"`
}

type AWSConfig struct {
    Region          string `mapstructure:"AWS_REGION"`
    AccessKeyID     string `mapstructure:"AWS_ACCESS_KEY_ID"`
    SecretAccessKey string `mapstructure:"AWS_SECRET_ACCESS_KEY"`
    S3Bucket        string `mapstructure:"AWS_S3_BUCKET"`
}

type TemporalConfig struct {
    HostPort    string `mapstructure:"TEMPORAL_HOST_PORT"`
    Namespace   string `mapstructure:"TEMPORAL_NAMESPACE"`
    TaskQueue   string `mapstructure:"TEMPORAL_TASK_QUEUE"`
}

type JWTConfig struct {
    Secret    string `mapstructure:"JWT_SECRET"`
    ExpiresIn int    `mapstructure:"JWT_EXPIRES_IN"`
}

func LoadConfig() (*Config, error) {
    viper.SetConfigName(".env")
    viper.SetConfigType("env")
    viper.AddConfigPath(".")
    
    viper.AutomaticEnv()

    if err := viper.ReadInConfig(); err != nil {
        return nil, err
    }

    config := &Config{}
    if err := viper.Unmarshal(config); err != nil {
        return nil, err
    }

    return config, nil
}
```