![Yopass-horizontal](https://user-images.githubusercontent.com/37777956/59544367-0867aa80-8f09-11e9-8d6a-02008e1bccc7.png)

# Yopass — Безопасный обмен секретами

[![Go Report Card](https://goreportcard.com/badge/github.com/Khovanskiy5/yopass)](https://goreportcard.com/report/github.com/Khovanskiy5/yopass)
[![codecov](https://codecov.io/gh/Khovanskiy5/yopass/branch/master/graph/badge.svg)](https://codecov.io/gh/Khovanskiy5/yopass)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/Khovanskiy5/yopass?sort=semver)

![demo](https://ydemo.netlify.com/yopass-demo.gif)

Yopass — это проект для быстрого и безопасного обмена конфиденциальной информацией (секретами).

Основная цель Yopass — свести к минимуму передачу паролей через системы управления тикетами, сообщения в Slack и электронную почту. Сообщения шифруются/дешифруются локально в браузере и отправляются в Yopass без ключа расшифровки, который виден только один раз в процессе шифрования. Затем Yopass возвращает одноразовую ссылку с указанным сроком действия.

Не существует идеального способа передачи секретов в интернете, и в каждой реализации есть свои компромиссы. Yopass спроектирован так, чтобы быть максимально простым и «глупым» без ущерба для безопасности. Между сгенерированным UUID и пользователем, отправившим зашифрованное сообщение, нет никакой связи. Всегда лучше отправлять весь контекст, за исключением самого пароля, по другому каналу связи.

**[Демо-версия доступна здесь](https://yopass.se)**. Если вы заботитесь о безопасности, рекомендуется развернуть Yopass самостоятельно.

- Сквозное (End-to-End) шифрование с использованием [OpenPGP](https://openpgpjs.org/)
- Секреты можно просмотреть только один раз
- Не требуются учетные записи или управление пользователями
- Секреты самоуничтожаются через X часов
- Возможность установки пользовательского пароля
- Ограниченная функциональность загрузки файлов

## История

Yopass был впервые выпущен в 2014 году и с тех пор поддерживается мной при участии этой фантастической группы [контрибьюторов](https://github.com/Khovanskiy5/yopass/graphs/contributors). Yopass используется многими крупными корпорациями, некоторые из которых перечислены ниже.

Если вы используете Yopass и хотите поддержать проект помимо вклада в код, вы можете выразить благодарность по электронной почте, сделать пожертвование или дать согласие на указание названия вашей компании в качестве пользователя Yopass в этом файле.

## Нам доверяют

- [Doddle LTD](https://doddle.com)
- [Spotify](https://spotify.com)
- [Gumtree Australia](https://www.gumtreeforbusiness.com.au/)

## Интерфейс командной строки

Основная идея Yopass заключается в том, чтобы каждый мог легко и быстро делиться секретами через простой веб-интерфейс. Также предоставляется интерфейс командной строки (CLI) для случаев, когда необходимо поделиться выводом программы.

```console
$ yopass --help
Yopass - Secure sharing for secrets, passwords and files

Flags:
      --api string          расположение API-сервера Yopass (по умолчанию "https://api.yopass.se")
      --decrypt string      URL для расшифровки секрета
      --expiration string   длительность, после которой секрет будет удален [1h, 1d, 1w] (по умолчанию "1h")
      --file string         прочитать секрет из файла вместо stdin
      --key string          вручную заданный ключ шифрования/расшифровки
      --one-time            одноразовая загрузка (по умолчанию true)
      --url string          публичный URL Yopass (по умолчанию "https://yopass.se")

Настройки считываются из флагов, переменных окружения или конфигурационного файла, расположенного по адресу
~/.config/yopass/defaults.<json,toml,yml,hcl,ini,...> в указанном порядке. Переменные окружения должны иметь префикс YOPASS_, а дефисы заменяются на подчеркивания.

Examples:
      # Зашифровать и отправить секрет из stdin
      printf 'secret message' | yopass

      # Зашифровать и отправить файл с секретом
      yopass --file /path/to/secret.conf

      # Делиться секретом несколько раз в течение всего дня
      cat secret-notes.md | yopass --expiration=1d --one-time=false

      # Расшифровать секрет в stdout
      yopass --decrypt https://yopass.se/#/...

Website: https://yopass.se
```

На данный момент доступны следующие варианты локальной установки CLI:

- Компиляция из исходного кода (требуется Go >= v1.21)

  ```console
  go install github.com/Khovanskiy5/yopass/cmd/yopass@latest
  ```

## Установка / Конфигурация

Параметры конфигурации сервера:

Флаги командной строки:

```console
$ yopass-server -h
      --address string             адрес прослушивания (по умолчанию 0.0.0.0)
      --allowed-expirations ints   допустимое время истечения срока действия в секундах (по умолчанию [3600,86400,604800])
      --database string            движок базы данных ('memcached' или 'redis') (default "memcached")
      --max-length int             максимальная длина зашифрованного секрета (по умолчанию 5242880)
      --memcached string           адрес Memcached (по умолчанию "localhost:11211")
      --metrics-port int           порт прослушивания сервера метрик (по умолчанию -1)
      --port int                   порт прослушивания (по умолчанию 1337)
      --redis string               URL Redis (по умолчанию "redis://localhost:6379/0")
      --tls-cert string            путь к TLS-сертификату
      --tls-key string             путь к TLS-ключу
      --cors-allow-origin string   настройка Access-Control-Allow-Origin CORS (по умолчанию "*")
      --force-onetime-secrets      запретить создание секретов, которые не являются одноразовыми
      --disable-upload             отключить эндпоинты загрузки /file
      --prefetch-secret            отображать информацию о том, что секрет может быть одноразовым (по умолчанию true)
      --disable-features           отключить раздел функций во фронтенде
      --no-language-switcher       отключить переключатель языков в интерфейсе
      --trusted-proxies strings    доверенные IP-адреса прокси или блоки CIDR для валидации заголовка X-Forwarded-For
      --privacy-notice-url string  URL страницы уведомления о конфиденциальности
      --imprint-url string         URL страницы с юридической информацией (imprint)
```

Зашифрованные секреты могут храниться в Memcached или Redis путем изменения флага `--database`. Списки допустимых сроков хранения (expiration) можно настроить с помощью флага `--allowed-expirations`.

### Настройка прокси

Когда Yopass развернут за обратным прокси-сервером или балансировщиком нагрузки (таким как Nginx, Caddy, Cloudflare или AWS ALB), вы можете захотеть логировать реальные IP-адреса клиентов вместо IP-адреса прокси. Yopass поддерживает настройку доверенных прокси для безопасной обработки заголовков `X-Forwarded-For`.

**Примечание по безопасности**: Заголовки X-Forwarded-For считаются доверенными только в том случае, если запросы поступают от явно настроенных доверенных прокси. Это предотвращает подмену IP-адресов из ненадежных источников.

#### Примеры:

```bash
# Доверять одному IP-адресу прокси
yopass-server --trusted-proxies 192.168.1.100

# Доверять нескольким IP-адресам прокси
yopass-server --trusted-proxies 192.168.1.100,10.0.0.50

# Доверять подсетям прокси (нотация CIDR)
yopass-server --trusted-proxies 192.168.1.0/24,10.0.0.0/8

# Переменная окружения (полезно для Docker)
export TRUSTED_PROXIES="192.168.1.0/24,10.0.0.0/8"
yopass-server
```

#### Типичные сценарии использования прокси:

- **Nginx/Apache**: Используйте IP-адрес вашего сервера обратного прокси.
- **Cloudflare**: Используйте диапазоны IP-адресов Cloudflare (доступны в их документации).
- **AWS ALB/ELB**: Используйте блок CIDR вашей VPC или подсеть балансировщика нагрузки.
- **Docker networks**: Используйте IP-адрес шлюза или подсеть сети Docker.

Без настройки доверенных прокси Yopass в целях безопасности всегда будет использовать IP-адрес прямого подключения, что является рекомендуемым поведением по умолчанию.

### Docker Compose

Используйте файл Docker Compose `deploy/with-nginx-proxy-and-letsencrypt/docker-compose.yml` для настройки экземпляра Yopass с шифрованием транспорта TLS и автоматическим продлением сертификатов с помощью [Let's Encrypt](https://letsencrypt.org/). Сначала направьте свой домен на хост, где вы хотите запустить Yopass. Затем замените значения-заполнители для `VIRTUAL_HOST`, `LETSENCRYPT_HOST` и `LETSENCRYPT_EMAIL` в файле docker-compose.yml вашими значениями. Перейдите в каталог развертывания и запустите контейнеры:

```console
docker-compose up -d
```

После этого Yopass будет доступен по домену, который вы указали в `VIRTUAL_HOST` / `LETSENCRYPT_HOST`.

Продвинутые пользователи, у которых уже есть обратный прокси-сервер, обрабатывающий TLS-соединения, могут использовать «небезопасную» (insecure) настройку:

```console
cd deploy/docker-compose/insecure
docker-compose up -d
```

Затем направьте ваш обратный прокси на `127.0.0.1:80`.

### Docker

С TLS-шифрованием:

```console
docker run --name memcached_yopass -d memcached
docker run -p 443:1337 -v /local/certs/:/certs \
    --link memcached_yopass:memcached -d Khovanskiy5/yopass --memcached=memcached:11211 --tls-key=/certs/tls.key --tls-cert=/certs/tls.crt
```

После этого Yopass будет доступен на порту 443 через все IP-адреса хоста, включая публичные. Чтобы ограничить доступ конкретным IP-адресом, используйте `-p 127.0.0.1:443:1337`.

Без TLS-шифрования (требуется обратный прокси для шифрования транспорта):

```console
docker run --name memcached_yopass -d memcached
docker run -p 127.0.0.1:80:1337 --link memcached_yopass:memcached -d Khovanskiy5/yopass --memcached=memcached:11211
```

Затем направьте ваш обратный прокси, обрабатывающий TLS-соединения, на `127.0.0.1:80`.

### Kubernetes

```console
kubectl apply -f deploy/yopass-k8.yaml
kubectl port-forward service/yopass 1337:1337
```

_Это предназначено для ознакомления, пожалуйста, настройте TLS при реальном использовании Yopass._

## Мониторинг

Yopass может опционально предоставлять метрики в текстовом формате [OpenMetrics][] / [Prometheus][]. Используйте флаг `--metrics-port <port>`, чтобы Yopass запустил второй HTTP-сервер на этом порту, делая метрики доступными по пути `/metrics`.

Поддерживаемые метрики:

- Базовые [метрики процесса][process metrics] с префиксом `process_` (например, использование процессора, памяти и дескрипторов файлов)
- Метрики среды выполнения Go с префиксом `go_` (например, использование памяти Go, статистика сборки мусора и т. д.)
- Метрики HTTP-запросов с префиксом `yopass_http_` (счетчик HTTP-запросов и гистограмма задержки HTTP-запросов)

[openmetrics]: https://openmetrics.io/
[prometheus]: https://prometheus.io/
[process metrics]: https://prometheus.io/docs/instrumenting/writing_clientlibs/#process-metrics

## Переводы

Yopass принимает переводы на дополнительные языки. Фронтенд поддерживает интернационализацию с использованием react-i18next, см. [текущие переводы](https://github.com/Khovanskiy5/yopass/blob/master/website/src/shared/lib/i18n.ts). Вклады в перевод приветствуются через pull requests, см. пример [здесь](https://github.com/Khovanskiy5/yopass/pull/3024) для добавления нового языка.
