Feature: Validate CLI configuration

    Background: Reset test data
        Given the test fixures have been reset

    @Acceptance
    Scenario: Prints help correctly
        When the app runs with parameters "-h"
        Then the app exits without error
        And the app output contains "Usage:"

    @Acceptance
    Scenario: Prints on bad config
        When the app runs with parameters "tmp/fixtures/inputs/bad-config.layli"
        Then the app exits with an error
        And the app output contains "creating config:"

    @Acceptance
    Scenario: Sets output file correctly
        When the app runs with parameters "--output tmp/another-file.svg tmp/fixtures/inputs/2-nodes.layli"
        Then the app exits without error
        And a file "tmp/another-file.svg" exists

    
    @Acceptance
    Scenario: Non-existent file returns error
        When the app runs with parameters "non-existant.layli"
        Then the app exits with an error
        And the app output contains "opening input: open non-existant.layli: no such file or directory"

    @Acceptance
    Scenario: Errors when cannot write output
        When the app runs with parameters "--output / tmp/fixtures/inputs/2-nodes.layli"
        Then the app exits with an error
        And the app output contains "drawing diagram"
