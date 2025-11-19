package models

type DuenoRecibo struct {
	Identificacion        string `json:"IDENTIFICACION"`
	CodTipoIdentificacion string `json:"COD_TIPO_IDENTIFICACION"`
	TipoIdentificacion    string `json:"TIPO_IDENTIFICACION"`
	Nombre                string `json:"NOMBRE"`
	CorreoElectronico     string `json:"CORREO_ELECTRONICO"`
	CodEstudiante         string `json:"COD_ESTUDIANTE"`
	CodCarrera            string `json:"COD_CARRERA"`
	Carrera               string `json:"CARRERA"`
	Nivel                 string `json:"NIVEL"`
}

type ConceptoRecibo struct {
	CodConcepto string `json:"COD_CONCEPTO"`
	Concepto    string `json:"CONCEPTO"`
	Valor       string `json:"VALOR"`
}
