Feature: Create new toggle
    
    In order to use toggle service
    I need to create it first

    Scenario: Invalid json request body (string)
        Given the toggle is empty
        When I create toggle with body
            | string |
        Then response status code must be 400
        And response must match json
            """
            {
                "code": 3,
                "message": "invalid character 's' looking for beginning of value",
                "details": []
            }
            """

    Scenario: Invalid json request body (integer)
        Given the toggle is empty
        When I create toggle with body
            | integer |
        Then response status code must be 400
        And response must match json
            """
            {
                "code": 3,
                "message": "invalid character 'i' looking for beginning of value",
                "details": []
            }
            """

    Scenario: Invalid json request body (double)
        Given the toggle is empty
        When I create toggle with body
            | double |
        Then response status code must be 400
        And response must match json
            """
            {
                "code": 3,
                "message": "invalid character 'd' looking for beginning of value",
                "details": []
            }
            """

    Scenario: Invalid json request body (key doesn't exist)
        Given the toggle is empty
        When I create toggle with body
            | {"toggle": "toggle"} |
        Then response status code must be 400
        And response must match json
            """
            {
                "code": 3,
                "message": "",
                "details": [
                    {
                        "@type": "type.googleapis.com/google.rpc.BadRequest",
                        "fieldViolations": [
                            {
                            "field": "key",
                            "description": "contain character outside of alphanumeric and dash"
                            }
                        ]
                    },
                    {
                        "@type": "type.googleapis.com/proto.indrasaputra.toggle.v1.ToggleError",
                        "errorCode": "TOGGLE_ERROR_CODE_INVALID_KEY"
                    }
                ]
            }
            """

    Scenario: Invalid json request body (key doesn't exist)
        Given the toggle is empty
        When I create toggle with body
            | {"value": "value"} |
        Then response status code must be 400
        And response must match json
            """
            {
                "code": 3,
                "message": "",
                "details": [
                    {
                        "@type": "type.googleapis.com/google.rpc.BadRequest",
                        "fieldViolations": [
                            {
                            "field": "key",
                            "description": "contain character outside of alphanumeric and dash"
                            }
                        ]
                    },
                    {
                        "@type": "type.googleapis.com/proto.indrasaputra.toggle.v1.ToggleError",
                        "errorCode": "TOGGLE_ERROR_CODE_INVALID_KEY"
                    }
                ]
            }
            """

    Scenario: Invalid json request body (key doesn't contain alphanumeric or dash)
        Given the toggle is empty
        When I create toggle with body
            | {"key": "___"} |
        Then response status code must be 400
        And response must match json
            """
            {
                "code": 3,
                "message": "",
                "details": [
                    {
                        "@type": "type.googleapis.com/google.rpc.BadRequest",
                        "fieldViolations": [
                            {
                            "field": "key",
                            "description": "contain character outside of alphanumeric and dash"
                            }
                        ]
                    },
                    {
                        "@type": "type.googleapis.com/proto.indrasaputra.toggle.v1.ToggleError",
                        "errorCode": "TOGGLE_ERROR_CODE_INVALID_KEY"
                    }
                ]
            }
            """

    Scenario: Invalid json request body (key doesn't contain alphanumeric or dash)
        Given the toggle is empty
        When I create toggle with body
            | {"key": "-!@"} |
        Then response status code must be 400
        And response must match json
            """
            {
                "code": 3,
                "message": "",
                "details": [
                    {
                        "@type": "type.googleapis.com/google.rpc.BadRequest",
                        "fieldViolations": [
                            {
                            "field": "key",
                            "description": "contain character outside of alphanumeric and dash"
                            }
                        ]
                    },
                    {
                        "@type": "type.googleapis.com/proto.indrasaputra.toggle.v1.ToggleError",
                        "errorCode": "TOGGLE_ERROR_CODE_INVALID_KEY"
                    }
                ]
            }
            """
    
    Scenario: Invalid json request body (key doesn't contain alphanumeric or dash)
        Given the toggle is empty
        When I create toggle with body
            | {"key": "%^&*!@#$%^&*()"} |
        Then response status code must be 400
        And response must match json
            """
            {
                "code": 3,
                "message": "",
                "details": [
                    {
                        "@type": "type.googleapis.com/google.rpc.BadRequest",
                        "fieldViolations": [
                            {
                            "field": "key",
                            "description": "contain character outside of alphanumeric and dash"
                            }
                        ]
                    },
                    {
                        "@type": "type.googleapis.com/proto.indrasaputra.toggle.v1.ToggleError",
                        "errorCode": "TOGGLE_ERROR_CODE_INVALID_KEY"
                    }
                ]
            }
            """

    Scenario: Valid json request body
        Given the toggle is empty
        When I create toggle with body
            | {"key": "toggle-1", "description": "description 1"} |
            | {"key": "toggle-2", "description": "description 2"} |
            | {"key": "toggle-3", "description": "description 3"} |
            | {"key": "toggle-4", "description": "description 4"} |
            | {"key": "toggle-5", "description": "description 5"} |
        Then response status code must be 200
        And response must match json
            """
            {}
            """
    
    Scenario: Create new toggle with already exist keys
        Given there are toggles with
            | {"key": "toggle-1"} |
            | {"key": "toggle-2"} |
            | {"key": "toggle-3"} |
            | {"key": "toggle-4"} |
            | {"key": "toggle-5"} |
        When I create toggle with body
            | {"key": "toggle-1"} |
            | {"key": "toggle-2"} |
            | {"key": "toggle-3"} |
            | {"key": "toggle-4"} |
            | {"key": "toggle-5"} |
        Then response status code must be 409
        And response must match json
            """
            {
                "code": 6,
                "message": "",
                "details": [
                    {
                        "@type": "type.googleapis.com/proto.indrasaputra.toggle.v1.ToggleError",
                        "errorCode": "TOGGLE_ERROR_CODE_ALREADY_EXISTS"
                    }
                ]
            }
            """