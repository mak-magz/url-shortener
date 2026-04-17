<template>
	<div>
		<UPageHero>
			<template #headline>
				<UBadge
					variant="subtle"
					color="primary"
					size="sm"
					class="mb-4 uppercase font-black tracking-wider "
				>
					Analytics Cloud 2.0
				</UBadge>
			</template>

			<template #title>
				<div class="lg:max-w-3xl mx-auto">
					Data that drives your brand's <span class="text-primary">growth</span>
				</div>
			</template>
			<template #description>
				Turn every click into a customer insight.
				Professional-grade link tracking designed for performance marketing and high-growth teams.
			</template>

			<template #links>
				<UButton
					label="Get Started Free"
					color="primary"
					size="xl"
					class="w-[160px] justify-center text-sm p-3"
				/>
				<UButton
					label="Watch Demo"
					size="xl"
					variant="outline"
					class="w-[160px] justify-center text-sm p-3"
				/>
			</template>
		</UPageHero>
	</div>
</template>
