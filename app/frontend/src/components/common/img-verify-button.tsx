import { Camera, ImageIcon, SwitchCamera } from 'lucide-react'
import { Button } from '../ui/button'
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from '../ui/dropdown-menu'
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle } from '../ui/dialog'
import { useRef, useState, useEffect } from 'react'
import Image from 'next/image'
import { showMessage } from '@/lib/utils/common'
import QrScanner from 'qr-scanner'

interface ImgVerifyButtonProps {
  // eslint-disable-next-line no-unused-vars
  onCodeDetected?: (code: string) => void
}

const ImgVerifyButton: React.FC<ImgVerifyButtonProps> = ({ onCodeDetected }) => {
  const fileInputRef = useRef<HTMLInputElement>(null)
  const videoRef = useRef<HTMLVideoElement>(null)
  const canvasRef = useRef<HTMLCanvasElement>(null)
  const [isCameraOpen, setIsCameraOpen] = useState(false)
  const [capturedImage, setCapturedImage] = useState<string | null>(null)
  const [stream, setStream] = useState<MediaStream | null>(null)
  const [facingMode, setFacingMode] = useState<'user' | 'environment'>('environment')
  const [isMobile, setIsMobile] = useState(false)

  // Detect if user is on mobile device
  useEffect(() => {
    const checkMobile = () => {
      const userAgent = navigator.userAgent.toLowerCase()
      const mobileKeywords = ['android', 'webos', 'iphone', 'ipad', 'ipod', 'blackberry', 'windows phone']
      const isMobileDevice = mobileKeywords.some((keyword) => userAgent.includes(keyword))
      setIsMobile(isMobileDevice)
      // Set default camera based on device type
      setFacingMode(isMobileDevice ? 'environment' : 'user')
    }
    checkMobile()
  }, [])

  const handleOpenCamera = async () => {
    // Mở camera thiết bị
    if (navigator.mediaDevices && navigator.mediaDevices.getUserMedia) {
      try {
        const mediaStream = await navigator.mediaDevices.getUserMedia({
          video: {
            width: { ideal: 1280 },
            height: { ideal: 720 },
            facingMode: facingMode // sử dụng camera trước hoặc sau tùy theo thiết bị
          }
        })

        setStream(mediaStream)
        setIsCameraOpen(true)

        // Đợi modal mở rồi mới gán stream cho video
        setTimeout(() => {
          if (videoRef.current) {
            videoRef.current.srcObject = mediaStream
            videoRef.current.play()
          }
        }, 100)
      } catch (error) {
        console.error('Error accessing camera:', error)

        showMessage('Không thể truy cập camera. Vui lòng kiểm tra quyền truy cập')
      }
    } else {
      showMessage('Trình duyệt không hỗ trợ camera')
    }
  }

  const handleSwitchCamera = async () => {
    // Dừng camera hiện tại
    if (stream) {
      stream.getTracks().forEach((track) => track.stop())
      setStream(null)
    }

    // Đổi hướng camera
    const newFacingMode = facingMode === 'user' ? 'environment' : 'user'
    setFacingMode(newFacingMode)

    // Mở camera với hướng mới
    try {
      const mediaStream = await navigator.mediaDevices.getUserMedia({
        video: {
          width: { ideal: 1280 },
          height: { ideal: 720 },
          facingMode: newFacingMode
        }
      })

      setStream(mediaStream)

      setTimeout(() => {
        if (videoRef.current) {
          videoRef.current.srcObject = mediaStream
          videoRef.current.play()
        }
      }, 100)
    } catch (error) {
      console.error('Error switching camera:', error)
      showMessage('Không thể chuyển đổi camera')
    }
  }

  const handleCloseCamera = () => {
    if (stream) {
      stream.getTracks().forEach((track) => track.stop())
      setStream(null)
    }
    setIsCameraOpen(false)
    setCapturedImage(null)
  }

  const handleCapturePhoto = () => {
    if (videoRef.current && canvasRef.current) {
      const video = videoRef.current
      const canvas = canvasRef.current
      const context = canvas.getContext('2d')

      if (context) {
        // Đặt kích thước canvas bằng với video
        canvas.width = video.videoWidth
        canvas.height = video.videoHeight

        // Vẽ frame hiện tại từ video lên canvas
        context.drawImage(video, 0, 0, canvas.width, canvas.height)

        // Chuyển canvas thành base64 image
        const imageDataUrl = canvas.toDataURL('image/jpeg', 0.8)
        setCapturedImage(imageDataUrl)
      }
    }
  }

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

  const handleQRCodeScan = async (imageSource: string | File) => {
    try {
      let result: string

      if (typeof imageSource === 'string') {
        // For base64 images (captured photos)
        result = await QrScanner.scanImage(imageSource)
      } else {
        // For file objects (selected images)
        result = await QrScanner.scanImage(imageSource)
      }

      // Check if the QR code contains a kma.edu.vn URL with code parameter
      if (result.includes('kma.edu.vn') && result.includes('code=')) {
        const code = extractCodeFromUrl(result)
        if (code) {
          onCodeDetected?.(code)
          showMessage('Đã quét thành công mã QR')
          handleCloseCamera()
          return
        }
      }

      showMessage('QR code không hợp lệ, vui lòng thử lại')
    } catch (error) {
      console.error('Error scanning QR code:', error)
      showMessage('Không tìm thấy QR code trong ảnh, vui lòng thử lại')
    }
  }

  const handleConfirmPhoto = () => {
    if (capturedImage) {
      // Scan for QR code in the captured image
      handleQRCodeScan(capturedImage)
    }
  }

  const handleRetakePhoto = () => {
    setCapturedImage(null)
    handleOpenCamera()
  }

  const handleSelectImage = () => {
    // Mở dialog chọn file
    fileInputRef.current?.click()
  }

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0]
    if (file) {
      // Check if selected file is an image
      if (file.type.startsWith('image/')) {
        // Scan for QR code in the selected image
        handleQRCodeScan(file)
      } else {
        showMessage('Vui lòng chọn file ảnh hợp lệ.')
      }

      // Reset the input value so the same file can be selected again
      event.target.value = ''
    }
  }

  return (
    <>
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button variant={'secondary'} className='flex items-center gap-2'>
            <Camera />
            Xác thực bằng hình ảnh
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

      {/* Camera Preview Modal */}
      <Dialog open={isCameraOpen} onOpenChange={handleCloseCamera}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Camera</DialogTitle>
          </DialogHeader>

          <div className='space-y-4'>
            {!capturedImage ? (
              // Camera Preview
              <>
                <div className='relative'>
                  <video
                    ref={videoRef}
                    autoPlay
                    playsInline
                    muted
                    className='h-96 w-full rounded-lg bg-black object-cover'
                  />
                  {/* Switch camera button for mobile devices */}
                  {isMobile && (
                    <div className='absolute right-4 top-4'>
                      <Button variant='secondary' size='icon' onClick={handleSwitchCamera} className='rounded-full'>
                        <SwitchCamera className='h-5 w-5' />
                      </Button>
                    </div>
                  )}
                </div>
                <DialogFooter>
                  <Button onClick={handleCapturePhoto} className='w-full'>
                    <Camera />
                    Chụp ảnh
                  </Button>
                </DialogFooter>
              </>
            ) : (
              // Captured Image Preview
              <div className='space-y-4'>
                <div className='relative h-96 w-full'>
                  <Image src={capturedImage} alt='Captured' fill className='rounded-lg object-cover' />
                </div>
                <DialogFooter className='flex flex-col justify-center gap-2 md:flex-row'>
                  <Button variant='outline' onClick={handleRetakePhoto} className='flex-1'>
                    <Camera />
                    Chụp lại
                  </Button>
                  <Button onClick={handleConfirmPhoto} className='flex-1'>
                    <ImageIcon />
                    Sử dụng ảnh này
                  </Button>
                </DialogFooter>
              </div>
            )}
          </div>
        </DialogContent>
      </Dialog>

      {/* Hidden canvas for image capture */}
      <canvas ref={canvasRef} className='hidden' />

      {/* Hidden file input */}
      <input ref={fileInputRef} type='file' accept='image/*' onChange={handleFileChange} className='hidden' />
    </>
  )
}

export default ImgVerifyButton
