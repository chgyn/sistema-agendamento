import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../services/api'

export interface PublicTenant {
  id: string
  slug: string
  name: string
  phone: string
  email: string
  address: string
  city: string
  state: string
  logo_url: string
  primary_color: string
}

export interface PublicService {
  id: string
  name: string
  description: string
  duration_minutes: number
  price: number
  color: string
}

export interface PublicProfessional {
  id: string
  name: string
  title: string
  specialty: string
  bio: string
  avatar_url: string
}

export interface TimeSlot {
  start_time: string
  end_time: string
  start_datetime: string
  end_datetime: string
  is_available: boolean
  professional_id: string
  professional_name: string
}

export const useBookingStore = defineStore('booking', () => {
  const currentStep = ref<number>(1) // 1: Serviço, 2: Profissional, 3: Data/Hora, 4: Dados do Cliente, 5: Sucesso
  const tenant = ref<PublicTenant | null>(null)
  const services = ref<PublicService[]>([])
  const professionals = ref<PublicProfessional[]>([])
  
  const selectedService = ref<PublicService | null>(null)
  const selectedProfessional = ref<PublicProfessional | null>(null) // null = Qualquer Profissional
  const selectedDate = ref<string>(new Date().toISOString().split('T')[0])
  const selectedSlot = ref<TimeSlot | null>(null)

  const availableSlots = ref<TimeSlot[]>([])
  const isLoadingSlots = ref<boolean>(false)
  const isSubmitting = ref<boolean>(false)
  const bookingError = ref<string | null>(null)
  const confirmedAppointment = ref<any | null>(null)

  async function loadTenant(slug: string) {
    bookingError.value = null
    try {
      const res = await api.get(`/public/${slug}`)
      if (res.data.success) {
        tenant.value = res.data.data
        await loadServices(slug)
      }
    } catch (err: any) {
      bookingError.value = err.response?.data?.error || 'Estabelecimento não encontrado'
    }
  }

  async function loadServices(slug: string) {
    try {
      const res = await api.get(`/public/${slug}/services`)
      if (res.data.success) {
        services.value = res.data.data
      }
    } catch (err) {
      console.error(err)
    }
  }

  async function loadProfessionalsForService(slug: string, serviceId: string) {
    try {
      const res = await api.get(`/public/${slug}/services/${serviceId}/professionals`)
      if (res.data.success) {
        professionals.value = res.data.data
      }
    } catch (err) {
      console.error(err)
    }
  }

  async function fetchAvailability(slug: string) {
    if (!selectedService.value || !selectedDate.value) return
    isLoadingSlots.value = true
    bookingError.value = null
    try {
      const proParam = selectedProfessional.value ? selectedProfessional.value.id : 'any'
      const res = await api.get(`/public/${slug}/availability`, {
        params: {
          service_id: selectedService.value.id,
          professional_id: proParam,
          date: selectedDate.value,
        },
      })
      if (res.data.success) {
        availableSlots.value = res.data.data.slots || []
      }
    } catch (err: any) {
      bookingError.value = err.response?.data?.error || 'Erro ao carregar horários disponíveis'
    } finally {
      isLoadingSlots.value = false
    }
  }

  async function createAppointment(slug: string, customerData: {
    name: string
    phone: string
    email?: string
    notes?: string
  }) {
    if (!selectedService.value || !selectedSlot.value) {
      throw new Error('Serviço ou horário não selecionado')
    }

    isSubmitting.value = true
    bookingError.value = null

    try {
      const payload = {
        service_id: selectedService.value.id,
        professional_id: selectedSlot.value.professional_id,
        start_at: selectedSlot.value.start_datetime,
        customer_name: customerData.name,
        customer_phone: customerData.phone,
        customer_email: customerData.email,
        notes: customerData.notes,
      }

      const res = await api.post(`/public/${slug}/appointments`, payload)
      if (res.data.success) {
        confirmedAppointment.value = res.data.data
        currentStep.value = 5 // Step final de confirmação
        return res.data.data
      }
    } catch (err: any) {
      const errorMsg = err.response?.data?.error || 'Falha ao agendar atendimento. Tente outro horário.'
      bookingError.value = errorMsg
      throw new Error(errorMsg)
    } finally {
      isSubmitting.value = false
    }
  }

  function resetBooking() {
    currentStep.value = 1
    selectedService.value = null
    selectedProfessional.value = null
    selectedSlot.value = null
    bookingError.value = null
    confirmedAppointment.value = null
  }

  return {
    currentStep,
    tenant,
    services,
    professionals,
    selectedService,
    selectedProfessional,
    selectedDate,
    selectedSlot,
    availableSlots,
    isLoadingSlots,
    isSubmitting,
    bookingError,
    confirmedAppointment,
    loadTenant,
    loadServices,
    loadProfessionalsForService,
    fetchAvailability,
    createAppointment,
    resetBooking,
  }
})
