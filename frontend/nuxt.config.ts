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

	runtimeConfig: {
		public: {
			backendApiBaseUrl: process.env.BACKEND_API_BASE_URL,
			backendApiVersion: process.env.BACKEND_API_VERSION,
			backendApiPrefix: process.env.BACKEND_API_PREFIX,
			appBaseUrl: process.env.APP_BASE_URL
		}
	},

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
	},

	icon: {
		clientBundle: {
			scan: true,
			sizeLimitKb: 256
		},
		serverBundle: 'local'
	}
})
