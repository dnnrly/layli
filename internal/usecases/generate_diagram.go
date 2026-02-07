package usecases

import "fmt"

// GenerateDiagram orchestrates the complete diagram generation workflow.
// Maps to a complete Gherkin scenario: Given → When → Then
type GenerateDiagram struct {
	configParser ConfigParser
	layoutEngine LayoutEngine
	pathfinder   Pathfinder
	renderer     Renderer
}

// NewGenerateDiagram creates a new GenerateDiagram use case.
func NewGenerateDiagram(
	parser ConfigParser,
	layout LayoutEngine,
	pathfinder Pathfinder,
	renderer Renderer,
) *GenerateDiagram {
	return &GenerateDiagram{
		configParser: parser,
		layoutEngine: layout,
		pathfinder:   pathfinder,
		renderer:     renderer,
	}
}

// Execute runs the complete diagram generation pipeline.
//
// Steps:
//
//	1. Parse configuration (Given)
//	2. Validate diagram (Given)
//	3. Arrange layout (When)
//	4. Calculate paths (When)
//	5. Render output (Then)
func (uc *GenerateDiagram) Execute(configPath, outputPath string) error {
	// Parse configuration
	diagram, err := uc.configParser.Parse(configPath)
	if err != nil {
		return fmt.Errorf("parse config: %w", err)
	}

	// Validate diagram
	if err := diagram.Validate(); err != nil {
		return fmt.Errorf("validate diagram: %w", err)
	}

	// Arrange layout
	if err := uc.layoutEngine.Arrange(diagram); err != nil {
		return fmt.Errorf("arrange layout: %w", err)
	}

	// Calculate paths
	if err := uc.pathfinder.FindPaths(diagram); err != nil {
		return fmt.Errorf("find paths: %w", err)
	}

	// Render output
	if err := uc.renderer.Render(diagram, outputPath); err != nil {
		return fmt.Errorf("render diagram: %w", err)
	}

	return nil
}
