package config

// Декорированная функция для опций конфигурации
type optionFunc func(v *Config)

// Декорация интерфейса опций
func (fn optionFunc) apply(v *Config) {
	fn(v)
}
