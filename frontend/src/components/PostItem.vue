<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { api } from '../api'
import type { Comment, Post } from '../types'

const props = defineProps<{
  post: Post
}>()

const emit = defineEmits<{
  changed: []
}>()

const comments = ref<Comment[]>([])
const draft = ref('')
const error = ref('')
const busy = ref(false)

async function loadComments() {
  try {
    comments.value = (await api.comments(props.post.id)) || []
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Failed to load comments'
  }
}

async function submitComment() {
  const content = draft.value.trim()
  if (!content || busy.value) return
  busy.value = true
  error.value = ''
  try {
    await api.createComment(props.post.id, content)
    draft.value = ''
    await loadComments()
    emit('changed')
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Failed to comment'
  } finally {
    busy.value = false
  }
}

function formatTime(value?: string) {
  if (!value) return ''
  try {
    return new Date(value).toLocaleString()
  } catch {
    return value
  }
}

onMounted(() => {
  void loadComments()
})
</script>

<template>
  <article class="post">
    <div class="row" style="justify-content: flex-start; gap: 0.55rem">
      <span class="post-author" :class="{ bot: post.user?.is_bot }">
        @{{ post.user?.username || 'unknown' }}
      </span>
      <span class="post-meta">{{ formatTime(post.created_at) }}</span>
    </div>
    <p class="post-body">{{ post.content }}</p>

    <div class="comments">
      <div v-for="c in comments" :key="c.id" class="comment">
        <strong
          class="post-author"
          :class="{ bot: c.user?.is_bot || c.user?.username === 'vibe_critic' }"
        >
          @{{ c.user?.username || 'unknown' }}
        </strong>
        <span class="post-meta"> · {{ formatTime(c.created_at) }}</span>
        <div>{{ c.content }}</div>
      </div>
    </div>

    <div class="comment-box" style="margin-top: 0.85rem">
      <textarea v-model="draft" maxlength="280" rows="2" placeholder="Reply…" />
      <div class="row" style="margin-top: 0.55rem">
        <span class="muted">{{ draft.length }}/280</span>
        <button
          class="btn btn-ghost"
          type="button"
          :disabled="!draft.trim() || busy"
          @click="submitComment"
        >
          Reply
        </button>
      </div>
      <p v-if="error" class="error">{{ error }}</p>
    </div>
  </article>
</template>
