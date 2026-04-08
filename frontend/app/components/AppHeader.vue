<script setup lang="ts">
import type { NavigationMenuItem } from '@nuxt/ui'

const route = useRoute()

const items = computed<NavigationMenuItem[]>(() => [
	{
		label: 'Features',
		to: '/features',
		active: route.path.startsWith('/features')
	},
	{
		label: 'Pricing',
		to: '/pricing',
		active: route.path.startsWith('/pricing')
	},
	{
		label: 'Analytics',
		to: '/analytics'
	},
	{
		label: 'Enterprise',
		to: '/enterprise'
	}
])
</script>

<template>
	<UHeader mode="slideover">
		<template #left>
			<NuxtLink to="/">
				<AppLogo class="w-auto h-6 shrink-0" />
			</NuxtLink>

			<!-- <TemplateMenu /> -->
		</template>

		<UNavigationMenu :items="items" />

		<template #right>
			<UColorModeButton />

			<UButton
				to="/login"
				label="Login"
				color="neutral"
				variant="ghost"
			/>

			<UButton
				to="https://github.com/nuxt-ui-templates/starter"
				target="_blank"
				icon="i-simple-icons-github"
				aria-label="GitHub"
				color="neutral"
				variant="ghost"
			/>
		</template>

		<template #body>
			<UNavigationMenu
				:items="items"
				orientation="vertical"
				class="-mx-2.5"
				slideover
			/>
		</template>
	</UHeader>
</template>
