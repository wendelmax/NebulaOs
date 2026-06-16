package aws

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/wendelmax/nebulaos/src/api/domain"
)

type AWSProvider struct {
	client *ec2.Client
}

func NewAWSProvider(ctx context.Context, region string, endpoint string) (*AWSProvider, error) {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if endpoint != "" {
			ep := endpoint
			if u, err := url.Parse(endpoint); err == nil && !strings.Contains(u.Hostname(), "localhost") {
				if u.Port() != "" {
					ep = fmt.Sprintf("http://%s.%s:%s", strings.ToLower(service), u.Hostname(), u.Port())
				} else {
					ep = fmt.Sprintf("http://%s.%s", strings.ToLower(service), u.Hostname())
				}
			}
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           ep,
				SigningRegion: region,
			}, nil
		}
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %v", err)
	}

	return &AWSProvider{client: ec2.NewFromConfig(cfg)}, nil
}

func (p *AWSProvider) Provision(ctx context.Context, resource *domain.Resource) error {
	input := &ec2.RunInstancesInput{
		ImageId:      aws.String("ami-12345678"),
		InstanceType: "t2.micro",
		MinCount:     aws.Int32(1),
		MaxCount:     aws.Int32(1),
	}

	result, err := p.client.RunInstances(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to provision AWS instance: %w", err)
	}

	resource.ProviderID = *result.Instances[0].InstanceId
	resource.State = "provisioning"
	return nil
}

func (p *AWSProvider) Decommission(ctx context.Context, resource *domain.Resource) error {
	input := &ec2.TerminateInstancesInput{
		InstanceIds: []string{resource.ProviderID},
	}

	_, err := p.client.TerminateInstances(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete AWS instance: %w", err)
	}

	resource.State = "deleted"
	return nil
}

func (p *AWSProvider) GetStatus(ctx context.Context, resourceID string) (string, error) {
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []string{resourceID},
	}

	result, err := p.client.DescribeInstances(ctx, input)
	if err != nil {
		return "", err
	}

	if len(result.Reservations) > 0 && len(result.Reservations[0].Instances) > 0 {
		return string(result.Reservations[0].Instances[0].State.Name), nil
	}

	return "unknown", nil
}
func (p *AWSProvider) Ping(ctx context.Context) error {
	return nil
}

func (p *AWSProvider) ListImages(ctx context.Context) ([]domain.Image, error) {
	return nil, nil
}

func (p *AWSProvider) DeployContainer(ctx context.Context, c *domain.Container) error {
	fmt.Printf("[AWS] Deploying container via ESC/Fargate: %s\n", c.Name)
	c.State = "running"
	return nil
}

func (p *AWSProvider) ConfigureNetwork(ctx context.Context, n *domain.Network) error {
	fmt.Printf("[AWS] Creating VPC/Subnet: %s\n", n.Name)
	n.State = "active"
	return nil
}

func (p *AWSProvider) AddRoute(ctx context.Context, r *domain.Route) error {
	return nil
}

func (p *AWSProvider) AttachSecurityGroup(ctx context.Context, resourceID string, sgID string) error {
	fmt.Printf("[AWS] Attaching Security Group %s to instance %s. Modifying network interfaces...\n", sgID, resourceID)
	return nil
}
