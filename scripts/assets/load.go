package assets

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/adm87/utilities/hashset"
	"github.com/adm87/utilities/linq"
)

var (
	cache       = make(map[AssetHandle]any)      // cache stores loaded assets mapped by their handles.
	filesystems = make(map[string]fs.FS)         // filesystems maps root paths to their corresponding filesystems.
	loading     = make(hashset.Set[AssetHandle]) // loading tracks assets currently being loaded to prevent duplicate loads.
	mu          sync.RWMutex                     // mu protects access to the cache and filesystems maps.
)

// Get retrieves a loaded asset by its handle and asserts it to the specified type T.
// It returns the asset and a boolean indicating whether the asset was found and of the correct type.
func Get[T any](handle AssetHandle) (T, bool) {
	mu.RLock()
	defer mu.RUnlock()

	var zero T

	asset, exists := cache[handle]
	if !exists {
		return zero, false
	}

	typedAsset, ok := asset.(T)
	if !ok {
		return zero, false
	}

	return typedAsset, true
}

func MustGet[T any](handle AssetHandle) T {
	asset, ok := Get[T](handle)
	if !ok {
		panic(fmt.Sprintf("asset not found or wrong type: %s", handle))
	}
	return asset
}

// RegisterFilesystems registers the the default filesystems for asset loading.
func RegisterFilesystem(root string, fsys fs.FS) {
	if _, exists := filesystems[root]; exists {
		panic("duplicate filesystems")
	}
	filesystems[root] = fsys
}

// Load loads the assets corresponding to the provided handles.
// It processes the handles in batches and supports concurrent loading.
func Load(handles ...AssetHandle) error {
	if len(handles) == 0 {
		return nil
	}

	batches := linq.Batch(linq.Distinct(handles), 100)

	if len(batches) == 0 {
		return nil
	}

	return loadBatches(batches)
}

// MustLoad is like Load but panics if any error occurs.
func MustLoad(handles ...AssetHandle) {
	if err := Load(handles...); err != nil {
		panic(err)
	}
}

func loadBatches(batches [][]AssetHandle) error {
	if len(batches) == 1 {
		return loadBatch(batches[0])
	}

	var wg sync.WaitGroup
	errCh := make(chan error, len(batches))

	for _, batch := range batches {
		wg.Add(1)
		go func(b []AssetHandle) {
			defer wg.Done()

			defer func() {
				if r := recover(); r != nil {
					errCh <- fmt.Errorf("panic: %v", r)
				}
			}()

			if err := loadBatch(b); err != nil {
				errCh <- err
			}
		}(batch)
	}

	wg.Wait()
	close(errCh)

	var errs []error
	for err := range errCh {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func loadBatch(batch []AssetHandle) error {
	tryLoad := func(handle AssetHandle) error {
		path := handle.String()
		root := handle.Root()
		ext := handle.Ext()

		if !CanImport(ext) {
			return fmt.Errorf("no importer for asset type: %s", ext)
		}

		mu.Lock()
		if _, exists := cache[handle]; exists {
			mu.Unlock()
			return nil
		}
		if loading.Contains(handle) {
			mu.Unlock()
			return nil
		}
		loading.Add(handle)
		mu.Unlock()

		defer func() {
			mu.Lock()
			loading.Remove(handle)
			mu.Unlock()
		}()

		mu.RLock()
		fsys, exists := filesystems[root]
		mu.RUnlock()

		var data []byte
		var err error

		if exists {
			if _, ok := fsys.(embed.FS); !ok {
				path = strings.TrimPrefix(path, root)
				path = strings.TrimPrefix(path, string(filepath.Separator))
			}
			data, err = fs.ReadFile(fsys, path)
			if err != nil {
				return fmt.Errorf("failed to read asset: %w", err)
			}
		} else {
			data, err = os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read asset: %w", err)
			}
		}

		asset, err := importers[ext].Import(handle, data)
		if err != nil {
			return fmt.Errorf("failed to import asset: %w", err)
		}

		mu.Lock()
		cache[handle] = asset
		mu.Unlock()

		return nil
	}
	for _, handle := range batch {
		if err := tryLoad(handle); err != nil {
			return err
		}
	}
	return nil
}
