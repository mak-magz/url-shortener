export const isValidUrl = (url: string) => {
	try {
		const parsed = new URL(url)
		return (parsed.protocol === 'http:' || parsed.protocol === 'https:') && parsed.host.includes('.')
	} catch {
		return false
	}
}
