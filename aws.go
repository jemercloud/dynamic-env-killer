package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
)

type CFClient struct {
	client *cloudformation.Client
}

func InitializeCFClient() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("error while loading config: %v", err)
	}

	cf = CFClient{
		client: cloudformation.NewFromConfig(cfg),
	}
}
