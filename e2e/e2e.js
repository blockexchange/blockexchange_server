const {Builder, By, until} = require('selenium-webdriver');

(async function example() {
  let driver = await new Builder()
    .usingServer('http://selenium:4444/wd/hub')
    .forBrowser('firefox')
    .build();
  try {
    await driver.get('http://blockexchange:8080/#/');
    await driver.wait(until.elementLocated(By.className('img-fluid')), 10000);
  } finally {
    await driver.quit();
    console.log("e2e tests done!");
  }
})();
