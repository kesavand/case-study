package eventshandler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/kesavand/case-study/internal/pkg/datahandler"
	"github.com/kesavand/case-study/internal/pkg/kafka"
	"github.com/kesavand/case-study/internal/pkg/utils"
)

const (
	eventDuration     = 5
	eventKafkaAddress = "KAFKA_ENDPOINT"
	kafkaTopic        = "test-topic"
)

type evntHandler struct {
}

type evntHandlerInterface interface {
	//ProduceEvnt(ctx context.Context, arams *utils.Parameter) error
	//	ConsumeEvnt(ctx context.Context, params *utils.Parameter) error
	Start(ctx context.Context, params *utils.Parameter, exitChannel chan error)
	StartConsumer(ctx context.Context, params *utils.Parameter, exitChannel chan error, dh datahandler.DataHandlerInterface)
	Stop(ctx context.Context) error
}

func NewEvntHandler() evntHandlerInterface {
	return &evntHandler{}
}

func (evh *evntHandler) ProduceEvnt(ctx context.Context, params *utils.Parameter) error {
	//Read kafka endpoint
	kafkaBroker, _ := params.ReadString(eventKafkaAddress, "kafka:9092")
	fmt.Printf("The kafka broker is %s", kafkaBroker)

	kp, err := kafka.NewKafkaProducer(kafkaBroker)
	if err != nil {
		return err
	}

	produceTick := time.NewTicker(eventDuration * time.Second)
	defer produceTick.Stop()
FORLOOP:
	for {
		select {
		case <-ctx.Done():
			fmt.Println("stopping producer")
			break FORLOOP
		case tick := <-produceTick.C:
			user := utils.GetUserInfo(tick)
			if data, err := json.Marshal(user); err != nil {
				fmt.Printf("marshal failed")
			} else {
				fmt.Println("producing event", data)
				kp.Produce(ctx, kafkaTopic, string(data))

			}

		}

	}
	return nil
}

func (evh *evntHandler) ConsumeEvnt(ctx context.Context, params *utils.Parameter, dh datahandler.DataHandlerInterface) error {
	kafkaBroker, _ := params.ReadString(eventKafkaAddress, "kafka:9092")
	fmt.Printf("The conusmer kafka broker is %s", kafkaBroker)

	kc, err := kafka.NewKafkaConsumer(kafkaBroker, dh, HandleMessage)
	if err != nil {
		return err
	}

	errChn := make(chan error, 100)
	kc.Read(ctx, kafkaTopic, errChn)

	for {
		fmt.Println("stopping consumer loop")
		select {
		case err := <-errChn:
			fmt.Printf("conusmer error %s", err)
			return err
		case <-ctx.Done():
			fmt.Printf("consumer closed")
			return nil
		}
	}

}

func (evh *evntHandler) Start(ctx context.Context, params *utils.Parameter, exitChannel chan error) {
	fmt.Printf("entering Start\n")
	go func() {
		<-ctx.Done()
		if err := evh.Stop(ctx); err != nil {
			fmt.Printf("Failed to stop event handler\n")
		}
	}()
	fmt.Printf("starting producer\n")
	exitChannel <- evh.ProduceEvnt(ctx, params)

}

func (evh *evntHandler) StartConsumer(ctx context.Context, params *utils.Parameter, exitChannel chan error, dh datahandler.DataHandlerInterface) {

	fmt.Printf("entering Start consumer \n")
	go func() {
		<-ctx.Done()
		if err := evh.Stop(ctx); err != nil {
			fmt.Printf("Failed to stop event handler\n")
		}
	}()

	fmt.Printf("starting producer\n")
	evh.ConsumeEvnt(ctx, params, dh)
	exitChannel <- evh.ConsumeEvnt(ctx, params, dh)

}

func (evh *evntHandler) Stop(ctx context.Context) error {
	return nil
}

func HandleMessage(msg string, dh datahandler.DataHandlerInterface) error {

	userInfo := &utils.UserInfo{}
	err := json.Unmarshal([]byte(msg), userInfo)
	if err != nil {
		fmt.Printf("Failed to unmarshal user inf")
	}

	fmt.Printf("Successfully recvd msg %v", msg)

	fmt.Printf("Trying to write Data %+v\n", userInfo)

	dh.Write(userInfo)

	return nil
}
