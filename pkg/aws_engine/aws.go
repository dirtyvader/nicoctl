package aws_engine

import (
	"bytes"
	"context"
	_ "embed"
	"html/template"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type AwsModel struct {
	// paramsOut *ec2.DescribeAccountAttributesOutput
	// tpl       *template.Template
	viewport viewport.Model
	err      error
}

func NewAws() (*AwsModel, error) {
	// Init viewport
	vp := viewport.New(80, 20)
	vp.Style = vpStyle

	// Load template
	// temp, err := template.ParseFiles("./templates/aws.md.tpl")

	temp, err := template.New("aws").Parse(awsTplStr)
	if err != nil {
		return nil, err
	}

	// Init AWS Data
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	client := ec2.NewFromConfig(cfg)

	paramsIn := &ec2.DescribeAccountAttributesInput{}
	result, err := client.DescribeAccountAttributes(context.TODO(), paramsIn)
	if err != nil {
		return nil, err
	}

	content, err := setViewportContent_Aws(temp, &vp, result)
	if err != nil {
		return nil, err
	}
	w, h := lipgloss.Size(content)
	vp.Height = h
	vp.Width = w

	return &AwsModel{
		// paramsOut: result,
		// tpl:       temp,
		viewport: vp,
		err:      nil,
	}, nil
}

func setViewportContent_Aws(temp *template.Template, vp *viewport.Model, in *ec2.DescribeAccountAttributesOutput) (string, error) {
	var attr = make(map[string]string)
	if in != nil {
		for i := 0; i < len(in.AccountAttributes); i++ {
			k := *in.AccountAttributes[i].AttributeName
			v := *in.AccountAttributes[i].AttributeValues[0].AttributeValue
			attr[k] = v
		}
	}

	var tpl bytes.Buffer
	if err := temp.Execute(&tpl, attr); err != nil {
		return "", err
	}
	content := tpl.String()
	// fmt.Printf("\n%v\n", content)

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

func (m AwsModel) helpView() string {
	return helpBar()
}

func (m AwsModel) Init() tea.Cmd {
	return nil
}

func (m AwsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		return m, nil

	case errMsg:
		m.err = msg
		return m, tea.Quit

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit

		default:
			var cmd tea.Cmd
			m.viewport, cmd = m.viewport.Update(msg)
			return m, cmd
		}

	default:
		return m, nil
	}
}

func (m AwsModel) View() string {
	return m.viewport.View() + m.helpView()
}
