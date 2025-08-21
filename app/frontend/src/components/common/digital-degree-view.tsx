'use client'

import useSWR from 'swr'
import DecriptionView from './description-view'
import {
  Book,
  BookOpen,
  Calendar,
  ChartAreaIcon,
  CheckCircleIcon,
  CircleX,
  Eye,
  FileTextIcon,
  Library,
  School,
  TagsIcon,
  Text,
  User
} from 'lucide-react'
import { Button } from '../ui/button'
import PDFView from './pdf-view'
import { Separator } from '../ui/separator'
import CertificateBlankButton from '../role/education-admin/certificate-management/certificate-blank-button'
import { Alert, AlertDescription, AlertTitle } from '../ui/alert'
import { showNotification } from '@/lib/utils/common'
import { cn } from '@/lib/utils/common'
import CertificateQrCode from './certificate-qr-code'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from '../ui/dialog'
import { getDigitalDegreeById, getDigitalDegreeFileById } from '@/lib/api/digital-degree'
import { Badge } from '../ui/badge'

interface Props {
  isBlockchain: boolean
  id: string
}

const DigitalDegreeView: React.FC<Props> = (props) => {
  const queryData = useSWR(
    props.isBlockchain ? undefined : `digital-degree-view-${props.id}`,
    () => getDigitalDegreeById(props.id),
    {
      onError: (error) => {
        showNotification('error', error.message || 'Không tải được dữ liệu')
      }
    }
  )

  const queryFile = useSWR(
    props.isBlockchain ? undefined : `digital-degree-file-${props.id}`,
    () => getDigitalDegreeFileById(props.id),
    {
      revalidateOnFocus: false,
      shouldRetryOnError: false,
      onError: () => {
        showNotification('error', 'Không tải được tệp PDF')
      }
    }
  )

  //   const queryBlockchainData = useSWR(
  //     props.isBlockchain ? `blockchain-data-${props.id}` : undefined,
  //     () => getBlockchainData(props.id),
  //     {
  //       revalidateOnFocus: false,
  //       shouldRetryOnError: false,
  //       onError: (error) => {
  //         showNotification('error', error.message || 'Không tải được dữ liệu')
  //       }
  //     }
  //   )
  //   const queryBlockchainFile = useSWR(
  //     props.isBlockchain ? `blockchain-file-${props.id}` : undefined,
  //     () => getBlockchainFile(props.id),
  //     {
  //       revalidateOnFocus: false,
  //       shouldRetryOnError: false,
  //       onError: () => {
  //         showNotification('error', 'Không tải được tệp PDF')
  //       }
  //     }
  //   )

  const currentDataQuery = queryData
  const currentFileQuery = queryFile

  const getDegreeItems = (data: any) => {
    return [
      {
        icon: <School className='h-5 w-5 text-gray-500' />,
        title: 'Trường đại học/Học viện',
        value: `${data?.university_code} - ${data?.university_name}`
      },
      {
        icon: <User className='h-5 w-5 text-gray-500' />,
        title: 'Sinh viên',
        value: `${data?.student_code} - ${data?.student_name}`
      },
      {
        icon: <Library className='h-5 w-5 text-gray-500' />,
        title: 'Ngành học',
        value: `${data?.faculty_code} - ${data?.faculty_name}`
      },
      {
        icon: <BookOpen className='h-5 w-5 text-gray-500' />,
        title: 'Hệ đào tạo',
        value: data?.education_type
      },
      {
        icon: <Book className='h-5 w-5 text-gray-500' />,
        title: 'Văn bằng',
        value: (
          <div>
            <Badge className='bg-blue-500 text-white hover:bg-blue-400'>{data?.certificate_type ?? '-'}</Badge>
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
        value: data?.issue_date
      },
      {
        icon: <TagsIcon className='h-5 w-5 text-gray-500' />,
        title: 'Số hiệu',
        value: data?.serial_number
      },
      {
        icon: <FileTextIcon className='h-5 w-5 text-gray-500' />,
        title: 'Số vào sổ gốc cấp văn bằng',
        value: data?.registration_number
      },
      {
        icon: <Text className='h-5 w-5 text-gray-500' />,
        title: 'Mô tả',
        value: data?.description
      }
      // {
      //   icon: <Key className='h-5 w-5 text-gray-500' />,
      //   title: 'Trạng thái ký',
      //   value: <Badge variant={data?.signed ? 'default' : 'outline'}>{data?.signed ? 'Đã ký' : 'Chưa ký'}</Badge>
      // }
    ]
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
            items={getDegreeItems(currentDataQuery?.data?.data)}
            description={`Thông tin chi tiết về văn bằng số`}
            extra={
              <div className='flex items-center gap-2'>
                <CertificateQrCode id={props.id} isIcon={false} />
                {currentDataQuery?.data?.certificate && (
                  <Dialog>
                    <DialogTrigger asChild>
                      <Button variant='outline'>
                        <Eye />
                        <span className='hidden md:block'>Xem trước</span>
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

export default DigitalDegreeView
