const MUTE_CRITIC_KEY = 'vibegopher_mute_critic'

export function isCriticMuted(): boolean {
  return localStorage.getItem(MUTE_CRITIC_KEY) === '1'
}

export function setCriticMuted(muted: boolean): void {
  if (muted) {
    localStorage.setItem(MUTE_CRITIC_KEY, '1')
  } else {
    localStorage.removeItem(MUTE_CRITIC_KEY)
  }
}
