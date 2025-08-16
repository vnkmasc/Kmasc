import { getSignDegreeConfig } from './handle-storage'

interface SignatureResponse {
  status: number
  signature?: string
  message?: string
  validationStatus?: {
    status: number
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
    console.log('🔍 DEBUG: Bắt đầu kết nối WebSocket đến:', endpoint)

    if (!('WebSocket' in window)) {
      const error = 'WebSocket không được hỗ trợ bởi trình duyệt này!'
      console.log('❌ ERROR:', error)
      reject(new Error(error))
      return
    }

    try {
      const jsonParams = JSON.stringify(params)
      console.log('🚀 DEBUG: Tham số gửi đi:', jsonParams)

      console.log('🔄 DEBUG: Đang tạo kết nối WebSocket...')
      const ws = new WebSocket(endpoint)

      ws.onopen = () => {
        console.log('✅ DEBUG: Kết nối WebSocket đã mở, đang gửi dữ liệu...')
        ws.send(jsonParams)
      }

      ws.onmessage = (evt) => {
        console.log('📩 DEBUG: Đã nhận phản hồi từ WebSocket:', evt.data)
        try {
          const response = JSON.parse(evt.data)
          console.log('🔍 DEBUG: Phản hồi đã parse:', response)
          resolve(handleResponse(response))
        } catch (error) {
          console.log('❌ DEBUG: Lỗi xử lý dữ liệu từ plugin:', error)
          reject(error)
        } finally {
          console.log('🔒 DEBUG: Đóng kết nối WebSocket')
          ws.close()
        }
      }

      ws.onclose = (event) => {
        console.log('🔒 DEBUG: Kết nối WebSocket đã đóng với mã:', event.code, 'lý do:', event.reason)
      }

      ws.onerror = (error) => {
        console.log('❌ DEBUG: Lỗi WebSocket:', error)
        reject(error)
      }

      // Thêm timeout để phát hiện kết nối không thành công
      setTimeout(() => {
        if (ws.readyState !== WebSocket.OPEN) {
          console.log('⏱️ DEBUG: Kết nối WebSocket timeout sau 5 giây')
          ws.close()
          reject(new Error('Kết nối WebSocket timeout - Plugin ký số có thể chưa được cài đặt hoặc không chạy'))
        }
      }, 5000)
    } catch (error) {
      console.log('❌ DEBUG: Lỗi khởi tạo kết nối:', error)
      reject(error)
    }
  })
}

export const signDigitalSignature = (
  hashValue: string,
  wsSign: WebSocketEndpoint = DEFAULT_SIGN_ENDPOINT
): Promise<string> => {
  console.log('🚀 DEBUG: Bắt đầu quá trình ký số với hash:', hashValue)

  // Kiểm tra plugin ký số đã được cài đặt chưa
  console.log('🔍 DEBUG: Kiểm tra plugin ký số...')

  const params = {
    HashValue: hashValue,
    HashAlg: 'SHA256'
  }
  console.log('🚀 DEBUG: Tham số ký số:', params)

  return createWebSocketConnection<string>(params, wsSign, (response: SignatureResponse) => {
    console.log('🚀 DEBUG: Kết quả ký số:', response)
    if (response.status === 0 && response.signature) {
      console.log('✅ DEBUG: Ký số thành công')
      return response.signature
    }
    const errorMsg = response.message || 'Lỗi không xác định khi ký số'
    console.log('❌ DEBUG: Ký số thất bại:', errorMsg)
    throw new Error(errorMsg)
  })
}

export const verifyDigitalSignature = (
  signature: string,
  hashValue: string,
  wsVerify: WebSocketEndpoint = DEFAULT_VERIFY_ENDPOINT
): Promise<boolean> => {
  console.log('🚀 DEBUG: Bắt đầu quá trình xác thực chữ ký')

  const params = {
    Signature: signature,
    Base64Content: hashValue
  }
  console.log('🚀 DEBUG: Tham số xác thực:', params)

  return createWebSocketConnection<boolean>(params, wsVerify, (response: SignatureResponse) => {
    console.log('🚀 DEBUG: Kết quả xác thực:', response)
    if (response.validationStatus?.status === 0) {
      console.log('✅ DEBUG: Xác thực thành công')
      return true
    }
    const errorMsg = response.message || 'Xác thực chữ ký không thành công'
    console.log('❌ DEBUG: Xác thực thất bại:', errorMsg)
    return false
  })
}

// Hàm kiểm tra plugin ký số đã được cài đặt chưa
export const checkVGCAPluginInstalled = (): Promise<boolean> => {
  console.log('🔍 DEBUG: Kiểm tra plugin VGCA...')

  return new Promise((resolve) => {
    try {
      const testWs = new WebSocket('wss://127.0.0.1:8987/GetVersion')

      testWs.onopen = () => {
        console.log('✅ DEBUG: Plugin VGCA đã được cài đặt và đang chạy')
        testWs.send('GetVersion')
        testWs.close()
        resolve(true)
      }

      testWs.onerror = () => {
        console.log('❌ DEBUG: Plugin VGCA chưa được cài đặt hoặc không chạy')
        resolve(false)
      }

      setTimeout(() => {
        if (testWs.readyState !== WebSocket.OPEN) {
          console.log('⏱️ DEBUG: Kiểm tra plugin VGCA timeout')
          testWs.close()
          resolve(false)
        }
      }, 2000)
    } catch (error) {
      console.log('❌ DEBUG: Lỗi kiểm tra plugin VGCA:', error)
      resolve(false)
    }
  })
}
