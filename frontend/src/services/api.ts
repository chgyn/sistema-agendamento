import axios from 'axios'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || '/api/v1',
  headers: {
    'Content-Type': 'application/json',
  },
})

// Interceptor para injetar JWT Token nas requisições autenticadas
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Interceptor para tratar 401 (token expirado)
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Se não for a rota de login ou registro, limpa o token
      if (!window.location.pathname.includes('/login') && !window.location.pathname.includes('/register') && !window.location.pathname.includes('/agendamento')) {
        localStorage.removeItem('token')
        localStorage.removeItem('user')
        localStorage.removeItem('tenant')
        window.location.href = '/login'
      }
    }
    return Promise.reject(error)
  }
)

export default api
