package authz

import (
	"net/http"
	"net/url"

	"github.com/auth4flow/auth4flow-core/pkg/service"
	"github.com/gorilla/mux"
<<<<<<< HEAD
	wookie "github.com/warrant-dev/warrant/pkg/authz/wookie"
	"github.com/warrant-dev/warrant/pkg/service"
=======
>>>>>>> main
)

func (svc ObjectService) Routes() ([]service.Route, error) {
	return []service.Route{
		// create
		service.WarrantRoute{
			Pattern: "/v1/objects",
			Method:  "POST",
			Handler: service.NewRouteHandler(svc, CreateHandler),
		},

		// get
		service.WarrantRoute{
			Pattern: "/v1/objects",
			Method:  "GET",
			Handler: service.ChainMiddleware(
				service.NewRouteHandler(svc, ListHandler),
				service.ListMiddleware[ObjectListParamParser],
			),
		},
		service.WarrantRoute{
			Pattern: "/v1/objects/{objectType}/{objectId}",
			Method:  "GET",
			Handler: service.NewRouteHandler(svc, GetHandler),
		},

		// delete
		service.WarrantRoute{
			Pattern: "/v1/objects/{objectType}/{objectId}",
			Method:  "DELETE",
			Handler: service.NewRouteHandler(svc, DeleteHandler),
		},
	}, nil
}

func CreateHandler(svc ObjectService, w http.ResponseWriter, r *http.Request) error {
	var newObject ObjectSpec
	err := service.ParseJSONBody(r.Body, &newObject)
	if err != nil {
		return err
	}

	createdObject, err := svc.Create(r.Context(), newObject)
	if err != nil {
		return err
	}

	service.SendJSONResponse(w, createdObject)
	return nil
}

func ListHandler(svc ObjectService, w http.ResponseWriter, r *http.Request) error {
	listParams := service.GetListParamsFromContext[ObjectListParamParser](r.Context())
	queryParams := r.URL.Query()
	objectType, err := url.QueryUnescape(queryParams.Get("objectType"))
	if err != nil {
		return service.NewInvalidParameterError("objectType", "")
	}

	filterOptions := FilterOptions{ObjectType: objectType}
	objects, err := svc.List(r.Context(), &filterOptions, listParams)
	if err != nil {
		return err
	}

	service.SendJSONResponse(w, objects)
	return nil
}

func GetHandler(svc ObjectService, w http.ResponseWriter, r *http.Request) error {
	objectType := mux.Vars(r)["objectType"]
	objectId := mux.Vars(r)["objectId"]
	object, err := svc.GetByObjectTypeAndId(r.Context(), objectType, objectId)
	if err != nil {
		return err
	}

	service.SendJSONResponse(w, object)
	return nil
}

func DeleteHandler(svc ObjectService, w http.ResponseWriter, r *http.Request) error {
	objectType := mux.Vars(r)["objectType"]
	objectId := mux.Vars(r)["objectId"]
	newWookie, err := svc.DeleteByObjectTypeAndId(r.Context(), objectType, objectId)
	if err != nil {
		return err
	}
	wookie.AddAsResponseHeader(w, newWookie)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	return nil
}
