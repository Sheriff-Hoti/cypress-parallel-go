import { defineConfig } from "cypress";

export default defineConfig({
  defaultCommandTimeout: 30000,
  pageLoadTimeout: 60000,
  viewportWidth: 1440,
  viewportHeight: 900,
  chromeWebSecurity: false,
  screenshotsFolder: "cypress/screenshots",
  screenshotOnRunFailure: true,

  reporter: "json",

  reporterOptions: {
    output: "filename.json",
  },

  retries: {
    runMode: 2,
    openMode: 3,
  },
  e2e: {
    supportFile: false,
    setupNodeEvents(on, config) {
      return config;
    },
    experimentalRunAllSpecs: true,
    video: false,
    experimentalOriginDependencies: true,
    defaultCommandTimeout: 30000,
    requestTimeout: 30000,
  },
});
