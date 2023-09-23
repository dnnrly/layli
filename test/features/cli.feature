Feature: Simple CLI commands

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
        When the app runs with parameters "tmp/fixtures/inputs/hello-world.layli"
        Then the app exits without error
        And a file "tmp/fixtures/inputs/hello-world.svg" exists
        And in the SVG file, all node text fits inside the node boundaries
        And in the SVG file, grid dots are not shown

    @Acceptance
    Scenario: Generates a 2 node diagram
        When the app runs with parameters "tmp/fixtures/inputs/2-nodes.layli"
        Then the app exits without error
        And a file "tmp/fixtures/inputs/2-nodes.svg" exists
        And the number of nodes is 2
        And the number of paths is 1
        And in the SVG file, all node text fits inside the node boundaries
        And in the SVG file, nodes do not overlap
        And in the SVG file, all nodes fit on the image

    @Acceptance
    Scenario: Generates an image with smallest area
        When the app runs with parameters "tmp/fixtures/inputs/14-nodes.layli"
        Then the app exits without error
        And a file "tmp/fixtures/inputs/14-nodes.svg" exists
        And the number of nodes is 14
        And in the SVG file, all node text fits inside the node boundaries
        And in the SVG file, nodes do not overlap
        And in the SVG file, all nodes fit on the image

    @Acceptance
    Scenario: Generates an image with a square number of nodes
        When the app runs with parameters "tmp/fixtures/inputs/9-nodes.layli"
        Then the app exits without error
        And a file "tmp/fixtures/inputs/9-nodes.svg" exists
        And the number of nodes is 9
        And in the SVG file, all node text fits inside the node boundaries
        And in the SVG file, nodes do not overlap
        And in the SVG file, all nodes fit on the image

    @Acceptance
    Scenario: Generates an image with topological sorted nodes
        When the app runs with parameters "tmp/fixtures/inputs/topological.layli"
        Then the app exits without error
        And a file "tmp/fixtures/inputs/topological.svg" exists
        And the number of nodes is 5
        And in the SVG file, all node text fits inside the node boundaries
        And in the SVG file, nodes do not overlap
        And in the SVG file, all nodes fit on the image

    @Acceptance
    Scenario: Sets output file correctly
        When the app runs with parameters "--output tmp/another-file.svg tmp/fixtures/inputs/2-nodes.layli"
        Then the app exits without error
        And a file "tmp/another-file.svg" exists

    @Acceptance
    Scenario: Shows path grid positions
        When the app runs with parameters "--show-grid --output tmp/fixtures/inputs/2-nodes-with-grid.svg tmp/fixtures/inputs/2-nodes.layli"
        Then the app exits without error
        And a file "tmp/fixtures/inputs/2-nodes-with-grid.svg" exists
        And in the SVG file, path grid dots are shown

    @Acceptance
    Scenario: Corrects crossing lines without moving nodes
        When the app runs with parameters "tmp/fixtures/inputs/crossing-lines.layli"
        Then the app exits without error
        And a file "tmp/fixtures/inputs/crossing-lines.svg" exists
        And the number of paths is 6
        And no paths cross

    @Acceptance
    Scenario: Errors when cannot find paths without crossing
        When the app runs with parameters "tmp/fixtures/inputs/impossible-paths.layli"
        Then the app exits with an error
        And the app output contains "cannot find a path between node2 and node3"

    @Acceptance
    Scenario: Errors when cannot write output
        When the app runs with parameters "--output / tmp/fixtures/inputs/2-nodes.layli"
        Then the app exits with an error
        And the app output contains "drawing diagram"
