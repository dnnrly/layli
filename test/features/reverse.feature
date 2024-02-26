Feature: Image reveral

    Background: Reset test data
        Given the test fixures have been reset

    @Acceptance
    Scenario: Flow generated image can be reversed into absolute
        When the app runs with parameters "tmp/fixtures/inputs/2-nodes.layli"
        Then the app exits without error
        And the app runs with parameters "to-absolute tmp/fixtures/inputs/2-nodes.layli -o tmp/absolute.layli"
        Then the app exits without error
        And a layli file "tmp/absolute.layli" exists
        And the layli file contains the following nodes:
        | id | x | y |
        | 1  | 3 | 3 |
        | 2  | 7 | 3 |
