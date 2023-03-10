package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"wgs/internal/config"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var CmdServer = &cobra.Command{
	Use:   "server",
	Short: "Run the wgs server",
	Run: func(*cobra.Command, []string) {
		//////////////////////////////////////////////////////////

		app := InitApp()
		log.Info().Msg("starting http server")
		go app.Listen(config.Get().Listen) //nolint:errcheck

		ctx, cancel := context.WithCancel(context.Background())
		// Start the app
		go func() {
			<-ctx.Done()
			if err := app.Shutdown(); err != nil {
				log.Fatal().Err(err).Msg("failed to shutdown")
			}
		}()

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		log.Info().Msg("initialize complete")
		<-sigs // Blocks here until interrupted
		// Handle shutdown
		fmt.Println("Shutdown signal received")
		cancel()
		fmt.Println("All workers done, shutting down!")
	},
}
