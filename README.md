# GoPay

[![Go Reference](https://pkg.go.dev/badge/github.com/Anton-Kraev/gopay.svg)](https://pkg.go.dev/github.com/Anton-Kraev/gopay)
[![Go Report Card](https://goreportcard.com/badge/github.com/Anton-Kraev/gopay)](https://goreportcard.com/report/github.com/Anton-Kraev/gopay)
[![CI Pipeline](https://github.com/Anton-Kraev/gopay/actions/workflows/ci.yml/badge.svg)](https://github.com/Anton-Kraev/gopay/actions/workflows/ci.yml)

Сервис для управления доступом к цифровым товарам, включающий прием платежей от клиентов и доставку им товаров после оплаты.

## Функциональность
- **Платежные провайдеры**
  - ЮKassa
    - Создание платежных ссылок и прием платежей
    - Обработка уведомлений о совершении платежей
- **Хранилища данных**
  - BoltDB
    - Хранение данных о пользователях и их платежах и товарах
  - MinIO
    - Хранение продаваемых файлов
- **Способы доставки**
  - По ссылке в браузере
    - Генерация уникальных ссылок для доступа к цифровым товарам
- **Виды цифровых товаров**
  - PDF-файлы

## Требования к окружению
- **Go 1.23+**
- GNU Make 3.80+ (опционально, для упрощения сборки)

## Запуск и тестирование локально
Ниже представлена инструкция по использованию основных команд, для получения списка всех доступных команд воспользуйтесь
`make help`.

### Тестирование
Генерация mock-объектов и запуск unit-тестов:
```shell
make test
```

### Запуск API
Ниже приведены доступные настройки запуска (флаг > переменная):

| Флаг Go CLI                | Переменная окружения     | Значение по умолчанию | Описание                        |
|----------------------------|--------------------------|-----------------------|---------------------------------|
| `--env`                    | `ENV`                    | `dev`                 | Окружение (dev/prod)            |
| `--gopay-host `            | `GOPAY_HOST`             | `localhost`           | Хост для HTTP-сервера           |
| `--gopay-port`/`-p`        | `GOPAY_PORT`             | `8080`                | Порт для HTTP-сервера           |
| `--db-file-path`           | `DB_FILE_PATH`           | `data.db`             | Путь к файлу базы данных        |
| `--db-open-timeout`        | `DB_OPEN_TIMEOUT`        | `10s`                 | Таймаут подключения к БД        |
| *`--yookassa-checkout-url` | *`YOOKASSA_CHECKOUT_URL` | -                     | URL для вебхука ЮKassa          |
| *`--yookassa-shop-id`      | *`YOOKASSA_SHOP_ID`      | -                     | Идентификатор магазина в ЮKassa |
| *`--yookassa-api-token`    | *`YOOKASSA_API_TOKEN`    | -                     | Секретный токен API ЮKassa      |

Пример сборки и запуска веб-сервера и API:
```shell
go run cmd/api/main.go --yookassa-checkout-url <gopay_checkout> --yookassa-shop-id <shop_id> --yookassa-api-token <api_token>
```

> при локальном запуске (серый IP-адрес) уведомления от платежного сервиса (ЮKassa) приходить не будут

Документация API будет доступна после запуска по адресу:
`http://<GOPAY_HOST>:<GOPAY_PORT>/swagger/index.html`

### Запуск бота
Ниже приведены доступные настройки запуска (флаг > переменная):

| Флаг Go CLI          | Переменная окружения | Значение по умолчанию    | Описание                                   |
|----------------------|----------------------|--------------------------|--------------------------------------------|
| `--env`              | `ENV`                | `dev`                    | Окружение (dev/prod)                       |
| `--gopay-server-url` | `GOPAY_SERVER_URL`   | `http://127.0.0.1:8080`  | Базовый URL сервера                        |
| *`--tg-bot-token`    | *`TG_BOT_TOKEN`      | -                        | Токен бота от BotFather                    |
| *`--tg-admin-ids`    | *`TG_ADMIN_IDS`      | -                        | Telegram ID администраторов через запятую  |

Пример сборки и запуска Telegram-бота для управления сервисом:
```shell
go run cmd/bot/main.go --tg-bot-token <token> --tg-admin-ids <id1>,<id2>
```

## Установка библиотеки
Также реализована библиотека для Go, которая предоставляет:
- **PaymentManager** 
  - Объект с бизнес-логикой GoPay
- **AdminClient** 
  - HTTP-клиент для GoPay API

Для установки библиотеки в другой проект на Go:
```shell
go get github.com/Anton-Kraev/gopay
```

Затем можно импортировать и использовать в своем проекте:
```go
package main

import "github.com/Anton-Kraev/gopay"

const serverURL = "http://localhost:8080"

func main() {
	client, err := gopay.NewAdminClient(serverURL)
	if err != nil {
		print(err)
		return
	}
	
	statuses, err := client.NewAllPaymentService().Do()
	if err != nil {
		print(err)
		return
	}
	
	for id, status := range statuses {
		println(id, status)
	}
}
```
