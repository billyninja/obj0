package tmx

import (
	"encoding/xml"
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Tile struct {
	Gid int32 `xml:"gid,attr"`
}

type Layer struct {
	Name   string `xml:"name,attr"`
	Tiles  []Tile `xml:"data>tile"`
	Height int32  `xml:"height,attr"`
	Width  int32  `xml:"width,attr"`
}

type TMX struct {
	Layers   []Layer   `xml:"layer"`
	Tilesets []TileSet `xml:"tileset"`

	XMLName     xml.Name `xml:"map"`
	HeightTiles int32    `xml:"height,attr"`
	WidthTiles  int32    `xml:"width,attr"`
	TileH       int32    `xml:"tileheight,attr"`
	TileW       int32    `xml:"tilewidth,attr"`
}

type TSImg struct {
	Src       string `xml:"source,attr"`
	SrcHeight int32  `xml:"height,attr"`
	SrcWidth  int32  `xml:"width,attr"`
}

type Property struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type TerrainType struct {
	Name       string     `xml:"name,attr"`
	Properties []Property `xml:"properties>property"`
}

type Terrain struct {
	Source       *sdl.Rect
	ColMode      int32
	TerrainTypes [4]*TerrainType
}

type SpecTerrain struct {
	Id         int32  `xml:"id,attr"`
	TerrainStr string `xml:"terrain,attr"`
	// Build manually
	TerrainRefs [4]int32
}

type Space struct {
	Source  *sdl.Rect
	ColMode int32
	Terrain *Terrain
}

type TileSet struct {
	Name  string `xml:"name,attr"`
	Image TSImg  `xml:"image"`
	TileH int32  `xml:"tileheight,attr"`
	TileW int32  `xml:"tilewidth,attr"`

	TerrainTypes []*TerrainType `xml:"terraintypes>terrain"`
	SpecTerrain  []*SpecTerrain `xml:"tile"`
	Gids         []*Terrain

	Txtr *sdl.Texture
}

func (ts *TileSet) GetGIDRect(gid int32) *sdl.Rect {
	w := ts.Image.SrcWidth / ts.TileW

	var x int32 = ((gid - 1) % w) * ts.TileW
	var y int32 = (gid / w) * ts.TileH

	return &sdl.Rect{x, y, int32(ts.TileW), int32(ts.TileH)}
}

func (ts *TileSet) Load(renderer *sdl.Renderer) {

	tilesetImg, err := img.Load("data/assets/" + ts.Image.Src)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load PNG: %s\n", err)
		os.Exit(3)
	}
	defer tilesetImg.Free()

	ts.Txtr, err = renderer.CreateTextureFromSurface(tilesetImg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create texture: %s\n", err)
		os.Exit(4)
	}
	var HSteps int32 = ts.Image.SrcHeight / ts.TileH
	var WSteps int32 = ts.Image.SrcWidth / ts.TileW
	ts.Gids = make([]*Terrain, (HSteps+1)*(WSteps+1))

	spMap := make(map[int32]*SpecTerrain)
	for _, spTerrain := range ts.SpecTerrain {
		for ti, terr := range strings.Split(spTerrain.TerrainStr, ",") {
			terrainTypeId, _ := strconv.Atoi(terr)
			spTerrain.TerrainRefs[ti] = int32(terrainTypeId)
			spMap[spTerrain.Id] = spTerrain
		}
	}

	var currentGid int32 = 1
	var h int32 = 0
	for ; h < HSteps; h++ {
		var w int32 = 0
		for ; w < WSteps; w++ {

			terrainSourceR := &sdl.Rect{w * ts.TileW, h * ts.TileH, ts.TileW, ts.TileH}
			if terrainSourceR.X > ts.Image.SrcWidth || terrainSourceR.Y > ts.Image.SrcHeight {
				os.Exit(9)
			}

			terr := &Terrain{
				Source:  terrainSourceR,
				ColMode: 0,
			}
			if val, ok := spMap[currentGid]; ok {
				for i, v := range val.TerrainRefs {
					terr.TerrainTypes[i] = ts.TerrainTypes[v]
					if terr.TerrainTypes[i].Name == "COLL_BLOCK" {
						terr.ColMode = 1
					}
				}
			}
			ts.Gids[currentGid] = terr

			currentGid++
		}
	}
}

func LoadTMXFile(mapname string, renderer *sdl.Renderer) (*TMX, [][]*Terrain) {
	f, err := os.Open(mapname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening the tmx map: %s\n", err)
		os.Exit(2)
	}
	output, _ := ioutil.ReadAll(f)
	_ = f.Close()

	tmx := &TMX{}
	err = xml.Unmarshal(output, tmx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing the tmx map: %s\n", err)
		os.Exit(3)
	}

	for i := 0; i < len(tmx.Tilesets); i++ {
		tmx.Tilesets[i].Load(renderer)
	}

	world := make([][]*Terrain, tmx.HeightTiles)

	ts := tmx.Tilesets[0]
	layer0 := tmx.Layers[0]

	var i int32 = 0
	for ; i < layer0.Height; i++ {
		world[i] = make([]*Terrain, tmx.WidthTiles)

		var j int32 = 0
		for ; j < layer0.Width; j++ {
			idx := (i * layer0.Width) + j
			tile := layer0.Tiles[idx]
			world[i][j] = ts.Gids[tile.Gid]
			//println("TRef", tile.TerrainRefStr)
		}
	}

	return tmx, world
}
