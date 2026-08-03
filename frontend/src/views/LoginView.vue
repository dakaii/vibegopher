<script setup lang="ts">
import { SignIn } from '@clerk/vue'
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const publishableKey = import.meta.env.VITE_CLERK_PUBLISHABLE_KEY || ''

const redirectUrl = computed(() => {
  const next = typeof route.query.next === 'string' ? route.query.next : '/'
  return next || '/'
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

      <div class="clerk-slot">
        <SignIn
          v-if="publishableKey"
          routing="hash"
          :fallback-redirect-url="redirectUrl"
          :force-redirect-url="redirectUrl"
        />
        <p v-else class="error">
          Set <code>VITE_CLERK_PUBLISHABLE_KEY</code> in <code>frontend/.env</code>.
        </p>
      </div>

      <p v-if="publishableKey" class="muted" style="margin-top: 1rem">
        Sign in with Clerk to join the feed. Password signup is legacy/API-only.
      </p>
    </section>
  </div>
</template>
