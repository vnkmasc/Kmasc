'use client'

import React, { useRef } from 'react'
import html2pdf from 'html2pdf.js'
import { Button } from '../ui/button'
import { Download } from 'lucide-react'

interface Props {
  degreeType: string
  major: string
  recipientName: string
  dateOfBirth: string
  graduationYear: string
  grade: string
  issueDate: string
}

// Component nhận dữ liệu qua props
const Certificate: React.FC<Props> = ({
  degreeType,
  major,
  recipientName,
  dateOfBirth,
  graduationYear,
  grade,
  issueDate
}) => {
  // Sử dụng useRef để tham chiếu đến DOM element của văn bằng
  const certificateRef = useRef(null)

  // Hàm xử lý việc xuất PDF
  const handleExportPdf = () => {
    const element = certificateRef.current
    const fileName = `VanBang_${recipientName.replace(/ /g, '_')}.pdf`

    const options = {
      margin: [0, 0, 0, 0], // không lề
      filename: fileName,
      image: { type: 'jpeg', quality: 0.98 },
      html2canvas: { scale: 3, letterRendering: true, useCORS: true },
      jsPDF: { unit: 'in', format: 'letter', orientation: 'landscape' }
    }

    // Tạo và tải file PDF
    html2pdf().from(element).set(options).save()
  }

  return (
    <div className='font-sans'>
      {/* Phần hiển thị văn bằng */}
      <div
        ref={certificateRef}
        className='certificate-container aspect-[1.414] w-full bg-[#fefbed] p-8 shadow-lg'
        style={{
          fontFamily: "'Times New Roman', Times, serif",
          border: '10px solid transparent',
          borderImage:
            "url(\"data:image/svg+xml,%3Csvg width='100' height='100' viewBox='0 0 100 100' xmlns='http://www.w3.org/2000/svg'%3E%3Cstyle%3E.s%7Bfill:none;stroke:%23e6c669;stroke-width:2%7D%3C/style%3E%3Cpath class='s' d='M0 50h25M75 50h25M50 0v25M50 75v25'/%3E%3Cpath class='s' d='M0 0h10v10H0zM90 0h10v10H90zM0 90h10v10H0zM90 90h10v10H90z'/%3E%3Cpath class='s' d='M25 25h50v50H25z'/%3E%3C/svg%3E\") 20 round",
          position: 'relative'
        }}
      >
        {/* Các element trang trí (watermark, seal) */}
        <div className="absolute left-1/2 top-1/2 z-0 h-3/5 w-3/5 -translate-x-1/2 -translate-y-1/2 bg-[url('data:image/svg+xml,%3Csvg%20xmlns=%27http://www.w3.org/2000/svg%27%20viewBox=%270%200%20200%20200%27%3E%3Ccircle%20cx=%27100%27%20cy=%27100%27%20r=%2790%27%20fill=%27none%27%20stroke=%27%23e6c669%27%20stroke-width=%272%27/%3E%3Ctext%20x=%27100%27%20y=%27105%27%20font-size=%2720%27%20text-anchor=%27middle%27%20fill=%27%23e6c669%27%20font-family=%27Arial%27%3EHVKTMM%3C/text%3E%3C/svg%3E')] bg-contain bg-center bg-no-repeat opacity-10"></div>
        <div className='-rotate-15 absolute bottom-24 right-32 flex h-[120px] w-[120px] items-center justify-center rounded-full border-4 border-red-700 text-center text-xs font-bold text-red-700 opacity-15'>
          HỌC VIỆN KỸ THUẬT MẬT MÃ
        </div>

        {/* Nội dung chính của văn bằng */}
        <div className='relative z-10 text-center text-gray-800'>
          <p className='text-sm tracking-widest'>CỘNG HÒA XÃ HỘI CHỦ NGHĨA VIỆT NAM</p>
          <p className='text-sm font-semibold tracking-wider'>Độc lập - Tự do - Hạnh phúc</p>
          <div className='mx-auto my-2 h-px w-24 bg-gray-600'></div>

          <p className='mt-8 text-lg font-semibold'>GIÁM ĐỐC HỌC VIỆN KỸ THUẬT MẬT MÃ</p>
          <p className='text-md mt-4'>cấp</p>

          <h1 className='my-4 text-5xl uppercase text-red-700' style={{ fontFamily: "'Playfair Display', serif" }}>
            {degreeType}
          </h1>

          <div className='mx-auto mt-8 max-w-md space-y-3 text-left text-lg'>
            <div className='flex items-center'>
              <span className='w-48 font-semibold'>Ngành:</span>
              <span className='flex-1 border-b border-dotted border-gray-500 pb-1 text-center font-bold'>{major}</span>
            </div>
            <div className='flex items-center'>
              <span className='w-48 font-semibold'>Cho:</span>
              <span className='flex-1 border-b border-dotted border-gray-500 pb-1 text-center font-bold'>{`Ông ${recipientName}`}</span>
            </div>
            <div className='flex items-center'>
              <span className='w-48 font-semibold'>Ngày sinh:</span>
              <span className='flex-1 border-b border-dotted border-gray-500 pb-1 text-center font-bold'>
                {dateOfBirth}
              </span>
            </div>
            <div className='flex items-center'>
              <span className='w-48 font-semibold'>Năm tốt nghiệp:</span>
              <span className='flex-1 border-b border-dotted border-gray-500 pb-1 text-center font-bold'>
                {graduationYear}
              </span>
            </div>
            <div className='flex items-center'>
              <span className='w-48 font-semibold'>Hạng tốt nghiệp:</span>
              <span className='flex-1 border-b border-dotted border-gray-500 pb-1 text-center font-bold'>{grade}</span>
            </div>
          </div>

          <div className='mt-16 flex justify-end'>
            <div className='mr-8 text-center'>
              <p>{issueDate}</p>
              <p className='mt-2 font-semibold'>GIÁM ĐỐC</p>
              <p className='mt-4 text-4xl text-red-700 opacity-75' style={{ fontFamily: "'Srisakdi', cursive" }}>
                N. H. Hùng
              </p>
              <p className='mt-2 font-bold'>TS. Nguyễn Hữu Hùng</p>
            </div>
          </div>
        </div>
      </div>

      {/* Nút Xuất PDF */}
      <div className='mt-8 text-center'>
        <Button onClick={handleExportPdf}>
          <Download />
          Xuất ra PDF
        </Button>
      </div>
    </div>
  )
}

export default Certificate
