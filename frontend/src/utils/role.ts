export type AccessRole = 'server' | 'client'

export function roleFromPath(path: string): AccessRole | null {
  if (path === '/server' || path.startsWith('/server/')) return 'server'
  if (path === '/client' || path.startsWith('/client/')) return 'client'
  return null
}

/**
 * Explicit /server and /client routes are authoritative. Hostname detection is
 * retained only as a compatibility fallback for old bookmarks and the root URL.
 */
export function currentRole(path = typeof window === 'undefined' ? '' : window.location.pathname): AccessRole {
  const routeRole = roleFromPath(path)
  if (routeRole) return routeRole
  if (typeof window === 'undefined') return 'client'
  const hostname = window.location.hostname
  return hostname === 'localhost' || hostname === '127.0.0.1' || hostname === '::1' || hostname === '[::1]' || hostname === ''
    ? 'server'
    : 'client'
}

export function isServerAccess(path?: string): boolean {
  return currentRole(path) === 'server'
}

export function defaultRouteFor(role: AccessRole): string {
  return role === 'server' ? '/server/devices' : '/client/inbox'
}
