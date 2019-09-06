# sga_mid
Mid del SGA 
## Requirements
Go version >= 1.8.
 ## Preparation:
    Para usar el API, usar el comando:
        - go get github.com/udistrital/sga_mid
 ## Run

## Example:
PERSONAS_SERVICE=localhost:8083/v1 CORE_SERVICE=localhost:8102/v1 PROGRAMA_ACADEMICO_SERVICE=localhost:8101/v1 EVENTOS_SERVICE:localhost:8083/v1 OIKOS_SERVICE:10.20.0.254/oikos_api/v1 SGA_MID_HTTP_PORT=8095 bee run -downdoc=true -gendoc=true