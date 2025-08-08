'use client'

import { getBlockchainData, getBlockchainFile, getCertificateDataById, getCertificateFile } from '@/lib/api/certificate'
import useSWR from 'swr'
import DecriptionView from './description-view'
import {
  Blocks,
  Book,
  BookOpen,
  Calendar,
  ChartAreaIcon,
  CheckCircleIcon,
  CircleX,
  Eye,
  FileTextIcon,
  Key,
  Library,
  School,
  TagsIcon,
  Text,
  User
} from 'lucide-react'
import { Badge } from '../ui/badge'
import { Button } from '../ui/button'
import PDFView from './pdf-view'
import { Separator } from '../ui/separator'
import CertificateBlankButton from './certificate-blank-button'
import { Alert, AlertDescription, AlertTitle } from '../ui/alert'
import { showNotification } from '@/lib/utils/common'
import { cn } from '@/lib/utils'
import CertificateQrCode from './certificate-qr-code'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from '../ui/dialog'

interface Props {
  isBlockchain: boolean
  id: string
}

const CertificateView: React.FC<Props> = (props) => {
  const queryData = useSWR(
    props.isBlockchain ? undefined : `certificate-view-${props.id}`,
    () => getCertificateDataById(props.id),
    {
      onError: (error) => {
        showNotification('error', error.message || 'Không tải được dữ liệu')
      }
    }
  )
  const isDegree = queryData.data?.certificate?.certificateType !== undefined
  const queryFile = useSWR(
    props.isBlockchain ? undefined : `certificate-file-${props.id}`,
    () => getCertificateFile(props.id),
    {
      revalidateOnFocus: false,
      shouldRetryOnError: false,
      onError: () => {
        showNotification('error', 'Không tải được tệp PDF')
      }
    }
  )

  const queryBlockchainData = useSWR(
    props.isBlockchain ? `blockchain-data-${props.id}` : undefined,
    () => getBlockchainData(props.id),
    {
      revalidateOnFocus: false,
      shouldRetryOnError: false,
      onError: (error) => {
        showNotification('error', error.message || 'Không tải được dữ liệu')
      }
    }
  )
  const queryBlockchainFile = useSWR(
    props.isBlockchain ? `blockchain-file-${props.id}` : undefined,
    () => getBlockchainFile(props.id),
    {
      revalidateOnFocus: false,
      shouldRetryOnError: false,
      onError: () => {
        showNotification('error', 'Không tải được tệp PDF')
      }
    }
  )

  const currentDataQuery = props.isBlockchain ? queryBlockchainData : queryData
  const currentFileQuery = props.isBlockchain ? queryBlockchainFile : queryFile

  const getCertificateItems = (data: any) => {
    return [
      {
        icon: <School className='h-5 w-5 text-gray-500' />,
        title: 'Trường đại học/Học viện',
        value: `${data?.universityCode} - ${data?.universityName}`
      },
      {
        icon: <User className='h-5 w-5 text-gray-500' />,
        title: 'Sinh viên',
        value: `${data?.studentCode} - ${data?.studentName}`
      },
      {
        icon: <Library className='h-5 w-5 text-gray-500' />,
        title: 'Ngành học',
        value: `${data?.facultyCode} - ${data?.facultyName}`
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
        value: data?.regNo
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
        value: `${data?.universityCode} - ${data?.universityName}`
      },
      {
        icon: <User className='h-5 w-5 text-gray-500' />,
        title: 'Sinh viên',
        value: `${data?.studentCode} - ${data?.studentName}`
      },
      {
        icon: <Library className='h-5 w-5 text-gray-500' />,
        title: 'Ngành học',
        value: `${data?.facultyCode} - ${data?.facultyName}`
      },
      {
        icon: <BookOpen className='h-5 w-5 text-gray-500' />,
        title: 'Hệ đào tạo',
        value: data?.educationType
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
        value: data?.regNo
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

  const getDecriptionViewItems = (data: any) => {
    const items = isDegree ? getDegreeItems(data?.certificate) : getCertificateItems(data?.certificate)

    if (props.isBlockchain) {
      return [
        {
          icon: <Blocks className='h-5 w-5 text-gray-500' />,
          title: 'Mã HASH',
          value: (
            <p className='max-w-[300px] truncate' title={data?.on_chain.cert_hash}>
              {data?.on_chain.cert_hash}
            </p>
          )
        },
        ...items
      ]
    }

    return items
  }

  return (
    <div>
      {currentDataQuery.isLoading ? null : !currentDataQuery.error ? (
        <>
          <Alert className={cn('mx-auto mb-4 max-w-[800px]', !props.isBlockchain && 'hidden')} variant='success'>
            <CheckCircleIcon />
            <AlertTitle>Thông báo</AlertTitle>
            <AlertDescription>{currentDataQuery.data?.message || 'Không tải được dữ liệu'}</AlertDescription>
          </Alert>
          <DecriptionView
            title={currentDataQuery?.data?.certificate?.name || 'Không có dữ liệu'}
            items={getDecriptionViewItems(currentDataQuery?.data)}
            description={`Thông tin chi tiết về ${isDegree ? 'văn bằng' : 'chứng chỉ'}`}
            extra={
              <div className='flex items-center gap-2'>
                <CertificateQrCode id={props.id} isIcon={false} />
                {isDegree && currentDataQuery?.data?.certificate && (
                  <Dialog>
                    <DialogTrigger asChild>
                      <Button variant='outline'>
                        <Eye />
                        Xem trước
                      </Button>
                    </DialogTrigger>
                    <DialogContent className='md:min-w-[600px]'>
                      <DialogHeader>
                        <DialogTitle>Văn bằng</DialogTitle>
                      </DialogHeader>
                      {/* <CertificatePreview {...getCertificatePreviewProps(currentDataQuery?.data)!} /> */}
                    </DialogContent>
                  </Dialog>
                )}
              </div>
            }
          />
        </>
      ) : (
        <>
          <Alert className='mx-auto mb-4 max-w-[800px]' variant='destructive'>
            <CircleX />
            <AlertTitle>Lỗi</AlertTitle>
            <AlertDescription>{currentDataQuery?.error?.message || 'Không tải được dữ liệu'}</AlertDescription>
          </Alert>
        </>
      )}

      {currentFileQuery.data ? (
        <>
          <Separator className='my-3' />
          <div className='flex items-center justify-between'>
            <h3 className='mb-3'>Tệp PDF</h3>
            <CertificateBlankButton action={() => currentFileQuery.mutate()} />
          </div>
          <div className='mt-4 h-[700px]'>
            <PDFView url={currentFileQuery?.data} loading={currentFileQuery.isLoading} />
          </div>{' '}
        </>
      ) : (
        <p className='mt-4 text-center text-red-500'>Không có tệp PDF</p>
      )}
    </div>
  )
}

export default CertificateView
