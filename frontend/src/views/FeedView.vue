<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { api, postCursor } from '../api'
import { clearToken } from '../auth'
import PostItem from '../components/PostItem.vue'
import { isCriticMuted, setCriticMuted } from '../preferences'
import type { Post, User } from '../types'

const PAGE_SIZE = 20

const router = useRouter()
const me = ref<User | null>(null)
const posts = ref<Post[]>([])
const draft = ref('')
const error = ref('')
const loading = ref(true)
const posting = ref(false)
const loadingMore = ref(false)
const hasMore = ref(true)
const muteCritic = ref(isCriticMuted())
const refreshTick = ref(0)
let pollTimer: ReturnType<typeof setInterval> | undefined

async function loadInitial() {
  error.value = ''
  try {
    const [user, feed] = await Promise.all([api.me(), api.posts({ limit: PAGE_SIZE })])
    me.value = user
    posts.value = feed || []
    hasMore.value = (feed || []).length >= PAGE_SIZE
  } catch (e) {
    const message = e instanceof Error ? e.message : 'Failed to load feed'
    error.value = message
    if (message.toLowerCase().includes('unauthorized')) {
      clearToken()
      await router.push('/login')
    }
  } finally {
    loading.value = false
  }
}

async function loadMore() {
  if (loadingMore.value || !hasMore.value || !posts.value.length) return
  loadingMore.value = true
  error.value = ''
  try {
    const last = posts.value[posts.value.length - 1]
    const page = (await api.posts({ limit: PAGE_SIZE, before: postCursor(last) })) || []
    posts.value = [...posts.value, ...page]
    hasMore.value = page.length >= PAGE_SIZE
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Failed to load more'
  } finally {
    loadingMore.value = false
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
    loading.value = true
    await loadInitial()
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Failed to post'
  } finally {
    posting.value = false
  }
}

async function signOut() {
  clearToken()
  await router.push('/login')
}

function toggleMuteCritic() {
  muteCritic.value = !muteCritic.value
  setCriticMuted(muteCritic.value)
}

onMounted(async () => {
  await loadInitial()
  // Refresh comments so @vibe_critic replies appear without a full remount.
  pollTimer = setInterval(() => {
    refreshTick.value += 1
  }, 8000)
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
        <button class="btn btn-ghost" type="button" @click="toggleMuteCritic">
          {{ muteCritic ? 'Show critic' : 'Mute critic' }}
        </button>
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
        :refresh-tick="refreshTick"
        :mute-critic="muteCritic"
        @changed="loadInitial"
      />
      <p v-if="!posts.length" class="muted" style="padding: 1rem 0">No posts yet — go first.</p>
      <div v-if="posts.length && hasMore" class="row" style="justify-content: center; padding: 1rem 0">
        <button class="btn btn-ghost" type="button" :disabled="loadingMore" @click="loadMore">
          {{ loadingMore ? 'Loading…' : 'Load more' }}
        </button>
      </div>
    </section>
  </div>
</template>
