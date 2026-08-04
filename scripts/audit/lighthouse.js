#!/usr/bin/env node
// Lighthouse audit — target 90+ per NFR outstanding UI, WCAG AA, SEO 90
// Best practice: CI check performance 90, accessibility 100, best-practices 100

const lighthouse = require('lighthouse');
const chromeLauncher = require('chrome-launcher');

const urls = [
  'http://localhost:3000', // merchant-web landing
  'http://localhost:3000/onboarding', // onboarding wizard 6 steps outstanding
  'http://localhost:3001', // checkout-web outstanding mobile 420px
  'http://localhost:3000/compliance', // compliance center RAG chat
];

async function run() {
  const chrome = await chromeLauncher.launch({ chromeFlags: ['--headless'] });
  const options = { logLevel: 'info', output: 'json', onlyCategories: ['performance','accessibility','best-practices','seo'], port: chrome.port };

  for (const url of urls) {
    console.log(`\nAuditing ${url}...`);
    const runnerResult = await lighthouse(url, options);
    const lhr = runnerResult.lhr;
    const perf = lhr.categories.performance.score * 100;
    const a11y = lhr.categories.accessibility.score * 100;
    const bp = lhr.categories['best-practices'].score * 100;
    const seo = lhr.categories.seo.score * 100;

    console.log(`Results for ${url}: Perf ${perf} | A11y ${a11y} | BP ${bp} | SEO ${seo}`);

    if (perf < 90) console.warn(`⚠ Performance ${perf} <90 for ${url} — need optimize images, dynamic import Lottie, code split`);
    if (a11y < 100) console.warn(`⚠ Accessibility ${a11y} <100 — check axe audit WCAG AA`);
    if (bp < 100) console.warn(`⚠ Best Practices ${bp} <100`);
    if (seo < 90) console.warn(`⚠ SEO ${seo} <90`);

    // Outstanding: axe audit via @axe-core/cli would run here too
  }

  await chrome.kill();
}

run().catch(e=> { console.error(e); process.exit(1); });
