package aws_engine

import (
	_ "embed"

	"github.com/charmbracelet/lipgloss"
)

var (
	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render
	vpStyle   = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("202"))

	//go:embed templates/aws.md.tpl
	awsTplStr string

	//go:embed templates/aws_instances.md.tpl
	awsInstTplStr string
)

type errMsg struct{ err error }
type instance struct {
	InstanceID     string
	InstanceState  string
	InstanceType   string
	PublicDNSName  string
	PublicIPv4     string
	SecurityGroups string
	KeyName        string
}

func helpBar() string {
	return helpStyle("\n  ↑/↓: Navigate • q: Quit\n")
}
