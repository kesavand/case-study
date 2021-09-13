package utils

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type UserInfo struct {
	gorm.Model    `json:"-"`
	Name          string `json:"name"`
	Age           int    `json:"age"`
	EmployeeID    string `json:"employeeid"`
	CardNo        string `json:"cardnumber"`
	CardSwipeTime string `json:"cardswipetime"`
}

var userIdx = 0
var maxUserCount = 10
var testUsers map[int]*UserInfo

//WaitForExit waits on an exit channel until error or an interrupt is received
func Wait(exitChannel chan error, cancel context.CancelFunc) error {
	err := <-exitChannel
	// cancel the context
	cancel()
	return err
}

//Parameter construction for holding configuration options
type Parameter struct {
	parameters *sync.Map
}

//New creates empty bag
func New() *Parameter {
	param := &Parameter{
		parameters: &sync.Map{},
	}

	return param

}

//Read reads from env variable s, returns default if not present in env
func (p *Parameter) Read(name string, defaultVal interface{}) (interface{}, bool) {
	valI, okI := p.parameters.Load(name)
	if okI {
		return valI, true
	}
	envVal, found := os.LookupEnv(name)
	if found {
		return envVal, true
	}
	return defaultVal, false
}

func (p *Parameter) ReadString(name string, defaultVal string) (string, bool) {
	val, found := p.Read(name, defaultVal)

	if !found {
		return defaultVal, false
	}

	return val.(string), true
}

func getRandomNo(min int, max int) int {
	rand.Seed(time.Now().UnixNano())

	return (rand.Intn(max-min+1) + min)
}

func GenerateTestUsers() {
	testUsers = make(map[int]*UserInfo)

	for i := 0; i <= maxUserCount; i++ {
		testUser := new(UserInfo)
		testUser.Name = fmt.Sprintf("USER-" + strconv.Itoa(i))
		testUser.Age = getRandomNo(20, 100)
		testUser.EmployeeID = strconv.Itoa(i)
		id := uuid.New()
		testUser.CardNo = id.String()
		testUsers[i] = testUser
		fmt.Printf("User is %v", testUser)

	}
}

func GetUserInfo(t time.Time) *UserInfo {
	fmt.Printf("K7>>>>the index is %d  %s\n", userIdx, t.String())
	if userIdx > maxUserCount {
		userIdx = 0
	}
	idx := userIdx
	testUsers[userIdx].CardSwipeTime = t.String()
	userIdx++
	return testUsers[idx]

}
