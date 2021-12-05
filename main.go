package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var c *Config

type Config struct {
	// Разделитель, разделяющий список ключей
	// Используется для доступа к вложенному значению за один раз
	keyDelim string
	// Имя файла, который нужно искать внутри пути
	configName string
	// Файл, который нужно искать внутри пути
	configFile string
	// Тип конфигурации, может быть файлом или переменными окружения 'json' | 'env'
	configType string
	// Права доступа к файлу конфигурации
	configPermissions os.FileMode
	// Префикс, который будут использовать переменные ОКРУЖЕНИЯ.
	envPrefix string
	// Основной список конфигурации
	config   map[string]interface{}
	defaults map[string]interface{}
	override map[string]interface{}
	allkeys  []string
	env      map[string]string
	aliases  map[string]string
	loaded   bool
}

// GetInstance создает единственную конфигурацию
func GetInstance() *Config {
	if c == nil {
		c = New()
	}
	return c
}

// New создает новую конфигурацию по умолчанию
func New() *Config {
	c = new(Config)
	c.keyDelim = "."
	c.configName = "config"
	c.configFile = "config.json"
	c.configType = "json"
	c.loaded = false
	c.configPermissions = os.FileMode(0644)
	c.config = make(map[string]interface{})
	c.defaults = make(map[string]interface{})
	c.override = make(map[string]interface{})
	c.allkeys = []string{}
	c.env = make(map[string]string)
	c.aliases = make(map[string]string)
	return c
}

// NewWithOptions создает новую конфигурацию с опциями.
func NewWithOptions(opts ...Option) *Config {
	v := New()
	for _, opt := range opts {
		opt.apply(v)
	}
	return v
}

// GetString Получение значения конфигурации тип строка
func GetString(key string, def interface{}) string { return c.GetString(key, def) }

// GetString Получение значения конфигурации тип строка
func (v *Config) GetString(key string, def interface{}) string {
	return ToString(v.Get(key, def))
}

// GetInt Получение значения конфигурации тип число
func GetInt(key string, def interface{}) int { return c.GetInt(key, def) }

// GetInt Получение значения конфигурации тип число
func (v *Config) GetInt(key string, def interface{}) int {
	return ToInt(v.Get(key, def))
}

// GetInt64 Получение значения конфигурации тип число 64
func GetInt64(key string, def interface{}) int64 { return c.GetInt64(key, def) }

// GetInt64 Получение значения конфигурации тип число 64
func (v *Config) GetInt64(key string, def interface{}) int64 {
	return ToInt64(v.Get(key, def))
}

// GetFloat Получение значения конфигурации тип число с плавоющей точкой
func GetFloat(key string, def interface{}) float64 { return c.GetFloat(key, def) }

// GetFloat Получение значения конфигурации тип число с плавоющей точкой
func (v *Config) GetFloat(key string, def interface{}) float64 {
	return ToFloat(v.Get(key, def))
}

// GetBool Получение значения конфигурации тип истина или ложь
func GetBool(key string, def interface{}) bool { return c.GetBool(key, def) }

// GetBool Получение значения конфигурации тип истина или ложь
func (v *Config) GetBool(key string, def interface{}) bool {
	return ToBool(v.Get(key, def))
}

// GetSlice Получение значения конфигурации тип истина или ложь
func GetSlice(key string, def interface{}) []string { return c.GetSlice(key, def) }

// GetSlice Получение значения конфигурации тип истина или ложь
func (v *Config) GetSlice(key string, def interface{}) []string {
	return ToStringSlice(v.Get(key, def))
}

// KeyDelimiter Устанавливает разделитель, используемый для определения частей ключа.
func (v *Config) KeyDelimiter(d string) {
	v.keyDelim = d
}

// SetConfigFile Устанавливает путь файла конфигурации.
func (v *Config) SetConfigFile(path string) {
	if !Empty(path) {
		file := strings.Split(filepath.Base(path), ".")
		v.configName = file[0]
		v.configFile = path
		if !FileExists(path) {
			v.createFileConf()
		}
	}
}

// SetConfigName Устанавливает путь файла конфигурации.
func (v *Config) SetConfigName(name string) {
	v.configName = name
}

// SetEnvPrefix Устанавливает префикс, который будут использовать переменные ОКРУЖЕНИЯ.
func (v *Config) SetEnvPrefix(in string) {
	if len(in) > 0 {
		v.envPrefix = in
	}
}

// UseConfigFile Устанавливает конфигурацию из файла
func (v *Config) UseConfigFile() {
	v.configType = "json"
	v.getConfigFile()
}

// UseEventSystem Устанавливает конфигурацию из файла
func (v *Config) UseEventSystem() {
	v.configType = "env"
}

// Get Возвращает значение переменной окружения либо значение по умолчанию
func (c *Config) GetEnvSystem(key string, def interface{}) string {
	if len(os.Getenv(key)) == 0 {
		return ToString(def)
	}
	return os.Getenv(key)
}

// Get Возвращает интерфейс. Для конкретного знач Get____ methods.
func Get(key string, def string) interface{} { return c.Get(key, def) }

