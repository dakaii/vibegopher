<script setup lang="ts">
import { useAuth, useClerk } from '@clerk/vue'
import { watchEffect } from 'vue'
import { bindClerkSession, setClerkSignOut } from '../auth'

const { isLoaded, isSignedIn, getToken } = useAuth()
const clerk = useClerk()

watchEffect(() => {
  bindClerkSession({
    isLoaded: Boolean(isLoaded.value),
    isSignedIn: Boolean(isSignedIn.value),
    getToken: async () => {
      const fn = getToken.value
      if (!fn) return null
      return fn()
    },
  })
  setClerkSignOut(async () => {
    if (clerk.value) {
      await clerk.value.signOut()
    }
  })
})
</script>

<template></template>
