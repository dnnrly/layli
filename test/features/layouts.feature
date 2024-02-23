Feature: Layout behaviour

    Background: Reset test data
        Given the test fixures have been reset

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
    Scenario: Generates an image with random shortest square nodes
        When the app runs with parameters "tmp/fixtures/inputs/random-shortest-square.layli"
        Then the app exits without error
        And a file "tmp/fixtures/inputs/random-shortest-square.svg" exists
        And in the SVG file, all node text fits inside the node boundaries
        And in the SVG file, nodes do not overlap
        And in the SVG file, all nodes fit on the image

    @Acceptance
    Scenario: Arranges paths to prevent blockages
        When the app runs with parameters "tmp/fixtures/inputs/blocked.layli"
        Then the app exits without error
        And a file "tmp/fixtures/inputs/blocked.svg" exists

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
    Scenario: Generates layout with absolute positions
        When the app runs with parameters "tmp/fixtures/inputs/absolute-layout.layli"
        Then the app exits without error
        And a file "tmp/fixtures/inputs/absolute-layout.svg" exists
        And in the SVG file, nodes do not overlap
        And no paths cross

    @Acceptance
    Scenario: Generates layout with styles
        When the app runs with parameters "tmp/fixtures/inputs/styles.layli"
        Then the app exits without error
        And a file "tmp/fixtures/inputs/styles.svg" exists
        And in the SVG file, style class "class-1" exists
        And in the SVG file, style class "class-2" exists
        And in the SVG file, element "b" has class "class-2"
        And in the SVG file, element "c" has class "class-1"
        And in the SVG file, element "a" has style "fill:cyan; stroke:red;"
        And in the SVG file, element "c" has style "fill:cyan; stroke:red;"
        And in the SVG file, element "edge-1" has class "path-line class-1"
        And in the SVG file, element "edge-2" has style "stroke:green;"

    @Acceptance
    Scenario: Embeds layout details in output
        When the app runs with parameters "tmp/fixtures/inputs/9-nodes.layli"
        Then the app exits without error
        And a file "tmp/fixtures/inputs/9-nodes.svg" exists
        And in the SVG file, element "a" has attribute "data-pos-x" with value "3"
        And in the SVG file, element "a" has attribute "data-pos-y" with value "3"
        And in the SVG file, element "a" has attribute "data-order" with value "0"
        And in the SVG file, element "a" has attribute "data-width" with value "7"
        And in the SVG file, element "a" has attribute "data-height" with value "4"
        And in the SVG file, element "g" has attribute "data-order" with value "6"
        And in the SVG file, element "edge-2" has attribute "data-order" with value "1"
        And in the SVG file, element "edge-2" has attribute "data-from" with value "b"
        And in the SVG file, element "edge-2" has attribute "data-to" with value "c"
