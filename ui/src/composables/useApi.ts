import { ShowbizClient } from '@showbiz/sdk'

const client = new ShowbizClient({
  baseURL: '',
  onTokenRefresh: (tokens) => {
    localStorage.setItem('accessToken', tokens.accessToken)
    localStorage.setItem('refreshToken', tokens.refreshToken)
  },
})

const savedAccessToken = localStorage.getItem('accessToken')
const savedRefreshToken = localStorage.getItem('refreshToken')
if (savedAccessToken && savedRefreshToken) {
  client.setTokens(savedAccessToken, savedRefreshToken)
} else if (savedAccessToken) {
  client.setTokens(savedAccessToken, '')
}

export function useApi() {
  return client
}
