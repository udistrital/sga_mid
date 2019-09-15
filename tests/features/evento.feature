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

Scenario Outline: To probe response route /evento        
    When I send "<method>" request to "<route>" where body is json "<bodyreq>"
    Then the response code should be "<codres>"      
    And the response should match json "<bodyres>"

    Examples: 
    |method |route                     |bodyreq                      |codres           |bodyres                                       |                                                    
    |GET    |/v1/evento/0|./assets/requests/empty.json               |200 OK           |./assets/responses/evento/get_empty.json      |
    |GET    |/v1/evento/1|./assets/requests/empty.json               |200 OK           |./assets/responses/evento/get_valid.json      |
    |POST   |/v1/evento  |./assets/requests/empty.json               |200 OK           |./assets/responses/invalid_post.json          |
    |POST   |/v1/evento  |./assets/requests/evento/post_invalid.json |200 OK           |./assets/responses/evento/post_invalid.json   |
    |POST   |/v1/evento  |./assets/requests/evento/post_valid.json   |200 OK           |./assets/responses/evento/post_valid.json     |
    |PUT    |/v1/evento/0|./assets/requests/empty.json               |200 OK           |./assets/responses/invalid_post.json          |
    |PUT    |/v1/evento/1|./assets/requests/evento/put_invalid.json  |200 OK           |./assets/responses/evento/put_invalid.json    |
    |PUT    |/v1/evento/1|./assets/requests/evento/put_valid.json    |200 OK           |./assets/responses/evento/put_valid.json      |
    |DELETE |/v1/evento/0|./assets/requests/empty.json               |200 OK           |./assets/responses/evento/delete_invalid.json |
    |DELETE |/v1/evento/1|./assets/requests/empty.json               |200 OK           |./assets/responses/evento/delete_valid.json   |
