export default (size) => ({
  cellSize: 0,
  ctx: null,
  sprite: null,

  init() {
    console.log('Editor init')

    const canvas = this.$refs.canvas
    this.ctx = this.$refs.canvas.getContext('2d')
    this.cellSize = canvas.width / size
    this.sprite = this.$store.sprites.selected()

    this.drawSprite()
    this.drawGrid()
  },

  drawGrid() {
    // Draws a grid over the canvas
    this.ctx.strokeStyle = '#222'
    this.ctx.lineWidth = 2

    for (let x = 0; x <= size; x++) {
      this.ctx.beginPath()
      this.ctx.moveTo(x * this.cellSize, 0)
      this.ctx.lineTo(x * this.cellSize, this.$refs.canvas.height)
      this.ctx.stroke()
    }

    for (let y = 0; y <= size; y++) {
      this.ctx.beginPath()
      this.ctx.moveTo(0, y * this.cellSize)
      this.ctx.lineTo(this.$refs.canvas.width, y * this.cellSize)
      this.ctx.stroke()
    }
  },

  drawSprite() {
    // Weirdly this makes reactivity work
    this.sprite = this.$store.sprites.selected()

    this.ctx.clearRect(0, 0, this.$refs.canvas.width, this.$refs.canvas.height)

    this.ctx.fillStyle = this.$store.transparent ? 'rgba(0, 0, 0, 0)' : 'black'
    this.ctx.fillRect(0, 0, this.$refs.canvas.width, this.$refs.canvas.height)

    for (let y = 0; y < size; y++) {
      for (let x = 0; x < size; x++) {
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
    this.ctx.fillStyle = this.$store.transparent ? 'rgba(0, 0, 0, 0)' : 'black'
    this.ctx.fillRect(x * this.cellSize, y * this.cellSize, this.cellSize, this.cellSize)
    this.$store.sprites.selected().data[y][x] = null
    this.drawGrid()
  },

  colourChange(event) {
    this.colour = event.target.value
  },

  toolClear() {
    this.sprite.data = this.sprite.data.map(() => Array(size).fill(null))
  },

  toolFlipX() {
    this.sprite.data = this.sprite.data.map((row) => row.reverse())
  },

  toolFlipY() {
    this.sprite.data = this.sprite.data.reverse()
  },

  toolColour() {
    this.sprite.data = this.sprite.data.map((row) => row.map((cell) => (cell === null ? null : this.$store.pal.selected())))
  },

  toolMoveDown() {
    this.sprite.data.unshift(this.sprite.data.pop())
  },

  toolMoveUp() {
    this.sprite.data.push(this.sprite.data.shift())
  },

  toolMoveRight() {
    this.sprite.data = this.sprite.data.map((row) => {
      row.unshift(row.pop())
      return row
    })
  },

  toolMoveLeft() {
    this.sprite.data = this.sprite.data.map((row) => {
      row.push(row.shift())
      return row
    })
  },
})
