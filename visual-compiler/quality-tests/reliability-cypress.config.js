const { defineConfig } = require("cypress");

module.exports = defineConfig({
    e2e: {
        setupNodeEvents(on, config) {
      
        },
        specPattern: 'reliability-tests/tests/**/*.cy.js',
        supportFile: 'reliability-tests/support/e2e.js',
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
