package data

type Tercero struct {
	Id                  int         `json:"Id,omitempty"`
	NombreCompleto      string      `json:"NombreCompleto,omitempty"`
	PrimerNombre        string      `json:"PrimerNombre,omitempty"`
	SegundoNombre       string      `json:"SegundoNombre,omitempty"`
	PrimerApellido      string      `json:"PrimerApellido,omitempty"`
	SegundoApellido     string      `json:"SegundoApellido,omitempty"`
	LugarOrigen         int         `json:"LugarOrigen,omitempty"`
	FechaNacimiento     string      `json:"FechaNacimiento,omitempty"`
	TipoContribuyenteId interface{} `json:"TipoContribuyenteId,omitempty"`
	UsuarioWSO2         string      `json:"UsuarioWSO2,omitempty"`
	Activo              bool        `json:"Activo,omitempty"`
	FechaCreacion       string      `json:"FechaCreacion,omitempty"`
	FechaModificacion   string      `json:"FechaModificacion,omitempty"`
}

type DatosIdentificacion struct {
	Id                 int           `json:"Id,omitempty"`
	TipoDocumentoId    TipoDocumento `json:"TipoDocumentoId,omitempty"`
	TerceroId          Tercero       `json:"TerceroId,omitempty"`
	Numero             string        `json:"Numero,omitempty"`
	DigitoVerificacion int           `json:"DigitoVerificacion,omitempty"`
	CiudadExpedicion   int           `json:"CiudadExpedicion,omitempty"`
	FechaExpedicion    string        `json:"FechaExpedicion,omitempty"`
	DocumentoSoporte   int           `json:"DocumentoSoporte,omitempty"`
	Activo             bool          `json:"Activo,omitempty"`
	FechaCreacion      string        `json:"FechaCreacion,omitempty"`
	FechaModificacion  string        `json:"FechaModificacion,omitempty"`
}

type TipoDocumento struct {
	Id                int    `json:"Id,omitempty"`
	Nombre            string `json:"Nombre,omitempty"`
	Descripcion       string `json:"Descripcion,omitempty"`
	CodigoAbreviacion string `json:"CodigoAbreviacion,omitempty"`
	NumeroOrden       int    `json:"NumeroOrden,omitempty"`
	Activo            bool   `json:"Activo,omitempty"`
	FechaCreacion     string `json:"FechaCreacion,omitempty"`
	FechaModificacion string `json:"FechaModificacion,omitempty"`
}

type InfoComplementariaTercero struct {
	Id                       int         `json:"Id,omitempty"`
	TerceroId                Tercero     `json:"TerceroId,omitempty"`
	InfoComplementariaId     interface{} `json:"InfoComplementariaId,omitempty"`
	Dato                     string      `json:"Dato,omitempty"`
	InfoCompleTerceroPadreId interface{} `json:"InfoCompleteroPadreId,omitempty"`
	Activo                   bool        `json:"Activo,omitempty"`
	FechaCreacion            string      `json:"FechaCreacion,omitempty"`
	FechaModificacion        string      `json:"FechaModificacion,omitempty"`
}

// TODO: Agregar lo que falta ...

type Vinculacion struct {
	Id int `json:"Id,omitempty"`

	Activo            bool   `json:"Activo,omitempty"`
	FechaCreacion     string `json:"FechaCreacion,omitempty"`
	FechaModificacion string `json:"FechaModificacion,omitempty"`
}
