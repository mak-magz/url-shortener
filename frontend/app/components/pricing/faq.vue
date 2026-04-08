<script setup lang="ts">
import type { AccordionItem } from '@nuxt/ui'

const items = ref<AccordionItem[]>([
	{
		label: 'Can I cancel anytime?',
		content: 'Yes, you can cancel your subscription at any time. There are no cancellation fees, and you can continue to use the service until the end of your billing cycle.'
	},
	{
		label: 'Do you offer a free trial?',
		content: 'Yes, all plans come with a 7-day free trial. You can try out all the features of the Pro plan for 7 days before deciding to upgrade.'
	},
	{
		label: 'Can I use custom domains?',
		content: 'Yes, you can use custom domains with SwiftLink. This allows you to create branded links that are more trustworthy and professional.'
	},
	{
		label: 'What is a custom alias?',
		content: 'A custom alias is a short, unique identifier that you can use to shorten your URLs. It allows you to create branded links that are more trustworthy and professional.'
	}
])
</script>

<template>
	<UContainer class="p-8">
		<h3 class="text-center text-3xl font-bold mb-16">
			Frequently Asked Questions
		</h3>
		<UCard>
			<UAccordion :items="items" />
		</UCard>
	</UContainer>
</template>
