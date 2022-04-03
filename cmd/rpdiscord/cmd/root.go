/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"time"

	"github.com/andrewstuart/rplace"
	"github.com/bwmarrin/discordgo"
	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rpdiscord",
	Short: "Discord bot for r/place",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		// Target image
		imgFlg, _ := cmd.Flags().GetString("image")
		imgF, err := os.OpenFile(imgFlg, os.O_RDONLY, 0400)
		if err != nil {
			log.Panic("error getting image file: ", err)
		}
		img, err := png.Decode(imgF)
		if err != nil {
			log.Panic("error decoding png image file: ", err)
		}
		imgF.Close()

		cli := rplace.Client{}
		x, _ := cmd.Flags().GetInt("x")
		y, _ := cmd.Flags().GetInt("y")

		ups, err := cli.NeededUpdatesFor(cmd.Context(), img, image.Point{X: x, Y: y})
		if err != nil {
			log.Panic("error getting updates for image: ", err)
		}

		// Encode and send example image
		example, err := cli.WithImage(img, image.Point{X: x, Y: y})
		if err != nil {
			log.Panic(err)
		}

		buf := &bytes.Buffer{}
		png.Encode(buf, example)

		tok, _ := cmd.Flags().GetString("token")
		disCli, err := discordgo.New("Bot " + tok)
		if err != nil {
			log.Panic("error connecting to discord: ", err)
		}

		disCli.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentDirectMessages | discordgo.IntentGuildMembers
		disCli.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
			spew.Dump(m)
			log.Println(m.Author.Username, ": ", m.Content)
		})

		err = disCli.Open()
		if err != nil {
			log.Panic(err)
		}

		guildID, _ := cmd.Flags().GetString("guild")
		for i := 0; ; i++ {
			if i > 0 {
				select {
				case <-cmd.Context().Done():
				case <-time.After(5 * time.Minute):
				}
			}

			ms, err := disCli.GuildMembers(guildID, "", 1000)
			if err != nil {
				log.Println("Error getting members", err)
				continue
			}

			for _, m := range ms {
				ch, err := disCli.UserChannelCreate(m.User.ID)
				if err != nil {
					log.Println(err)
					continue
				}
				up := <-ups

				const h = 25
				img := image.NewPaletted(image.Rect(0, 0, h, h), rplace.StdPalette)
				for i := 0; i < h; i++ {
					for j := 0; j < h; j++ {
						img.Set(i, j, up.Color.Color)
					}
				}
				buf := &bytes.Buffer{}
				png.Encode(buf, img)

				disCli.ChannelFileSendWithMessage(ch.ID, fmt.Sprintf("Please update %s to %s. (See color swatch)", up.Link(), up.Color.Name), "color.png", buf)
				if err != nil {
					log.Println(err)
				}
			}
		}
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

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.rpdiscord.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().StringP("image", "i", "gopher.png", "An image, in png format")
	rootCmd.Flags().StringP("token", "t", "", "The discord bot token")
	rootCmd.Flags().StringP("guild", "g", "", "The discord guild (server)")
	rootCmd.Flags().Int("x", 0, "The X coordinate")
	rootCmd.Flags().Int("y", 0, "The Y coordinate")
}