// Get Возвращает интерфейс. Для конкретного знач Get____ methods.
func (v *Config) Get(key string, def interface{}) (value interface{}) {
	lcaseKey := strings.ToLower(key)
	path := strings.Split(lcaseKey, v.keyDelim)
	value = v.searchMap(v.config, path)
	if value == nil {
		ucaseKey := strings.ToUpper(key)
		value = v.GetEnvSystem(ucaseKey, def)
	}
	return value
}

// AllKeys Получение всех ключей настройках
func AllKeys() []string { return c.AllKeys() }

// AllKeys Получение всех ключей настройках
func (v *Config) AllKeys() (keys []string) {
	return v.allkeys
}

// SetConfigPermissions Установка прав доступа файлу конфигурации
func SetConfigPermissions(perm os.FileMode) { c.SetConfigPermissions(perm) }

// SetConfigPermissions Установка прав доступа файлу конфигурации
func (v *Config) SetConfigPermissions(perm os.FileMode) {
	v.configPermissions = perm.Perm()
}

// Unmarshal Парсит файл конфигурации
func (v *Config) Unmarshal(data []byte) {
	if err := json.Unmarshal(data, &v.defaults); err != nil {
		panic(err)
	}
}

// Поиск ключ значение во вложеннии
func (v *Config) searchMap(source map[string]interface{}, path []string) interface{} {
	if len(path) == 0 {
		return source
	}

	next, ok := source[path[0]]
	if ok {
		// Fast path
		if len(path) == 1 {
			return next
		}

		// Nested case
		switch s := next.(type) {
		case map[interface{}]interface{}:
			return v.searchMap(ToStringMap(s), path[1:])
		case map[string]interface{}:
			// Type assertion is safe here since it is only reached
			// if the type of `next` is the same as the type being asserted
			return v.searchMap(next.(map[string]interface{}), path[1:])
		default:
			// got a value but nested key expected, return "nil" for not found
			return nil
		}
	}
	return nil
}

// Чтение файла конфигурации
func (v *Config) getConfigFile() {
	if !v.loaded && v.configType == "json" {
		jsonFile, err := os.Open(v.configFile)
		// Если файла нет, создать файл и снова его прочесть
		if err != nil {
			v.createFileConf()
			v.loaded = true
			v.getConfigFile()
			return
		}
		// Закрываем файл после чтения
		defer jsonFile.Close()
		byteValue, _ := ioutil.ReadAll(jsonFile)
		if err := json.Unmarshal([]byte(byteValue), &v.defaults); err == nil {
			// Заполнение всей конфигурации
			v.concatConfigFile()
		}
	}
}

// Запись конфигурации в файл
func (v *Config) createFileConf() {
	f, err := os.OpenFile(v.configFile, os.O_APPEND|os.O_WRONLY, v.configPermissions)
	if err != nil {
		f, err = os.Create(v.configFile)
		if err != nil {
			panic(err)
		}
	}
	defer f.Close()
	conf, err := json.Marshal(v.config)
	if err != nil {
		panic(err)
	}
	if _, err = f.Write(conf); err != nil {
		panic(err)
	}
}

// Заполнение всей конфигурации
func (v *Config) concatConfigFile() {
	m := map[string]bool{}
	// Копирование карты конфигурации не чувствительная к регистру
	v.config = copyAndInsensitiviseMap(v.defaults)
	// Сглаживание карты
	v.flattenAndMergeMap(m, v.defaults, "")
}

// Перебор файла конфигурации для сглаживания вложений
func (v *Config) flattenAndMergeMap(shadow map[string]bool, m map[string]interface{}, prefix string) map[string]bool {
	// Выход из рекурсии
	if shadow != nil && prefix != "" && shadow[prefix] {
		return shadow
	}
	// Инициация карты если ее нет
	if shadow == nil {
		shadow = make(map[string]bool)
	}
	// Карта для вложенных данных
	var m2 map[string]interface{}
	if prefix != "" {
		// Склеивание ключа
		prefix += v.keyDelim
	}
	// Перебор вложенной карты
	for k, val := range m {
		// Формирование полного ключа
		fullKey := strings.ToLower(prefix + k)
		switch m := val.(type) {
		case map[string]interface{}:
			m2 = m
		case map[interface{}]interface{}:
			m2 = ToStringMap(val)
		default:
			// Обработка последнего вложенного типа
			shadow[fullKey] = true
			v.allkeys = append(v.allkeys, fullKey)
			v.override[fullKey] = val
			continue
		}
		// Рекурсивное слияние карты
		shadow = v.flattenAndMergeMap(shadow, m2, fullKey)
	}
	return shadow
}

// copyAndInsensitiviseMap behaves like insensitiviseMap, but creates a copy of
// any map it makes case insensitive.
func copyAndInsensitiviseMap(m map[string]interface{}) map[string]interface{} {
	nm := make(map[string]interface{})

	for key, val := range m {
		lkey := strings.ToLower(key)
		switch v := val.(type) {
		case map[interface{}]interface{}:
			nm[lkey] = copyAndInsensitiviseMap(ToStringMap(v))
		case map[string]interface{}:
			nm[lkey] = copyAndInsensitiviseMap(v)
		default:
			nm[lkey] = v
		}
	}

	return nm
}
