Feature: Get toggle
    
    In order to use toggle service
    I need to get them/it first

    Scenario: Toggle is empty
        Given the toggle is empty
        When I get all toggles
        Then response status code must be 200
        And response must match json
            """
            {
                "toggles": []
            }
            """

    Scenario: Toggles are filled
        Given there are toggles with
            | {"key": "toggle-1", "description": "description 1"} |
            | {"key": "toggle-2", "description": "description 2"} |
            | {"key": "toggle-3", "description": "description 3"} |
            | {"key": "toggle-4", "description": "description 4"} |
            | {"key": "toggle-5", "description": "description 5"} |
        When I get all toggles
        Then response status code must be 200
        And response toggles should match
            """
            {
                "toggles": [
                    {
                        "key": "toggle-1",
                        "is_enabled": false,
                        "description": "description 1"
                    },
                    {
                        "key": "toggle-2",
                        "is_enabled": false,
                        "description": "description 2"
                    },
                    {
                        "key": "toggle-3",
                        "is_enabled": false,
                        "description": "description 3"
                    },
                    {
                        "key": "toggle-4",
                        "is_enabled": false,
                        "description": "description 4"
                    },
                    {
                        "key": "toggle-5",
                        "is_enabled": false,
                        "description": "description 5"
                    }
                ]
            }
            """

    Scenario: Get single toggle when toggle is empty (not found)
        Given the toggle is empty
        When I get single toggle with key "toggle-1"
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

    Scenario: Get single toggle when toggle exists
        Given there are toggles with
            | {"key": "toggle-1", "description": "description 1"} |
        When I get single toggle with key "toggle-1"
        Then response status code must be 200
        And response single toggle should match
            """
            {
                "toggle": {
                    "key": "toggle-1",
                    "is_enabled": false,
                    "description": "description 1"
                }
            }
            """
