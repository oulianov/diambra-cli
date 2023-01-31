/*
 * Copyright 2022 The DIAMBRA Authors
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package agent

import (
	"fmt"
	"os"

	"github.com/diambra/cli/pkg/diambra"
	"github.com/diambra/cli/pkg/log"
	"github.com/go-kit/log/level"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func NewSubmitCmd(logger *log.Logger) *cobra.Command {
	dump := false
	submissionConfig := diambra.NewSubmissionConfig(logger)
	c, err := diambra.NewConfig(logger)
	if err != nil {
		level.Error(logger).Log("msg", err.Error())
		os.Exit(1)
	}

	cmd := &cobra.Command{
		Use:   "submit [--manifest submission-manifest.yaml | docker-image]",
		Short: "Submits an agent for evaluation",
		Long:  `This takes a docker image or submission manifest and submits it for evaluation.`,
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				submissionConfig.Image = args[0]
			} else if submissionConfig.ManifestPath == "" {
				level.Error(logger).Log("msg", "either image or manifest path must be provided")
				os.Exit(1)
			}
			submission, err := submissionConfig.Submission()
			if err != nil {
				level.Error(logger).Log("msg", "failed to configure manifest", "err", err.Error())
				os.Exit(1)
			}
			if dump {
				b, err := yaml.Marshal(submission)
				if err != nil {
					level.Error(logger).Log("msg", "failed to marshal manifest", "err", err.Error())
					os.Exit(1)
				}
				fmt.Println(string(b))
				return
			}
			id, err := diambra.Submit(logger, c.CredPath, submission)
			if err != nil {
				level.Error(logger).Log("msg", "failed to submit agent", "err", err.Error())
				os.Exit(1)
			}
			level.Info(logger).Log("msg", fmt.Sprintf("Agent submitted: https://diambra.ai/submission/%d", id), "id", id)
		},
	}
	submissionConfig.AddFlags(cmd.Flags())
	cmd.Flags().BoolVar(&dump, "dump", false, "Dump the manifest to stdout instead of submitting")
	return cmd
}
