<script setup lang="ts">
import type { PricingPlanProps, PricingTableSection, PricingTableTier } from '@nuxt/ui'

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

const tiers = ref<PricingTableTier[]>([
	{
		id: 'free',
		title: 'Free'
	},
	{
		id: 'pro',
		title: 'Pro',
		highlight: true
	},
	{
		id: 'enterprise',
		title: 'Enterprise'
	}
])

const sections = ref<PricingTableSection[]>([
	{
		title: 'Features',
		features: [
			{
				title: 'Monthly Links',
				tiers: {
					free: '20',
					pro: 'Unlimited',
					enterprise: 'Unlimited'
				}
			},
			{
				title: 'Link History',
				tiers: {
					free: '30 days',
					pro: 'Unlimited',
					enterprise: 'Unlimited'
				}
			},
			{
				title: 'Realtime Analytics',
				tiers: {
					free: false,
					pro: true,
					enterprise: true
				}
			},
			{
				title: 'QR Code Generation',
				tiers: {
					free: true,
					pro: true,
					enterprise: true
				}
			},
			{
				title: 'Branded domains',
				tiers: {
					free: false,
					pro: false,
					enterprise: true
				}
			},
			{
				title: 'Team collaboration',
				tiers: {
					free: false,
					pro: false,
					enterprise: true
				}
			},
			{
				title: 'API access',
				tiers: {
					free: false,
					pro: false,
					enterprise: true
				}
			},
			{
				title: 'Dedicated support',
				tiers: {
					free: false,
					pro: false,
					enterprise: true
				}
			}
		]
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

		<UPageSection>
			<h3 class="text-center text-3xl font-bold">
				Compare Plans
			</h3>

			<UCard>
				<UPricingTable
					:tiers="tiers"
					:sections="sections"
				/>
			</UCard>
		</UPageSection>

		<PricingFaq />
	</div>
</template>
