<script setup lang="ts">
import { useMutation } from '@pinia/colada'

type UrlInfo = {
	shortCode: string
}

const url = ref('')
const shortUrl = ref('')

const { mutate: shortenUrl, status, asyncStatus, error, data } = useMutation({
	mutation: (url: string) => {
		return $fetch<{ data: UrlInfo }>('http://localhost:8080/api/v1/shorten', {
			method: 'POST',
			body: {
				originalUrl: url
			}
		})
	},
	onSuccess: (response) => {
		console.log(response)
		// response is correctly typed as { shortUrl: string }
		shortUrl.value = 'http://localhost:8080/' + response.data.shortCode
	},
	onError: (error) => {
		console.error(error)
	}
})

watch(data, (newData) => {
	console.log(newData)
})

watch(error, (newError) => {
	console.log(newError)
})

watch(status, (newStatus) => {
	console.log(newStatus)
})

watch(asyncStatus, (newAsyncStatus) => {
	console.log(newAsyncStatus)
})
</script>

<template>
	<div class="flex flex-col gap-2">
		<UInput
			v-model="url"
			placeholder="Enter your URL"
			size="xl"
			class="w-full"
			variant="soft"
			icon="i-lucide-link"
			:ui="{
				base: 'overflow-x-hidden text-ellipsis whitespace-nowrap truncate'
			}"
		>
			<template #trailing>
				<UButton
					label="Shorten"
					color="primary"
					size="md"
					@click="shortenUrl(url)"
				/>
			</template>
		</UInput>

		<UCard
			v-if="shortUrl"
			:ui="{
				body: 'p-2 sm:p-2'
			}"
			variant="soft"
		>
			<div class="flex h-12 justify-between items-center">
				<div class="w-full">
					<div class="font-bold">
						Original URL
					</div>
					<div class="text-ellipsis overflow-hidden whitespace-nowrap max-w-[200px]">
						{{ url }}
					</div>
				</div>
				<USeparator
					orientation="vertical"
					class="h-full px-4"
					size="md"
				/>
				<div class="flex flex-col  w-full">
					<span class="font-bold"> Short URL</span>
					<span>{{ shortUrl }}</span>
				</div>
				<UButton
					label="Copy"
					color="primary"
					size="md"
				/>
			</div>
		</UCard>
	</div>
</template>
