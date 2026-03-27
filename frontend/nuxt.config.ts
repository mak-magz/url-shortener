// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
	modules: [
		'@nuxt/eslint',
		'@nuxt/ui',
		'@nuxt/test-utils',
		'@pinia/colada-nuxt',
		'@pinia/nuxt'
	],

	ssr: true,

	devtools: {
		enabled: true
	},

	css: ['~/assets/css/main.css'],

	routeRules: {
		'/': { prerender: true }
	},

	compatibilityDate: '2025-01-15',

	eslint: {
		config: {
			stylistic: {
				indent: 'tab',
				commaDangle: 'never',
				braceStyle: '1tbs'
			}
		}
	}
})
