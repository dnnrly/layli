package test_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/antchfx/xmlquery"
	"github.com/cucumber/godog"
	"github.com/otiai10/copy"
	"github.com/stretchr/testify/assert"
)

// nolint: unused
type testContext struct {
	err      error
	cmdInput struct {
		parameters string
	}
	cmdResult struct {
		Output string
		Err    error
	}
	svgOutput struct {
		name     string
		contents []byte
		doc      *xmlquery.Node
	}
}

// Errorf is used by the called assertion to report an error and is required to
// make testify assertions work
func (c *testContext) Errorf(format string, args ...interface{}) {
	c.err = fmt.Errorf(format, args...)
}

func (c *testContext) theTestFixuresHaveBeenReset() error {
	// err := os.RemoveAll("tmp/fixtures")
	// if err != nil {
	// 	return fmt.Errorf("deleting old fixtures: %w", err)
	// }

	err := copy.Copy("fixtures/", "tmp/fixtures")
	if err != nil {
		return fmt.Errorf("copying new fixtures: %w", err)
	}

	return nil
}

func (c *testContext) theAppRunsWithoutArgs() error {
	cmd := exec.Command("../layli")
	output, err := cmd.CombinedOutput()
	c.cmdResult.Output = string(output)
	c.cmdResult.Err = err

	return nil
}

func (c *testContext) theAppRunsWithParameters(args string) error {
	c.cmdInput.parameters = args
	cmdArgs := strings.Split(args, " ")
	cmd := exec.Command("../layli", cmdArgs...)
	output, err := cmd.CombinedOutput()
	c.cmdResult.Output = string(output)
	c.cmdResult.Err = err

	return nil
}

func (c *testContext) theAppExitsWithoutError() error {
	assert.NoError(c, c.cmdResult.Err)
	return c.err
}

func (c *testContext) theAppExitsWithAnError() error {
	assert.Error(c, c.cmdResult.Err)
	return c.err
}

func (c *testContext) theAppOutputContains(expected string) error {
	expected = strings.ReplaceAll(expected, "\\\"", "\"")
	assert.Contains(c, c.cmdResult.Output, expected)
	return c.err
}

func (c *testContext) theAppOutputDoesNotContain(unexpected string) error {
	unexpected = strings.ReplaceAll(unexpected, "\\\"", "\"")
	assert.NotContains(c, c.cmdResult.Output, unexpected)
	return c.err
}

func (c *testContext) aFileExists(file string) error {
	c.svgOutput.name = file
	f, err := os.Open(file)
	if err != nil {
		assert.NoError(c, err)
		return c.err
	}
	c.svgOutput.contents, err = ioutil.ReadAll(f)
	if err != nil {
		assert.NoError(c, err)
		return c.err
	}

	c.svgOutput.doc, err = xmlquery.Parse(bytes.NewBuffer(c.svgOutput.contents))
	if err != nil {
		assert.NoError(c, err)
		return c.err
	}

	return nil
}

func (c *testContext) inTheSVGFileAllNodeTextFitsInsideTheNodeBoundaries() error {
	ids := getNodeIds(c.svgOutput.doc)
	set := map[string]bool{}

	for _, v := range ids {
		id, _, _ := strings.Cut(v, "-")
		set[id] = true
	}

	for n := range set {
		rect := xmlquery.FindOne(c.svgOutput.doc, "//rect[starts-with(@id, '"+n+"')]")
		rectX, _ := strconv.Atoi(rect.SelectAttr("x"))
		rectY, _ := strconv.Atoi(rect.SelectAttr("y"))
		rectWidth, _ := strconv.Atoi(rect.SelectAttr("width"))
		rectHeight, _ := strconv.Atoi(rect.SelectAttr("height"))

		text := xmlquery.FindOne(c.svgOutput.doc, "//text[starts-with(@id, '"+n+"')]")
		textX, _ := strconv.Atoi(text.SelectAttr("x"))
		textY, _ := strconv.Atoi(text.SelectAttr("y"))
		textWidth, _ := strconv.Atoi(text.SelectAttr("width"))
		textHeight, _ := strconv.Atoi(text.SelectAttr("height"))

		assert.Greater(c, textX, rectX, fmt.Sprintf(`"%s" text "%s" does not fit inside node`, n, text.InnerText()))
		assert.Greater(c, textY, rectY, fmt.Sprintf(`"%s" text "%s" does not fit inside node`, n, text.InnerText()))
		assert.Less(c, textX+textWidth, rectX+rectWidth, fmt.Sprintf(`"%s" text "%s" does not fit inside node`, n, text.InnerText()))
		assert.Less(c, textY+textHeight, rectY+rectHeight, fmt.Sprintf(`"%s" text "%s" does not fit inside node`, n, text.InnerText()))
	}

	return c.err
}

func (c *testContext) theNumberOfNodesIs(expected int) error {
	ids := getNodeIds(c.svgOutput.doc)
	set := map[string]bool{}

	for _, v := range ids {
		id, _, _ := strings.Cut(v, "-")
		set[id] = true
	}

	assert.Len(c, set, 2)

	return c.err
}

// nolint: unused
func InitializeTestSuite(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {})
}

// nolint: unused
func InitializeScenario(ctx *godog.ScenarioContext) {
	tc := testContext{}
	ctx.BeforeScenario(func(*godog.Scenario) {})
	ctx.AfterScenario(func(s *godog.Scenario, err error) {
		if err != nil {
			fmt.Printf(
				"Command line output for \"%s\"\nUsing parameters: %s\n%s",
				s.GetName(),
				tc.cmdInput.parameters,
				tc.cmdResult.Output,
			)
		}
	})
	ctx.Step(`^the test fixures have been reset$`, tc.theTestFixuresHaveBeenReset)
	ctx.Step(`^the app runs without args$`, tc.theAppRunsWithoutArgs)
	ctx.Step(`^the app runs with parameters "(.*)"$`, tc.theAppRunsWithParameters)
	ctx.Step(`^the app exits without error$`, tc.theAppExitsWithoutError)
	ctx.Step(`^the app exits with an error$`, tc.theAppExitsWithAnError)
	ctx.Step(`^the app output contains "(.*)"$`, tc.theAppOutputContains)
	ctx.Step(`^the app output does not contain "(.*)"$`, tc.theAppOutputDoesNotContain)
	ctx.Step(`^a file "([^"]*)" exists$`, tc.aFileExists)
	ctx.Step(`^in the SVG file, all node text fits inside the node boundaries$`, tc.inTheSVGFileAllNodeTextFitsInsideTheNodeBoundaries)
	ctx.Step(`^the number of nodes is (\d+)$`, tc.theNumberOfNodesIs)
}
