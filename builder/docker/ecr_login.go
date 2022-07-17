//go:generate packer-sdc struct-markdown

package docker

import (
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go/service/ecrpublic"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	awsCredentials "github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	awsbase "github.com/hashicorp/aws-sdk-go-base"
	"github.com/hashicorp/go-cleanhttp"
)

type AwsAccessConfig struct {
	// The AWS access key used to communicate with
	// AWS. Learn how to set
	// this.
	AccessKey string `mapstructure:"aws_access_key" required:"false"`
	// The AWS secret key used to communicate with
	// AWS. Learn how to set
	// this.
	SecretKey string `mapstructure:"aws_secret_key" required:"false"`
	// The AWS access token to use. This is different from
	// the access key and secret key. If you're not sure what this is, then you
	// probably don't need it. This will also be read from the AWS_SESSION_TOKEN
	// environmental variable.
	Token string `mapstructure:"aws_token" required:"false"`
	// The AWS shared credentials profile used to
	// communicate with AWS. Learn how to set
	// this.
	Profile string `mapstructure:"aws_profile" required:"false"`
}

type ECRType string

const (
	Public  ECRType = "public"
	Private ECRType = "private"
	Invalid ECRType = "invalid"
)

const EcrPublicHost = "public.ecr.aws"

// EcrPublicApiRegion : The Amazon ECR Public registry requires authentication in the us-east-1 Region,
// so you need to specify --region us-east-1 each time you authenticate
const EcrPublicApiRegion = "us-east-1"

// GetEcrType : Get ECR type (Public or Private) based on the given URL.
// If the URL can't be parsed the function returns Invalid.
func (c *AwsAccessConfig) GetEcrType(ecrUrl string) (ECRType, error) {
	_, err := url.ParseRequestURI(ecrUrl)
	if err != nil {
		return Invalid, err
	}

	u, err := url.Parse(ecrUrl)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return Invalid, err
	}

	if u.Host == EcrPublicHost {
		return Public, nil
	}
	return Private, nil
}

// PublicEcrLogin : Get a login token for Amazon AWS ECR Public. Returns username and password
// or an error.
func (c *AwsAccessConfig) PublicEcrLogin(ecrUrl string) (string, string, error) {
	config := aws.NewConfig().WithCredentialsChainVerboseErrors(true)
	config = config.WithRegion(EcrPublicApiRegion)

	config = config.WithHTTPClient(cleanhttp.DefaultClient())
	transport := config.HTTPClient.Transport.(*http.Transport)
	transport.Proxy = http.ProxyFromEnvironment

	// Figure out which possible credential providers are valid; test that we
	// can get credentials via the selected providers, and set the providers in
	// the config.
	creds, err := c.GetCredentials(config)
	if err != nil {
		return "", "", fmt.Errorf(err.Error())
	}
	config.WithCredentials(creds)

	// Create session options based on our AWS config
	opts := session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config:            *config,
	}

	if c.Profile != "" {
		opts.Profile = c.Profile
	}

	sess, err := session.NewSessionWithOptions(opts)
	if err != nil {
		return "", "", err
	}
	session := sess

	cp, err := session.Config.Credentials.Get()
	if err != nil {
		return "", "", fmt.Errorf("failed to create session: %s", err)
	}
	log.Printf("[INFO] AWS authentication used: %q", cp.ProviderName)

	service := ecrpublic.New(session)
	params := &ecrpublic.GetAuthorizationTokenInput{}

	resp, err := service.GetAuthorizationToken(params)
	if err != nil {
		return "", "", fmt.Errorf(err.Error())
	}

	auth, err := base64.StdEncoding.DecodeString(*resp.AuthorizationData.AuthorizationToken)
	if err != nil {
		return "", "", fmt.Errorf("error decoding ECR Public AuthorizationToken: %s", err)
	}

	authParts := strings.SplitN(string(auth), ":", 2)
	log.Printf("Successfully got login for ECR Public: %s", ecrUrl)

	username := authParts[0]
	password := authParts[1]

	return username, password, nil
}

