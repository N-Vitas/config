## Config Проект конфигурации приложения на Go 
> Может кому то показаться что это велосипед. Ведь уже есть готовые решения бери да пользуйся  
> Частично вы правы. Я до этого проекта пользовался различными пакетами конфигураций проектов. И в целом они отличные пакеты, но в них очень много лишнего.  
> Я решил сделать версию урезанного пакета, так как весь функционал предоставляемых возможностей в сторонних пакетах мне не нужен. 
> Мне нужно было только чтение конфигурации из файла json и переменные среды для конфигурации прода в докер контейнерах.

## Установка
``` 
go get github.com/n-vitas/config 
```  


## Пример как использовать
``` golang
package main

import (
	"github.com/n-vitas/config"
	"fmt"
)

func main() {
	con := config.NewWithOptions(
		config.SetConfigFile("./main.json"),
		config.UseConfigFile(),
	)
	fmt.Printf("%v\n", con.GetBool("auth.ismart", false))
	fmt.Printf("%d\n", con.GetInt64("db.test.port", 3306))
	fmt.Printf("%v\n", con.GetSlice("auth.maps", []string{}))
	fmt.Printf("%v", con.GetString("go111module", "off"))
}
```
