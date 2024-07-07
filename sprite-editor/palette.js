export default () => ({
  editColour(e, i) {
    this.$store.pal.colours[i] = e.target.value
    this.$store.pal.select(i)

    this.$dispatch('palette', {})
  },

  selectColour(e, i) {
    e.preventDefault()
    this.$store.pal.select(i)
  },
})
