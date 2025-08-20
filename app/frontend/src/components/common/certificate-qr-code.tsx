'use client'

import { useRef } from 'react'
import { QrCode, Copy, Download, ArrowRight } from 'lucide-react'
import { QRCodeSVG } from 'qrcode.react'
import { Button } from '../ui/button'
import { Input } from '../ui/input'
import { Dialog, DialogTrigger, DialogContent, DialogHeader, DialogTitle, DialogDescription } from '../ui/dialog'
import { showMessage } from '@/lib/utils/common'
import { Label } from '../ui/label'
import UseBreakpoint from '@/lib/hooks/use-breakpoint'
import Link from 'next/link'

interface QRCodeDialogProps {
  id: string
  isIcon?: boolean
}

const QRCodeDialog: React.FC<QRCodeDialogProps> = (props) => {
  const qrRef = useRef<SVGSVGElement>(null)
  const { md } = UseBreakpoint()
  // Generate the URL for QR code
  const certificateUrl = `${window.location.origin}?code=${props.id}`

  const copyToClipboard = async (copyText: string) => {
    try {
      await navigator.clipboard.writeText(copyText)
      showMessage('Đã sao chép vào clipboard', {
        duration: 1000
      })
    } catch {
      showMessage('Không thể sao chép', {
        duration: 1000
      })
    }
  }

  const downloadQRCode = () => {
    if (qrRef.current) {
      const svg = qrRef.current
      const canvas = document.createElement('canvas')
      const ctx = canvas.getContext('2d')
      const img = new Image()

      // Convert SVG to data URL
      const svgData = new XMLSerializer().serializeToString(svg)
      const svgBlob = new Blob([svgData], { type: 'image/svg+xml;charset=utf-8' })
      const svgUrl = URL.createObjectURL(svgBlob)

      img.onload = () => {
        canvas.width = 250
        canvas.height = 250

        // Add white background
        if (ctx) {
          ctx.fillStyle = 'white'
          ctx.fillRect(0, 0, canvas.width, canvas.height)
          ctx.drawImage(img, 25, 25, 200, 200)
        }

        // Download
        const link = document.createElement('a')
        link.download = `certificate-qr-${props.id}.png`
        link.href = canvas.toDataURL('image/png')
        link.click()

        URL.revokeObjectURL(svgUrl)
        showMessage('Mã QR đã được tải xuống', {
          duration: 1000
        })
      }

      img.src = svgUrl
    }
  }

  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button
          size={props.isIcon ? 'icon' : md ? 'default' : 'icon'}
          variant={props.isIcon ? 'outline' : undefined}
          title='Hiển thị mã QR'
        >
          <QrCode /> {!props.isIcon && md ? 'Mã QR' : ''}
        </Button>
      </DialogTrigger>
      <DialogContent className='sm:max-w-md'>
        <DialogHeader>
          <DialogTitle>Mã QR Chứng chỉ</DialogTitle>
          <DialogDescription>Quét mã QR để truy cập thông tin chứng chỉ</DialogDescription>
        </DialogHeader>
        <div className='flex flex-col items-center space-y-4 p-4'>
          {/* QR Code Display */}
          <div className='rounded-lg border bg-white p-4 shadow-sm'>
            <QRCodeSVG ref={qrRef} value={certificateUrl} size={200} level='M' />
          </div>

          <div className='flex flex-col gap-2 sm:flex-row'>
            <Button onClick={downloadQRCode}>
              <Download />
              Tải xuống
            </Button>
            <Link href={certificateUrl}>
              <Button variant='outline'>
                <ArrowRight />
                Chuyển hướng
              </Button>
            </Link>
          </div>

          {/* URL Input with Copy Button */}
          <div className='w-full'>
            <Label>ID chứng chỉ</Label>
            <div className='mt-1 flex gap-2'>
              <Input value={props.id} readOnly className='flex-1 text-xs' />
              <Button size='icon' variant='outline' onClick={() => copyToClipboard(props.id)}>
                <Copy className='h-4 w-4' />
              </Button>
            </div>
          </div>
          <div className='w-full'>
            <Label>Đường dẫn chứng chỉ</Label>
            <div className='mt-1 flex gap-2'>
              <Input value={certificateUrl} readOnly className='flex-1 text-xs' />
              <Button size='icon' variant='outline' onClick={() => copyToClipboard(certificateUrl)}>
                <Copy className='h-4 w-4' />
              </Button>
            </div>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  )
}

export default QRCodeDialog
