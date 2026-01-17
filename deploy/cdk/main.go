package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/akrylysov/algnhsa"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/Khovanskiy5/yopass/internal/config"
	"github.com/Khovanskiy5/yopass/internal/secret/domain"
	"github.com/Khovanskiy5/yopass/internal/secret/handler"
	"github.com/Khovanskiy5/yopass/internal/secret/service"
	"github.com/Khovanskiy5/yopass/internal/server"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.SetDefault("cors-allow-origin", "*")
	viper.SetDefault("prefetch-secret", true)
	viper.SetDefault("max-length", 10000)
	viper.SetDefault("force-onetime-secrets", false)

	logger := configureZapLogger(zapcore.InfoLevel)
	registry := prometheus.NewRegistry()

	cfg := &config.Config{
		MaxLength:           viper.GetInt("max-length"),
		PrefetchSecret:      viper.GetBool("prefetch-secret"),
		CORSAllowOrigin:     viper.GetString("cors-allow-origin"),
		ForceOneTimeSecrets: viper.GetBool("force-onetime-secrets"),
	}

	repo := NewDynamo(os.Getenv("TABLE_NAME"))
	
	secretService := service.NewSecretService(
		repo,
		cfg.MaxLength,
		cfg.ForceOneTimeSecrets,
		[]int32{3600, 86400, 604800},
	)

	secretHandler := handler.NewSecretHandler(secretService, logger)
	configHandler := handler.NewConfigHandler(cfg, logger)

	router := server.NewRouter(cfg, secretHandler, configHandler, registry)

	algnhsa.ListenAndServe(router, nil)
}

// Dynamo Database implementation
type Dynamo struct {
	tableName string
	svc       *dynamodb.DynamoDB
}

// NewDynamo returns a database client
func NewDynamo(tableName string) domain.Repository {
	sess, _ := session.NewSession()
	return &Dynamo{tableName: tableName, svc: dynamodb.New(sess)}
}

// Get item from dynamo
func (d *Dynamo) Get(key string) (domain.Secret, error) {
	var s domain.Secret
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(key),
			},
		},
		TableName: aws.String(d.tableName),
	}
	result, err := d.svc.GetItem(input)
	if err != nil {
		return s, err
	}
	if len(result.Item) == 0 {
		return s, domain.ErrNotFound
	}

	if *result.Item["one_time"].BOOL {
		if err := d.deleteItem(key); err != nil {
			return s, err
		}
	}
	s.Message = *result.Item["secret"].S
	s.OneTime = *result.Item["one_time"].BOOL
	return s, nil
}

// Delete item
func (d *Dynamo) Delete(key string) (bool, error) {
	err := d.deleteItem(key)

	if errors.Is(err, &dynamodb.ResourceNotFoundException{}) {
		return false, nil
	}

	return err == nil, err
}

func (d *Dynamo) deleteItem(key string) error {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(key),
			},
		},
		TableName:    aws.String(d.tableName),
		ReturnValues: aws.String("ALL_OLD"),
	}

	_, err := d.svc.DeleteItem(input)
	return err
}

// Put item in Dynamo
func (d *Dynamo) Put(key string, secret domain.Secret) error {
	input := &dynamodb.PutItemInput{
		// TABLE GENERATED NAME
		Item: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(key),
			},
			"secret": {
				S: aws.String(secret.Message),
			},
			"one_time": {
				BOOL: aws.Bool(secret.OneTime),
			},
			"ttl": {
				N: aws.String(
					fmt.Sprintf(
						"%d", time.Now().Unix()+int64(secret.Expiration))),
			},
		},
		TableName: aws.String(d.tableName),
	}
	_, err := d.svc.PutItem(input)
	return err
}

// Status returns the OneTime status of a secret without retrieving or deleting it
func (d *Dynamo) Status(key string) (bool, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(key),
			},
		},
		TableName:            aws.String(d.tableName),
		ProjectionExpression: aws.String("one_time"),
	}
	result, err := d.svc.GetItem(input)
	if err != nil {
		return false, err
	}
	if len(result.Item) == 0 {
		return false, domain.ErrNotFound
	}

	oneTime := *result.Item["one_time"].BOOL
	return oneTime, nil
}

func configureZapLogger(logLevel zapcore.Level) *zap.Logger {
	loggerCfg := zap.NewProductionConfig()
	loggerCfg.Level.SetLevel(logLevel)
	logger, err := loggerCfg.Build()
	if err != nil {
		log.Fatalf("Unable to build logger %v", err)
	}
	zap.ReplaceGlobals(logger)
	return logger
}
