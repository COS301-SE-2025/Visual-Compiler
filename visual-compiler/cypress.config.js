import { defineConfig } from 'cypress';

export default defineConfig({
  e2e: {
    supportFile: './e2e-tests/support/e2e.js',
    specPattern: './e2e-tests/e2e/**/*.cy.{js,ts}',
    fixturesFolder: './e2e-tests/fixtures',
    videosFolder: './e2e-tests/videos',
    screenshotsFolder: './e2e-tests/screenshots',
    env: {
    frontendUrl: 'http://localhost:5173',
    backendUrl: 'http://localhost:8080',
    }
  },
});