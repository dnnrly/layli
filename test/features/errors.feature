Feature: Error handling

    Background: Reset test data
        Given the test fixures have been reset

    @Acceptance
    Scenario: Errors when cannot find paths without crossing
        When the app runs with parameters "tmp/fixtures/inputs/impossible-paths.layli"
        Then the app exits with an error
        And the app output contains "finding path between node2 and node3: no path found"

    @Acceptance
    Scenario: Errors when overlapping absolute nodes defined
        When the app runs with parameters "tmp/fixtures/inputs/absolute-layout-overlap.layli"
        Then the app exits with an error
        And the app output contains "arranging nodes: nodes a and b margins overlap"

    @Acceptance
    Scenario: Errors when not specifying output for to-absolute
        When the app runs with parameters "to-absolute tmp/fixtures/inputs/2-nodes.svg"
        Then the app exits with an error

    @Acceptance
    Scenario: Errors when cannot find input for to-absolute
        When the app runs with parameters "to-absolute non-existent.svg -o tmp/absolute.layli"
        Then the app exits with an error

    @Acceptance
    Scenario: Errors when not specifying invalid input for to-absolute
        When the app runs with parameters "to-absolute . -o tmp/absolute.layli"
        Then the app exits with an error

    @Acceptance
    Scenario: Errors when cannot parse input for to-absolute
        When the app runs with parameters "to-absolute tmp/fixtures/inputs/2-nodes.layli -o tmp/absolute.layli"
        Then the app exits with an error
