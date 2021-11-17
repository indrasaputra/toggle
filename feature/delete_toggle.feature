Feature: Delete toggle
    
    In order to remove toggle permanently
    I need to delete it first

    Scenario: Non-exists toggle can't be deleted
        Given the toggle is empty
        When I delete toggle with key "toggle-1"
        Then response status code must be 404
        And response must match json
            """
            {
                "code": 5,
                "message": "",
                "details": [
                    {
                        "@type": "type.googleapis.com/proto.indrasaputra.toggle.v1.ToggleError",
                        "errorCode": "TOGGLE_ERROR_CODE_NOT_FOUND"
                    }
                ]
            }
            """

    Scenario: Enabled toggle can't be deleted
        Given there are toggles with
            | {"key": "toggle-1"} |
        And I enable toggle with key "toggle-1"
        When I delete toggle with key "toggle-1"
        Then response status code must be 400
        And response must match json
            """
            {
                "code": 9,
                "message": "toggle's is ENABLED hence it can't be deleted",
                "details": [
                    {
                        "@type": "type.googleapis.com/proto.indrasaputra.toggle.v1.ToggleError",
                        "errorCode": "TOGGLE_ERROR_CODE_PROHIBITED_TO_DELETE"
                    }
                ]
            }
            """

    Scenario: Newly created toggle can be deleted directly
        Given there are toggles with
            | {"key": "toggle-1"} |
        When I delete toggle with key "toggle-1"
        Then response status code must be 200
        And response must match json
            """
            {}
            """

    Scenario: Disabled toggle can be deleted
        Given there are toggles with
            | {"key": "toggle-1"} |
        And I disable toggle with key "toggle-1"
        When I delete toggle with key "toggle-1"
        Then response status code must be 200
        And response must match json
            """
            {}
            """