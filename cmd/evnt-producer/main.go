package main

import (
	"context"

	"github.com/joho/godotenv"
	"github.com/kesavand/case-study/internal/pkg/eventshandler"
	"github.com/kesavand/case-study/internal/pkg/logger"
	"github.com/kesavand/case-study/internal/pkg/utils"
)

func main() {
	logger := logger.NewLogger()
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	logger.Info.Println("Starting process")

	err = startProcess()
	if err != nil {
		panic(err)
	}
}

func startProcess() error {
	exitChannel := make(chan error, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	utils.GenerateTestUsers()

	opts := utils.New()

	event := eventshandler.NewEvntHandler()
	go event.Start(ctx, opts, exitChannel)

	return utils.Wait(exitChannel, cancel)

}
