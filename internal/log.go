package internal

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

func LoggingSetup(cmd *cobra.Command, args []string) {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	level, err := log.ParseLevel(viper.GetString(EnvVarLogLevel))
	if err != nil {
		log.Warnf("Failed to parse log level from config: %v", err)
		level = DefaultLogLevel
	}
	log.SetLevel(level)
	level = log.GetLevel()
	log.Debugf("Log level set to %s", level)
}
