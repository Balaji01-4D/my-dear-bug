export function navigate(path: string) {
	if (window.location.pathname !== path) {
		window.history.pushState({}, '', path)
		window.dispatchEvent(new PopStateEvent('popstate'))
	}
}

export function normalizeConfessionPath() {
	const m = window.location.pathname.match(/^\/confessions\/(\d+)$/)
	if (m) {
		navigate(`/c${window.location.pathname}`)
	}
}
