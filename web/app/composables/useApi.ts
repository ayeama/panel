export function useApi() {
    const httpScheme = window.location.protocol === "https:" ? "https:" : "http:"
    const wsScheme = window.location.protocol === "https:" ? "wss:" : "ws:"

    function http(path: string): string {
        if (import.meta.dev) {
            return `http://localhost:8000${path}`
        }
        return `${httpScheme}//${window.location.host}/api${path}`
    }

    function ws(path: string): string {
        if (import.meta.dev) {
            return `ws://localhost:8000${path}`
        }
        return `${wsScheme}//${window.location.host}/api${path}`
    }

    return {
        http,
        ws
    }
}
