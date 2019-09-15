Feature: Validate API responses
    /v1/evento/
    probe JSON responses

Scenario Outline: To probe route code response /evento
    When I send "<method>" request to "<route>" where body is json "<bodyreq>"
    Then the response code should be "<codres>"

    Examples:
    |method |route       |bodyreq                       |codres       |
    |GET    |/v1/evento  |./assets/requests/empty.json  |404 Not Found|
    |GET    |/v1/evento/0|./assets/requests/empty.json  |200 OK       |
    |POST   |/v1/evento/0|./assets/requests/empty.json  |404 Not Found|
    |POST   |/v1/evento  |./assets/requests/empty.json  |200 OK       |
    |PUT    |/v1/evento  |./assets/requests/empty.json  |404 Not Found|
    |PUT    |/v1/evento/0|./assets/requests/empty.json  |200 OK       |
    |DELETE |/v1/evento  |./assets/requests/empty.json  |404 Not Found|
    |DELETE |/v1/evento/0|./assets/requests/empty.json  |200 OK       |
    |GET    |/v1/eventos |./assets/requests/empty.json  |404 Not Found|
    |POST   |/v1/eventos |./assets/requests/empty.json  |404 Not Found|
    |PUT    |/v1/eventos |./assets/requests/empty.json  |404 Not Found|
    |DELETE |/v1/eventos |./assets/requests/empty.json  |404 Not Found|
