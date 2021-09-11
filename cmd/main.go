package main

import (
	"context"
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kesavand/case-study/internal/pkg/eventshandler"
	"github.com/kesavand/case-study/internal/pkg/utils"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	fmt.Printf("starting process\n")
	err = startProcess()
	if err != nil {
		panic(err)
	}
}

func startProcess() error {
	exitChannel := make(chan error, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	opts := utils.New()

	event := eventshandler.NewEvntHandler()
	go event.StartConsumer(ctx, opts, exitChannel)
	go event.Start(ctx, opts, exitChannel)

	return utils.Wait(exitChannel, cancel)

}
