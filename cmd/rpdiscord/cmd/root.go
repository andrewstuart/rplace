/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"os"

	"github.com/andrewstuart/rplace"
	"github.com/bwmarrin/discordgo"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rpdiscord",
	Short: "Discord bot for r/place",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
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

		disCli.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentDirectMessages | discordgo.IntentGuildMembers

		chid, _ := cmd.Flags().GetString("channel")

		disCli.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
			fmt.Printf("m = %#v\n", m)
			// chid = m.ChannelID
			fmt.Printf("m.ChannelID = %+v\n", m.ChannelID)
			fmt.Printf("m.Content = %+v\n", m.Content)
			if m.Thread != nil {
				fmt.Printf("m.Thread.ID = %+v\n", m.Thread.ID)
			}
		})

		ch, err := disCli.Channel(chid)
		if err != nil {
			log.Panic(err)
		}

		ms, err := disCli.GuildMembers(ch.GuildID, "", 100)
		if err != nil {
			log.Panic(err)
		}
		fmt.Printf("ms = %+v\n", ms)
		for _, m := range ms {
			fmt.Printf("m.Nick = %+v\n", m.Nick)
			fmt.Printf("m.User.Username = %+v\n", m.User.Username)
			fmt.Printf("m.User.IED = %+v\n", m.User.ID)
		}
		// return

		err = disCli.Open()
		if err != nil {
			log.Panic(err)
		}

		fmt.Printf("ups = %+v\n", ups)
		// go func() {
		// 	for {
		// 		select {
		// 		case <-cmd.Context().Done():
		// 			return
		// 		case up := <-ups:
		// 			if chid != "" {
		// 				disCli.ChannelMessageSend(chid, up.Link())
		// 			}
		// 		}
		// 	}
		// }()

		// fmt.Printf("ups = %+v\n", ups)
		// fmt.Printf("disCli = %+v\n", disCli)

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
	rootCmd.Flags().StringP("image", "i", "gopher.png", "An image, in png format")
	rootCmd.Flags().StringP("token", "t", "", "The discord bot token")
	rootCmd.Flags().StringP("channel", "c", "", "The discord bot channel")
	rootCmd.Flags().Int("x", 0, "The X coordinate")
	rootCmd.Flags().Int("y", 0, "The Y coordinate")
}
