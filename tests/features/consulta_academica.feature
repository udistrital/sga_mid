Feature: Validate API responses
    /v1/consulta_academica/
    probe JSON responses

Scenario Outline: To probe route code response /consulta_academica
    When I send "<method>" request to "<route>" where body is json "<bodyreq>"
    Then the response code should be "<codres>"

    Examples:
    |method |route       |bodyreq                       |codres       |
    |GET    |/v1/consulta_academica  |./assets/requests/empty.json  |200 OK       |
    |GET    |/v1/consulta_academica/0|./assets/requests/empty.json  |200 OK       |
    |POST   |/v1/consulta_academica/0|./assets/requests/empty.json  |404 Not Found|
    |POST   |/v1/consulta_academica  |./assets/requests/empty.json  |404 Not Found|
    |PUT    |/v1/consulta_academica  |./assets/requests/empty.json  |404 Not Found|
    |PUT    |/v1/consulta_academica/0|./assets/requests/empty.json  |404 Not Found|
    |DELETE |/v1/consulta_academica  |./assets/requests/empty.json  |404 Not Found|
    |DELETE |/v1/consulta_academica/0|./assets/requests/empty.json  |404 Not Found|
    |GET    |/v1/consulta_academicas |./assets/requests/empty.json  |404 Not Found|
    |POST   |/v1/consulta_academicas |./assets/requests/empty.json  |404 Not Found|
    |PUT    |/v1/consulta_academicas |./assets/requests/empty.json  |404 Not Found|
    |DELETE |/v1/consulta_academicas |./assets/requests/empty.json  |404 Not Found|

Scenario Outline: To probe response route /evento        
    When I send "<method>" request to "<route>" where body is json "<bodyreq>"
    Then the response code should be "<codres>"      
    And the response should match json "<bodyres>"

    Examples: 
    |method |route                   |bodyreq                      |codres           |bodyres                                                     |                                                    
    |GET    |/v1/consulta_academica/0|./assets/requests/empty.json |200 OK           |./assets/responses/consulta_academica/get_invalid.json      |
