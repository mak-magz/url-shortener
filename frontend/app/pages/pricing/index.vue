<script setup lang="ts">
import type { PricingPlanProps } from '@nuxt/ui'

const plans = ref<PricingPlanProps[]>([
	{
		title: 'FREE TIER',
		description: 'Essential tools for personal projects.',
		price: '$0',
		features: [
			'Basic link shortening',
			'20 links per month',
			'Basic analytics'
		],
		button: {
			label: 'Start for free',
			color: 'neutral',
			variant: 'outline'
		},
		billingCycle: '/mo'
	},
	{
		title: 'PRO PLAN',
		description: 'Scale your reach with advanced tools.',
		price: '$15',
		scale: true,
		features: [
			'Unlimited links',
			'Custom aliases',
			'Advanced analytics',
			'Branded QR codes'
		],
		button: {
			label: 'Buy now',
			color: 'primary',
			variant: 'solid'
		},
		billingCycle: '/mo',
		badge: {
			label: 'Most Popular',
			color: 'primary'
		},
		highlight: true
	},
	{
		title: 'ENTERPRISE PLAN',
		description: 'Ideal for larger teams and organizations.',
		price: 'Custom',
		features: [
			'Branded domains',
			'Team collaboration',
			'API access',
			'Dedicated support'
		],
		button: {
			label: 'Contact Sales',
			color: 'neutral',
			variant: 'outline'
		}
	}
])
</script>

<template>
	<div>
		<UPageHero
			title="Simple, transparent pricing"
			description="Choose the plan that's right for you. All plans come with a 7-day free trial. Save 20% with annual billing."
		/>

		<UContainer>
			<UPricingPlans
				:plans="plans"
				scale
			/>
		</UContainer>
	</div>
</template>
