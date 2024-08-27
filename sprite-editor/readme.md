# Sprite Editor

This is very simple sprite editor which is run as a standalone web app.

Is designed for retro games so is based on the approach where all sprites share a common and editable palette (i.e. indexed) which was common with older systems. It's also optimized for single colour sprites

Features:

- Configurable sprite pixel dimensions, pallet colour count, and "bank" size
- Scratch area to preview how sprites would look together on a tiled grid
- Save projects as JSON files
- Sprites can be exported for use in a game, to a single image PNG sprite sheet along with some JSON metadata

## Running

- From root of repo, run `make editor`
- Go to http://localhost:8000
