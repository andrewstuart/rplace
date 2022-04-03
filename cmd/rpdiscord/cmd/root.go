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

		if chid, err := cmd.Flags().GetString("channel"); err == nil {
			ch, err := disCli.Channel(chid)
			if err != nil {
				log.Panic("Channel details error: ", err)
			}

			if testu, err := cmd.Flags().GetString("testuser"); err == nil {
				ms, err := disCli.GuildMembers(ch.GuildID, "", 1000)
				if err != nil {
					log.Println("Error getting members", err)
				}
				for _, m := range ms {
					if m.User.Username == testu {
						ch, err := disCli.UserChannelCreate(m.User.ID)

						disCli.ChannelFileSend(ch.ID, "example.png", buf)
						const h = 25
						for up := range ups {
							img := image.NewPaletted(image.Rect(0, 0, h, h), rplace.StdPalette)
							for i := 0; i < h; i++ {
								for j := 0; j < h; j++ {
									img.Set(i, j, up.Color.Color)
								}
							}
							buf := &bytes.Buffer{}
							png.Encode(buf, img)
							disCli.ChannelFileSendWithMessage(ch.ID, fmt.Sprintf("Please update %s to %s. (See attached 25x25 color swatch)", up.Link(), up.Color.Name), "color.png", buf)
							if err != nil {
								log.Println(err)
							}
							time.Sleep(5 * time.Minute)
						}
						return
					}
				}
				ch, err := disCli.Channel(chid)
				if err != nil {
					log.Panic(err)
				}
				go func() {
					for {
						select {
						case <-cmd.Context().Done():
							return
						case up := <-ups:
							if chid != "" {
								time.Sleep(5 * time.Second)
								disCli.ChannelMessageSend(ch.ID, up.Link())
							}
						}
					}
				}()
			}

			disCli.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
				spew.Dump(m)
				log.Println(m.Author.Username, ": ", m.Content)
			})

			// return

			err = disCli.Open()
			if err != nil {
				log.Panic(err)
			}
		}
		select {
		case <-cmd.Context().Done():
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
	rootCmd.Flags().StringP("testuser", "u", "", "The user ID")
	rootCmd.Flags().StringP("image", "i", "gopher.png", "An image, in png format")
	rootCmd.Flags().StringP("token", "t", "", "The discord bot token")
	rootCmd.Flags().StringP("channel", "c", "", "The discord bot channel")
	rootCmd.Flags().Int("x", 0, "The X coordinate")
	rootCmd.Flags().Int("y", 0, "The Y coordinate")
}
