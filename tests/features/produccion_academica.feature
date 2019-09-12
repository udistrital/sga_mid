Feature: Validate API responses
    /v1/produccion_academica/
    probe JSON responses

Scenario Outline: To probe route code response /produccion_academica
    When I send "<method>" request to "<route>" where body is json "<bodyreq>"
    Then the response code should be "<codres>"

    Examples:
    |method |route                     |bodyreq                       |codres       |
    |GET    |/v1/produccion_academica  |./assets/requests/empty.json  |404 Not Found|
    |GET    |/v1/produccion_academica/0|./assets/requests/empty.json  |200 OK       |
    |POST   |/v1/produccion_academica/0|./assets/requests/empty.json  |404 Not Found|
    |POST   |/v1/produccion_academica  |./assets/requests/empty.json  |200 OK       |
    |PUT    |/v1/produccion_academica  |./assets/requests/empty.json  |404 Not Found|
    |PUT    |/v1/produccion_academica/0|./assets/requests/empty.json  |200 OK       |
    |DELETE |/v1/produccion_academica  |./assets/requests/empty.json  |404 Not Found|
    |DELETE |/v1/produccion_academica/0|./assets/requests/empty.json  |200 OK       |
    |GET    |/v1/produccion_academicas |./assets/requests/empty.json  |404 Not Found|
    |POST   |/v1/produccion_academicas |./assets/requests/empty.json  |404 Not Found|
    |PUT    |/v1/produccion_academicas |./assets/requests/empty.json  |404 Not Found|
    |DELETE |/v1/produccion_academicas |./assets/requests/empty.json  |404 Not Found|
    |GET    |/v1/produccion_academica/estado_autor_produccion  |./assets/requests/empty.json  |200 OK       |
    |POST   |/v1/produccion_academica/estado_autor_produccion  |./assets/requests/empty.json  |404 Not Found|
    |DELETE |/v1/produccion_academica/estado_autor_produccion  |./assets/requests/empty.json  |200 OK       |
    |PUT    |/v1/produccion_academica/estado_autor_producciones|./assets/requests/empty.json  |200 OK       |
    |PUT    |/v1/produccion_academica/estado_autor_produccion  |./assets/requests/empty.json  |200 OK       |
    |PUT    |/v1/produccion_academica/estado_autor_produccion/0|./assets/requests/empty.json  |200 OK       |

Scenario Outline: To probe response route /produccion_academica        
    When I send "<method>" request to "<route>" where body is json "<bodyreq>"
    Then the response code should be "<codres>"      
    And the response should match json "<bodyres>"

    Examples: 
    |method |route                     |bodyreq                                                  |codres           |bodyres                                                     |                                                    
    |POST   |/v1/produccion_academica  |./assets/requests/empty.json                             |200 OK           |./assets/responses/invalid_post.json                        |
    |POST   |/v1/produccion_academica  |./assets/requests/produccion_academica/post_invalid.json |200 OK           |./assets/responses/produccion_academica/post_invalid.json   |
    |POST   |/v1/produccion_academica  |./assets/requests/produccion_academica/post_valid.json   |200 OK           |./assets/responses/produccion_academica/post_valid.json     |
    |PUT    |/v1/produccion_academica/1|./assets/requests/empty.json                             |200 OK           |./assets/responses/invalid_post.json                        |
    |PUT    |/v1/produccion_academica/1|./assets/requests/produccion_academica/post_invalid.json |200 OK           |./assets/responses/produccion_academica/post_invalid.json   |
    |PUT    |/v1/produccion_academica/1|./assets/requests/produccion_academica/post_valid.json   |200 OK           |./assets/responses/produccion_academica/post_valid.json     |
    |GET    |/v1/produccion_academica/0|./assets/requests/empty.json                             |200 OK           |./assets/responses/produccion_academica/get_1.json   |
    |GET    |/v1/produccion_academica/0|./assets/requests/empty.json                             |200 OK           |./assets/responses/produccion_academica/get_2.json   |
