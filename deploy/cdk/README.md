# Конфигурация CDK для развертывания yopass в AWS.

Это предназначено __прежде всего__ для самого проекта и требует изменений для использования в вашей собственной настройке.

```
npx cdk deploy
```

```
GOOS=linux GOARCH=arm64 go build -o ./bootstrap -tags lambda.norpc
zip deployment.zip bootstrap
```


Файл `cdk.json` сообщает CDK Toolkit, как запускать ваше приложение.

## Полезные команды

* `npm run build`   компиляция typescript в js
* `npm run watch`   отслеживание изменений и компиляция
* `npm run test`    выполнение юнит-тестов jest
* `npx cdk deploy`  развертывание этого стека в вашу учетную запись/регион AWS по умолчанию
* `npx cdk diff`    сравнение развернутого стека с текущим состоянием
* `npx cdk synth`   генерация синтезированного шаблона CloudFormation
