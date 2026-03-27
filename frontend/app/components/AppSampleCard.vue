<template>
	<UCard>
		<template #header>
			<div class="flex gap-2 items-center">
				<div class="h-3 w-3 bg-red-500 rounded-full" />
				<div class="h-3 w-3 bg-yellow-500 rounded-full" />
				<div class="h-3 w-3 bg-green-500 rounded-full" />
				<USkeleton class="ml-2 h-4 w-full" />
			</div>
		</template>
		<template #default>
			<UCard
				variant="soft"
				class="bg-purple-200"
			>
				<div class="h-[250px] flex flex-col gap-4 items-center justify-center">
					<div class="w-full flex flex-row gap-4 items-center justify-center">
						<USkeleton
							class="
								h-10 max-w-[150px]
								overflow-x-hidden
								text-ellipsis
								whitespace-nowrap
								bg-white
								rounded-lg
								text-gray-700
								items-center
								p-2
								shadow-sm"
						>
							https://your-very-long-url.com/this-is-a-very-long-url-that-needs-to-be-shortened
						</USkeleton>

						<UButton
							color="primary"
							variant="link"
							icon="i-lucide-arrow-right"
						/>

						<USkeleton
							class="
								h-10
								bg-primary
								rounded-lg
								text-white
								items-center
								p-2
								w-fit
								shadow-sm"
						>
							swft.lk/growth
						</USkeleton>
					</div>

					<UCard
						class="w-[80%]"
						variant="soft"
					>
						<div class="flex items-end justify-start gap-2  h-[70px]">
							<USkeleton class="h-6 w-1/5 bg-purple-300" />
							<USkeleton class="h-10 w-1/5 bg-purple-400" />
							<USkeleton class="h-14 w-1/5 bg-purple-500" />
							<USkeleton class="h-12 w-1/5 bg-purple-600" />
							<USkeleton class="h-10 w-1/5 bg-purple-700" />
						</div>
					</UCard>
				</div>
			</UCard>
		</template>
	</UCard>
</template>
