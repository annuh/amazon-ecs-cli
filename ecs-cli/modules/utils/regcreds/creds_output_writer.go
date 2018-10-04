// Copyright 2015-2018 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//	http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package readers

import (
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

const (
	// ECSCredFileTimeFmt is the timestamp format to use on 'registry-creds up' outputs
	ECSCredFileTimeFmt  = "20060102T150405Z"
	ecsCredFileBaseName = "ecs-registry-creds"
)

// CredsOutputEntry contains the credential ARN, key, and associated container names for a registry
type CredsOutputEntry struct {
	CredentialARN  string   `yaml:"secret_manager_arn"`
	KMSKeyID       string   `yaml:"kms_key_id,omitempty"`
	ContainerNames []string `yaml:"container_names"`
}

// RegistryCredsOutput contains the content of the output file
type RegistryCredsOutput struct {
	Version             string
	CredentialResources CredResources `yaml:"registry_credential_outputs"`
}

// CredResources contains the
type CredResources struct {
	TaskExecutionRole    string                      `yaml:"task_execution_role"`
	ContainerCredentials map[string]CredsOutputEntry `yaml:"container_credentials"`
}

// GenerateCredsOutput marshals credential output JSON into YAML and outputs it to a file
func GenerateCredsOutput(creds map[string]CredsOutputEntry, roleName, outputDir string, policyCreatTime *time.Time) error {
	outputResources := CredResources{
		ContainerCredentials: creds,
		TaskExecutionRole:    roleName,
	}
	regOutput := RegistryCredsOutput{
		Version:             "1",
		CredentialResources: outputResources,
	}
	credBytes, err := yaml.Marshal(regOutput)
	if err != nil {
		return err
	}

	outputFileDir := outputDir
	if outputFileDir == "" {
		wdir, err := os.Getwd()
		if err != nil {
			return err
		}
		outputFileDir = wdir
	}

	timeStamp := time.Now().UTC()
	if policyCreatTime != nil {
		timeStamp = *policyCreatTime
	}
	timestampedSuffix := fmt.Sprintf("_%s.yml", timeStamp.Format(ECSCredFileTimeFmt))

	file, err := os.Create(outputFileDir + string(os.PathSeparator) + ecsCredFileBaseName + timestampedSuffix)
	if err != nil {
		return err
	}
	defer file.Close()

	log.Info("Writing registry credential output to new file " + file.Name())
	_, err = file.Write(credBytes)
	if err != nil {
		return err
	}

	return nil
}

// BuildOutputEntry returns a CredsOutputEntry with the provided parameters
func BuildOutputEntry(arn string, key string, containers []string) CredsOutputEntry {
	return CredsOutputEntry{
		CredentialARN:  arn,
		KMSKeyID:       key,
		ContainerNames: containers,
	}
}