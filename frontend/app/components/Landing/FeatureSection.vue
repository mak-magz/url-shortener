<script lang="ts" setup>
import type { PageFeatureProps } from '@nuxt/ui'

const features = ref<PageFeatureProps[]>([
	{
		title: 'Real-time Analytics',
		description: 'Track every click in real-time with detailed geographic data, referral sources, and device information to optimize your marketing.',
		icon: 'i-lucide-chart-column-stacked',
		to: '/docs/getting-started/integrations/icons',
		orientation: 'vertical'
	},
	{
		title: 'Custom Branded Links',
		description: 'Replace generic domains with your own vanity URL. Increase click-through rates by up to 34% by building brand trust and recognition.',
		icon: 'i-lucide-a-large-small',
		to: '/docs/getting-started/integrations/fonts',
		orientation: 'vertical'
	},
	{
		title: 'Secure & Private',
		description: 'We prioritize security with SSL encryption, link expiration options, and private analytics dashboards to keep your data safe.',
		icon: 'i-lucide-shield-check',
		to: '/docs/getting-started/integrations/color-mode',
		orientation: 'vertical'
	}
])
</script>

<template>
	<UPageSection
		title="Everything you need to manage your links"
		description="Our suite of tools gives you total control over your shared content with industry-leading technology."
		headline="Platform Benefits"
		:features="features"
	/>
</template>
