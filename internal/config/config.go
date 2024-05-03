// Package config is defined application settings.
package config

import (
	"fmt"
	"reflect"

	"github.com/caarlos0/env/v10"
)

// Config アプリケーション設定を表す構造体。
type Config struct {
	// Port ポート
	// 0 が指定された場合は動的にポートを割り当てる。
	Port uint16 `env:"PORT" envDefault:"8080"`

	// DBHost データベースホスト
	DBHost string `env:"DB_HOST" envDefault:"localhost"`
	// DBPort データベースポート
	DBPort uint16 `env:"DB_PORT" envDefault:"3306"`
	// DBName データベース名
	DBName string `env:"DB_NAME,required"`
	// DBUsername データベースユーザー名
	DBUsername string `env:"DB_USERNAME,required"`
	// DBPassword データベースパスワード
	DBPassword string `env:"DB_PASSWORD"`
	DBUrl      string `env:"DB_URL,required"`

	// AuthSecret 認証トークンの署名用シークレット
	AuthSecret   string `env:"AUTH_SECRET,required"`
	SecretIssuer string `env:"SECRET_ISSUER,required"`

	// ClientOrigin クライアントのオリジン
	ClientOrigin ClientOrigin `env:"CLIENT_ORIGIN"`

	// DebugCORS CORS デバッグモード
	DebugCORS bool `env:"DEBUG_CORS"`

	AppDebug bool `env:"APP_DEBUG"`
	// development, staging, production
	AppEnv EnvironmentMode `env:"APP_ENV" envDefault:"production"`
	// FakeTime Fake time mode setting
	// If a time is specified, fix to that time.
	// If a truthy value is specified, fix to the default time.
	FakeTime FakeTimeMode `env:"FAKE_TIME"`
	LogLevel LogLevel     `env:"LOG_LEVEL,required"`
}

var parseFuncMap = map[reflect.Type]env.ParserFunc{
	reflect.TypeOf(ProductionEnv):  parseEnvironmentMode,
	reflect.TypeOf(FakeTimeMode{}): parseFakeTimeMode,
	reflect.TypeOf(InfoLevel):      parseLogLevel,
	reflect.TypeOf(ClientOrigin{}): parseClientOrigin,
}

// Get Get application settings from environment variables.
func Get() (*Config, error) {
	cfg := &Config{}
	if err := env.ParseWithOptions(cfg, env.Options{FuncMap: parseFuncMap}); err != nil {
		return nil, fmt.Errorf("parse env: %w", err)
	}

	return cfg, nil
}
