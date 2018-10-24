package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/kinvolk/karydia/pkg/server"
	"github.com/kinvolk/karydia/pkg/util"
)

var runserverCmd = &cobra.Command{
	Use:   "runserver",
	Short: "Run the karydia server",
	Run:   runserverFunc,
}

func init() {
	rootCmd.AddCommand(runserverCmd)

	runserverCmd.Flags().String("addr", "0.0.0.0:33333", "Address to listen on")
	runserverCmd.Flags().String("tls-cert", "cert.pem", "Path to TLS certificate file")
	runserverCmd.Flags().String("tls-key", "key.pem", "Path to TLS private key file")
}

func runserverFunc(cmd *cobra.Command, args []string) {
	tlsConfig, err := util.CreateTLSConfig(viper.GetString("tls-cert"), viper.GetString("tls-key"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create TLS config: %v\n", err)
		os.Exit(1)
	}
	s, err := server.New(&server.Config{
		Addr:      viper.GetString("addr"),
		TLSConfig: tlsConfig,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load server: %v\n", err)
		os.Exit(1)
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, os.Kill)

		<-sigChan

		ctx, cancelCtx := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelCtx()

		if err := s.Shutdown(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "HTTP Shutdown error: %v\n", err)
		}

		close(idleConnsClosed)
	}()

	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		fmt.Fprintf(os.Stderr, "HTTP ListenAndServe error: %v", err)
	}

	<-idleConnsClosed
}
