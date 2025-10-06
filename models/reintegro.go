package models

type Reintegro struct {
	Id                int    `json:"id"`
	CodigoEstudiante  string `json:"codigo_estudiante"`
	MotivoRetiro      string `json:"motivo_retiro"`
	Activo            bool   `json:"activo"`
	FechaCreacion     string `json:"fecha_creacion"`
	FechaModificacion string `json:"fecha_modificacion"`
	InscripcionId     int    `json:"inscripcion_id"`
	Telefono1         string `json:"telefono_1"`
	Telefono2         string `json:"telefono_2"`
}
