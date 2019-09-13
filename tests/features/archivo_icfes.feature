Feature: Validate API responses
    /v1/archivo_icfes/
    probe JSON responses

Scenario Outline: To probe response route /archivo_icfes        
    When I send "<method>" request to "<route>" where  where body is multipart/form-data with this params "<bodyreq>" and the file "<bodyfile>"
    Then the response code should be "<codres>"      
    And the response should match json "<bodyres>"

    Examples: 
    |method |route             |bodyreq                      |bodyfile                     |codres           |bodyres                                                     |                                                    
    |POST   |/v1/archivo_icfes |./assets/requests/empty.json |./assets/requests/empty.json |200 OK           |./assets/responses/produccion_academica/get_empty.json      |
