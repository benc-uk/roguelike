import Alpine from 'https://cdn.jsdelivr.net/npm/alpinejs@3.14.0/dist/module.esm.min.js'

Alpine.data('palette', () => ({
  editColour(e, i) {
    this.$store.pal.colours[i] = e.target.value
    this.$store.pal.select(i)

    this.$dispatch('palette', {})
  },

  selectColour(e, i) {
    e.preventDefault()
    this.$store.pal.select(i)
  },
}))
