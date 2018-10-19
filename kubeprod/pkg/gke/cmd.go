/*
 * Bitnami Kubernetes Production Runtime - A collection of services that makes it
 * easy to run production workloads in Kubernetes.
 *
 * Copyright 2018 Bitnami
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package gke

import (
	"os"

	"github.com/spf13/cobra"

	kubeprodcmd "github.com/bitnami/kube-prod-runtime/kubeprod/cmd"
)

const (
	flagEmail             = "email"
	flagDNSSuffix         = "dns-zone"
	flagProject           = "project"
	flagOauthClientId     = "oauth-client-id"
	flagOauthClientSecret = "oauth-client-secret"
	flagOauthGoogleGroups = "oauth-google-groups"
)

var gkeCmd = &cobra.Command{
	Use:   "gke",
	Short: "Install Bitnami Production Runtime for GKE",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := kubeprodcmd.NewInstallSubcommand(cmd)
		if err != nil {
			return err
		}

		conf := GKEConfig{}
		c.PlatformConfig = &conf
		if err := c.ReadPlatformConfig(&conf); err != nil {
			return err
		}
		if err := config(cmd, &conf); err != nil {
			return err
		}
		if err := c.WritePlatformConfig(&conf); err != nil {
			return err
		}

		return c.Run(cmd.OutOrStdout())
	},
}

func init() {
	kubeprodcmd.InstallCmd.AddCommand(gkeCmd)

	gkeCmd.PersistentFlags().String(flagEmail, os.Getenv("EMAIL"), "Contact email for cluster admin")
	gkeCmd.PersistentFlags().String(flagDNSSuffix, "", "External DNS zone for public endpoints")
	gkeCmd.PersistentFlags().String(flagProject, "", "GCP project to use for managed resources")
	gkeCmd.PersistentFlags().String(flagOauthClientId, "", "Client ID to use for OAuth")
	gkeCmd.PersistentFlags().String(flagOauthClientSecret, "", "Client secret to use for OAuth")
	gkeCmd.PersistentFlags().StringSlice(flagOauthGoogleGroups, []string{}, "Google groups used to restrict OAuth access")
}
