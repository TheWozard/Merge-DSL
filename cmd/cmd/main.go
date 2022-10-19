package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	cmd  *cobra.Command
	vars *viper.Viper
}

func (c Config) Bind(name string, env []string) {
	if len(env) > 0 {
		c.vars.MustBindEnv(append([]string{name}, env...)...)
	}
	// c.cmd.
}

func main() {
	cmd := &cobra.Command{
		Use:   os.Args[0],
		Short: "An example of cobra",
		Long: `This application shows how to create modern CLI
	applications in go using Cobra CLI library`,
		Run: func(cmd *cobra.Command, args []string) {},
	}
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	provider := viper.New()
	provider.BindEnv("num", "TEST_NUM")
	fmt.Println(provider.GetInt("num"))
}
