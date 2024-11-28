package clis

import (
	"github.com/gelleson/changescout/changescout/pkg/flags"
	"time"
)

var (
	FlagsLogLevel = flags.NewStringFlag("log-level",
		flags.WithDefaultValue[string]("info"),
		flags.WithCategory[string]("logging"),
		flags.WithUsage[string]("The log level"),
		flags.WithEnvVars[string]("CS_LOG_LEVEL"))

	FlagsPort = flags.NewStringFlag("port",
		flags.WithDefaultValue[string](":3311"),
		flags.WithAlias[string]("p"),
		flags.WithCategory[string]("http"),
		flags.WithEnvVars[string]("CS_PORT"),
		flags.WithUsage[string]("The port to listen on"))

	FlagsSecret = flags.NewStringFlag("secret",
		flags.WithRequired[string](true),
		flags.WithAlias[string]("s"),
		flags.WithCategory[string]("security"),
		flags.WithEnvVars[string]("CS_JWT_SECRET"),
		flags.WithUsage[string]("The secret to use for JWT"))

	FlagsDBUrl = flags.NewStringFlag("db-url",
		flags.WithDefaultValue[string](":memory:?_pragma=foreign_keys(1)"),
		flags.WithAlias[string]("d"),
		flags.WithCategory[string]("db"),
		flags.WithEnvVars[string]("CS_DB_URL"),
		flags.WithUsage[string]("The database url"))

	FlagsDBEngine = flags.NewStringFlag("db-engine",
		flags.WithDefaultValue[string]("sqlite"),
		flags.WithAlias[string]("e"),
		flags.WithCategory[string]("db"),
		flags.WithEnvVars[string]("CS_DB_ENGINE"),
		flags.WithUsage[string]("The database engine"))

	FlagsBrokerEnabled = flags.NewBoolFlag("broker-enabled",
		flags.WithDefaultValue[bool](true),
		flags.WithCategory[bool]("broker"),
		flags.WithEnvVars[bool]("CS_BROKER_ENABLED"),
		flags.WithUsage[bool]("Enable the broker"))

	FlagsSchedulerEnabled = flags.NewBoolFlag("scheduler-enabled",
		flags.WithDefaultValue[bool](true),
		flags.WithCategory[bool]("scheduler"),
		flags.WithEnvVars[bool]("CS_SCHEDULER_ENABLED"),
		flags.WithUsage[bool]("Enable the scheduler"))

	FlagsSchedulerInterval = flags.NewDurationFlag("scheduler-interval",
		flags.WithDefaultValue[time.Duration](time.Second*5),
		flags.WithAlias[time.Duration]("si"),
		flags.WithCategory[time.Duration]("scheduler"),
		flags.WithEnvVars[time.Duration]("CS_SCHEDULER_INTERVAL"),
		flags.WithUsage[time.Duration]("The interval to check for due websites"))
)
