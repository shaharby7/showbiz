<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import InputText from 'primevue/inputtext'
import Password from 'primevue/password'
import Button from 'primevue/button'
import Message from 'primevue/message'
import { useApi } from '@/composables/useApi'

const router = useRouter()
const api = useApi()

const email = ref('')
const password = ref('')
const confirmPassword = ref('')
const displayName = ref('')
const error = ref('')
const loading = ref(false)

const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

function validate(): string | null {
  if (!displayName.value.trim()) return 'Display name is required.'
  if (!email.value.trim()) return 'Email is required.'
  if (!emailRegex.test(email.value)) return 'Please enter a valid email address.'
  if (password.value.length < 8) return 'Password must be at least 8 characters.'
  if (password.value !== confirmPassword.value) return 'Passwords do not match.'
  return null
}

async function handleRegister() {
  error.value = ''
  const validationError = validate()
  if (validationError) {
    error.value = validationError
    return
  }

  loading.value = true
  try {
    await api.auth.register({
      email: email.value,
      password: password.value,
      displayName: displayName.value,
    })
    router.push({ name: 'login', query: { registered: 'true' } })
  } catch (err: unknown) {
    if (err instanceof Error) {
      error.value = err.message || 'Registration failed. Please try again.'
    } else {
      error.value = 'Registration failed. Please try again.'
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
        <img src="/showbiz.svg" alt="Showbiz" class="h-12 w-12 mx-auto mb-3" />
        <h1 class="text-4xl font-bold text-primary mb-2">Showbiz</h1>
        <p class="text-muted-color">Create your account</p>
      </div>

      <div class="bg-surface-card rounded-xl shadow-md p-8 border border-surface-border">
        <Message v-if="error" severity="error" class="mb-4">
          {{ error }}
        </Message>

        <form class="flex flex-col gap-5" @submit.prevent="handleRegister">
          <div class="flex flex-col gap-2">
            <label for="displayName" class="text-sm font-medium text-color">Display Name</label>
            <InputText
              id="displayName"
              v-model="displayName"
              placeholder="Jane Smith"
              class="w-full"
              autocomplete="name"
            />
          </div>

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
              toggle-mask
              class="w-full"
              input-class="w-full"
              autocomplete="new-password"
            />
          </div>

          <div class="flex flex-col gap-2">
            <label for="confirmPassword" class="text-sm font-medium text-color">
              Confirm Password
            </label>
            <Password
              v-model="confirmPassword"
              input-id="confirmPassword"
              :feedback="false"
              toggle-mask
              class="w-full"
              input-class="w-full"
              autocomplete="new-password"
            />
          </div>

          <Button
            type="submit"
            label="Create Account"
            icon="pi pi-user-plus"
            class="w-full"
            :loading="loading"
          />
        </form>

        <div class="text-center mt-6 text-sm text-muted-color">
          Already have an account?
          <router-link to="/login" class="text-primary font-medium hover:underline">
            Sign in
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>
