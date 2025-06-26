'use client'

import { getCertificateDataById, getCertificateFile } from '@/lib/api/certificate'
import useSWR from 'swr'
import DecriptionView from './description-view'
import { Book, Calendar, FileTextIcon, Key, Library, School, TagsIcon, User } from 'lucide-react'
import { Badge } from '../ui/badge'
import PDFView from './pdf-view'
import { Separator } from '../ui/separator'
import CertificateBlankButton from './certificate-blank-button'

interface Props {
  id: string
}

const CertificateView: React.FC<Props> = (props) => {
  const queryData = useSWR(`certificate-view-${props.id}`, () => getCertificateDataById(props.id))
  const isDegree = queryData.data?.certificateType !== undefined
  const queryFile = useSWR(`certificate-file-${props.id}`, () => getCertificateFile(props.id), {
    revalidateOnFocus: false,
    shouldRetryOnError: false
  })

  const certificateItems = [
    {
      icon: <School className='h-5 w-5 text-gray-500' />,
      title: 'Trường đại học/Học viện',
      value: `${queryData.data?.universityCode} - ${queryData.data?.universityName}`
    },
    {
      icon: <User className='h-5 w-5 text-gray-500' />,
      title: 'Sinh viên',
      value: `${queryData.data?.studentCode} - ${queryData.data?.studentName}`
    },
    {
      icon: <Library className='h-5 w-5 text-gray-500' />,
      title: 'Ngành học',
      value: `${queryData.data?.facultyCode} - ${queryData.data?.facultyName}`
    },
    {
      icon: <Book className='h-5 w-5 text-gray-500' />,
      title: 'Chứng chỉ',
      value: queryData.data?.name
    },
    {
      icon: <Calendar className='h-5 w-5 text-gray-500' />,
      title: 'Ngày cấp',
      value: queryData.data?.date
    },
    {
      icon: <Key className='h-5 w-5 text-gray-500' />,
      title: 'Trạng thái ký',
      value: (
        <Badge variant={queryData.data?.signed ? 'default' : 'outline'}>
          {queryData.data?.signed ? 'Đã ký' : 'Chưa ký'}
        </Badge>
      )
    }
  ]

  const degreeItems = [
    {
      icon: <School className='h-5 w-5 text-gray-500' />,
      title: 'Trường đại học/Học viện',
      value: `${queryData.data?.universityCode} - ${queryData.data?.universityName}`
    },
    {
      icon: <User className='h-5 w-5 text-gray-500' />,
      title: 'Sinh viên',
      value: `${queryData.data?.studentCode} - ${queryData.data?.studentName}`
    },
    {
      icon: <Library className='h-5 w-5 text-gray-500' />,
      title: 'Ngành học',
      value: `${queryData.data?.facultyCode} - ${queryData.data?.facultyName}`
    },
    {
      icon: <Book className='h-5 w-5 text-gray-500' />,
      title: 'Văn bằng',
      value: (
        <div>
          <Badge className='bg-blue-500 text-white hover:bg-blue-400'>{queryData.data?.certificateType ?? '-'}</Badge>
          {' - '}
          <span>{queryData.data?.name}</span>
        </div>
      )
    },
    {
      icon: <Calendar className='h-5 w-5 text-gray-500' />,
      title: 'Ngày cấp',
      value: queryData.data?.date
    },
    {
      icon: <TagsIcon className='h-5 w-5 text-gray-500' />,
      title: 'Số hiệu',
      value: queryData.data?.serialNumber
    },
    {
      icon: <FileTextIcon className='h-5 w-5 text-gray-500' />,
      title: 'Số vào sổ gốc cấp văn bằng',
      value: queryData.data?.regNo
    },
    {
      icon: <Key className='h-5 w-5 text-gray-500' />,
      title: 'Trạng thái ký',
      value: (
        <Badge variant={queryData.data?.signed ? 'default' : 'outline'}>
          {queryData.data?.signed ? 'Đã ký' : 'Chưa ký'}
        </Badge>
      )
    }
  ]

  return (
    <div>
      <DecriptionView
        title={queryData.data?.name || 'Không có dữ liệu'}
        items={isDegree ? degreeItems : certificateItems}
        description={`Thông tin chi tiết về ${isDegree ? 'văn bằng' : 'chứng chỉ'}`}
      />
      {queryFile.data ? (
        <>
          <Separator className='my-3' />
          <div className='flex items-center justify-between'>
            <h3 className='mb-3'>Tệp PDF</h3>
            <CertificateBlankButton action={() => queryFile.mutate()} />
          </div>
          <div className='mt-4 h-[700px]'>
            <PDFView url={queryFile.data} loading={queryFile.isLoading} />
          </div>{' '}
        </>
      ) : (
        <p className='mt-4 text-center text-red-500'>Không có tệp PDF</p>
      )}
    </div>
  )
}

export default CertificateView
