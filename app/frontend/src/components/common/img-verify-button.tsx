import { Camera, ImageIcon, SwitchCamera } from 'lucide-react'
import { Button } from '../ui/button'
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from '../ui/dropdown-menu'
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle } from '../ui/dialog'
import { useRef, useState, useEffect } from 'react'
import { showMessage } from '@/lib/utils/common'
import QrScanner from 'qr-scanner'

interface ImgVerifyButtonProps {
  // eslint-disable-next-line no-unused-vars
  onCodeDetected?: (code: string) => void
}

const ImgVerifyButton: React.FC<ImgVerifyButtonProps> = ({ onCodeDetected }) => {
  const fileInputRef = useRef<HTMLInputElement>(null)
  const videoRef = useRef<HTMLVideoElement>(null)
  const scannerRef = useRef<QrScanner | null>(null)
  const [isCameraOpen, setIsCameraOpen] = useState(false)
  const [stream, setStream] = useState<MediaStream | null>(null)
  const [facingMode, setFacingMode] = useState<'user' | 'environment'>('environment')
  const [isMobile, setIsMobile] = useState(false)

  // Detect mobile device
  useEffect(() => {
    const userAgent = navigator.userAgent.toLowerCase()
    const mobileKeywords = ['android', 'webos', 'iphone', 'ipad', 'ipod', 'blackberry', 'windows phone']
    const isMobileDevice = mobileKeywords.some((k) => userAgent.includes(k))
    setIsMobile(isMobileDevice)
    setFacingMode(isMobileDevice ? 'environment' : 'user')
  }, [])

  const extractCodeFromUrl = (url: string): string | null => {
    try {
      const urlObj = new URL(url)
      const code = urlObj.searchParams.get('code')
      return code
    } catch (error) {
      console.error('Error parsing URL:', error)
      return null
    }
  }

  const handleOpenCamera = async () => {
    try {
      const mediaStream = await navigator.mediaDevices.getUserMedia({
        video: {
          width: { ideal: 1920 },
          height: { ideal: 1080 },
          facingMode: facingMode
        }
      })

      setStream(mediaStream)
      setIsCameraOpen(true)

      setTimeout(() => {
        if (videoRef.current) {
          videoRef.current.srcObject = mediaStream

          // Khởi tạo scanner
          scannerRef.current = new QrScanner(
            videoRef.current,
            (result) => {
              if (result) {
                const code = extractCodeFromUrl(result.data)
                if (code) {
                  onCodeDetected?.(code)
                  showMessage('Đã quét thành công QR code')
                } else {
                  showMessage('QR code không hợp lệ')
                }
                handleCloseCamera()
              }
            },
            {
              highlightScanRegion: true,
              highlightCodeOutline: true,
              maxScansPerSecond: 10 // scan nhiều frame/giây
            }
          )

          scannerRef.current.start()
        }
      }, 200)
    } catch (error) {
      console.error('Error accessing camera:', error)
      showMessage('Không thể truy cập camera. Vui lòng kiểm tra quyền.')
    }
  }

  const handleSwitchCamera = async () => {
    if (stream) {
      stream.getTracks().forEach((track) => track.stop())
      setStream(null)
    }
    if (scannerRef.current) {
      scannerRef.current.stop()
      scannerRef.current.destroy()
      scannerRef.current = null
    }

    const newFacingMode = facingMode === 'user' ? 'environment' : 'user'
    setFacingMode(newFacingMode)
    handleOpenCamera()
  }

  const handleCloseCamera = () => {
    if (stream) {
      stream.getTracks().forEach((track) => track.stop())
      setStream(null)
    }
    if (scannerRef.current) {
      scannerRef.current.stop()
      scannerRef.current.destroy()
      scannerRef.current = null
    }
    setIsCameraOpen(false)
  }

  const handleSelectImage = () => {
    fileInputRef.current?.click()
  }

  const handleFileChange = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0]
    if (file) {
      try {
        const result = await QrScanner.scanImage(file)
        if (result) {
          onCodeDetected?.(result)
          showMessage('Đã quét thành công QR code từ ảnh')
        } else {
          showMessage('Không tìm thấy QR code trong ảnh')
        }
      } catch (err) {
        console.error('Error scanning QR from file:', err)
        showMessage('Không tìm thấy QR code trong ảnh')
      }
    }
    event.target.value = ''
  }

  return (
    <>
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button variant='secondary' className='flex items-center gap-2'>
            <Camera />
            Quét QR
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent>
          <DropdownMenuItem onClick={handleOpenCamera}>
            <Camera />
            Mở camera
          </DropdownMenuItem>
          <DropdownMenuItem onClick={handleSelectImage}>
            <ImageIcon />
            Chọn ảnh có sẵn
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>

      {/* Camera Modal */}
      <Dialog open={isCameraOpen} onOpenChange={handleCloseCamera}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Quét QR Code</DialogTitle>
          </DialogHeader>
          <div className='relative'>
            <video ref={videoRef} autoPlay playsInline muted className='h-96 w-full rounded-lg bg-black object-cover' />
            {isMobile && (
              <div className='absolute right-4 top-4'>
                <Button variant='secondary' size='icon' onClick={handleSwitchCamera} className='rounded-full'>
                  <SwitchCamera className='h-5 w-5' />
                </Button>
              </div>
            )}
          </div>
          <DialogFooter>
            <Button variant='outline' onClick={handleCloseCamera} className='w-full'>
              Đóng
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Hidden file input */}
      <input ref={fileInputRef} type='file' accept='image/*' onChange={handleFileChange} className='hidden' />
    </>
  )
}

export default ImgVerifyButton
