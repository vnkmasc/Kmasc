'use client'

import { getBlockchainData, getBlockchainFile, getCertificateDataById, getCertificateFile } from '@/lib/api/certificate'
import useSWR from 'swr'
import DecriptionView from './description-view'
import { Book, Calendar, ChartAreaIcon, FileTextIcon, Key, Library, School, TagsIcon, Text, User } from 'lucide-react'
import { Badge } from '../ui/badge'
import PDFView from './pdf-view'
import { Separator } from '../ui/separator'
import CertificateBlankButton from './certificate-blank-button'

interface Props {
  isBlockchain: boolean
  id: string
}

const CertificateView: React.FC<Props> = (props) => {
  const queryData = useSWR(props.isBlockchain ? undefined : `certificate-view-${props.id}`, () =>
    getCertificateDataById(props.id)
  )
  const isDegree = queryData.data?.certificateType !== undefined
  const queryFile = useSWR(
    props.isBlockchain ? undefined : `certificate-file-${props.id}`,
    () => getCertificateFile(props.id),
    {
      revalidateOnFocus: false,
      shouldRetryOnError: false
    }
  )

  const queryBlockchainData = useSWR(
    props.isBlockchain ? `blockchain-data-${props.id}` : undefined,
    () => getBlockchainData(props.id),
    {
      revalidateOnFocus: false,
      shouldRetryOnError: false
    }
  )
  const queryBlockchainFile = useSWR(
    props.isBlockchain ? `blockchain-file-${props.id}` : undefined,
    () => getBlockchainFile(props.id),
    {
      revalidateOnFocus: false,
      shouldRetryOnError: false
    }
  )

  const getCertificateItems = (data: any) => {
    return [
      {
        icon: <School className='h-5 w-5 text-gray-500' />,
        title: 'Trường đại học/Học viện',
        value: `${data?.universityCode ?? 'KMA'} - ${data?.universityName ?? 'Học viện Kỹ thuật Mật mã'}`
      },
      {
        icon: <User className='h-5 w-5 text-gray-500' />,
        title: 'Sinh viên',
        value: `${data?.studentCode ?? '20200000'} - ${data?.studentName ?? 'Nguyễn Ngọc Tuyền'}`
      },
      {
        icon: <Library className='h-5 w-5 text-gray-500' />,
        title: 'Ngành học',
        value: `${data?.facultyCode ?? 'CNTT'} - ${data?.facultyName ?? 'Công nghệ thông tin'}`
      },

      {
        icon: <Book className='h-5 w-5 text-gray-500' />,
        title: 'Chứng chỉ',
        value: data?.name
      },
      {
        icon: <Calendar className='h-5 w-5 text-gray-500' />,
        title: 'Ngày cấp',
        value: data?.date
      },
      {
        icon: <TagsIcon className='h-5 w-5 text-gray-500' />,
        title: 'Số hiệu',
        value: data?.serialNumber
      },
      {
        icon: <FileTextIcon className='h-5 w-5 text-gray-500' />,
        title: 'Số vào sổ gốc cấp văn bằng',
        value: data?.regNo ?? 512
      },
      {
        icon: <Text className='h-5 w-5 text-gray-500' />,
        title: 'Mô tả',
        value: data?.description
      },
      {
        icon: <Key className='h-5 w-5 text-gray-500' />,
        title: 'Trạng thái ký',
        value: <Badge variant={data?.signed ? 'default' : 'outline'}>{data?.signed ? 'Đã ký' : 'Chưa ký'}</Badge>
      }
    ]
  }

  const getDegreeItems = (data: any) => {
    return [
      {
        icon: <School className='h-5 w-5 text-gray-500' />,
        title: 'Trường đại học/Học viện',
        value: `${data?.universityCode ?? 'KMA'} - ${data?.universityName ?? 'Học viện Kỹ thuật Mật mã'}`
      },
      {
        icon: <User className='h-5 w-5 text-gray-500' />,
        title: 'Sinh viên',
        value: `${data?.studentCode} - ${data?.studentName ?? 'Nguyễn Ngọc Tuyền'}`
      },
      {
        icon: <Library className='h-5 w-5 text-gray-500' />,
        title: 'Ngành học',
        value: `${data?.facultyCode ?? 'CNTT'} - ${data?.facultyName ?? 'Công nghệ thông tin'}`
      },
      {
        icon: <Book className='h-5 w-5 text-gray-500' />,
        title: 'Văn bằng',
        value: (
          <div>
            <Badge className='bg-blue-500 text-white hover:bg-blue-400'>{data?.certificateType ?? '-'}</Badge>
            {' - '}
            <span>{data?.name}</span>
          </div>
        )
      },
      {
        icon: <ChartAreaIcon className='h-5 w-5 text-gray-500' />,
        title: 'GPA',
        value: data?.gpa
      },
      {
        icon: <Calendar className='h-5 w-5 text-gray-500' />,
        title: 'Ngày cấp',
        value: data?.date
      },
      {
        icon: <TagsIcon className='h-5 w-5 text-gray-500' />,
        title: 'Số hiệu',
        value: data?.serialNumber
      },
      {
        icon: <FileTextIcon className='h-5 w-5 text-gray-500' />,
        title: 'Số vào sổ gốc cấp văn bằng',
        value: data?.regNo ?? 512
      },
      {
        icon: <Text className='h-5 w-5 text-gray-500' />,
        title: 'Mô tả',
        value: data?.description
      },
      {
        icon: <Key className='h-5 w-5 text-gray-500' />,
        title: 'Trạng thái ký',
        value: <Badge variant={data?.signed ? 'default' : 'outline'}>{data?.signed ? 'Đã ký' : 'Chưa ký'}</Badge>
      }
    ]
  }

  return (
    <div>
      <DecriptionView
        title={queryData.data?.name || 'Không có dữ liệu'}
        items={
          isDegree
            ? getDegreeItems(props.isBlockchain ? queryBlockchainData.data : queryData.data)
            : getCertificateItems(props.isBlockchain ? queryBlockchainData.data : queryData.data)
        }
        description={`Thông tin chi tiết về ${isDegree ? 'văn bằng' : 'chứng chỉ'}`}
      />
      {(props.isBlockchain ? queryBlockchainFile.data : queryFile.data) ? (
        <>
          <Separator className='my-3' />
          <div className='flex items-center justify-between'>
            <h3 className='mb-3'>Tệp PDF</h3>
            <CertificateBlankButton
              action={() => (props.isBlockchain ? queryBlockchainFile.mutate() : queryFile.mutate())}
            />
          </div>
          <div className='mt-4 h-[700px]'>
            <PDFView
              url={props.isBlockchain ? queryBlockchainFile.data : queryFile.data}
              loading={queryFile.isLoading}
            />
          </div>{' '}
        </>
      ) : (
        <p className='mt-4 text-center text-red-500'>Không có tệp PDF</p>
      )}
    </div>
  )
}

export default CertificateView
