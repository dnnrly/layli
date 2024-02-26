//go:build acceptance

package test_test

import (
	"fmt"
	"testing"

	"github.com/cucumber/godog"
)

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format: "pretty",
			Paths:  []string{"features"},
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
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
	ctx.Step(`^a layli file "([^"]*)" exists$`, tc.aLayliFileExists)
	ctx.Step(`^in the SVG file, all node text fits inside the node boundaries$`, tc.inTheSVGFileAllNodeTextFitsInsideTheNodeBoundaries)
	ctx.Step(`^the number of nodes is (\d+)$`, tc.theNumberOfNodesIs)
	ctx.Step(`^the number of paths is (\d+)$`, tc.theNumberOfPathsIs)
	ctx.Step(`^no paths cross$`, tc.noPathsCross)
	ctx.Step(`^in the SVG file, nodes do not overlap$`, tc.inTheSVGFileNodesDoNotOverlap)
	ctx.Step(`^the image has a width less than (\d+)$`, tc.theImageHasAWidthLessThan)
	ctx.Step(`^the image has a height less than (\d+)$`, tc.theImageHasAHeightLessThan)
	ctx.Step(`^in the SVG file, all nodes fit on the image$`, tc.inTheSVGFileAllNodesFitOnTheImage)
	ctx.Step(`^in the SVG file, grid dots are not shown$`, tc.inTheSVGFileGridDotsAreNotShown)
	ctx.Step(`^in the SVG file, path grid dots are shown$`, tc.inTheSVGFilePathGridDotsAreShown)
	ctx.Step(`^in the SVG file, style class "([^"]*)" exists$`, tc.inTheSVGFileStyleClassExists)
	ctx.Step(`^in the SVG file, element "([^"]*)" has class "([^"]*)"$`, tc.inTheSVGFileElementHasClass)
	ctx.Step(`^in the SVG file, element "([^"]*)" has style "([^"]*)"$`, tc.inTheSVGFileElementHasStyle)
	ctx.Step(`^in the SVG file, element "([^"]*)" has attribute "([^"]*)" with value "([^"]*)"$`, tc.inTheSVGFileElementHasAttrWithVal)
	ctx.Step(`^the layli file contains the following nodes:$`, tc.theLayliFileContainsTheFollowingNodes)
}
