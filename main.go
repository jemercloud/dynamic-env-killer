package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation"

	yaml "sigs.k8s.io/yaml/goyaml.v2"
)

type Config struct {
	Environments []struct {
		ID         string `yaml:"id"`
		Containers []struct {
			Container    string `yaml:"container"`
			TargetBranch string `yaml:"targetBranch,omitempty"`
		} `yaml:"containers"`
	} `yaml:"environments"`
}

var cf CFClient
var stackPrefix string = "Samply-ecs-service-dynamic-"

func getDeployedDynamicEnvs() []string {
	response, err := cf.client.DescribeStacks(context.TODO(), &cloudformation.DescribeStacksInput{})
	var deployedDynamicEnvs []string

	if err != nil {
		log.Fatalf("error while describing stacks: %v", err)
	}

	for _, stack := range response.Stacks {
		if strings.HasPrefix(*stack.StackName, stackPrefix) {
			deployedDynamicEnvs = append(deployedDynamicEnvs, *stack.StackName)
		}
	}

	return deployedDynamicEnvs
}

func getIDsInConfigFile(fileName string) []string {
	configYml := Config{}
	data, err := os.ReadFile(fileName)

	if err != nil {
		log.Fatalf("error while reading yaml: %v", err)
	}

	err = yaml.Unmarshal([]byte(data), &configYml)

	if err != nil {
		log.Fatalf("error while parsing yaml: %v", err)
	}

	var ids []string

	for _, env := range configYml.Environments {
		ids = append(ids, env.ID)
	}

	return ids
}

func main() {
	fileName := os.Getenv("INPUT_FILENAME")
	log.Println("using file: ", fileName)
	InitializeCFClient()

	deployedDynamicEnvs := getDeployedDynamicEnvs()
	idsOfDynamicEnvs := getIDsInConfigFile(fileName)
	var stackNamesOfDeployedEnvs []string

	for _, env := range idsOfDynamicEnvs {
		stackName := fmt.Sprintf("%s%s", stackPrefix, env)
		stackNamesOfDeployedEnvs = append(stackNamesOfDeployedEnvs, stackName)
	}

	log.Println("Theoretical stack names: ", stackNamesOfDeployedEnvs)
	log.Println("Deployed stacks: ", deployedDynamicEnvs)

	for _, stack := range deployedDynamicEnvs {
		if !contains(stackNamesOfDeployedEnvs, stack) {
			log.Println("Deleting stack: ", stack)
			_, err := cf.client.DeleteStack(context.TODO(), &cloudformation.DeleteStackInput{
				StackName: &stack,
			})

			if err != nil {
				log.Fatalf("error while deleting stack: %v", err)
			}
		}
	}
}
