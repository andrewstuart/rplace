/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"git.stuart.fun/andrew/rester/v2"
	"git.stuart.fun/andrew/rester/v2/iou"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rpbot",
	Short: "A brief description of your application",
	Run: func(cmd *cobra.Command, args []string) {
		oc := oauth2.Config{
			ClientID:     "M3ecb_lk04u8L-S8Pb9-wg",
			ClientSecret: "mJfHhBMBJPXwTAzMWqABOMQh5zAx9w",
			Endpoint: oauth2.Endpoint{
				AuthURL:   "https://www.reddit.com/api/v1/authorize",
				TokenURL:  "https://www.reddit.com/api/v1/access_token",
				AuthStyle: oauth2.AuthStyleInParams,
			},
			RedirectURL: "http://localhost:8080/auth/oauth",
			Scopes:      []string{"read", "history"},
		}

		r := gin.Default()

		r.GET("/", func(c *gin.Context) {
			c.Redirect(http.StatusTemporaryRedirect, oc.AuthCodeURL("bar", oauth2.SetAuthURLParam("duration", "permanent")))
		})

		cli := &http.Client{
			Transport: rester.All{
				rester.MaxStatus(399),
				rester.Logging(logrus.WithField("cli", "reddit")),
				rester.MergeHeaders{"User-Agent": []string{"web:testcli:0.1.0 (by u/10gistic)"}},
				rester.RequestFunc(func(req *http.Request) {
					req.Body = iou.ReadCloser{
						Reader: io.TeeReader(req.Body, os.Stdout),
						Closer: req.Body,
					}
				}),
			}.Wrap(http.DefaultTransport),
		}

		r.GET("/auth/oauth", func(c *gin.Context) {
			cliCtx := context.WithValue(c.Request.Context(), oauth2.HTTPClient, cli)
			code := c.Query("code")
			log.Println("code ", code)
			tok, err := oc.Exchange(cliCtx, code)
			if err != nil {
				c.Writer.Header().Set("Content-Type", "text/html")
				c.Writer.WriteString(err.Error())
				return
			}
			log.Println("here")

			fmt.Printf("tok = %+v\n", tok)

			c.Writer.WriteString(tok.AccessToken)
		})

		r.Run()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.rpbot.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