// EcrGetLogin Get a login token for Amazon AWS ECR. Returns username and password
// or an error.
func (c *AwsAccessConfig) EcrGetLogin(ecrUrl string) (string, string, error) {

	// Check ECR Type
	ecrType, parsingErr := c.GetEcrType(ecrUrl)
	if parsingErr != nil {
		errMsg := "failed to parse the ECR URL: %v" +
			"\n%v" +
			"\nit should be either on the form `public.ecr.aws/<registry_alias>/<registry_name>` or " +
			"`<account number>.dkr.ecr.<region>.amazonaws.com`"
		return "", "", fmt.Errorf(errMsg, ecrUrl, parsingErr)
	}

	if ecrType == Public {
		return c.PublicEcrLogin(ecrUrl)
	}

	exp := regexp.MustCompile(`(?:http://|https://|)([0-9]*)\.dkr\.ecr\.(.*)\.amazonaws\.com.*`)
	splitUrl := exp.FindStringSubmatch(ecrUrl)
	if len(splitUrl) != 3 {
		return "", "", fmt.Errorf("Failed to parse the ECR URL: %s it should be on the form <account number>.dkr.ecr.<region>.amazonaws.com", ecrUrl)
	}
	accountId := splitUrl[1]
	region := splitUrl[2]

	log.Println(fmt.Sprintf("Getting ECR token for account: %s in %s..", accountId, region))

	// Create new AWS config
	config := aws.NewConfig().WithCredentialsChainVerboseErrors(true)
	config = config.WithRegion(region)

	config = config.WithHTTPClient(cleanhttp.DefaultClient())
	transport := config.HTTPClient.Transport.(*http.Transport)
	transport.Proxy = http.ProxyFromEnvironment

	// Figure out which possible credential providers are valid; test that we
	// can get credentials via the selected providers, and set the providers in
	// the config.
	creds, err := c.GetCredentials(config)
	if err != nil {
		return "", "", fmt.Errorf(err.Error())
	}
	config.WithCredentials(creds)

	// Create session options based on our AWS config
	opts := session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config:            *config,
	}

	if c.Profile != "" {
		opts.Profile = c.Profile
	}

	sess, err := session.NewSessionWithOptions(opts)
	if err != nil {
		return "", "", err
	}
	log.Printf("Found region %s", *sess.Config.Region)
	session := sess

	cp, err := session.Config.Credentials.Get()

	if err != nil {
		return "", "", fmt.Errorf("failed to create session: %s", err)
	}

	log.Printf("[INFO] AWS authentication used: %q", cp.ProviderName)

	service := ecr.New(session)
	params := &ecr.GetAuthorizationTokenInput{
		RegistryIds: []*string{
			aws.String(accountId),
		},
	}
	resp, err := service.GetAuthorizationToken(params)
	if err != nil {
		return "", "", fmt.Errorf(err.Error())
	}

	auth, err := base64.StdEncoding.DecodeString(*resp.AuthorizationData[0].AuthorizationToken)
	if err != nil {
		return "", "", fmt.Errorf("Error decoding ECR AuthorizationToken: %s", err)
	}

	authParts := strings.SplitN(string(auth), ":", 2)
	log.Printf("Successfully got login for ECR: %s", ecrUrl)

	return authParts[0], authParts[1], nil
}

// GetCredentials gets credentials from the environment, shared credentials,
// the session (which may include a credential process), or ECS/EC2 metadata
// endpoints. GetCredentials also validates the credentials and the ability to
// assume a role or will return an error if unsuccessful.
func (c *AwsAccessConfig) GetCredentials(config *aws.Config) (*awsCredentials.Credentials, error) {
	// Reload values into the config used by the Packer-Terraform shared SDK
	awsbaseConfig := &awsbase.Config{
		AccessKey:    c.AccessKey,
		DebugLogging: false,
		Profile:      c.Profile,
		SecretKey:    c.SecretKey,
		Token:        c.Token,
	}

	return awsbase.GetCredentials(awsbaseConfig)
}
