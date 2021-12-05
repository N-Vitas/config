package config

// KeyDelimiter устанавливает разделитель, используемый для определения частей ключа.
func KeyDelimiter(d string) Option {
	return optionFunc(func(v *Config) {
		v.KeyDelimiter(d)
	})
}

// SetConfigFile Устанавливает путь файла конфигурации.
func SetConfigFile(in string) Option {
	return optionFunc(func(v *Config) {
		v.SetConfigFile(in)
	})
}

// SetConfigName Устанавливает путь файла конфигурации.
func SetConfigName(in string) Option {
	return optionFunc(func(v *Config) {
		v.SetConfigName(in)
	})
}

// SetEnvPrefix Устанавливает префикс, который будут использовать переменные ОКРУЖЕНИЯ.
func SetEnvPrefix(in string) Option {
	return optionFunc(func(v *Config) {
		v.SetEnvPrefix(in)
	})
}

// UseConfigFile Устанавливает конфигурацию из файла
func UseConfigFile() Option {
	return optionFunc(func(v *Config) {
		v.UseConfigFile()
	})
}

// UseEventSystem Устанавливает конфигурацию из файла
func UseEventSystem() Option {
	return optionFunc(func(v *Config) {
		v.UseEventSystem()
	})
}
