export class Sprite {
  constructor(index = 0, size = 16) {
    this.name = `Sprite ${index}`
    this.index = index
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
}
