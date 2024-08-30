import { isHexColor } from './lib/colours'

export default () => ({
  bulkColours: '',
  bulkEdit: false,

  init() {
    console.log('Palette init')

    for (const c of this.$store.pal.colours) {
      // Trim the # from the colour
      this.bulkColours += c.substring(1) + '\n'
    }

    this.bulkEdit = false
  },

  colourUpdated(colour, i) {
    this.$store.pal.colours[i] = colour
    this.$dispatch('palette', {})
  },

  bulkSave() {
    this.bulkEdit = false

    const colourStrings = this.bulkColours.split('\n')
    for (let i = 0; i < this.$store.pal.colours.length; i++) {
      if (!colourStrings[i]) {
        continue
      }

      let colour = colourStrings[i].trim()
      if (colour[0] !== '#') {
        colour = '#' + colour
      }

      if (!isHexColor(colour)) {
        colour = '#000000'
      }

      this.$store.pal.colours[i] = colour
    }

    // Hideous hack to restore the colour picker on each button
    this.$nextTick(() => {
      const btns = document.querySelectorAll('.colour-but')
      for (let i = 0; i < btns.length; i++) {
        const btn = btns[i]
        if (!btn.jscolor) {
          new JSColor(btn, {
            showOnClick: false,
            preset: 'dark large',
            value: this.$store.pal.colours[i],
            onChange: function () {
              this.targetElement.dispatchEvent(new CustomEvent('col-change', { detail: { colour: this.toHEXString() } }))
            },
          })
        }
      }
    })

    this.$dispatch('palette', {})
  },
})
