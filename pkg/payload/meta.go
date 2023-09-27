package payload

type Dependency uint8

const (
	PackageName Dependency = iota
	SharedLibrary
	PkgConfig
	Interpreter
	CMake
	Python
	Binary
	SystemBinary
	PkgConfig32
)

type Tag uint16

const (
	// Name of the package
	Name Tag = 1
	// Architecture of the package
	Architecture = 2
	// Version of the package
	Version = 3
	// Summary of the package
	Summary = 4
	// Description of the package
	Description = 5
	// Homepage for the package
	Homepage = 6
	// ID for the source package, used for grouping
	SourceID = 7
	// Runtime dependencies
	Depends = 8
	// Provides some capability or name
	Provides = 9
	// Conflicts with some capability or name
	Conflicts = 10
	// Release number for the package
	Release = 11
	// SPDX license identifier
	License = 12
	// Currently recorded build number
	BuildRelease = 13
	// Repository index specific (relative URI)
	PackageURI = 14
	// Repository index specific (Package hash)
	PackageHash = 15
	// Repository index specific (size on disk)
	PackageSize = 16
	// A Build Dependency
	BuildDepends = 17
	// Upstream URI for the source
	SourceURI = 18
	// Relative path for the source within the upstream URI
	SourcePath = 19
	// Ref/commit of the upstream source
	SourceRef = 20
)
