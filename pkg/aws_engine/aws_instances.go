package aws_engine

import (
	"bytes"
	"context"
	"html/template"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

func NewAwsInst() (*AwsModel, error) {
	vp := viewport.New(200, 20)
	vp.Style = vpStyle

	temp, err := template.New("aws_instances").Parse(awsInstTplStr)
	if err != nil {
		return nil, err
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	client := ec2.NewFromConfig(cfg)

	paramsIn := &ec2.DescribeInstancesInput{}
	result, err := client.DescribeInstances(context.TODO(), paramsIn)
	if err != nil {
		return nil, err
	}

	content, err := setViewportContent_AwsInst(temp, &vp, result)
	if err != nil {
		return nil, err
	}
	_, vp.Height = lipgloss.Size(content)

	return &AwsModel{
		viewport: vp,
		err:      nil,
	}, nil

}

func setViewportContent_AwsInst(temp *template.Template, vp *viewport.Model, in *ec2.DescribeInstancesOutput) (string, error) {
	attr := make([]instance, 0)
	if in != nil {
		for i := 0; i < len(in.Reservations); i++ {
			for j := 0; j < len(in.Reservations[i].Instances); j++ {
				inst := in.Reservations[i].Instances[j]
				instanceId := *inst.InstanceId
				instanceState := inst.State.Name
				instanceType := inst.InstanceType
				publicIpV4 := *inst.PublicIpAddress
				keyName := *inst.KeyName
				attr = append(attr, instance{
					InstanceID:    instanceId,
					InstanceState: string(instanceState),
					InstanceType:  string(instanceType),
					PublicIPv4:    publicIpV4,
					KeyName:       keyName,
				})
			}
		}
	}

	var tpl bytes.Buffer
	if err := temp.Execute(&tpl, attr); err != nil {
		return "", err
	}
	content := tpl.String()

	renderer, err := glamour.NewTermRenderer(glamour.WithStylePath("dark"))
	if err != nil {
		return "", err
	}

	str, err := renderer.Render(content)
	if err != nil {
		return "", err
	}

	vp.SetContent(str)
	return str, nil
}
