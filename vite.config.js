// Hacky plugin to force full reload when files change, as HMR will not work with WASM
const fullReloadAlways = {
  name: "full-reload",
  handleHotUpdate({ server }) {
    server.ws.send({ type: "full-reload" });
    return [];
  },
};

export default {
  root: "./web",
  server: {
    port: 3000,
  },
  build: {
    watch: {
      include: "./web/**",
    },
    hmr: false,
  },
  plugins: [fullReloadAlways],
};
