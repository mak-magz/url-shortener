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

		<UContainer>
			<div class="grid grid-cols-1 lg:grid-cols-3 gap-4 ">
				<div class="col-span-1 lg:col-span-2">
					<UCard class="h-full">
						<template #default>
							<h6 class="text-lg font-semibold flex items-center gap-2">
								<UIcon
									name="i-lucide-chart-no-axes-combined"
									class="text-primary"
								/>
								<span class="uppercase font-black text-sm text-muted">Live Pulse</span>
							</h6>
							<h3 class="text-xl font-black mt-2">
								Real-time tracking
							</h3>

							<p class="text-muted mt-2 max-w-sm md:max-w-lg lg:max-w-sm">
								Every interaction is recorded in sub-milliseconds.
								Watch your campaign spread across the globe live.
							</p>

							<UCard
								class="mt-12"
								variant="subtle"
							>
								<template #default>
									<div class="flex items-center justify-between">
										<div>
											<h6 class="text-xs font-semibold uppercase">
												Active now
											</h6>
											<h3 class="text-4xl font-black tracking-tight">
												1,234
											</h3>
										</div>
										<div class="h-[50px] flex items-end gap-[2px] justify-between">
											<div class="h-3/10 w-[6px] bg-primary/10 rounded-t-sm" />
											<div class="h-5/10 w-[6px] bg-primary/30 rounded-t-sm" />
											<div class="h-4/10 w-[6px] bg-primary/20 rounded-t-sm" />
											<div class="h-8/10 w-[6px] bg-primary rounded-t-sm" />
											<div class="h-6/10 w-[6px] bg-primary/80 rounded-t-sm" />
											<div class="h-7/10 w-[6px] bg-primary/90 rounded-t-sm" />
										</div>
									</div>
									<div class="text-muted mt-2">
										<div class="flex items-center justify-between text-xs py-2">
											<div>/black-friday-sale</div>
											<div class="text-primary">
												+42 clicks
											</div>
										</div>
										<div class="flex items-center justify-between text-xs py-2">
											<div>/twitter-bio-link</div>
											<div class="text-primary">
												+18 clicks
											</div>
										</div>
									</div>
								</template>
							</UCard>
						</template>
					</UCard>
				</div>

				<div class="col-span-1 lg:col-span-1 h-full">
					<UCard
						class="h-full"
						variant="solid"
						:ui="{ body: 'flex flex-col h-full gap-8' }"
					>
						<div>
							<h6 class="text-lg font-semibold flex items-center gap-2">
								<UIcon
									name="i-lucide-globe"
								/>
								<span class="uppercase font-black text-sm text-muted">Global Reach</span>
							</h6>

							<h3 class="text-xl font-black mt-2">
								Geographical insights
							</h3>
							<p class="text-muted mt-2 text-sm">
								Heatmaps and regional breakdowns to pinpoint where your audience lives.
							</p>
						</div>

						<div class="h-[150px] w-full mt-auto card border border-neutral-700 rounded-lg overflow-hidden">
							<img
								src="https://lh3.googleusercontent.com/aida-public/AB6AXuCIVJtyR9ZOTWV7hRCym4f4ROTT97HDsokLLGG-zUwTUCq-rSFxs5VyesG_FHE47-pcf-p2nOvdBY9bw-PRjFXR2D_RxLcsA_3QrnxRquDF-_2bT5-PTy2jowecncNdrsTstKlkJIMu2DtLxXfyeVnsrqq68rda0Hk-yPPn7dKEOQ7c8qxNl0AChdUEJ9i5JkuQL6MnwfgutPwb1T8ReRvbWwgw3M-74lr_hH4IeNWcJsXbMBXJqbcZvk1hSHx80UowPlU7fOYuvRQ"
								alt="Map"
								class="w-full h-full object-cover"
							>
						</div>
					</UCard>
				</div>
			</div>
		</UContainer>
	</div>
</template>
