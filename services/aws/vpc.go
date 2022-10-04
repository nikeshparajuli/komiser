package aws

import (
	"context"

	awsConfig "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	. "github.com/mlabouardy/komiser/models/aws"
)

type AWSVPC struct {
	Name   string   `json:"name"`
	Region string   `json:"region"`
	ARN    string   `json:"arn"`
	Tags   []string `json:"tags"`
}

func (aws AWS) DescribeVPCsTotal(cfg awsConfig.Config) ([]AWSVPC, error) {
	vpcs := make([]AWSVPC, 0)
	regions, err := aws.getRegions(cfg)
	if err != nil {
		return vpcs, err
	}
	for _, region := range regions {
		vpcsRes, err := aws.getVPCs(cfg, region.Name)
		if err != nil {
			return vpcs, err
		}

		for _, vpc := range vpcsRes {
			vpcs = append(vpcs, AWSVPC{
				Name:   vpc.ID,
				Region: region.Name,
				ARN:    vpc.ID,
				Tags:   vpc.Tags,
			})
		}
	}
	return vpcs, nil
}

func (aws AWS) DescribeSubnets(cfg awsConfig.Config) (int64, error) {
	var sum int64
	regions, err := aws.getRegions(cfg)
	if err != nil {
		return sum, err
	}
	for _, region := range regions {
		cfg.Region = region.Name
		svc := ec2.NewFromConfig(cfg)
		res, err := svc.DescribeSubnets(context.Background(), &ec2.DescribeSubnetsInput{})
		if err != nil {
			return sum, err
		}
		sum += int64(len(res.Subnets))
	}
	return sum, nil
}

func (aws AWS) getVPCs(cfg awsConfig.Config, region string) ([]VPC, error) {
	cfg.Region = region
	svc := ec2.NewFromConfig(cfg)
	result, err := svc.DescribeVpcs(context.Background(), &ec2.DescribeVpcsInput{})
	if err != nil {
		return []VPC{}, err
	}
	listOfVPCs := make([]VPC, 0)
	for _, vpc := range result.Vpcs {
		vpcTags := make([]string, 0)
		for _, tag := range vpc.Tags {
			vpcTags = append(vpcTags, *tag.Value)
		}
		listOfVPCs = append(listOfVPCs, VPC{
			ID:        *vpc.VpcId,
			State:     string(vpc.State),
			CidrBlock: *vpc.CidrBlock,
			Tags:      vpcTags,
		})
	}
	return listOfVPCs, nil
}
