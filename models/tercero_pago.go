package models

type TerceroPago struct {
	TERPA_ANO_PAGO         int     `json:"TERPA_ANO_PAGO"`
	TERPA_DIGITO_CHEQUEO   *int    `json:"TERPA_DIGITO_CHEQUEO"`
	TERPA_DIRECCION        string  `json:"TERPA_DIRECCION"`
	TERPA_EMAIL            string  `json:"TERPA_EMAIL"`
	TERPA_ESTADO_REGISTRO  string  `json:"TERPA_ESTADO_REGISTRO"`
	TERPA_FECHA_REGISTRO   string  `json:"TERPA_FECHA_REGISTRO"`
	TERPA_NATURALEZA       string  `json:"TERPA_NATURALEZA"`
	TERPA_NRO_DOCUMENTO    int     `json:"TERPA_NRO_DOCUMENTO"`
	TERPA_PRIMER_APELLIDO  string  `json:"TERPA_PRIMER_APELLIDO"`
	TERPA_PRIMER_NOMBRE    string  `json:"TERPA_PRIMER_NOMBRE"`
	TERPA_RAZON_SOCIAL     *string `json:"TERPA_RAZON_SOCIAL"`
	TERPA_SECUENCIA        int     `json:"TERPA_SECUENCIA"`
	TERPA_SEGUNDO_APELLIDO string  `json:"TERPA_SEGUNDO_APELLIDO"`
	TERPA_SEGUNDO_NOMBRE   string  `json:"TERPA_SEGUNDO_NOMBRE"`
	TERPA_TDO_CODVAR       string  `json:"TERPA_TDO_CODVAR"`
	TERPA_TELEFONO         int     `json:"TERPA_TELEFONO"`
	// TERPA_DATOS_ADICIONALES string  `json:"TERPA_DATOS_ADICIONALES"`
}

type TerceroPagoRequest struct {
	PostTerceroPago TerceroPago `json:"_posttercero_pago"`
}
