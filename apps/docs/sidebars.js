/**
 * Sidebars for ApexPay Docs — Day 6 docs site like ApexPay/ApexPay docs.ApexPay.co
 * Outstanding modern: 6 lines copy-paste, merchant onboarding NBE + Fayda, bank list, error codes, SDKs
 */

/** @type {import('@docusaurus/plugin-content-docs').SidebarsConfig} */
const sidebars = {
  tutorialSidebar: [
    'intro',
    'quickstart',
    {
      type: 'category',
      label: 'Merchant Onboarding NBE-Grade',
      items: ['onboarding', 'fayda', 'merchant-guide'],
    },
    {
      type: 'category',
      label: 'Payments + Smart Routing',
      items: ['payments', 'routing', 'refunds', 'payouts', 'payroll'],
    },
    {
      type: 'category',
      label: 'Intelligence — RAG + Swarm',
      items: ['compliance-rag', 'swarm', 'recon'],
    },
    {
      type: 'category',
      label: 'API Reference',
      items: ['api-reference', 'postman', 'comparison'],
    },
  ],
};

module.exports = sidebars;
