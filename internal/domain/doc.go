// Package domain contains pure business entities for layli.
//
// This package defines the core domain model with NO external dependencies
// (only Go stdlib). These types represent the "ubiquitous language" from
// our Gherkin feature files.
//
// Key principles:
//   - No dependencies on external packages (except Go stdlib)
//   - Value objects are immutable where possible
//   - Business logic lives here, not in adapters
//   - Types map directly to concepts in feature files
//
// Domain entities:
//   - Diagram: Complete diagram specification
//   - Node: A box/component in the diagram
//   - Edge: A connection between nodes
//   - Position: X/Y coordinates on the grid
//   - Path: A series of positions forming an edge route
package domain
