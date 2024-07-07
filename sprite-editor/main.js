import Alpine from 'https://cdn.jsdelivr.net/npm/alpinejs@3.14.0/dist/module.esm.min.js'

import './assets/style.css'

import './editor.js'
import './palette.js'
import './sheet.js'
import { Sprite } from './sprite.js'

export const SIZE = 12

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

const sprites = []

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

// Global store for sprite data
Alpine.store('sprites', {
  sprites,
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
})

for (let i = 0; i < 80; i++) {
  const s = new Sprite(i, SIZE)
  sprites.push(s)
}
