package data

type Parametro struct {
	Id                int         `json:"Id,omitempty"`
	Nombre            string      `json:"Nombre,omitempty"`
	Descripcion       string      `json:"Descripcion,omitempty"`
	CodigoAbreviacion string      `json:"CodigoAbreviacion,omitempty"`
	TipoParametroId   interface{} `json:"TipoParametroId,omitempty"`
	ParametroPadreId  interface{} `json:"ParametroPadreId,omitempty"`
	Activo            bool        `json:"Activo,omitempty"`
	FechaCreacion     string      `json:"FechaCreacion,omitempty"`
	FechaModificacion string      `json:"FechaModificacion,omitempty"`
}

type Periodo struct {
	Id                int    `json:"Id,omitempty"`
	Nombre            string `json:"Nombre,omitempty"`
	Descripcion       string `json:"Descripcion,omitempty"`
	Year              int    `json:"Year,omitempty"`
	Ciclo             string `json:"Ciclo,omitempty"`
	CodigoAbreviacion string `json:"CodigoAbreviacion,omitempty"`
	AplicacionId      int    `json:"AplicacionId,omitempty"`
	InicioVigencia    string `json:"InicioVigencia,omitempty"`
	FinVigencia       string `json:"FinVigencia,omitempty"`
	Activo            bool   `json:"Activo,omitempty"`
	FechaCreacion     string `json:"FechaCreacion,omitempty"`
	FechaModificacion string `json:"FechaModificacion,omitempty"`
}

type ParametroPeriodo struct {
	Id                int       `json:"Id,omitempty"`
	ParametroId       Parametro `json:"ParametroId,omitempty"`
	PeriodoId         Periodo   `json:"PeriodoId,omitempty"`
	Valor             string    `json:"Valor,omitempty"`
	Activo            bool      `json:"Activo,omitempty"`
	FechaCreacion     string    `json:"FechaCreacion,omitempty"`
	FechaModificacion string    `json:"FechaModificacion,omitempty"`
}

// ? Hay más en modelo pero no se han usado en mid
