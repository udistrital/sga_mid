package models

type Reintegro struct {
	Id                int    `json:"id"`
	CodigoEstudiante  int64  `json:"codigo_estudiante"`
	MotivoRetiro      string `json:"motivo_retiro"`
	Activo            bool   `json:"activo"`
	FechaCreacion     string `json:"fecha_creacion"`
	FechaModificacion string `json:"fecha_modificacion"`
	InscripcionId     int    `json:"inscripcion_id"`
	Telefono1         int64  `json:"telefono_1"`
	Telefono2         int64  `json:"telefono_2"`
}
