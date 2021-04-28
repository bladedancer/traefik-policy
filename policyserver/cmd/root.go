package cmd

import (
	"strings"

	"github.com/bladedancer/traefik-policy/policyserver/pkg/log"
	"github.com/bladedancer/traefik-policy/policyserver/pkg/policyservice"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "policyserver",
	Short: "Policy Server.",
	Long:  `Implementation of the Policy as a Service interface.`,
	RunE:  run,
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initRootConfig)

	rootCmd.PersistentFlags().String("logLevel", "info", "log level")
	rootCmd.PersistentFlags().String("logFormat", "json", "line or json")
	bindOrPanic("log.level", rootCmd.PersistentFlags().Lookup("logLevel"))
	bindOrPanic("log.format", rootCmd.PersistentFlags().Lookup("logFormat"))

	rootCmd.PersistentFlags().String("host", "", "host interface to bind to")
	rootCmd.PersistentFlags().Uint32("port", 9000, "port to listen on")
	bindOrPanic("host", rootCmd.PersistentFlags().Lookup("host"))
	bindOrPanic("port", rootCmd.PersistentFlags().Lookup("port"))

	rootCmd.PersistentFlags().String("key", "", "the path to the private key")
	rootCmd.PersistentFlags().String("cert", "", "the path to the cert")
	bindOrPanic("key", rootCmd.PersistentFlags().Lookup("key"))
	bindOrPanic("cert", rootCmd.PersistentFlags().Lookup("cert"))

	rootCmd.PersistentFlags().Bool("headerPolicy", false, "Apply the header policy")
	rootCmd.PersistentFlags().Bool("bodyPolicy", false, "Apply the body policy")
	bindOrPanic("headerPolicy", rootCmd.PersistentFlags().Lookup("headerPolicy"))
	bindOrPanic("bodyPolicy", rootCmd.PersistentFlags().Lookup("bodyPolicy"))

}

func initRootConfig() {
	viper.SetTypeByDefaultValue(true)
	viper.SetEnvPrefix("pe")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	err := log.SetLogLevel(viper.GetString("log.level"), viper.GetString("log.format"))
	if err != nil {
		panic(err)
	}
}

func bindOrPanic(key string, flag *pflag.Flag) {
	if err := viper.BindPFlag(key, flag); err != nil {
		panic(err)
	}
}

func run(cmd *cobra.Command, args []string) error {
	// Create the file source and run the server command
	ps := &policyservice.PolicyServer{
		Config: &policyservice.PolicyServiceConfig{
			Host:         viper.GetString("host"),
			Port:         viper.GetUint32("port"),
			CertFile:     viper.GetString("cert"),
			KeyFile:      viper.GetString("key"),
			HeaderPolicy: viper.GetBool("headerPolicy"),
			BodyPolicy:   viper.GetBool("bodyPolicy"),
		},
	}

	return ps.Run()
}
