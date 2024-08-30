const transparentPixelImg =
  'data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNkYAAAAAYAAjCB0C8AAAAASUVORK5CYII='

export const WIDTH = 14
export const HEIGHT = 8

export default () => ({
  map: [],
  mapData: [],

  init() {
    console.log('Map init')

    // This is important hack to set the grid columns to the correct width
    this.$refs.mapGrid.style.gridTemplateColumns = `repeat(${WIDTH}, fit-content(100%))`
    this.mapData = this.$store.map

    this.$refs.mapBg.style.backgroundColor = this.$store.transparent ? 'rgba(0, 0, 0, 0)' : 'black'

    // Create a new map with images, from the raw mapdata
    this.map = []
    for (let i = 0; i < this.mapData.length; i++) {
      let image = transparentPixelImg
      if (this.mapData[i] !== -1) {
        image = this.$store.sprites.get(this.mapData[i]).toImageSrc(this.$store.pal.colours)
      }

      this.map.push({
        index: this.mapData[i],
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
    this.map[index].image = transparentPixelImg
  },

  updateStore() {
    const mapData = this.map.map((cell) => cell.index)
    this.$store.map = mapData
  },

  clearMap() {
    for (let i = 0; i < this.map.length; i++) {
      this.map[i].index = -1
      this.map[i].image = transparentPixelImg
      this.$store.map[i] = -1
    }
  },
})
