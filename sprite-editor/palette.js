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

  editColour(e, i) {
    this.$store.pal.colours[i] = e.target.value
    this.$store.pal.select(i)

    this.$dispatch('palette', {})
  },

  selectColour(e, i) {
    e.preventDefault()
    this.$store.pal.select(i)
  },

  bulkSave() {
    this.bulkEdit = false

    const colours = this.bulkColours.split('\n')
    for (let i = 0; i < this.$store.pal.colours.length; i++) {
      if (!colours[i]) {
        continue
      }

      let colour = colours[i].trim()
      if (colour[0] !== '#') {
        colour = '#' + colour
      }

      if (!isHexColor(colour)) {
        colour = '#000000'
      }

      // Add the # back to the colour!
      this.$store.pal.colours[i] = colour
    }

    this.$dispatch('palette', {})
  },

  createColourPicker(e) {
    console.log('Create colour picker', this.$store.pal.selected())

    // Create a colour picker
    const input = document.createElement('input')
    input.type = 'color'
    // input.value = this.$store.pal.selected()
    input.style.position = 'absolute'
    input.style.left = e.clientX + 'px'
    input.style.top = e.clientY + 'px'
    input.style.zIndex = 1000

    document.body.appendChild(input)

    input.addEventListener('input', (e) => {
      this.$store.pal.colours[this.$store.pal.selected()] = e.target.value
      this.$dispatch('palette', {})
    })

    input.addEventListener('change', () => {
      document.body.removeChild(input)
    })

    input.click()
  },
})
