<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { api } from '../api'
import { setToken } from '../auth'

const router = useRouter()
const route = useRoute()
const error = ref('')
const ready = ref(false)
const clientId = import.meta.env.VITE_GOOGLE_CLIENT_ID || ''

async function handleCredential(response: GoogleCredentialResponse) {
  error.value = ''
  try {
    const result = await api.googleAuth(response.credential)
    setToken(result.token)
    const next = typeof route.query.next === 'string' ? route.query.next : '/'
    await router.replace(next)
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Google sign-in failed'
  }
}

function renderButton() {
  if (!window.google || !clientId) {
    ready.value = true
    return
  }
  window.google.accounts.id.initialize({
    client_id: clientId,
    callback: handleCredential,
    ux_mode: 'popup',
  })
  const el = document.getElementById('google-btn')
  if (el) {
    window.google.accounts.id.renderButton(el, {
      theme: 'outline',
      size: 'large',
      shape: 'pill',
      text: 'signin_with',
      width: 320,
    })
  }
  ready.value = true
}

onMounted(() => {
  if (!clientId) {
    ready.value = true
    return
  }
  if (window.google?.accounts.id) {
    renderButton()
    return
  }
  const timer = window.setInterval(() => {
    if (window.google?.accounts.id) {
      window.clearInterval(timer)
      renderButton()
    }
  }, 100)
  window.setTimeout(() => {
    window.clearInterval(timer)
    ready.value = true
  }, 5000)
})
</script>

<template>
  <div class="login-stage">
    <section class="panel">
      <p class="muted">Welcome to</p>
      <h1 class="brand">Vibe<span>Gopher</span></h1>
      <p class="lede">
        A small square for sharp takes. An AI critic replies — ideas get poked, not people.
      </p>

      <div class="google-slot">
        <div v-if="clientId" id="google-btn" />
        <p v-else class="error">
          Set <code>VITE_GOOGLE_CLIENT_ID</code> in <code>frontend/.env</code>.
        </p>
      </div>

      <p v-if="error" class="error">{{ error }}</p>
      <p v-if="ready && clientId" class="muted" style="margin-top: 1rem">
        Sign in with Google to join the feed. Password signup is legacy/API-only.
      </p>
    </section>
  </div>
</template>
