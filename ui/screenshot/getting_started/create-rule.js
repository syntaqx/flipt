const { capture } = require('../../screenshot.js');

(async () => {
  await capture('getting_started', 'create_rule', async (page) => {
    await page.getByRole('link', { name: 'new-login' }).click();
    await page.getByRole('link', { name: 'Evaluation' }).click();
    await page.getByRole('button', { name: 'New Rule' }).click();
    await page.locator('#segmentKey-select-button').click();
    await page.getByText('all-users').click();
    await page.getByLabel('Multi-Variate').check();
  });
})();
