import { ref } from 'vue'

const TOKEN_KEY = 'vibegopher_token'

/** Reactive Clerk session flags for the router (bound from App.vue). */
export const clerkLoaded = ref(false)
export const clerkSignedIn = ref(false)

type TokenGetter = () => Promise<string | null | undefined>
type SignOutFn = () => Promise<void>

let clerkTokenGetter: TokenGetter | null = null
let clerkSignOut: SignOutFn | null = null

export function bindClerkSession(opts: {
  isLoaded: boolean
  isSignedIn: boolean
  getToken: TokenGetter
}): void {
  clerkLoaded.value = opts.isLoaded
  clerkSignedIn.value = opts.isSignedIn
  clerkTokenGetter = opts.getToken
}

export function setClerkSignOut(fn: SignOutFn): void {
  clerkSignOut = fn
}

export async function signOutClerk(): Promise<void> {
  if (clerkSignOut) {
    await clerkSignOut()
  }
}

/** Thin localStorage shim for password-auth API JWTs during local/manual testing. */
export function getLegacyToken(): string {
  return localStorage.getItem(TOKEN_KEY) || ''
}

export function setLegacyToken(token: string): void {
  localStorage.setItem(TOKEN_KEY, token)
}

export function clearLegacyToken(): void {
  localStorage.removeItem(TOKEN_KEY)
}

export async function getAccessToken(): Promise<string> {
  if (clerkTokenGetter) {
    const token = await clerkTokenGetter()
    if (token) return token
  }
  return getLegacyToken()
}

export function isSignedIn(): boolean {
  if (clerkLoaded.value) {
    return clerkSignedIn.value || Boolean(getLegacyToken())
  }
  return Boolean(getLegacyToken())
}

export async function waitForClerk(timeoutMs = 10000): Promise<void> {
  if (clerkLoaded.value) return
  const started = Date.now()
  await new Promise<void>((resolve) => {
    const tick = () => {
      if (clerkLoaded.value || Date.now() - started > timeoutMs) {
        resolve()
        return
      }
      window.setTimeout(tick, 50)
    }
    tick()
  })
}
