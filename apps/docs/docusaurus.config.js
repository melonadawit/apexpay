// Docusaurus config — Day 6 docs site like ApexPay/ApexPay docs.ApexPay.co + developer guides outstanding
// Best practice: dark/light, search, i18n AM/EN, OpenAPI Swagger embedded

/** @type {import('@docusaurus/types').Config} */
const config = {
  title: 'ApexPay Docs — AI-native Payment Gateway for Ethiopia',
  tagline: 'Collect, Disburse, Payroll, RAG Compliance, Swarm AI — ETB-first, NBE PSO Gateway Operator ONPS/02/2020, Fayda ID front/back <2MB OTP consent id.gov.et',
  favicon: 'img/favicon.ico',

  url: 'https://docs.apexpay.et',
  baseUrl: '/',

  organizationName: 'zitadave',
  projectName: 'ApexPay',

  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'warn',

  i18n: {
    defaultLocale: 'en',
    locales: ['en', 'am'],
    localeConfigs: {
      en: { label: 'English' },
      am: { label: 'አማርኛ' },
    },
  },

  presets: [
    [
      'classic',
      /** @type {import('@docusaurus/preset-classic').Options} */
      ({
        docs: {
          sidebarPath: './sidebars.js',
          editUrl: 'https://github.com/zitadave/ApexPay/tree/main/apps/docs/',
        },
        blog: {
          showReadingTime: true,
        },
        theme: {
          customCss: './src/css/custom.css',
        },
      }),
    ],
  ],

  themeConfig:
    /** @type {import('@docusaurus/preset-classic').ThemeConfig} */
    ({
      image: 'img/apexpay-social-card.png',
      navbar: {
        title: 'ApexPay Docs',
        logo: { alt: 'ApexPay Logo', src: 'img/logo.png' },
        items: [
          { type: 'docSidebar', sidebarId: 'tutorialSidebar', position: 'left', label: 'Docs' },
          { to: '/api', label: 'API Reference (OpenAPI)', position: 'left' },
          { type: 'localeDropdown', position: 'right' },
          { href: 'https://github.com/zitadave/ApexPay', label: 'GitHub', position: 'right' },
        ],
      },
      footer: {
        style: 'dark',
        links: [
          {
            title: 'Docs',
            items: [
              { label: 'Quickstart 6 Lines', to: '/docs/quickstart' },
              { label: 'Merchant Onboarding NBE + Fayda', to: '/docs/onboarding' },
              { label: 'Fayda ID Front/Back OTP', to: '/docs/fayda' },
            ],
          },
          {
            title: 'API',
            items: [
              { label: 'OpenAPI 21 Paths', to: '/docs/api-reference' },
              { label: 'Postman Collection', to: '/docs/postman' },
              { label: 'ApexPay vs ApexPay vs ApexPay', to: '/docs/comparison' },
            ],
          },
          {
            title: 'Community',
            items: [
              { label: 'Merchant Guides EN/AM Outstanding', to: '/docs/merchant-guide' },
              { label: 'NBE Compliance RAG Citations', to: '/docs/compliance-rag' },
              { label: 'Swarm AI Agents', to: '/docs/swarm' },
            ],
          },
        ],
        copyright: `Copyright © ${new Date().getFullYear()} ApexPay — collect, pay, prove. Built with Docusaurus.`,
      },
      prism: {
        theme: require('prism-react-renderer').themes.github,
        darkTheme: require('prism-react-renderer').themes.dracula,
      },
    }),
};

module.exports = config;
