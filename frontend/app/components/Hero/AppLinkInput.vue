<script setup lang="ts">
import { useMutation } from '@pinia/colada'
import * as v from 'valibot'

import type { FormErrorEvent, FormSubmitEvent } from '#ui/types'

const schema = v.object({
	url: v.pipe(v.string(), v.url('Invalid URL'))
})

type Schema = v.InferOutput<typeof schema>

type UrlInfo = {
	shortCode: string
}

const config = useRuntimeConfig()
const toast = useToast()

const state = reactive<Schema>({
	url: ''
})

const originalUrl = ref('')
const shortUrl = ref('')

const { mutate: shortenUrl, status, asyncStatus, error, data } = useMutation({
	mutation: (url: string) => {
		return $fetch<{ data: UrlInfo }>(`${config.public.backendApiBaseUrl}/${config.public.backendApiPrefix}/${config.public.backendApiVersion}/shorten`, {
			method: 'POST',
			body: {
				originalUrl: url
			}
		})
	},
	onSuccess: (response) => {
		console.log(response)
		// response is correctly typed as { shortUrl: string }
		shortUrl.value = config.public.appBaseUrl?.replace('https://', '').replace('http://', '').replace('www.', '') + '/' + response.data.shortCode
		originalUrl.value = state.url
		state.url = ''
	},
	onError: (error) => {
		console.error(error)
	}
})

const onSubmit = async (event: FormSubmitEvent<Schema>) => {
	shortenUrl(event.data.url)
}

const onError = async (event: FormErrorEvent) => {
	console.log(event)
}

const copyToClipboard = () => {
	navigator.clipboard.writeText(shortUrl.value)
	toast.add({
		title: 'Copied to clipboard',
		color: 'success'
	})
}

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
		<UForm
			:schema="schema"
			:state="state"
			class="space-y-4"
			@submit="onSubmit"
			@error="onError"
		>
			<UFormField
				name="url"
			>
				<UInput
					v-model="state.url"
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
							type="submit"
						/>
					</template>
				</UInput>
			</UFormField>
		</UForm>
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
						{{ originalUrl }}
					</div>
				</div>
				<USeparator
					orientation="vertical"
					class="h-full px-4"
					size="md"
				/>
				<div class="flex flex-col  w-full">
					<span class="font-bold"> Short URL</span>
					<span class="text-ellipsis overflow-hidden whitespace-nowrap max-w-[200px]">{{ shortUrl }}</span>
				</div>
				<UButton
					label="Copy"
					color="primary"
					size="md"
					icon="i-lucide-copy"
					@click="copyToClipboard"
				/>
			</div>
		</UCard>
	</div>
</template>
