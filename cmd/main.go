package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/spf13/cobra"

	"github.com/adm87/deepdown/data"
	"github.com/adm87/deepdown/scripts/assets"
	"github.com/adm87/deepdown/scripts/deepdown"
	"github.com/adm87/deepdown/scripts/game"

	_ "net/http/pprof"

	assetcmd "github.com/adm87/deepdown/cmd/assets"
)

// TASK: Setup build tags

func main() {
	var (
		root    string
		profile bool
	)

	ctx := deepdown.NewContext()

	cmd := &cobra.Command{
		Use:   "deepdown",
		Short: "Deepdown Game Engine",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			root, err := filepath.Abs(root)
			if err != nil {
				ctx.Logger().Error("error", slog.Any("err", err))
				os.Exit(1)
			}

			ctx.Set(deepdown.CtxApplicationRoot, root)

			assets.RegisterFilesystem("assets", os.DirFS(path.Join(root, "data", "assets")))
			assets.RegisterFilesystem("embedded", data.EmbeddedFS)

			assets.RegisterImporters(ctx)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx.Logger().Info("Starting Deepdown...")

			if profile {
				ctx.Logger().Info("Profiling enabled")

				go func() {
					log.Println("Profiling server at http://localhost:6060/debug/pprof/")
					log.Println(http.ListenAndServe("localhost:6060", nil))
				}()
			}

			return ebiten.RunGame(game.NewGame(ctx))
		},
	}

	cmd.PersistentFlags().StringVar(&root, "root", ".", "Root directory of the application")
	cmd.Flags().BoolVar(&profile, "profile", false, "Enable profiling")

	cmd.AddCommand(assetcmd.GenerateHandles(ctx))

	if err := cmd.ExecuteContext(ctx.Ctx()); err != nil {
		ctx.Logger().Error("Command execution failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
