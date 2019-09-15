Feature: Validate API responses
    /v1/archivo_icfes/
    probe JSON responses

Scenario Outline: To probe route code response  /archivo_icfes        
    When I send "<method>" request to "<route>" where body is multipart/form-data with this params "<bodyreq>" and the file "<filename>" located at "<bodyfile>"
    Then the response code should be "<codres>"      

    Examples: 
    |method |route             |bodyreq                                         |filename         |bodyfile                                                   |codres           |bodyres                                                    |                                                    
    |GET    |/v1/archivo_icfes |./assets/requests/archivo_icfes/post_valid.json |archivo_icfes    |./assets/requests/archivo_icfes/post_valid.txt             |404 Not Found    |./assets/responses/archivo_icfes/post_valid.json           |
    |PUT    |/v1/archivo_icfes |./assets/requests/archivo_icfes/post_valid.json |archivo_icfes    |./assets/requests/archivo_icfes/post_valid.txt             |404 Not Found    |./assets/responses/archivo_icfes/post_valid.json           |
    |DELETE |/v1/archivo_icfes |./assets/requests/archivo_icfes/post_valid.json |archivo_icfes    |./assets/requests/archivo_icfes/post_valid.txt             |404 Not Found    |./assets/responses/archivo_icfes/post_valid.json           |
    |POST   |/v1/archivo_icfe  |./assets/requests/archivo_icfes/post_valid.json |archivo_icfes    |./assets/requests/archivo_icfes/post_valid.txt             |404 Not Found    |./assets/responses/archivo_icfes/post_valid.json           |
    |POST   |/v1/archivo_icfes |./assets/requests/archivo_icfes/post_valid.json |archivo_icfes    |./assets/requests/archivo_icfes/post_valid.txt             |200 OK           |./assets/responses/archivo_icfes/post_valid.json           |
    |POST   |/v1/archivo_icfes |./assets/requests/archivo_icfes/post_valid.json |archivo_icfesbad |./assets/requests/archivo_icfes/post_valid.txt             |200 OK           |./assets/responses/archivo_icfes/post_invalid_file.json    |

Scenario Outline: To probe response route /archivo_icfes        
    When I send "<method>" request to "<route>" where body is multipart/form-data with this params "<bodyreq>" and the file "<filename>" located at "<bodyfile>"
    Then the response code should be "<codres>"      
    And the response should match json "<bodyres>"

    Examples: 
    |method |route             |bodyreq                                         |filename         |bodyfile                                                   |codres           |bodyres                                                    |                                                    
    |POST   |/v1/archivo_icfes |./assets/requests/archivo_icfes/post_valid.json |archivo_icfes    |./assets/requests/archivo_icfes/post_valid.txt             |200 OK           |./assets/responses/archivo_icfes/post_valid.json           |
    |POST   |/v1/archivo_icfes |./assets/requests/archivo_icfes/post_valid.json |archivo_icfesbad |./assets/requests/archivo_icfes/post_valid.txt             |200 OK           |./assets/responses/archivo_icfes/post_invalid_file.json    |
    |POST   |/v1/archivo_icfes |./assets/requests/archivo_icfes/post_valid.json |archivo_icfes    |./assets/requests/archivo_icfes/post_invalid_content_1.txt |200 OK           |./assets/responses/archivo_icfes/post_invalid_content.json |
    |POST   |/v1/archivo_icfes |./assets/requests/archivo_icfes/post_valid.json |archivo_icfes    |./assets/requests/archivo_icfes/post_invalid_content_2.txt |200 OK           |./assets/responses/archivo_icfes/post_invalid_content.json |
