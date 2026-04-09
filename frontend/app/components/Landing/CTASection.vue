<script setup lang="ts">
import type { ButtonProps } from '@nuxt/ui'

const links = ref<ButtonProps[]>([
	{
		label: 'Create Free Account',
		color: 'neutral'
	},
	{
		label: 'View Pricing',
		color: 'neutral',
		variant: 'subtle'
	}
])
</script>

<template>
	<UPageCTA
		title="Ready to start shortening?"
		description="Join over 10,000+ creators and businesses who trust SwiftLink for their digital link management."
		:links="links"
		variant="subtle"
	/>
</template>
