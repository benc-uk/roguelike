export function randomHexColor() {
  const letters = '0123456789ABCDEF'
  let color = '#'
  for (let i = 0; i < 6; i++) {
    color += letters[Math.floor(Math.random() * 16)]
  }
  return color
}

export function isHexColor(str) {
  return /^#[0-9A-F]{6}$/i.test(str)
}
