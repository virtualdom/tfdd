package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/pkg/errors"
)

const (
	noCreds = iota
	profileCreds
	accessKeyCreds
	accessKeyWithSessionCreds
)

// NewSTSClient returns an STS client initialized using the creds passed
// in. If creds is empty, depend on the default AWS SDK behavior (env,
// instance profile, etc.). If creds has one element, it's the profile.
// If it has two or three elements, the first is an IAM access key ID, the second
// is an IAM secret access key, and the optional third is the session token.
func NewSTSClient(creds ...string) (*sts.STS, error) {
	var sess *session.Session
	var err error

	switch credentialType := len(creds); credentialType {
	case noCreds:
		sess, err = session.NewSession()

	case profileCreds:
		sess, err = session.NewSessionWithOptions(session.Options{
			Profile: creds[0],
		})

	case accessKeyCreds:
		creds = append(creds, "")
		fallthrough

	case accessKeyWithSessionCreds:
		sess, err = session.NewSession(&aws.Config{
			Credentials: credentials.NewStaticCredentials(creds[0], creds[1], creds[2]),
		})

	default:
		err = errors.New("incorrect value for AWS credentials used")
	}

	if err != nil {
		return nil, errors.Wrap(err, "failed to init AWS session")
	}

	return sts.New(sess), nil
}
