import { getSignDegreeConfig } from './handle-storage'

interface SignatureResponse {
  Status: number
  Signature?: string
  Message?: string
  ValidationStatus?: {
    Status: number
  }
}

type WebSocketEndpoint = string

const signDegreeConfig = getSignDegreeConfig()
const DEFAULT_SIGN_ENDPOINT = signDegreeConfig.signService
const DEFAULT_VERIFY_ENDPOINT = signDegreeConfig.verifyService

const createWebSocketConnection = <T>(
  params: object,
  endpoint: WebSocketEndpoint,
  // eslint-disable-next-line no-unused-vars
  handleResponse: (data: any) => T
): Promise<T> => {
  return new Promise((resolve, reject) => {
    if (!('WebSocket' in window)) {
      reject(new Error('WebSocket không được hỗ trợ bởi trình duyệt này!'))
      return
    }

    try {
      const jsonParams = JSON.stringify(params)
      const ws = new WebSocket(endpoint)

      ws.onopen = () => {
        ws.send(jsonParams)
      }

      ws.onmessage = (evt) => {
        try {
          const response = JSON.parse(evt.data)
          resolve(handleResponse(response))
        } catch (error) {
          reject(error)
        } finally {
          ws.close()
        }
      }

      ws.onclose = () => {}

      ws.onerror = (error) => {
        reject(error)
      }

      setTimeout(() => {
        if (ws.readyState !== WebSocket.OPEN) {
          ws.close()
          reject(new Error('Kết nối WebSocket timeout - Plugin ký số có thể chưa được cài đặt hoặc không chạy'))
        }
      }, 5000)
    } catch (error) {
      reject(error)
    }
  })
}

export const signDigitalSignature = (
  hashValue: string,
  wsSign: WebSocketEndpoint = DEFAULT_SIGN_ENDPOINT
): Promise<string> => {
  const params = {
    HashValue: hashValue,
    HashAlg: 'SHA256'
  }

  return createWebSocketConnection<string>(params, wsSign, (response: SignatureResponse) => {
    if (response.Status === 0 && response.Signature) {
      return response.Signature
    }
    const errorMsg = response.Message || 'Lỗi không xác định khi ký số'
    throw new Error(errorMsg)
  })
}

export const verifyDigitalSignature = (
  signature: string,
  hashValue: string,
  wsVerify: WebSocketEndpoint = DEFAULT_VERIFY_ENDPOINT
): Promise<boolean> => {
  const params = {
    Signature: signature,
    Base64Content: hashValue
  }

  return createWebSocketConnection<boolean>(params, wsVerify, (response: SignatureResponse) => {
    if (response.ValidationStatus?.Status === 0) {
      return true
    }
    return false
  })
}
