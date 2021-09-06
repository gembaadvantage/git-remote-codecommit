/*
Copyright (c) 2021 Gemba Advantage

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package main

import (
	"context"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/gembaadvantage/codecommit-sign/pkg/awsv4"
	"github.com/gembaadvantage/codecommit-sign/pkg/translate"
	"github.com/spf13/cobra"
)

var (
	grcRgx = regexp.MustCompile(`^codecommit::(.+)://(.+)$`)
)

type gitOptions struct {
	Cmd       string
	RemoteURL string
}

func newRootCmd(out io.Writer, args []string) *cobra.Command {
	opts := gitOptions{}

	cmd := &cobra.Command{
		Use:          "git-remote-codecommit",
		Short:        "TODO",
		Args:         cobra.ExactArgs(2),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Cmd = args[0]

			// Prefix protocol, as it will have been stipped off by the git client
			opts.RemoteURL = "codecommit::" + args[1]

			return opts.Run(out)
		},
	}

	cmd.AddCommand(newVersionCmd(out))
	return cmd
}

func (o gitOptions) Run(out io.Writer) error {
	// Load named profile if one has been provided
	opts := []func(*config.LoadOptions) error{}
	if profile := identifyProfile(o.RemoteURL); profile != "" {
		opts = append(opts, config.WithSharedConfigProfile(profile))
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), opts...)
	if err != nil {
		return err
	}

	creds, err := cfg.Credentials.Retrieve(context.TODO())
	if err != nil {
		return err
	}

	// Translate and then sign the RemoteURL
	grcURL, err := translate.FromGrc(o.RemoteURL)
	if err != nil {
		return err
	}

	signer := awsv4.NewSigner(creds)
	url, err := signer.Sign(grcURL)
	if err != nil {
		return err
	}

	cmd := exec.Command("git", "remote-http", o.Cmd, url)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func identifyProfile(url string) string {
	m := grcRgx.FindStringSubmatch(url)
	if len(m) < 3 {
		return ""
	}

	// A GRC url will contain an optional profile with a trailing @
	profile := m[len(m)-1]
	if strings.Contains(profile, "@") {
		return strings.Split(profile, "@")[0]
	}

	return ""
}
