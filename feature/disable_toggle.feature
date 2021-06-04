Feature: Disable toggle
    
    In order to delete toggle
    I need to disable it first

    Scenario: Non-exists toggle can't be disabled
        Given the toggle is empty
        When I disable toggle with key "toggle-1"
        Then response status code must be 404
        And response must match json
            """
            {
                "code": 5,
                "message": "",
                "details": [
                    {
                        "@type": "type.googleapis.com/proto.indrasaputra.toggle.v1.ToggleError",
                        "errorCode": "NOT_FOUND"
                    }
                ]
            }
            """

    Scenario: existing toggle can be disabled
        Given there are toggles with
            | {"key": "toggle-1", "description": "description 1"} |
        When I disable toggle with key "toggle-1"
        Then response status code must be 200
        And response must match json
            """
            {}
            """