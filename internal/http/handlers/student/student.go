package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/himanshukandari14/go-server/internal/storage"
	"github.com/himanshukandari14/go-server/internal/types"
	response "github.com/himanshukandari14/go-server/internal/utils/responses"
)


func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("creating student")

		var student types.Student

		err:= json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err,io.EOF) {
			response.WriteJSON(w,http.StatusBadRequest,response.GeneralError(fmt.Errorf("empty body")))
			return
		}

		if err!=nil {
			response.WriteJSON(w,http.StatusBadRequest,response.GeneralError(err))
			return
		}

		//validate request
		if err:= validator.New().Struct(student);err!=nil{
			validateErrs :=err.(validator.ValidationErrors)
			response.WriteJSON(w,http.StatusBadRequest,response.ValidationError(validateErrs))
			return
		}

		lastId,err:=storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
		)
		slog.Info("user created succesfully",slog.String("userId",fmt.Sprint(lastId)))
		if err!=nil {
			response.WriteJSON(w,http.StatusInternalServerError, err)
			return
		}

		response.WriteJSON(w,http.StatusCreated,map[string]int64{"id":lastId})
	}
}