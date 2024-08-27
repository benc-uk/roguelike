export class Sprite {
  constructor(name, size = 16) {
    this.name = name ? name : `Sprite ${Math.random().toString(36).substring(7)}`
    this.data = []
    this.size = size

    for (let y = 0; y < size; y++) {
      this.data[y] = []
      for (let x = 0; x < size; x++) {
        // random data
        this.data[y][x] = null //Math.floor(Math.random() * 16)
      }
    }
  }

  toImageSrc(palette = []) {
    const canvas = document.createElement('canvas')
    canvas.width = this.size
    canvas.height = this.size
    const ctx = canvas.getContext('2d')

    for (let y = 0; y < this.size; y++) {
      for (let x = 0; x < this.size; x++) {
        const pixel = palette[this.data[y][x]]
        if (pixel) {
          ctx.fillStyle = pixel
          ctx.fillRect(x, y, 1, 1)
        }
      }
    }

    // return image data URL
    return canvas.toDataURL()
  }

  loadData(data) {
    if (data.length !== this.size) {
      throw new Error('Data is not the correct size')
    }

    this.data = data
  }

  drawOnCanvas(ctx, x, y, palette = []) {
    for (let sy = 0; sy < this.size; sy++) {
      for (let sx = 0; sx < this.size; sx++) {
        const pixel = palette[this.data[sy][sx]]
        if (pixel) {
          ctx.fillStyle = pixel
          ctx.fillRect(x + sx, y + sy, 1, 1)
        }
      }
    }
  }
}
