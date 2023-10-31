package data

type ColocacionEspacioAcademico struct {
	Id                             int    `json:"Id,omitempty"`
	EspacioAcademicoId             string `json:"EspacioAcademicoId,omitempty"`
	EspacioFisicoId                int    `json:"EspacioFisicoId,omitempty"`
	ColocacionEspacioAcademico     string `json:"ColocacionEspacioAcademico,omitempty"`
	ResumenColocacionEspacioFisico string `json:"ResumenColocacionEspacioFisico,omitempty"`
	Activo                         bool   `json:"Activo,omitempty"`
	FechaCreacion                  string `json:"FechaCreacion,omitempty"`
	FechaModificacion              string `json:"FechaModificacion,omitempty"`
}

type ResumenColocacion struct {
	Colocacion    string        `json:"colocacion,omitempty"`
	EspacioFisico EspacioFisico `json:"espacio_fisico,omitempty"`
}

type EspacioFisico struct {
	SedeId     string `json:"sede_id"`
	EdificioId string `json:"edificio_id"`
	SalonId    string `json:"salon_id"`
}
