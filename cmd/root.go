/*
Copyright © 2021 Jose Ramon Mañes jr.mb47@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "axectl",
	Short: "axectl CLI tool",
	Long: `
+--------------------------------------------------+
|                                                  |
|              .dx:.                               |
|             cO00Ok;                              |
|           .x00OOOkkx.                            |
|          .d00OOOkkkxl                            |
|         :kOOOOkkkxl;'                            |
|      .l00OOOkkkxd:;cc'                           |
|    'd00OOOOkkkdc',:lol,                          |
|  ,d000OOOkkkxd,..',:loo:                         |
|.d000OOOkkkxxxl   ..,:codl.                       |
|xO0OOOkkkxxxxd,      .;cldo;                      |
|:kkOOkkkxxxddo.        ':codl'                    |
| ;xxkkxxxdddd:          .,:lddl.                  |
|  .cxddddddoo.            .;codd;                 |
|    .:ddooooc               ':lddl'               |
|       .,clo,                .,coddc.             |
|                               .:lddo,.           |
|                                 ,clddl,          |
|                                  .;lodo:.        |
|                                    ,clddl,       |
|                                     .;codoc.     |
|                                       ':lddo:.   |
|                                        .;codoc,  |
|                                         .,:lool:.|
|                                           .;cloo:|
|                                            .,:c;.|
|                                                  |
|                                                  |
+--------------------------------------------------+
--------------------------------------------------
axectl is a set of tools for DevOps/SRE.
-------------------------------------------------`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.piktocºtl.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	viper.SetDefault("author", "Jose Ramon Mañes")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		configHome := filepath.Join(home, "/.axectl/")
		configName := "config"
		configType := "yml"
		configPath := filepath.Join(configHome, configName+"."+configType)

		err = CreateFileInPath(configHome, configPath)
		if err != nil {
			log.Fatal(err)
		}

		viper.AddConfigPath(configHome)
		viper.SetConfigName(configName)
		viper.SetConfigType(configType)

		if _, err := os.Stat(configHome); os.IsNotExist(err) {
			err = os.MkdirAll(configHome, 0764)
			log.Fatal(err)
		}
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	//if err := viper.ReadInConfig(); err == nil {
	//	fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	//}
}

// CreateFileInPath Create a file in a path
func CreateFileInPath(configHome, configPath string) error {
	// Check if the path does not exist
	if _, err := os.Stat(configHome); os.IsNotExist(err) {
		// Create the path
		err := os.MkdirAll(configHome, 0764)
		if err != nil {
			log.Fatal("[ERROR] ", err)
			return err
		}
	}

	// Check if the file exists or not, create if not exists
	_, err := os.OpenFile(configPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0764)
	if err != nil {
		log.Fatal("[ERROR] ", err)
		return err
	}
	return nil
}
