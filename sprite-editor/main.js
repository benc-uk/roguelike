import Alpine from 'https://cdn.jsdelivr.net/npm/alpinejs@3.14.0/dist/module.esm.min.js'

import './assets/style.css'

// Sub components of the sprite editor
import App from './app.js'
import Editor from './editor.js'
import Palette from './palette.js'
import Bank from './bank.js'
import Map from './map.js'

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

Alpine.store('scratch', [])

Alpine.data('app', App)
Alpine.data('palette', Palette)
Alpine.data('bank', Bank)
Alpine.data('palette', Palette)
Alpine.data('editor', Editor)
Alpine.data('map', Map)
