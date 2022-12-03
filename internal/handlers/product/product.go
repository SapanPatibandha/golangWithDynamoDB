package product

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"

	"github.com/SapanPatibandha/golangWithDynamoDB/internal/repository/adapter"
	"github.com/SapanPatibandha/golangWithDynamoDB/internal/controllers/product"
	EntityProduct "github.com/SapanPatibandha/golangWithDynamoDB/internal/entities/product"
	"github.com/SapanPatibandha/golangWithDynamoDB/internal/handlers"
	Rules "github.com/SapanPatibandha/golangWithDynamoDB/internal/rules"
	RulesProduct "github.com/SapanPatibandha/golangWithDynamoDB/internal/rules/product"
	HttpStatus "github.com/SapanPatibandha/golangWithDynamoDB/utils/http"
)

type Handler struct {
	handler.Interface
	Controller product.Interface
	Rules      Rules.Interface
}

func NewHandler(repository adaptor.Interface) handlers.Interface {
	return &Handler{
		Controller: product.NewController(repository),
		Rules:      RulesProduct.NewRules(),
	}
}

func (h *Handler)Get(w http.ResponseWriter, r http.Request) {
	if chi.URLParm(r, "ID") != ""{
		h.getOne()(w, r)
	}else {
		h.getAll(w, r)
	}
}

func (h *Handler) getOne(w http.ResponseWriter, r http.Request) {
	ID, err := uuid.Parse(chi.URLParm(r, "ID"))

	if err!= nil{
		HttpStatus.StatusBadRequest(w, r, errors.New("ID is not uuid values"))
		return
	}

	response, err := h.Controller.ListOne(ID)
	if err != nil{
		HttpStatus.StatusInternalServerError(w, r, err)
		return
	}

	HttpStatus.StatusOK(w, r, response)
}

func (h *Handler) getAll(w http.ResponseWriter, r http.Request) {
	
	response, err := h.Controller.ListAll()

	if err != nil{
		HttpStatus.StatusInternalServerError(w, r, err)
		return
	}

	HttpStatus.StatusOK(w, r, response)
}

func (h *Handler) Post(w http.ResponseWriter, r *http.Request) {

	productBody, err := h.getBodyAndValidate(r, uuid.Nil)

	if err != nil {
		HttpStatus.StatusBadRequest(w, r, err)
		return
	}

	ID, err := h.Controller.Create(productBody)

	if err!=nil{
		HttpStatus.StatusInternalServerError(w, r, err)
		return
	}

	HttpStatus.StatusOK(w, r, map[string]interface{"id:", ID.string()})
}

func (h *Handler) Put(w http.ResponseWriter, r *http.Request) {
	ID, err := uuid.Parse(chi.URLParm(r, "ID"))

	if err!= nil{
		HttpStatus.StatusBadRequest(w, r, errors.New("ID is not uuid values"))
		return
	}

	productBody, err := h.getBodyAndValidate(r, ID)

	if err != nil{
		HttpStatus.StatusBadRequest(w, r, err)
		return
	}

	_, err := h.Controller.Update(ID, productBody)
	if err != nil{
		HttpStatus.StatusBadRequest(w, r, err)
		return
	}

	HttpStatus.StatusNoContent(w, r)
}

func (h *Handler)Delete(w http.ResponseWriter, r *http.Response) {
	ID, err := uuid.Parse(chi.URLParm(r, "ID"))

	if err!= nil{
		HttpStatus.StatusBadRequest(w, r, errors.New("ID is not uuid values"))
		return
	}

	_, err := h.Controller.Remove(ID)
	if err != nil{
		HttpStatus.StatusInternalServerError(w, r, err)
		return
	}

	HttpStatus.StatusNoContent(w, r)
}

func (h *Handler)Options(w http.ResponseWriter, r *http.Request) {
	
	HttpStatus.StatusNoContent(w, r)

}

func (h *Handler) getBodyAndValidate(r *http.Request, ID uuid.UUID)(*EntityProduct.Product, error) {
	productBody := &EntityProduct.Product{}
	body, err := h.Rules.ConvertIoReaderToStruct(r.Body, productBody)

	if err != nil{
		return &EntityProduct.Product{}, errors.New("body is required")
	}

	productParsed, err := EntityProduct.InterfaceToModel(body)
	if err != nil{
		return &EntityProduct.Product{}, errors.New("error on converting body to modal")
	}

	setDefaultValues(productParsed, ID)
	return productParsed, h.Role.Validate(productParsed)
}

func setDefaultValues(product *EntityProduct.Product, ID uuid.UUID)  {
	
	product.UpdateAt = time.Now()
	if ID == uuid.Nil{
		product.ID = uuid.New()
		product.CreatedAt = time.Now()
	}else{
		product.ID = ID
	}

}