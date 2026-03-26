<script lang="ts" setup>
import { useMutation } from '@pinia/colada'

type UrlInfo = {
  shortCode: string
}

const _url = ref('')
const shortUrl = ref('')

const { mutate: _shortenUrl, status, asyncStatus, error, data } = useMutation({
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
  <div>
    <!-- <UContainer class="flex flex-col justify-center items-center mt-30">
      <div class="flex gap-4 mb-4">
        <UInput
          v-model="url"
          placeholder="Enter URL"
          label="URL"
        />
        <UButton
          label="Shorten"
          color="primary"
          variant="solid"
          @click="shortenUrl(url)"
        />
      </div>

      <div>
        <p v-if="status === 'pending' && asyncStatus === 'loading'">
          Shortening URL...
        </p>
        <p v-if="status === 'success'">
          Shortened URL: {{ shortUrl }}
        </p>
        <p v-if="status === 'error'">
          Error: {{ error }}
        </p>
      </div>
    </UContainer> -->
    <UPageHero
      orientation="horizontal"
      title="Shorten Your Links, Expand Your Reach"
      description="SwiftLink helps you create short, powerful links that are easy to share and track. Boost your click-through rates and manage your digital footprint with ease."
      :links="[{
        label: 'Get started',
        to: 'https://ui.nuxt.com/docs/getting-started/installation/nuxt',
        target: '_blank',
        trailingIcon: 'i-lucide-arrow-right',
        size: 'xl'
      }, {
        label: 'Use this template',
        to: 'https://github.com/nuxt-ui-templates/starter',
        target: '_blank',
        icon: 'i-simple-icons-github',
        size: 'xl',
        color: 'neutral',
        variant: 'subtle'
      }]"
    >
      <template #title>
        <div class="flex flex-col justify-center items-center">
          <p>Shorten Your Links, <span class="text-primary">Expand Your Reach</span></p>
        </div>
      </template>

      <AppSampleCard />
    </UPageHero>
    <!--
    <UPageSection
      id="features"
      title="Everything you need to build modern Nuxt apps"
      description="Start with a solid foundation. This template includes all the essentials for building production-ready applications with Nuxt UI's powerful component system."
      :features="[{
        icon: 'i-lucide-rocket',
        title: 'Production-ready from day one',
        description: 'Pre-configured with TypeScript, ESLint, Tailwind CSS, and all the best practices. Focus on building features, not setting up tooling.'
      }, {
        icon: 'i-lucide-palette',
        title: 'Beautiful by default',
        description: 'Leveraging Nuxt UI\'s design system with automatic dark mode, consistent spacing, and polished components that look great out of the box.'
      }, {
        icon: 'i-lucide-zap',
        title: 'Lightning fast',
        description: 'Optimized for performance with SSR/SSG support, automatic code splitting, and edge-ready deployment. Your users will love the speed.'
      }, {
        icon: 'i-lucide-blocks',
        title: '100+ components included',
        description: 'Access Nuxt UI\'s comprehensive component library. From forms to navigation, everything is accessible, responsive, and customizable.'
      }, {
        icon: 'i-lucide-code-2',
        title: 'Developer experience first',
        description: 'Auto-imports, hot module replacement, and TypeScript support. Write less boilerplate and ship more features.'
      }, {
        icon: 'i-lucide-shield-check',
        title: 'Built for scale',
        description: 'Enterprise-ready architecture with proper error handling, SEO optimization, and security best practices built-in.'
      }]"
    />

    <UPageSection>
      <UPageCTA
        title="Ready to build your next Nuxt app?"
        description="Join thousands of developers building with Nuxt and Nuxt UI. Get this template and start shipping today."
        variant="subtle"
        :links="[{
          label: 'Start building',
          to: 'https://ui.nuxt.com/docs/getting-started/installation/nuxt',
          target: '_blank',
          trailingIcon: 'i-lucide-arrow-right',
          color: 'neutral'
        }, {
          label: 'View on GitHub',
          to: 'https://github.com/nuxt-ui-templates/starter',
          target: '_blank',
          icon: 'i-simple-icons-github',
          color: 'neutral',
          variant: 'outline'
        }]"
      />
    </UPageSection> -->
  </div>
</template>
