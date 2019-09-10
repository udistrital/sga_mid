Feature: Validate API responses
    SGA_MID
    probe JSON responses

Scenario Outline: To probe route code response /evento
    When I send "<method>" request to "<route>" where body is json "<bodyreq>"
    Then the response code should be "<codres>"

    Examples:
    |method |route       |bodyreq                       |codres       |
    |GET    |/v1/evento  |./assets/requests/empty.json  |200 OK       |
    |GET    |/v1/eventos |./assets/requests/empty.json  |404 Not Found|
    |POST   |/v1/eventos |./assets/requests/empty.json  |404 Not Found|
    |PUT    |/v1/eventos |./assets/requests/empty.json  |404 Not Found|
    |DELETE |/v1/eventos |./assets/requests/empty.json  |404 Not Found|

Scenario Outline: To probe response route /tipo_periodo       
    When I send "<method>" request to "<route>" where body is json "<bodyreq>"
    Then the response code should be "<codres>"      
    And the response should match json "<bodyres>"

    Examples: 
    |method |route       |bodyreq                      |codres           |bodyres                   |
    |GET    |/v1/evento  |./assets/requests/empty.json |200 OK           |./assets/responses/eventos/1.json    |
    |POST   |/v1/evento  |./assets/requests/empty.json |400 Bad Request  |./assets/responses/eventos/1.json   |
