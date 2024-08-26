import Alpine from 'https://cdn.jsdelivr.net/npm/alpinejs@3.14.0/dist/module.esm.min.js'

import './assets/style.css'

import { Sprite } from './lib/sprite.js'
import { randomHexColor } from './lib/colours.js'

// Sub components of the sprite editor
import Editor from './editor.js'
import Palette from './palette.js'
import Bank from './bank.js'
import Scratch from './scratch.js'

// Global store for palette data
Alpine.store('pal', {
  colours: [],
  _selected: 0,

  selected() {
    return this._selected
  },

  colour() {
    return this.colours[this._selected]
  },

  select(index) {
    this._selected = index
  },
})

// Global store for sprite data
Alpine.store('sprites', {
  sprites: [],
  _selected: 0,

  selected() {
    return this.sprites[this._selected]
  },

  selectedIndex() {
    return this._selected
  },

  select(index) {
    this._selected = index
  },

  get(index) {
    return this.sprites[index]
  },
})

Alpine.data('palette', Palette)
Alpine.data('bank', Bank)
Alpine.data('palette', Palette)
Alpine.data('editor', Editor)
Alpine.data('scratch', Scratch)

// Main Alpine app with global data
Alpine.data('app', () => ({
  size: 16,
  projectLoaded: false,
  newSpriteSize: 12,
  newPaletteSize: 16,
  newBankSize: 128,

  init() {
    console.log('App initializing')

    try {
      this.loadFromStorage()
    } catch (e) {
      console.error('Error loading project', e)
    }
  },

  async loadFromStorage() {
    const projectData = await localStorage.getItem('project')
    if (projectData) {
      console.log('Loading stored project')
      try {
        const proj = await JSON.parse(projectData)

        const loadSprites = []
        for (const s of proj.sprites) {
          const sprite = new Sprite(s.id, s.size)
          sprite.loadData(s.data)
          loadSprites.push(sprite)
        }

        const loadPal = []
        for (const c of proj.palette) {
          loadPal.push(c)
        }

        Alpine.store('pal').colours = loadPal
        Alpine.store('sprites').sprites = loadSprites
        this.size = proj.size

        this.projectLoaded = true
      } catch (e) {
        console.error('Error loading & parsing project', e)
      }
    }
  },

  saveToStorage() {
    console.log('Saving project to storage')
    localStorage.setItem(
      'project',
      JSON.stringify({
        size: this.size,
        sprites: Alpine.store('sprites').sprites,
        palette: Alpine.store('pal').colours,
      })
    )
  },

  newProject() {
    try {
      const spriteBank = []
      for (let i = 0; i < this.newBankSize; i++) {
        spriteBank.push(new Sprite(i, this.newSpriteSize))
      }
      Alpine.store('sprites').sprites = spriteBank

      const palette = []
      for (let i = 0; i < this.newPaletteSize; i++) {
        palette.push(randomHexColor())
      }
      Alpine.store('pal').colours = palette

      this.size = this.newSpriteSize

      this.saveToStorage()
      this.projectLoaded = true
      console.log('New project created and stored')
    } catch (e) {
      console.error('Error creating new project', e)
    }
  },

  eraseProject() {
    console.log('Erasing project')
    const resp = prompt('Are you sure you want to erase the project? Enter "yes" to confirm')
    if (resp === 'yes') {
      localStorage.removeItem('project')
      this.projectLoaded = false
    }
  },
}))
