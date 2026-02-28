package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type AppConfig struct {
	Env                 string
	Port                string
	ReadTimeoutSeconds  int
	WriteTimeoutSeconds int
}

type DBConfig struct {
	Host                   string
	Port                   int
	User                   string
	Password               string
	Name                   string
	SSLMode                string
	MaxOpenConns           int
	MaxIdleConns           int
	ConnMaxLifetimeMinutes int
}

type AuthConfig struct {
	JWTSecret        string
	JWTIssuer        string
	JWTExpiryMinutes int
}

type CORSConfig struct {
	AllowedOrigins []string
}

type StorageConfig struct {
	UploadBaseDir  string
	UploadMaxBytes int64
	Provider       string
	SupabaseURL    string
	SupabaseBucket string
	SupabaseKey    string
	PublicBaseURL  string
}

type BootstrapConfig struct {
	AdminEmail    string
	AdminPassword string
}

type GithubConfig struct {
	Token           string
	User            string
	Topic           string
	CacheTTLSeconds int
	PerPage         int
}

type Config struct {
	App       AppConfig
	DB        DBConfig
	Auth      AuthConfig
	CORS      CORSConfig
	Storage   StorageConfig
	Bootstrap BootstrapConfig
	Github    GithubConfig
}

func Load() (*Config, error) {
	portFallback := getEnv("PORT", "8080")
	cfg := &Config{
		App: AppConfig{
			Env:                 getEnv("APP_ENV", "development"),
			Port:                getEnv("APP_PORT", portFallback),
			ReadTimeoutSeconds:  getEnvAsInt("APP_READ_TIMEOUT_SECONDS", 15),
			WriteTimeoutSeconds: getEnvAsInt("APP_WRITE_TIMEOUT_SECONDS", 15),
		},
		DB: DBConfig{
			Host:                   getEnv("DB_HOST", "localhost"),
			Port:                   getEnvAsInt("DB_PORT", 5432),
			User:                   getEnv("DB_USER", "portfolio"),
			Password:               getEnv("DB_PASSWORD", "portfolio"),
			Name:                   getEnv("DB_NAME", "portfolio"),
			SSLMode:                getEnv("DB_SSLMODE", "disable"),
			MaxOpenConns:           getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:           getEnvAsInt("DB_MAX_IDLE_CONNS", 10),
			ConnMaxLifetimeMinutes: getEnvAsInt("DB_CONN_MAX_LIFETIME_MINUTES", 30),
		},
		Auth: AuthConfig{
			JWTSecret:        getEnv("JWT_SECRET", "replace-this-in-production"),
			JWTIssuer:        getEnv("JWT_ISSUER", "portfolio-platform"),
			JWTExpiryMinutes: getEnvAsInt("JWT_EXPIRY_MINUTES", 120),
		},
		CORS: CORSConfig{
			AllowedOrigins: splitCSV(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000")),
		},
		Storage: StorageConfig{
			UploadBaseDir:  getEnv("UPLOAD_BASE_DIR", "./uploads"),
			UploadMaxBytes: int64(getEnvAsInt("UPLOAD_MAX_BYTES", 10*1024*1024)),
			Provider:       strings.ToLower(getEnv("STORAGE_PROVIDER", "local")),
			SupabaseURL:    getEnv("SUPABASE_URL", ""),
			SupabaseBucket: getEnv("SUPABASE_STORAGE_BUCKET", ""),
			SupabaseKey:    getEnv("SUPABASE_SERVICE_ROLE_KEY", ""),
			PublicBaseURL:  getEnv("SUPABASE_STORAGE_PUBLIC_BASE_URL", ""),
		},
		Bootstrap: BootstrapConfig{
			AdminEmail:    getEnv("ADMIN_BOOTSTRAP_EMAIL", "admin@example.com"),
			AdminPassword: getEnv("ADMIN_BOOTSTRAP_PASSWORD", "ChangeMe123!"),
		},
		Github: GithubConfig{
			Token:           getEnv("GITHUB_TOKEN", ""),
			User:            getEnv("GITHUB_USER", "naodEthiop"),
			Topic:           getEnv("GITHUB_TOPIC", "portfolio"),
			CacheTTLSeconds: getEnvAsInt("GITHUB_CACHE_TTL_SECONDS", 300),
			PerPage:         getEnvAsInt("GITHUB_PER_PAGE", 30),
		},
	}

	if cfg.Auth.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}
	if cfg.Storage.Provider == "supabase" {
		if strings.TrimSpace(cfg.Storage.SupabaseURL) == "" {
			return nil, fmt.Errorf("SUPABASE_URL is required when STORAGE_PROVIDER=supabase")
		}
		if strings.TrimSpace(cfg.Storage.SupabaseKey) == "" {
			return nil, fmt.Errorf("SUPABASE_SERVICE_ROLE_KEY is required when STORAGE_PROVIDER=supabase")
		}
		if strings.TrimSpace(cfg.Storage.SupabaseBucket) == "" {
			return nil, fmt.Errorf("SUPABASE_STORAGE_BUCKET is required when STORAGE_PROVIDER=supabase")
		}
	}

	return cfg, nil
}

func (c *Config) PostgresDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.DB.Host,
		c.DB.Port,
		c.DB.User,
		c.DB.Password,
		c.DB.Name,
		c.DB.SSLMode,
	)
}

func (c *Config) JWTExpiryDuration() time.Duration {
	return time.Duration(c.Auth.JWTExpiryMinutes) * time.Minute
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

func getEnvAsInt(key string, fallback int) int {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	num, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}
	return num
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			out = append(out, trimmed)
		}
	}
	if len(out) == 0 {
		return []string{"http://localhost:3000"}
	}
	return out
}
