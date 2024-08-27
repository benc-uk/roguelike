import { Sprite } from './lib/sprite.js'
import { randomHexColor } from './lib/colours.js'

export default () => ({
  size: 16,
  projectLoaded: false,
  newSpriteSize: 12,
  newPaletteSize: 16,
  newBankSize: 128,
  copySprite: null,

  init() {
    console.log('App initializing')

    try {
      this.loadFromStorage()
    } catch (e) {
      console.error('Error loading project', e)
    }
  },

  async loadFromStorage() {
    console.log('Loading project from browser storage')

    const projectData = await localStorage.getItem('project')
    if (projectData) {
      console.log('Loading stored project')
      try {
        const proj = await JSON.parse(projectData)

        const loadSprites = []
        for (const s of proj.sprites) {
          const sprite = new Sprite(s.name, s.size)
          sprite.loadData(s.data)
          loadSprites.push(sprite)
        }

        const loadPal = []
        for (const c of proj.palette) {
          loadPal.push(c)
        }

        this.$store.pal.colours = loadPal
        this.$store.sprites.sprites = loadSprites
        this.size = proj.size
        this.$store.sprites.select(proj.selected)
        this.$store.map = proj.map

        this.projectLoaded = true
      } catch (e) {
        console.error('Error loading & parsing project', e)
      }
    }
  },

  saveToStorage() {
    localStorage.setItem(
      'project',
      JSON.stringify({
        size: this.size,
        sprites: this.$store.sprites.sprites,
        palette: this.$store.pal.colours,
        selected: this.$store.sprites.selectedIndex(),
        map: this.$store.map,
      })
    )
  },

  newProject() {
    console.log('Creating new project')

    try {
      const spriteBank = []
      for (let i = 0; i < this.newBankSize; i++) {
        const name = `Sprite ${i}`
        spriteBank.push(new Sprite(name, this.newSpriteSize))
      }
      this.$store.sprites.sprites = spriteBank

      const palette = []
      for (let i = 0; i < this.newPaletteSize; i++) {
        palette.push(randomHexColor())
      }
      this.$store.pal.colours = palette

      // Create a blank map 12*6
      // TODO: Fix hard coded map size
      this.$store.map = []
      for (let i = 0; i < 12 * 6; i++) {
        this.$store.map.push(-1)
      }

      this.size = this.newSpriteSize

      this.saveToStorage()
      this.projectLoaded = true
      console.log('New project created and stored')
    } catch (e) {
      console.error('Error creating new project', e)
    }
  },

  eraseProject() {
    const resp = prompt('Are you sure you want to erase the project? Enter "yes" to confirm')
    if (resp === 'yes') {
      console.log('Erasing & resetting project')

      localStorage.removeItem('project')
      this.projectLoaded = false
    }
  },

  exportProject() {
    console.log('Exporting project to file')

    const data = localStorage.getItem('project')
    const blob = new Blob([data], { type: 'application/json' })
    const url = URL.createObjectURL(blob)

    const a = document.createElement('a')
    a.href = url
    a.download = 'project.json'

    a.click()
  },

  importProject() {
    console.log('Importing project from file')

    const input = document.createElement('input')
    input.type = 'file'
    input.accept = 'application/json'

    input.onchange = async (e) => {
      const file = e.target.files[0]
      const reader = new FileReader()
      reader.onload = async (e) => {
        const data = e.target.result
        localStorage.setItem('project', data)
        this.loadFromStorage()
        location.reload()
      }
      reader.readAsText(file)
    }

    input.click()
  },

  exportAllSprites() {
    console.log('Exporting sprite sheet')

    // Create a big canvas to draw all sprites on
    const canvas = document.createElement('canvas')
    const spriteCount = this.$store.sprites.sprites.length
    const cols = 8
    const rows = Math.ceil(spriteCount / cols)
    canvas.width = this.size * cols
    canvas.height = this.size * rows

    console.log('Exporting sprites', spriteCount, 'in', cols, 'x', rows, 'grid')

    // Prepare spritesheet meta data
    const meta = {
      size: this.size,
      count: spriteCount,
      sprites: [],
    }

    // Process all sprites into canvas and meta data
    const ctx = canvas.getContext('2d')
    let i = 0
    for (let y = 0; y < rows; y++) {
      for (let x = 0; x < cols; x++) {
        if (i < spriteCount) {
          const sprite = this.$store.sprites.sprites[i]
          sprite.drawOnCanvas(ctx, x * this.size, y * this.size, this.$store.pal.colours)

          // Add sprite meta data
          meta.sprites.push({
            name: sprite.name,
            x: x * this.size,
            y: y * this.size,
          })

          i++
        }
      }
    }

    // Download image
    const imageA = document.createElement('a')
    imageA.href = canvas.toDataURL()
    imageA.download = 'sprites.png'
    imageA.click()

    // Download meta data
    const metaBlob = new Blob([JSON.stringify(meta)], { type: 'application/json' })
    const metaA = document.createElement('a')
    metaA.href = URL.createObjectURL(metaBlob)
    metaA.download = 'sprites.json'
    metaA.click()
  },

  toolCopy() {
    if (this.$store.sprites.selected()) {
      this.copySprite = this.$store.sprites.selected()
      console.log('Copied sprite', this.copySprite)
    }
  },

  toolPaste() {
    if (this.copySprite) {
      // Low level clone sprite data
      const data = []
      for (let y = 0; y < this.copySprite.size; y++) {
        data[y] = []
        for (let x = 0; x < this.copySprite.size; x++) {
          data[y][x] = this.copySprite.data[y][x]
        }
      }
      this.$store.sprites.selected().loadData(data)
      this.$store.sprites.selected().name = this.copySprite.name

      console.log('Pasted sprite', this.copySprite)
    }
  },
})
