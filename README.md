![Yopass-horizontal](https://user-images.githubusercontent.com/37777956/59544367-0867aa80-8f09-11e9-8d6a-02008e1bccc7.png)

# Yopass — Безопасный обмен секретами

[![Go Report Card](https://goreportcard.com/badge/github.com/jhaals/yopass)](https://goreportcard.com/report/github.com/jhaals/yopass)
[![codecov](https://codecov.io/gh/jhaals/yopass/branch/master/graph/badge.svg)](https://codecov.io/gh/jhaals/yopass)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/jhaals/yopass?sort=semver)

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

Yopass был впервые выпущен в 2014 году и с тех пор поддерживается мной при участии этой фантастической группы [контрибьюторов](https://github.com/jhaals/yopass/graphs/contributors). Yopass используется многими крупными корпорациями, некоторые из которых перечислены ниже.

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
      --api string          Yopass API server location (default "https://api.yopass.se")
      --decrypt string      Decrypt secret URL
      --expiration string   Duration after which secret will be deleted [1h, 1d, 1w] (default "1h")
      --file string         Read secret from file instead of stdin
      --key string          Manual encryption/decryption key
      --one-time            One-time download (default true)
      --url string          Yopass public URL (default "https://yopass.se")

Settings are read from flags, environment variables, or a config file located at
~/.config/yopass/defaults.<json,toml,yml,hcl,ini,...> in this order. Environment
variables have to be prefixed with YOPASS_ and dashes become underscores.

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
  go install github.com/jhaals/yopass/cmd/yopass@latest
  ```

## Установка / Конфигурация

Параметры конфигурации сервера:

Флаги командной строки:

```console
$ yopass-server -h
      --address string             listen address (default 0.0.0.0)
      --database string            database backend ('memcached' or 'redis') (default "memcached")
      --max-length int             max length of encrypted secret (default 5242880)
      --memcached string           Memcached address (default "localhost:11211")
      --metrics-port int           metrics server listen port (default -1)
      --port int                   listen port (default 1337)
      --redis string               Redis URL (default "redis://localhost:6379/0")
      --tls-cert string            path to TLS certificate
      --tls-key string             path to TLS key
      --cors-allow-origin          Access-Control-Allow-Origin CORS setting (default *)
      --force-onetime-secrets      reject non onetime secrets from being created
      --disable-upload             disable the /file upload endpoints
      --prefetch-secret            display information that the secret might be one time use (default true)
      --disable-features           disable features section on frontend
      --no-language-switcher       disable the language switcher in the UI
      --trusted-proxies strings    trusted proxy IP addresses or CIDR blocks for X-Forwarded-For header validation
      --privacy-notice-url string  URL to privacy notice page
      --imprint-url string         URL to imprint/legal notice page
```

Зашифрованные секреты могут храниться в Memcached или Redis путем изменения флага `--database`.

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
    --link memcached_yopass:memcached -d jhaals/yopass --memcached=memcached:11211 --tls-key=/certs/tls.key --tls-cert=/certs/tls.crt
```

После этого Yopass будет доступен на порту 443 через все IP-адреса хоста, включая публичные. Чтобы ограничить доступ конкретным IP-адресом, используйте `-p 127.0.0.1:443:1337`.

Без TLS-шифрования (требуется обратный прокси для шифрования транспорта):

```console
docker run --name memcached_yopass -d memcached
docker run -p 127.0.0.1:80:1337 --link memcached_yopass:memcached -d jhaals/yopass --memcached=memcached:11211
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

Yopass принимает переводы на дополнительные языки. Фронтенд поддерживает интернационализацию с использованием react-i18next, см. [текущие переводы](https://github.com/jhaals/yopass/blob/master/website/src/shared/lib/i18n.ts). Вклады в перевод приветствуются через pull requests, см. пример [здесь](https://github.com/jhaals/yopass/pull/3024) для добавления нового языка.
