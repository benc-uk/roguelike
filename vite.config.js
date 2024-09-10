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
