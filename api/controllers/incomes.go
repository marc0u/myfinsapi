package controllers

import (
	// "encoding/json"
	// "errors"
	// "fmt"
	// "io/ioutil"
	// "net/http"
	// "strconv"

	// "gitlab.com/marco.urriola/apifinances/api/models"
	// "gitlab.com/marco.urriola/apifinances/api/responses"
	// "gitlab.com/marco.urriola/apifinances/api/utils/formaterror"

	"github.com/gofiber/fiber"
	// "github.com/gorilla/mux"
)

func (server *Server) CreateIncome(c *fiber.Ctx) {
	c.JSON(fiber.Map{
		"ok":  true,
		"msg": "Create Income",
	})
	// // Reading body http request
	// body, err := ioutil.ReadAll(r.Body)
	// if err != nil {
	// 	responses.ERROR(w, http.StatusUnprocessableEntity, err)
	// 	return
	// }
	// item := models.Income{}
	// err = json.Unmarshal(body, &item)
	// if err != nil {
	// 	responses.ERROR(w, http.StatusUnprocessableEntity, err)
	// 	return
	// }
	// // Preparing and validating data
	// item.Prepare(uint32(uid))
	// err = item.Validate()
	// if err != nil {
	// 	responses.ERROR(w, http.StatusUnprocessableEntity, err)
	// 	return
	// }
	// // Verifying User ID
	// if uid != item.UserID {
	// 	responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
	// 	return
	// }
	// // Saving data
	// itemCreated, err := item.SaveIncome(server.DB)
	// if err != nil {
	// 	formattedError := formaterror.FormatError(err.Error())
	// 	responses.ERROR(w, http.StatusInternalServerError, formattedError)
	// 	return
	// }
	// // Http response
	// w.Header().Set("Lacation", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, itemCreated.ID))
	// responses.JSON(w, http.StatusCreated, itemCreated)
}

func (server *Server) GetIncomes(c *fiber.Ctx) {
	c.JSON(fiber.Map{
		"ok":  true,
		"msg": "Get Incomes",
	})
	// // Getting data
	// item := models.Income{}
	// items, err := item.FindAllIncomes(server.DB, uint32(uid))
	// if err != nil {
	// 	responses.ERROR(w, http.StatusInternalServerError, err)
	// 	return
	// }
	// // Http response
	// responses.JSON(w, http.StatusOK, items)
}

func (server *Server) GetIncome(c *fiber.Ctx) {
	c.JSON(fiber.Map{
		"ok":  true,
		"msg": "Get Income",
	})
	// // Getting URL parameter ID
	// vars := mux.Vars(r)
	// pid, err := strconv.ParseUint(vars["id"], 10, 64)
	// if err != nil {
	// 	responses.ERROR(w, http.StatusBadRequest, err)
	// 	return
	// }
	// // Getting data
	// item := models.Income{}
	// itemReceived, err := item.FindIncomeByID(server.DB, uint32(pid), uint32(uid))
	// if err != nil {
	// 	responses.ERROR(w, http.StatusInternalServerError, err)
	// 	return
	// }
	// // Http response
	// responses.JSON(w, http.StatusOK, itemReceived)
}

func (server *Server) UpdateIncome(c *fiber.Ctx) {
	c.JSON(fiber.Map{
		"ok":  true,
		"msg": "Update Income",
	})
	// // Getting URL parameter ID
	// vars := mux.Vars(r)
	// pid, err := strconv.ParseUint(vars["id"], 10, 64)
	// if err != nil {
	// 	responses.ERROR(w, http.StatusBadRequest, err)
	// 	return
	// }
	// // Getting data
	// item := models.Income{}
	// err = server.DB.Debug().Model(models.Income{}).Where("id = ?", pid).Take(&item).Error
	// if err != nil {
	// 	responses.ERROR(w, http.StatusNotFound, errors.New("Income not found"))
	// 	return
	// }
	// // Verifying User ID
	// if uid != item.UserID {
	// 	responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
	// 	return
	// }
	// // Reading body http request
	// body, err := ioutil.ReadAll(r.Body)
	// if err != nil {
	// 	responses.ERROR(w, http.StatusUnprocessableEntity, err)
	// 	return
	// }
	// itemUpdate := models.Income{}
	// err = json.Unmarshal(body, &itemUpdate)
	// if err != nil {
	// 	responses.ERROR(w, http.StatusUnprocessableEntity, err)
	// 	return
	// }
	// // Preparing and validating data
	// itemUpdate.PrepareUpdate(&item)
	// err = itemUpdate.Validate()
	// if err != nil {
	// 	responses.ERROR(w, http.StatusUnprocessableEntity, err)
	// 	return
	// }
	// // Updating data
	// itemUpdated, err := itemUpdate.UpdateAIncome(server.DB)
	// if err != nil {
	// 	formattedError := formaterror.FormatError(err.Error())
	// 	responses.ERROR(w, http.StatusInternalServerError, formattedError)
	// 	return
	// }
	// // Http responce
	// responses.JSON(w, http.StatusOK, itemUpdated)
}

func (server *Server) DeleteIncome(c *fiber.Ctx) {
	c.JSON(fiber.Map{
		"ok":  true,
		"msg": "Delete Income",
	})
	// // Getting URL parameter ID
	// vars := mux.Vars(r)
	// pid, err := strconv.ParseUint(vars["id"], 10, 64)
	// if err != nil {
	// 	responses.ERROR(w, http.StatusBadRequest, err)
	// 	return
	// }
	// // Getting data
	// item := models.Income{}
	// err = server.DB.Debug().Model(models.Income{}).Where("id = ?", pid).Take(&item).Error
	// if err != nil {
	// 	responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
	// 	return
	// }
	// // Verifying User ID
	// if uid != item.UserID {
	// 	responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
	// 	return
	// }
	// // Deleting item
	// _, err = item.DeleteAIncome(server.DB, uint32(pid))
	// if err != nil {
	// 	responses.ERROR(w, http.StatusBadRequest, err)
	// 	return
	// }
	// // Http response
	// w.Header().Set("Entity", fmt.Sprintf("%d", pid))
	// responses.JSON(w, http.StatusNoContent, "")
}
