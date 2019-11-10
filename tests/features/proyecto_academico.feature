Feature: Validate API responses
    /v1/proyecto_academico/
    probe JSON responses

Scenario Outline: To probe route code response /consulta_proyecto_academico
    When I send "<method>" request to "<route>" where body is json "<bodyreq>"
    Then the response code should be "<codres>"

    Examples:
    |method |route                            |bodyreq                       |codres       |
    # get_all proyectos
    |GET    |/v1/consulta_proyecto_academico  |./assets/requests/empty.json  |200 OK       |
    |GET    |/v1/consulta_proyecto_academica  |./assets/requests/empty.json  |404 Not Found|
    # get_by_id
    |GET    |/v1/consulta_proyecto_academico/0|./assets/requests/empty.json  |200 OK       |
    |GET    |/v1/consulta_proyecto_academica/0|./assets/requests/empty.json  |404 Not Found|

Scenario Outline: To probe response route /consulta_proyecto_academico        
    When I send "<method>" request to "<route>" where body is json "<bodyreq>"
    Then the response code should be "<codres>"      
    And the response should match json "<bodyres>"

    Examples: 
    |method |route                                  |bodyreq                              |codres           |bodyres                                                   |                                                    
    # get_all proyectos
    # |GET    |/v1/consulta_proyecto_academico        |./assets/requests/empty.json         |200 OK           |./assets/responses/proyecto_academico/get_empty.json      |
    # get_by_id
    |GET    |/v1/consulta_proyecto_academico/0      |./assets/requests/empty.json         |200 OK           |./assets/responses/proyecto_academico/get_empty.json      |
    |GET    |/v1/consulta_proyecto_academico/9999999|./assets/requests/empty.json         |200 OK           |./assets/responses/proyecto_academico/get_valid.json      |


Scenario Outline: To probe route code response /proyecto_academico
    When I send "<method>" request to "<route>" where body is json "<bodyreq>"
    Then the response code should be "<codres>"

    Examples:
    |method |route                     |bodyreq                     |codres       |
    |POST   |/v1/proyecto_academico/0|./assets/requests/empty.json  |404 Not Found|
    |POST   |/v1/proyecto_academico  |./assets/requests/empty.json  |200 OK       |
    # coordinador
    |POST   |/v1/proyecto_academico/coordinador    |./assets/requests/empty.json  |404 Not Found|
    |POST   |/v1/proyecto_academico/coordinador/0  |./assets/requests/empty.json  |200 OK       |
    # registro_calificado
    |POST   |/v1/proyecto_academico/registro_calificado    |./assets/requests/empty.json  |404 Not Found|
    |POST   |/v1/proyecto_academico/registro_calificado/0  |./assets/requests/empty.json  |200 OK       |
    # registro_alta_calidad
    |POST   |/v1/proyecto_academico/registro_alta_calidad    |./assets/requests/empty.json  |404 Not Found|
    |POST   |/v1/proyecto_academico/registro_alta_calidad/0  |./assets/requests/empty.json  |200 OK       |

Scenario Outline: To probe response route /proyecto_academico        
    When I send "<method>" request to "<route>" where body is json "<bodyreq>"
    Then the response code should be "<codres>"      
    And the response should match json "<bodyres>"

        Examples: 
    |method |route                   |bodyreq                                                |codres           |bodyres                                                   |                                                    
    |POST   |/v1/proyecto_academico  |./assets/requests/empty.json                           |200 OK           |./assets/responses/invalid_post.json                      |
    |POST   |/v1/proyecto_academico  |./assets/requests/proyecto_academico/post_invalid.json |200 OK           |./assets/responses/proyecto_academico/post_invalid.json   |
    |POST   |/v1/proyecto_academico  |./assets/requests/proyecto_academico/post_valid.json   |200 OK           |./assets/responses/proyecto_academico/post_valid.json     |
    # coordinador
    |POST   |/v1/proyecto_academico/coordinador/1  |./assets/requests/empty.json                                       |200 OK           |./assets/responses/invalid_post.json                                  |
    |POST   |/v1/proyecto_academico/coordinador/1  |./assets/requests/proyecto_academico/post_coordinador_invalid.json |200 OK           |./assets/responses/proyecto_academico/post_coordinador_invalid.json   |
    |POST   |/v1/proyecto_academico/coordinador/1  |./assets/requests/proyecto_academico/post_coordinador_valid.json   |200 OK           |./assets/responses/proyecto_academico/post_coordinador_valid.json     |
    # registro_calificado
    |POST   |/v1/proyecto_academico/registro_calificado/1  |./assets/requests/empty.json                                               |200 OK           |./assets/responses/invalid_post.json                                          |
    |POST   |/v1/proyecto_academico/registro_calificado/1  |./assets/requests/proyecto_academico/post_registro_calificado_invalid.json |200 OK           |./assets/responses/proyecto_academico/post_registro_calificado_invalid.json   |
    |POST   |/v1/proyecto_academico/registro_calificado/1  |./assets/requests/proyecto_academico/post_registro_calificado_valid.json   |200 OK           |./assets/responses/proyecto_academico/post_registro_calificado_valid.json     |
    # registro_alta_calidad
    |POST   |/v1/proyecto_academico/registro_alta_calidad/1  |./assets/requests/empty.json                                                 |200 OK           |./assets/responses/invalid_post.json                                            |
    |POST   |/v1/proyecto_academico/registro_alta_calidad/1  |./assets/requests/proyecto_academico/post_registro_alta_calidad_invalid.json |200 OK           |./assets/responses/proyecto_academico/post_registro_alta_calidad_invalid.json   |
    |POST   |/v1/proyecto_academico/registro_alta_calidad/1  |./assets/requests/proyecto_academico/post_registro_alta_calidad_valid.json   |200 OK           |./assets/responses/proyecto_academico/post_registro_alta_calidad_valid.json     |
