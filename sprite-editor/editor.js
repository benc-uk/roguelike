import Alpine from 'https://cdn.jsdelivr.net/npm/alpinejs@3.14.0/dist/module.esm.min.js'

import { SIZE } from './main.js'

Alpine.data('editor', () => ({
  cellSize: 0,
  ctx: null,
  sprite: null,
  pallet: null,

  init() {
    const canvas = this.$refs.canvas
    this.ctx = this.$refs.canvas.getContext('2d')
    this.cellSize = canvas.width / SIZE
    this.sprite = this.$store.sprites.selected()
    this.pallet = this.$store.pal

    this.drawSprite()
    this.drawGrid()
  },

  drawGrid() {
    // Draws a grid over the canvas
    this.ctx.strokeStyle = '#222'
    this.ctx.lineWidth = 2

    for (let x = 0; x <= SIZE; x++) {
      this.ctx.beginPath()
      this.ctx.moveTo(x * this.cellSize, 0)
      this.ctx.lineTo(x * this.cellSize, this.$refs.canvas.height)
      this.ctx.stroke()
    }

    for (let y = 0; y <= SIZE; y++) {
      this.ctx.beginPath()
      this.ctx.moveTo(0, y * this.cellSize)
      this.ctx.lineTo(this.$refs.canvas.width, y * this.cellSize)
      this.ctx.stroke()
    }
  },

  drawSprite() {
    // Weirdly this makes reactivity work
    this.sprite = this.$store.sprites.selected()

    this.ctx.fillStyle = '#000'
    this.ctx.fillRect(0, 0, this.$refs.canvas.width, this.$refs.canvas.height)

    for (let y = 0; y < SIZE; y++) {
      for (let x = 0; x < SIZE; x++) {
        if (this.sprite.data[y][x] !== null) {
          this.ctx.fillStyle = this.$store.pal.colours[this.sprite.data[y][x]]
          this.ctx.fillRect(x * this.cellSize, y * this.cellSize, this.cellSize, this.cellSize)
        }
      }
    }
  },

  handleClick(event) {
    const x = Math.floor(event.offsetX / this.cellSize)
    const y = Math.floor(event.offsetY / this.cellSize)

    if (event.buttons === 1) {
      this.fillCell(x, y)
    } else if (event.buttons === 2) {
      this.clearCell(x, y)
    }

    if (event.type === 'click') {
      this.fillCell(x, y)
    }

    if (event.type === 'contextmenu') {
      this.clearCell(x, y)
    }
  },

  fillCell(x, y) {
    this.ctx.fillStyle = this.$store.pal.colour()
    this.ctx.fillRect(x * this.cellSize, y * this.cellSize, this.cellSize, this.cellSize)
    this.$store.sprites.selected().data[y][x] = this.$store.pal.selected()
    this.drawGrid()
  },

  clearCell(x, y) {
    this.ctx.clearRect(x * this.cellSize, y * this.cellSize, this.cellSize, this.cellSize)
    this.$store.sprites.selected().data[y][x] = null
    this.drawGrid()
  },

  colourChange(event) {
    this.colour = event.target.value
  },
}))
