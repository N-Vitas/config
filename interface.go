package config

// Интерфейс для применения опций конфигурации
type Option interface {
	apply(v *Config)
}
