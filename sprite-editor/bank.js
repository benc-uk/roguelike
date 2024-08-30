export default () => ({
  spriteImages: [],
  activeSprite: null,
  transparent: false,

  init() {
    console.log('Bank init')

    this.activeSprite = this.$store.sprites.selected()
    for (const s of this.$store.sprites.sprites) {
      this.spriteImages.push(s.toImageSrc(this.$store.pal.colours))
    }

    this.transparent = this.$store.transparent
    this.$refs.bankBg.style.backgroundColor = this.transparent ? 'rgba(0, 0, 0, 0)' : 'black'

    window.addEventListener('palette', this.paletteUpdated.bind(this))
  },

  paletteUpdated() {
    let index = 0
    for (const sprite of this.$store.sprites.sprites) {
      this.spriteImages[index] = sprite.toImageSrc(this.$store.pal.colours)
      index++
    }
  },

  updateSprite() {
    this.spriteImages[this.$store.sprites.selectedIndex()] = this.$store.sprites.selected().toImageSrc(this.$store.pal.colours)
  },
})
