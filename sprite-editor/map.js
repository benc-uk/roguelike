const blackPixelImg =
  'data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAIAAAACCAYAAABytg0kAAAAAXNSR0IArs4c6QAAAAlwSFlzAAAWJQAAFiUBSVIk8AAAABNJREFUCB1jZGBg+A/EDEwgAgQADigBA//q6GsAAAAASUVORK5CYII%3D'

export default () => ({
  map: [],
  mapData: [],
  width: 12,
  height: 6,

  init() {
    console.log('Scratch init')

    this.$refs.scratchGrid.style.gridTemplateColumns = `repeat(${this.width}, fit-content(100%))`
    let total = this.width * this.height
    this.mapData = this.$store.map

    // Create a new map from the scratch data
    this.map = []
    for (let i = 0; i < total; i++) {
      let image = blackPixelImg
      if (this.mapData[i] !== -1) {
        image = this.$store.sprites.get(this.mapData[i]).toImageSrc(this.$store.pal.colours)
      }

      this.map.push({
        index: this.mapData[i] || -1,
        image,
      })
    }
  },

  updateSprite() {
    const spriteIndex = this.$store.sprites.selectedIndex()
    const imgSrc = this.$store.sprites.get(spriteIndex).toImageSrc(this.$store.pal.colours)

    // Update the map
    for (let i = 0; i < this.map.length; i++) {
      if (this.map[i].index === spriteIndex) {
        this.map[i].image = imgSrc
      }
    }
  },

  clickCell(event, index) {
    if (event.type === 'contextmenu') {
      this.clearCell(index)
      return
    }

    this.map[index].index = this.$store.sprites.selectedIndex()
    this.map[index].image = this.$store.sprites.selected().toImageSrc(this.$store.pal.colours)
  },

  clearCell(index) {
    this.map[index].index = -1
    this.map[index].image = blackPixelImg
  },

  updateStore() {
    // scratch data is stored in the store
    const mapData = this.map.map((cell) => cell.index)
    this.$store.map = mapData
  },
})
