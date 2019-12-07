package auth

import (
	"fmt"
	"os"
	"path"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/go-ini/ini"
	"github.com/pkg/errors"
	"github.com/virtualdom/tfdd/pkg/config"
)

var fileMutex = &sync.Mutex{}

const (
	accessKeyID     = "aws_access_key_id"
	secretAccessKey = "aws_secret_access_key" // nolint: gosec
	sessionToken    = "aws_session_token"
)

// AuthI is meant to be run as a goroutine that authenticates after pausing for
// the given interval in seconds.
func AuthI(cfg *config.Config, sts stsiface.STSAPI, roleName string, interval int) {
	for {
		time.Sleep(time.Duration(interval) * time.Second)
		Auth(cfg, sts, roleName)
	}
}

// Auth loops through each profile in the Config and fetches temporary IAM
// credentials before writing them all to the AWS credentials file.
func Auth(cfg *config.Config, stsClient stsiface.STSAPI, roleName string) {
	creds := make(map[string]*sts.Credentials)

	for _, profile := range cfg.Profiles {
		roleArn := fmt.Sprintf("arn:aws:iam::%s:role/%s", profile.AccountID, roleName)
		cred, err := getCreds(roleArn, stsClient)
		if err != nil {
			fmt.Printf("failed to get credentials for %s: %+v\n", roleArn, err)
		} else {
			creds[profile.Name] = cred
		}
	}

	writeCreds(creds)
}

func getCreds(roleArn string, svc stsiface.STSAPI) (*sts.Credentials, error) {
	input := &sts.AssumeRoleInput{
		RoleArn:         aws.String(roleArn),
		RoleSessionName: aws.String("tfdd"),
	}

	output, err := svc.AssumeRole(input)
	if err != nil {
		return nil, errors.Wrap(err, "failed to assume role "+roleArn)
	}

	return output.Credentials, nil
}

func writeCreds(creds map[string]*sts.Credentials) error {
	fileMutex.Lock()
	defer fileMutex.Unlock()

	credsFilepath := path.Join(os.Getenv("HOME"), ".aws", "credentials")

	file, err := ini.Load(credsFilepath)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		file = ini.Empty()
	}

	for profileName, cred := range creds {
		section := file.Section(profileName)
		for _, k := range section.Keys() {
			section.DeleteKey(k.Name())
		}

		_, err = section.NewKey(accessKeyID, *cred.AccessKeyId)
		if err != nil {
			return err
		}
		_, err = section.NewKey(secretAccessKey, *cred.SecretAccessKey)
		if err != nil {
			return err
		}
		_, err = section.NewKey(sessionToken, *cred.SessionToken)
		if err != nil {
			return err
		}
	}

	err = file.SaveTo(credsFilepath)
	if err != nil {
		return err
	}

	return nil
}
