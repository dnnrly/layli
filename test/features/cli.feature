Feature: Simple CLI commands

    Background: Reset test data
        Given the test fixures have been reset

    @Acceptance
    Scenario: Prints help correctly
        When the app runs with parameters "-h"
        Then the app exits without error
        And the app output contains "Usage:"

    @Acceptance
    Scenario: Non-existent file returns error
        When the app runs with parameters "non-existant.layli"
        Then the app exits with an error
        And the app output contains "opening input: open non-existant.layli: no such file or directory"

    @Acceptance
    Scenario: Returns error when file not specified
        When the app runs without args
        Then the app exits with an error
        And the app output contains "Error: accepts 1 arg(s), received 0"

    @Acceptance
    Scenario: Generates a single node diagram
        When the app runs with parameters "-o tmp/outputs/ tmp/fixtures/inputs/hello-world.layli"
        Then the app exits without error
        And a file "tmp/outputs/hello-world.svg" exists
