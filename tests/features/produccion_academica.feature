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
    |POST   |/v1/produccion_academica  |./assets/requests/empty.json  |200 OK       |
    |GET    |/v1/produccion_academica/1|./assets/requests/empty.json  |200 OK       |
    |PUT    |/v1/produccion_academica/1|./assets/requests/empty.json  |200 OK       |
    |DELETE |/v1/produccion_academica/1|./assets/requests/empty.json  |200 OK       |
    |GET    |/v1/produccion_academicas |./assets/requests/empty.json  |404 Not Found|
    |POST   |/v1/produccion_academicas |./assets/requests/empty.json  |404 Not Found|
    |PUT    |/v1/produccion_academicas |./assets/requests/empty.json  |404 Not Found|
    |DELETE |/v1/produccion_academicas |./assets/requests/empty.json  |404 Not Found|

Scenario Outline: To probe response route /tipo_periodo       
    When I send "<method>" request to "<route>" where body is json "<bodyreq>"
    Then the response code should be "<codres>"      
    And the response should match json "<bodyres>"

    Examples: 
    |method |route                     |bodyreq                      |codres           |bodyres                                           |
    |GET    |/v1/produccion_academica/0|./assets/requests/empty.json |200 OK           |./assets/responses/produccion_academica/get_1.json   |
    |POST   |/v1/produccion_academica  |./assets/requests/empty.json |400 Bad Request  |./assets/responses/produccion_academica/1.json   |
