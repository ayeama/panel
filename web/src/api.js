export const API_URL = `${window.location.origin}/api`
export const API_WS = `${window.location.protocol === 'https:' ? 'wss' : 'ws'}://${window.location.host}/api`
export const API_DAV = `${window.location.protocol === 'https:' ? 'davs' : 'dav'}://${window.location.host}/api`

export const API_ERROR_NOT_FOUND = 'API_ERROR_NOT_FOUND'
