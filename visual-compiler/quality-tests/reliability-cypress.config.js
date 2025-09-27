const { defineConfig } = require("cypress");

module.exports = defineConfig({
    e2e: {
        setupNodeEvents(on, config) {
      
        },
        specPattern: 'reliability-tests/tests/**/*.cy.js',
        supportFile: false,
        fixturesFolder: 'reliability-tests/fixtures',
        screenshotsFolder: 'reliability-tests/screenshots',
        videosFolder: 'reliability-tests/videos',
        env: {
            frontend_url: 'http://localhost:5173',
            backend_url: 'http://localhost:8080',
        },
        viewportWidth: 1290,
        viewportHeight: 1290,
    },
});
