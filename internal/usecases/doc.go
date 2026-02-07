// Package usecases contains the application use cases (orchestration logic).
//
// This layer defines the business workflows that coordinate domain entities
// with external concerns (file I/O, rendering, etc.). Use cases do NOT
// contain business logic—that belongs in the domain layer.
//
// Each use case:
//   - Defines one complete workflow (e.g., "Generate a diagram")
//   - Uses ports (interfaces) to depend on external concerns
//   - Never depends directly on specific implementations
//   - Maps to Gherkin scenarios: "Given → When → Then"
//
// Ports (interfaces) defined here:
//   - ConfigParser: Read and parse configuration files
//   - LayoutEngine: Arrange nodes using layout algorithms
//   - Pathfinder: Calculate paths between nodes
//   - Renderer: Generate output (SVG, PNG, etc.)
//   - FileReader/FileWriter: Abstract file system access
//
// Adapters (in internal/adapters/) implement ports
//
// Use cases map to Gherkin "When" steps.
package usecases
