package assets

import (
	"strings"

	"github.com/adm87/deepdown/scripts/deepdown"
)

// AssetHandle represents a unique identifier for an asset.
type AssetHandle string

// String returns a string representation of the asset handle (for debugging).
func (ah AssetHandle) String() string {
	return string(ah)
}

// Ext returns the file extension of the asset handle (without the dot).
func (ah AssetHandle) Ext() string {
	s := string(ah)
	if i := strings.LastIndexByte(s, '.'); i >= 0 && i < len(s)-1 {
		return s[i+1:]
	}
	return ""
}

// Root returns the root directory of the asset handle.
func (ah AssetHandle) Root() string {
	s := string(ah)
	if i := strings.IndexByte(s, '/'); i >= 0 {
		return s[:i]
	}
	return s
}

// AssetImporter is an interface for importing different types of assets.
type AssetImporter interface {
	AssetTypes() []string                                // AssetTypes returns the list of asset types the importer can handle.
	Import(handle AssetHandle, data []byte) (any, error) // Import imports the asset data for the given handle.
}

var (
	importers = make(map[string]AssetImporter) // importers maps asset types to their respective importers.
)

// RegisterImporter registers all supported asset importers.
func RegisterImporters(ctx deepdown.Context) {
	addImporter(ImageImporter(ctx))
	addImporter(TmxImporter(ctx))
	addImporter(TsxImporter(ctx))
	addImporter(TxImporter(ctx))
}

// CanImport checks if there is an importer registered for the given file extension.
func CanImport(ext string) bool {
	_, ok := importers[ext]
	return ok
}

func addImporter(ai AssetImporter) {
	for _, ext := range ai.AssetTypes() {
		importers[ext] = ai
	}
}
