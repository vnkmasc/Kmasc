'use client'

import useSWR from 'swr'
import DecriptionView from './description-view'
import {
  Book,
  BookOpen,
  Calendar,
  ChartAreaIcon,
  CheckCircleIcon,
  FileTextIcon,
  Library,
  School,
  TagsIcon,
  Text,
  User
} from 'lucide-react'
import PDFView from './pdf-view'
import { Separator } from '../ui/separator'
import CertificateBlankButton from '../role/education-admin/certificate-management/certificate-blank-button'
import { Alert, AlertDescription, AlertTitle } from '../ui/alert'
import { showNotification } from '@/lib/utils/common'
import { cn } from '@/lib/utils/common'
import {
  getDigitalDegreeById,
  getDigitalDegreeFileById,
  verifyDigitalDegreeDataBlockchain,
  verifyDigitalDegreeFileBlockchain
} from '@/lib/api/digital-degree'
import { Badge } from '../ui/badge'
import useSWRMutation from 'swr/mutation'
import { useEffect } from 'react'
import CertificateQrCode from './certificate-qr-code'
import { encodeJSON } from '@/lib/utils/lz-string'

interface Props {
  isBlockchain: boolean
  id: string
  universityCode?: string
  universityId?: string
  facultyId?: string
  certificateType?: string
  course?: string
}

const DigitalDegreeView: React.FC<Props> = (props) => {
  const queryData = useSWR(
    !props.isBlockchain && props.id ? `digital-degree-view-${props.id}` : undefined,
    () => getDigitalDegreeById(props.id),
    {
      onError: (error) => {
        showNotification('error', error.message || 'Không tải được dữ liệu')
      }
    }
  )

  const queryFile = useSWR(
    !props.isBlockchain && props.id ? `digital-degree-file-${props.id}` : undefined,
    () => getDigitalDegreeFileById(props.id),
    {
      revalidateOnFocus: false,
      shouldRetryOnError: false,
      onError: () => {
        showNotification('error', 'Không tải được tệp PDF')
      }
    }
  )

  const mutateVerifyDigitalDegreeDataBlockchain = useSWRMutation(
    'verify-digital-degree-data-blockchain',
    () =>
      verifyDigitalDegreeDataBlockchain(
        props?.universityId || '',
        props?.facultyId || '',
        props?.certificateType || '',
        props?.course || '',
        props.id
      ),
    {
      onError: (error: any) => {
        showNotification('error', error.message || 'Lỗi khi xác minh dữ liệu trên blockchain')
      },
      onSuccess: (data) => {
        if (!data.verified) {
          showNotification('error', data.message || 'Dữ liệu không hợp lệ')
        } else {
          showNotification('success', 'Xác minh dữ liệu trên blockchain thành công')
        }
      }
    }
  )

  useEffect(() => {
    if (props.isBlockchain) {
      mutateVerifyDigitalDegreeDataBlockchain.trigger()
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [props.isBlockchain])

  const queryBlockchainFile = useSWR(
    props.isBlockchain && props.id && props.universityCode && mutateVerifyDigitalDegreeDataBlockchain.data?.verified
      ? `digital-degree-blockchain-file-${props.id}`
      : undefined,
    () => verifyDigitalDegreeFileBlockchain(props.universityCode ?? '', props.id),
    {
      revalidateOnFocus: false,
      shouldRetryOnError: false,
      onError: () => {
        showNotification('error', 'Không tải được tệp PDF')
      }
    }
  )

  const currentFileQuery = props.isBlockchain ? queryBlockchainFile : queryFile

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
      {props.isBlockchain ? (
        <>
          <Alert className={cn('mx-auto mb-4 max-w-[800px]', !props.isBlockchain && 'hidden')} variant='success'>
            <CheckCircleIcon />
            <AlertTitle>Thông báo</AlertTitle>
            <AlertDescription>
              {mutateVerifyDigitalDegreeDataBlockchain.data?.message || 'Xác minh dữ liệu trên blockchain thành công'}
            </AlertDescription>
          </Alert>
          <DecriptionView
            title={mutateVerifyDigitalDegreeDataBlockchain?.data?.data?.name || 'Không có dữ liệu'}
            items={getDegreeItems(mutateVerifyDigitalDegreeDataBlockchain?.data?.data)}
            description={`Thông tin chi tiết về văn bằng số`}
            extra={
              <CertificateQrCode
                id={
                  encodeJSON({
                    university_id: props.universityId,
                    university_code: props.universityCode,
                    faculty_id: props.facultyId,
                    certificate_type: props.certificateType,
                    course: props.course,
                    ediploma_id: props.id
                  }) ?? ''
                }
                isIcon={false}
              />
            }
          />
        </>
      ) : (
        <>
          <Alert className={cn('mx-auto mb-4 max-w-[800px]', !props.isBlockchain && 'hidden')} variant='success'>
            <CheckCircleIcon />
            <AlertTitle>Thông báo</AlertTitle>
            <AlertDescription>{queryData.data?.message || 'Không tải được dữ liệu'}</AlertDescription>
          </Alert>

          <DecriptionView
            title={queryData?.data?.data?.name || 'Không có dữ liệu'}
            items={getDegreeItems(queryData?.data?.data)}
            description={`Thông tin chi tiết về văn bằng số`}
          />
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
