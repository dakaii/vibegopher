<script setup>
import { onMounted, onUnmounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../api'
import { clearToken } from '../auth'
import PostItem from '../components/PostItem.vue'

const router = useRouter()
const me = ref(null)
const posts = ref([])
const draft = ref('')
const error = ref('')
const loading = ref(true)
const posting = ref(false)
let pollTimer

async function load() {
  error.value = ''
  try {
    const [user, feed] = await Promise.all([api.me(), api.posts()])
    me.value = user
    posts.value = feed || []
  } catch (e) {
    error.value = e.message
    if (String(e.message).toLowerCase().includes('unauthorized')) {
      clearToken()
      router.push('/login')
    }
  } finally {
    loading.value = false
  }
}

async function submitPost() {
  const content = draft.value.trim()
  if (!content || posting.value) return
  posting.value = true
  error.value = ''
  try {
    await api.createPost(content)
    draft.value = ''
    await load()
  } catch (e) {
    error.value = e.message
  } finally {
    posting.value = false
  }
}

function signOut() {
  clearToken()
  router.push('/login')
}

onMounted(async () => {
  await load()
  // Poll so vibe_critic replies appear without a manual refresh.
  pollTimer = setInterval(load, 8000)
})

onUnmounted(() => {
  if (pollTimer) clearInterval(pollTimer)
})
</script>

<template>
  <div class="app-shell">
    <header class="row">
      <div>
        <h1 class="brand">Vibe<span>Gopher</span></h1>
        <p class="lede">Post freely. @vibe_critic will argue with the idea, not erase it.</p>
      </div>
      <div class="row" style="gap: 0.75rem">
        <span v-if="me" class="muted">@{{ me.username }}</span>
        <button class="btn btn-ghost" type="button" @click="signOut">Sign out</button>
      </div>
    </header>

    <section class="panel composer">
      <textarea
        v-model="draft"
        maxlength="280"
        placeholder="What's the take? Links welcome — the critic reads them."
      />
      <div class="row" style="margin-top: 0.85rem">
        <span class="muted">{{ draft.length }}/280</span>
        <button class="btn" type="button" :disabled="!draft.trim() || posting" @click="submitPost">
          {{ posting ? 'Posting…' : 'Post' }}
        </button>
      </div>
    </section>

    <p v-if="error" class="error">{{ error }}</p>
    <p v-if="loading" class="muted">Loading feed…</p>

    <section v-else class="panel" style="margin-top: 1rem; padding-top: 0.35rem">
      <PostItem
        v-for="post in posts"
        :key="post.id"
        :post="post"
        @changed="load"
      />
      <p v-if="!posts.length" class="muted" style="padding: 1rem 0">No posts yet — go first.</p>
    </section>
  </div>
</template>
