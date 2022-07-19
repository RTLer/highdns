/*
Copyright © 2022 Hooman Naghiee <me@hoomanna.com>

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
	"github.com/RTLer/highdns/dns"
	"github.com/RTLer/highdns/dns/cache"
	"github.com/RTLer/highdns/dns/client"
	"github.com/spf13/cobra"
)

// servCmd represents the serv command
var servCmd = &cobra.Command{
	Use:   "serv",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		pc := &client.DNS{}
		pc.Init()

		c := &cache.Cache{
			Acquire: pc.Acquire,
		}
		if err := c.Open(); err != nil {
			panic(err)
		}
		go c.Queue.Run(c)
		server := dns.DNSServer{
			Cache: c,
		}
		server.Serve()
	},
}

func init() {
	rootCmd.AddCommand(servCmd)
}
