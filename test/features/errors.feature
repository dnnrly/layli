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
        And the app output contains "arranging nodes: nodes a and c overlap"
