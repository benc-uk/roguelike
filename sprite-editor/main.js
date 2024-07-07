import Alpine from 'https://cdn.jsdelivr.net/npm/alpinejs@3.14.0/dist/module.esm.min.js'

import './assets/style.css'

import { Sprite } from './lib/sprite.js'

// Sub components of the sprite editor
import Editor from './editor.js'
import Palette from './palette.js'
import Bank from './bank.js'
import Scratch from './scratch.js'

// The colours available to the sprite editor
const colours = []
colours.push('#fefefe')
colours.push('#e5d5d0')
colours.push('#808080')
colours.push('#404040')
colours.push('#895826')
colours.push('#ff8000')
colours.push('#da0000')
colours.push('#774009')

colours.push('#d3bc87')
colours.push('#ffe228')
colours.push('#6dff66')
colours.push('#009900')
colours.push('#192ee8')
colours.push('#44d0ff')
colours.push('#6666cc')
colours.push('#a82cea')

// Global store for palette data
Alpine.store('pal', {
  colours,
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

const SIZE = 12
const BANK_SIZE = 32

// Bank of sprites
const spriteBank = []
for (let i = 0; i < BANK_SIZE; i++) {
  spriteBank.push(new Sprite(i, SIZE))
}

// Global store for sprite data
Alpine.store('sprites', {
  sprites: spriteBank,
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
  size: SIZE,
}))
