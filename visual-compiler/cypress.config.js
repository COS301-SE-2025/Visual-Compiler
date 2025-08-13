const { defineConfig } = require("cypress");

module.exports = defineConfig({
  e2e: {
    setupNodeEvents(on, config) {
  
    },
    env: {
        frontend_url: 'http://localhost:5173',
        backend_url: 'http://localhost:8080',
      }
  },
});
