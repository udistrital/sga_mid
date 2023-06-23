package data

type EspacioAcademico struct {
	Id                       string `json:"_id,omitempty"`
	Nombre                   string `json:"nombre,omitempty"`
	Codigo_abreviacion       string `json:"codigo_abreviacion,omitempty"`
	Codigo                   string `json:"codigo,,omitempty"`
	Plan_estudio_id          string `json:"plan_estudio_id,omitempty"`
	Proyecto_academico_id    int64  `json:"proyecto_academico_id,omitempty"`
	Creditos                 int64  `json:"creditos,omitempty"`
	Distribucion_horas       string `json:"distribucion_horas,omitempty"`
	Tipo_espacio_id          int64  `json:"tipo_espacio_id,omitempty"`
	Clasificacion_espacio_id int64  `json:"clasificacion_espacio_id,omitempty"`
	Enfoque_id               int64  `json:"enfoque_id,omitempty"`
	Espacios_requeridos      string `json:"espacios_requeridos,omitempty"`
	Grupo                    string `json:"grupo,omitempty"`
	Inscritos                int64  `json:"inscritos,omitempty"`
	Periodo_id               int64  `json:"periodo_id,omitempty"`
	Docente_id               int64  `json:"docente_id,omitempty"`
	Horario_id               string `json:"horario_id,omitempty"`
	Espacio_academico_padre  string `json:"espacio_academico_padre,omitempty"`
	Soporte_documental       string `json:"soporte_documental,omitempty"`
	Estado_aprobacion_id     string `json:"estado_aprobacion_id,omitempty"`
	Observacion              string `json:"observacion,omitempty"`
	Activo                   bool   `json:"activo,omitempty"`
	Fecha_creacion           string `json:"fecha_creacion,omitempty"`
	Fecha_modificacion       string `json:"fecha_modificacion,omitempty"`
}

type EstadoAprobacion struct {
	Id                 string `json:"_id,omitempty"`
	Nombre             string `json:"nombre,omitempty"`
	Descripcion        string `json:"descripcion,omitempty"`
	Codigo_abreviacion string `json:"codigo_abreviacion,omitempty"`
	Activo             bool   `json:"activo,omitempty"`
	Fecha_creacion     string `json:"fecha_creacion,omitempty"`
	Fecha_modificacion string `json:"fecha_modificacion,omitempty"`
}

type EspacioAcademicoEstudiantes struct {
	Id                      string `json:"_id,omitempty"`
	Espacio_academico_padre string `json:"espacio_academico_padre,omitempty"`
	Estudiante_id           int64  `json:"estudiante_id,omitempty"`
	Periodo_id              int64  `json:"periodo_id,omitempty"`
	Activo                  bool   `json:"activo,omitempty"`
	Fecha_creacion          string `json:"fecha_creacion,omitempty"`
	Fecha_modificacion      string `json:"fecha_modificacion,omitempty"`
}
