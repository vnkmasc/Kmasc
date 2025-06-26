import { getSession } from '../auth/session'

const BASE_URL = process.env.NEXT_PUBLIC_API_URL || ''

const defaultHeaders = {
  'Content-Type': 'application/json',

  Accept: 'application/json'
}

const apiService = async (
  method: 'POST' | 'PUT' | 'GET' | 'DELETE' | 'PATCH',
  url: string,
  data?: Record<string, string | number | null | undefined> | FormData,
  isAuth: boolean = true,
  headers?: HeadersInit,
  isBlob?: boolean
) => {
  const fullUrl = `${BASE_URL}/api/v1/${url}`

  const isFormData = data instanceof FormData

  const headersAPI = {
    ...(isFormData ? {} : defaultHeaders),
    ...headers,
    ...(isAuth ? { Authorization: `Bearer ${(await getSession())?.access_token}` } : {})
  }

  if (isFormData && headersAPI['Content-Type']) {
    delete headersAPI['Content-Type']
  }

  try {
    const response = await fetch(fullUrl, {
      method,
      headers: headersAPI,
      body: data ? (isFormData ? data : JSON.stringify(data)) : null
    })

    console.log(method + ' ' + fullUrl)

    const contentType = response.headers.get('content-type')

    if (!response.ok) {
      let errorData
      let errorMessage = `API Error: ${response.status} - ${response.statusText}`

      if (contentType && contentType.includes('application/json')) {
        try {
          errorData = await response.json()
          console.error('Error response:', errorData)

          if (errorData.error) {
            errorMessage = errorData.error
          } else if (errorData.message) {
            errorMessage = errorData.message
          }
        } catch (parseError) {
          console.error('Failed to parse JSON error response:', parseError)
        }
      } else {
        try {
          const errorText = await response.text()
          console.error('Error response (text):', errorText)
          errorMessage = errorText || errorMessage
        } catch (parseError) {
          console.error('Failed to read text error response:', parseError)
        }
      }

      const error = new Error(errorMessage) as any
      error.status = response.status
      error.data = errorData
      throw error
    }

    if (contentType && contentType.includes('application/json')) {
      return await response.json()
    } else if (isBlob) {
      return await response.blob()
    } else {
      const text = await response.text()
      console.warn('Response is not JSON:', text)
      return { error: 'Response is not JSON', data: text }
    }
  } catch (error) {
    console.error('API Service Error:', error)
    throw error
  }
}

export default apiService
