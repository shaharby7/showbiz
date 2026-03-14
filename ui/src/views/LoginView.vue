<script setup lang="ts">
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import InputText from 'primevue/inputtext'
import Password from 'primevue/password'
import Button from 'primevue/button'
import Message from 'primevue/message'
import { useAuthStore } from '@/stores/auth'
import { useOrganizationStore } from '@/stores/organization'
import { useProjectStore } from '@/stores/project'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const orgStore = useOrganizationStore()
const projectStore = useProjectStore()

const email = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)

const successMessage = ref(
  route.query.registered === 'true' ? 'Account created successfully. Please log in.' : '',
)

async function handleLogin() {
  error.value = ''
  successMessage.value = ''

  if (!email.value || !password.value) {
    error.value = 'Please enter both email and password.'
    return
  }

  loading.value = true
  try {
    await authStore.login(email.value, password.value)
    try {
      await orgStore.initialize()
      if (orgStore.currentOrg) {
        await projectStore.fetchProjects(orgStore.currentOrg.id)
        await projectStore.initialize()
      }
    } catch {
      // org/project init may fail if none exist yet
    }
    router.push({ name: 'dashboard' })
  } catch (err: unknown) {
    if (err instanceof Error) {
      error.value = err.message || 'Invalid credentials. Please try again.'
    } else {
      error.value = 'Invalid credentials. Please try again.'
    }
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="flex items-center justify-center min-h-screen bg-surface-ground px-4">
    <div class="w-full max-w-md">
      <div class="text-center mb-8">
        <h1 class="text-4xl font-bold text-primary mb-2">Showbiz</h1>
        <p class="text-muted-color">Sign in to your account</p>
      </div>

      <div class="bg-surface-card rounded-xl shadow-md p-8 border border-surface-border">
        <Message v-if="successMessage" severity="success" class="mb-4">
          {{ successMessage }}
        </Message>
        <Message v-if="error" severity="error" class="mb-4">
          {{ error }}
        </Message>

        <form class="flex flex-col gap-5" @submit.prevent="handleLogin">
          <div class="flex flex-col gap-2">
            <label for="email" class="text-sm font-medium text-color">Email</label>
            <InputText
              id="email"
              v-model="email"
              type="email"
              placeholder="you@example.com"
              class="w-full"
              autocomplete="email"
            />
          </div>

          <div class="flex flex-col gap-2">
            <label for="password" class="text-sm font-medium text-color">Password</label>
            <Password
              v-model="password"
              input-id="password"
              :feedback="false"
              toggle-mask
              class="w-full"
              input-class="w-full"
              autocomplete="current-password"
            />
          </div>

          <Button
            type="submit"
            label="Login"
            icon="pi pi-sign-in"
            class="w-full"
            :loading="loading"
          />
        </form>

        <div class="text-center mt-6 text-sm text-muted-color">
          Don't have an account?
          <router-link to="/register" class="text-primary font-medium hover:underline">
            Create one
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>
