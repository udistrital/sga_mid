Feature: Validate API responses
    /v1/archivo_icfes/
    probe JSON responses

Scenario Outline: To probe response route /archivo_icfes        
    When I send "<method>" request to "<route>" where body is multipart/form-data with this params "<bodyreq>" and the file "<filename>" located at "<bodyfile>"
    Then the response code should be "<codres>"      
    And the response should match json "<bodyres>"

    Examples: 
    |method |route             |bodyreq                                         |filename         |bodyfile                                       |codres           |bodyres                                               |                                                    
    |POST   |/v1/archivo_icfes |./assets/requests/archivo_icfes/post_valid.json |archivo_icfes    |./assets/requests/archivo_icfes/post_valid.txt |200 OK           |./assets/responses/archivo_icfes/post_valid.json      |
