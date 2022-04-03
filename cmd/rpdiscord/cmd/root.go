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

		tok, _ := cmd.Flags().GetString("token")
		disCli, err := discordgo.New("Bot " + tok)
		if err != nil {
			log.Panic("error connecting to discord: ", err)
		}

		disCli.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentDirectMessages | discordgo.IntentGuildMembers | discordgo.IntentGuildMessageReactions | discordgo.IntentDirectMessageReactions
		disCli.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
			spew.Dump(m)
			log.Println(m.Author.Username, ": ", m.Content)
		})

		acked := map[string]rplace.Update{}
		disCli.AddHandler(func(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
			spew.Dump(m)
			if m.MessageReaction.Emoji.Name == "✅" {
				delete(acked, m.MessageID)
			}
		})
		disCli.AddHandler(func(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
			spew.Dump(m)
			if m.MessageReaction.Emoji.ID == "" {
				delete(acked, m.MessageID)
			}
		})

		err = disCli.Open()
		if err != nil {
			log.Panic(err)
		}

		if chid, err := cmd.Flags().GetString("channel"); err == nil {
			// Encode and send example image
			example, err := cli.WithImage(img, image.Point{X: x, Y: y})
			if err != nil {
				log.Panic(err)
			}

			buf := &bytes.Buffer{}
			png.Encode(buf, example)

			disCli.ChannelFileSendWithMessage(chid, "Preview of the intended canvas state", "example.png", buf)
			buf.Reset()
		}

		guildID, _ := cmd.Flags().GetString("guild")
		for i := 0; ; i++ {
			if i > 0 {
				select {
				case <-cmd.Context().Done():
				case <-time.After(6 * time.Minute):
				}
			}

			ms, err := disCli.GuildMembers(guildID, "", 1000)
			if err != nil {
				log.Println("Error getting members", err)
				continue
			}

			for _, m := range ms {
				if m.User.ID == disCli.State.User.ID {
					continue
				}
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

				msg, err := disCli.ChannelFileSendWithMessage(ch.ID, fmt.Sprintf("Please update %s to %s. (See color swatch). Please react ✅ when done or the request will be requeued.", up.Link(), up.Color.Name), "color.png", buf)
				if err != nil {
					log.Println(err)
				}

				acked[msg.ID] = up
				go func() {
					defer func() {
						err := recover()
						if err != nil {
							log.Printf("Panicked trying to update: %v\n", err)
						}
					}()
					time.Sleep(6 * time.Minute)
					if up, ok := acked[msg.ID]; ok {
						log.Printf("requeueing: %+v\n", up)
						ups <- up
					}
				}()

				fmt.Printf("msg.ID = %+v\n", msg.ID)
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
	rootCmd.Flags().String("channel", "", "The channel to post a preview in.")
	rootCmd.Flags().Int("x", 0, "The X coordinate")
	rootCmd.Flags().Int("y", 0, "The Y coordinate")
}
