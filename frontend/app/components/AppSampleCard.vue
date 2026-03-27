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
			<UCard variant="soft">
				<div class="h-[200px]">
					<div class="w-full flex flex-col gap-2">
						<div
							class="
								h-10 max-w-[150px]
								overflow-x-hidden
								text-ellipsis
								whitespace-nowrap
								bg-gray-200
								rounded-lg
								text-gray-700
								items-center
								p-2"
						>
							https://your-very-long-url.com/this-is-a-very-long-url-that-needs-to-be-shortened
						</div>
						<div
							class="
								h-10
								bg-primary
								rounded-lg
								text-white
								items-center
								p-2
								w-fit"
						>
							swft.lk/growth
						</div>
					</div>
				</div>
			</UCard>
		</template>
	</UCard>
</template>
