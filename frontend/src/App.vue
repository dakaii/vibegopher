<script setup lang="ts">
import { RouterView } from 'vue-router'
import { clerkLoaded } from './auth'
import ClerkSessionBridge from './components/ClerkSessionBridge.vue'

const hasClerk = Boolean(import.meta.env.VITE_CLERK_PUBLISHABLE_KEY)
if (!hasClerk) {
  clerkLoaded.value = true
}
</script>

<template>
  <ClerkSessionBridge v-if="hasClerk" />
  <RouterView />
</template>
