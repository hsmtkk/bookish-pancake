package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/hsmtkk/bookish-pancake/openweather"
	"github.com/hsmtkk/bookish-pancake/util"
	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use:  "openweather",
		Args: cobra.ExactArgs(1),
		Run:  run,
	}
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func run(cmd *cobra.Command, args []string) {
	apiKey, err := util.RequiredEnvVar("API_KEY")
	if err != nil {
		log.Fatal(err)
	}
	city := args[0]
	result, err := openweather.New(http.DefaultClient).GetCurrentWeather(context.Background(), apiKey, city)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("result: %v\n", result)
}
