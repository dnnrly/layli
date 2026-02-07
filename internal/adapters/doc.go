// Package adapters contains implementations of port interfaces.
//
// Adapters are the concrete implementations that interact with
// external systems or provide specific algorithms. They implement
// the port interfaces defined in the usecases package.
//
// Structure:
//   - config/     : Configuration file parsers (YAML, JSON, etc.)
//   - layout/     : Layout algorithms (FlowSquare, TopoSort, etc.)
//   - pathfinding/: Pathfinding algorithms (Dijkstra, A*, etc.)
//   - rendering/  : Output renderers (SVG, PNG, etc.)
//
// Each adapter is independent and can be tested in isolation.
package adapters
