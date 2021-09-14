package apihandler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kesavand/case-study/internal/pkg/datahandler"
	"github.com/kesavand/case-study/internal/pkg/utils"
)

type apihandler struct {
	dh datahandler.DataHandlerInterface
}

func StartApiHandler(ctx context.Context, params *utils.Parameter, exitChannel chan error) {
	fmt.Println("starting api handler")

	aph := &apihandler{
		dh: datahandler.NewDatahandler(params),
	}

	router := mux.NewRouter()
	router.Path("/users").HandlerFunc(aph.getUsers)
	http.ListenAndServe(":8000", router)

	exitChannel <- nil
}

func (aph *apihandler) getUsers(w http.ResponseWriter, r *http.Request) {
	var users []utils.UserInfo
	name := r.URL.Query().Get("name")
	age, _ := strconv.Atoi(r.URL.Query().Get("age"))
	if name != "" && age != 0 {
		aph.dh.GetUsersByNameAndAge(w, &users, name, age)
	} else if name != "" {
		fmt.Println("Name is ", name)
		aph.dh.GetUsersByName(w, &users, name)
	} else if age != 0 {
		fmt.Println("Age is ", age)
		aph.dh.GetUsersByAge(w, &users, age)
	} else {

		aph.dh.GetAllusers(w, &users)
	}

}
