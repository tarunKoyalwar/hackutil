package burpconfig

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tarunKoyalwar/hackutil/pkg/utils"
)

type Options struct {
	Protocols  string
	Ports      string
	Inscope    string
	Outofscope string
	Extra      string
	Output     string
}

var defaultOutOfScope = []string{
	"mozilla",
	"postog",
	"sentry\\.io",
	"stripe\\.com",
	"stripe\\.network",
}

var defaultExtra = []string{
	"firebaseio\\.com",
	"amazonaws\\.com",
}

func NewBurpConfigCmd() *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:   "burpconfig",
		Short: "Generate Burp Suite configuration file",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runBurpConfig(opts)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&opts.Protocols, "proto", "https", "Protocols (comma-separated)")
	flags.StringVar(&opts.Ports, "port", "443", "Ports (comma-separated)")
	flags.StringVar(&opts.Inscope, "inscope", "", "In-scope domains (file path or comma-separated values)")
	flags.StringVar(&opts.Outofscope, "outofscope", strings.Join(defaultOutOfScope, ","), "Out-of-scope domains")
	flags.StringVar(&opts.Extra, "extra", strings.Join(defaultExtra, ","), "Extra in-scope domains")
	flags.StringVar(&opts.Output, "output", "burp-config.json", "Output file path")

	return cmd
}

func runBurpConfig(opts *Options) error {
	// Parse protocols
	protocols := strings.Split(opts.Protocols, ",")

	// Parse ports
	ports := strings.Split(opts.Ports, ",")

	// Get inscope domains
	inscope, err := utils.GetInputList(opts.Inscope)
	if err != nil {
		return fmt.Errorf("failed to process inscope input: %w", err)
	}

	// Get outofscope domains
	outofscope, err := utils.GetInputList(opts.Outofscope)
	if err != nil {
		return fmt.Errorf("failed to process outofscope input: %w", err)
	}

	// Get extra domains and append to inscope
	extra, err := utils.GetInputList(opts.Extra)
	if err != nil {
		return fmt.Errorf("failed to process extra input: %w", err)
	}
	inscope = append(inscope, extra...)

	// Create config
	config := createBurpConfig(protocols, ports, inscope, outofscope)

	// Write to file
	data, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(opts.Output, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func createBurpConfig(protocols, ports, inscope, outofscope []string) BurpConfig {
	config := BurpConfig{
		Target: Target{
			Scope: Scope{
				AdvancedMode: true,
				Include:      make([]Rule, 0),
				Exclude:      make([]Rule, 0),
			},
		},
	}

	// Add include rules
	for _, domain := range inscope {
		for _, proto := range protocols {
			for _, port := range ports {
				config.Target.Scope.Include = append(config.Target.Scope.Include, Rule{
					Enabled:  true,
					File:     "^/.*",
					Host:     fmt.Sprintf("^.*%s.*$", domain),
					Port:     fmt.Sprintf("^%s$", port),
					Protocol: proto,
				})
			}
		}
	}

	// Add exclude rules
	for _, domain := range outofscope {
		for _, proto := range protocols {
			for _, port := range ports {
				config.Target.Scope.Exclude = append(config.Target.Scope.Exclude, Rule{
					Enabled:  true,
					File:     "^/.*",
					Host:     fmt.Sprintf("^.*%s.*$", domain),
					Port:     fmt.Sprintf("^%s$", port),
					Protocol: proto,
				})
			}
		}
	}

	return config
}
