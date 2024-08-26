const blackPixelImg =
  'data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAIAAAACCAYAAABytg0kAAAAAXNSR0IArs4c6QAAAAlwSFlzAAAWJQAAFiUBSVIk8AAAABNJREFUCB1jZGBg+A/EDEwgAgQADigBA//q6GsAAAAASUVORK5CYII%3D'

export default () => ({
  map: null,
  width: 12,
  height: 6,

  init() {
    console.log('Scratch init')

    this.$refs.scratchGrid.style.gridTemplateColumns = `repeat(${this.width}, fit-content(100%))`

    let total = this.width * this.height

    // Create a new map
    this.map = []
    for (let i = 0; i < total; i++) {
      this.map.push({
        index: -1,
        image: blackPixelImg,
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
      console.log('right click')
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
})
