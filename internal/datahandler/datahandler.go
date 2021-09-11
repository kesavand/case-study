package datahandler

import (
	"encoding/json"
	"fmt"

	"github.com/kesavand/case-study/internal/pkg/utils"
)

func HandleMessage(msg string) error {
	userInfo := utils.UserInfo{}
	err := json.Unmarshal([]byte(msg), &userInfo)
	if err != nil {
		fmt.Printf("Failed to unmarshal user inf")
	}

	fmt.Printf("Successfully recvd msg %v", userInfo)
	return nil
}
