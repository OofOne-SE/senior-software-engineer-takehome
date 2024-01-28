package weather

import (
	"net/http"
	"senior-software-engineer-takehome/models"
	"senior-software-engineer-takehome/restapi/operations/weather_assets"

	"github.com/ascarter/requestid"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
)

type database interface {
	AddUnit(u *models.WeatherUnit) error
	GetUnits(from, to *strfmt.Date) ([]*models.WeatherUnit, error)
}

type dispatcher interface {
	Notify(event *models.WeatherUnit)
}

type Weather struct {
	Db         database
	Dispatcher dispatcher
	Log        func(string, ...any)
}

func (w Weather) GetWeatherUnits(params weather_assets.GetDataWeatherParams) middleware.Responder {
	parentCtx := params.HTTPRequest.Context()
	rID, _ := requestid.FromContext(parentCtx)
	if units, err := w.Db.GetUnits(&params.From, &params.To); err != nil {
		w.Log("GetWeatherUnits: %s", err)
		return weather_assets.NewGetDataWeatherDefault(http.StatusInternalServerError).WithPayload(&models.Error{
			Code:      "data.weather.get.internal_server_error",
			Error:     "Internal server error",
			RequestID: rID,
		})
	} else {
		return weather_assets.NewGetDataWeatherOK().WithPayload(units)
	}
}

func (w Weather) AddWeatherUnit(params weather_assets.PostDataWeatherParams) middleware.Responder {
	parentCtx := params.HTTPRequest.Context()
	rID, _ := requestid.FromContext(parentCtx)
	if err := w.Db.AddUnit(params.Body); err != nil {
		w.Log("AddUnit: %s", err)
		return weather_assets.NewPostDataWeatherDefault(http.StatusInternalServerError).WithPayload(&models.Error{
			Code:      "data.weather.post.internal_server_error",
			Error:     "Internal server error",
			RequestID: rID,
		})
	} else {
		w.Dispatcher.Notify(params.Body)
		return weather_assets.NewPostDataWeatherOK()
	}
}
