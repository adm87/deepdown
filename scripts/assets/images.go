package assets

import (
	"bytes"

	"github.com/adm87/deepdown/scripts/deepdown"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type imageImporter struct {
	ctx deepdown.Context
}

func (ii *imageImporter) AssetTypes() []string {
	return []string{"png", "jpg", "jpeg"}
}

func (ii *imageImporter) Import(handle AssetHandle, data []byte) (any, error) {
	img, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return img, nil
}

func ImageImporter(ctx deepdown.Context) AssetImporter {
	return &imageImporter{ctx: ctx}
}
