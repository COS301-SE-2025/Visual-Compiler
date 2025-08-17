const { defineConfig } = require("cypress");

module.exports = defineConfig({
    e2e: {
        setupNodeEvents(on, config) {
      
        },
        specPattern: 'e2e-tests/e2e/**/*.cy.{js,jsx,ts,tsx}',
        supportFile: 'e2e-tests/support/e2e.{js,jsx,ts,tsx}',
        fixturesFolder: 'e2e-tests/fixtures',
        screenshotsFolder: 'e2e-tests/screenshots',
        videosFolder: 'e2e-tests/videos',
        env: {
            frontend_url: 'http://localhost:5173',
            backend_url: 'http://localhost:8080',
        }
    },
});
