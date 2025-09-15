'use client'

import PageHeader from '@/components/common/page-header'
import CommonPagination from '@/components/common/pagination'
import { UseData } from '@/components/providers/data-provider'
import DetailDialog from '@/components/role/education-admin/detail-dialog'
import Filter from '@/components/common/filter'
import TableList from '@/components/common/table-list'
import UploadButton, { UploadButtonRef } from '@/components/role/education-admin/upload-button'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger
} from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import {
  CERTIFICATE_TYPE_OPTIONS,
  GRADUATION_RANK_OPTIONS,
  PAGE_SIZE,
  STUDENT_CODE_SEARCH_SETTING
} from '@/constants/common'
import CertificateQRCode from '@/components/common/certificate-qr-code'
import {
  createCertificate,
  createDegree,
  getCertificateList,
  importCertificateExcel,
  uploadCertificate,
  uploadCertificatesBlockchain,
  uploadDegree
} from '@/lib/api/certificate'
import { findLabel, formatResponseImportExcel, showNotification } from '@/lib/utils/common'
import { formatFacultyOptions } from '@/lib/utils/format-api'
import { validateGPA, validateNoEmpty } from '@/lib/utils/validators'
import { AlertCircleIcon, Blocks, CheckCircle2Icon, EyeIcon, FileUpIcon, Grid2X2Plus, PlusIcon } from 'lucide-react'
import { useRef, useState } from 'react'

import { useCallback } from 'react'
import useSWR from 'swr'
import useSWRMutation from 'swr/mutation'
import {
  AlertDialog,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogContent,
  AlertDialogTrigger,
  AlertDialogFooter,
  AlertDialogCancel,
  AlertDialogAction
} from '@/components/ui/alert-dialog'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import Link from 'next/link'
import { encodeJSON } from '@/lib/utils/lz-string'

const CertificateManagementPage = () => {
  const [openCreateDegreeDialog, setOpenCreateDegreeDialog] = useState(false)
  const [openCreateCertificateDialog, setOpenCreateCertificateDialog] = useState(false)
  const [typeUpload, setTypeUpload] = useState<'degree' | 'certificate'>('degree')
  const [certificateName, setCertificateName] = useState<string>('')
  const uploadButtonRef = useRef<UploadButtonRef>(null)
  const [openUploadDialog, setOpenUploadDialog] = useState(false)
  const [filter, setFilter] = useState<any>({})
  const facultyOptions = UseData().facultyList

  const queryCertificates = useSWR('certificates-list' + JSON.stringify(filter), () =>
    getCertificateList({
      student_code: filter.studentCode || undefined,
      certificate_type: filter.certificateType || undefined,
      page: filter.page || 1,
      page_size: PAGE_SIZE,
      faculty_code: filter.faculty || undefined,
      // signed: filter.signed || undefined,
      course: filter.course || undefined
    })
  )
  console.log('🚀 ~ CertificateManagementPage ~ queryCertificates:', queryCertificates.data?.data)

  const mutateCreateCertificate = useSWRMutation('create-certificate', (_, { arg }: any) => createCertificate(arg), {
    onSuccess: () => {
      showNotification('success', 'Cấp chứng chỉ thành công')
      queryCertificates.mutate()
      setOpenCreateCertificateDialog(false)
    },
    onError: (error) => {
      showNotification('error', error.message || 'Cấp chứng chỉ thất bại')
    }
  })

  const mutateCreateDegree = useSWRMutation('create-degree', (_, { arg }: any) => createDegree(arg), {
    onSuccess: () => {
      showNotification('success', 'Cấp văn bằng thành công')
      queryCertificates.mutate()
      setOpenCreateDegreeDialog(false)
    },
    onError: (error) => {
      showNotification('error', error.message || 'Cấp văn bằng thất bại')
    }
  })

  const mutateUploadFile = useSWRMutation('upload-certificate', (_, { arg }: { arg: FormData }) => uploadDegree(arg), {
    onSuccess: () => {
      showNotification('success', 'Tải tệp lên thành công')
      queryCertificates.mutate()
      setOpenUploadDialog(false)
      setCertificateName('')
      setTypeUpload('degree')
    },
    onError: (error) => {
      showNotification('error', error.message || 'Lỗi khi tải tệp lên')
    }
  })

  const mutateUploadCertificateFile = useSWRMutation(
    'upload-certificate',
    (_, { arg }: { arg: FormData }) => uploadCertificate(arg, certificateName),
    {
      onSuccess: () => {
        showNotification('success', 'Tải tệp lên thành công')
        queryCertificates.mutate()
        setOpenUploadDialog(false)
        setCertificateName('')
        setTypeUpload('certificate')
      },
      onError: (error) => {
        showNotification('error', error.message || 'Lỗi khi tải tệp lên')
      }
    }
  )

  const handleUploadPDF = useCallback(() => {
    uploadButtonRef.current?.triggerUpload()
  }, [uploadButtonRef])

  const mutateImportCertificateExcel = useSWRMutation(
    'import-certificate-excel',
    (_, { arg }: { arg: FormData }) => importCertificateExcel(arg),
    {
      onSuccess: (data) => {
        const formatData = formatResponseImportExcel(data)

        if (data.error_count === 0) {
          showNotification('success', `Thêm ${data.success_count} văn bằng/chứng chỉ thành công`)
          queryCertificates.mutate()
          return
        }

        if (data.success_count === 0) {
          formatData.error.forEach((item) => {
            showNotification('error', `Văn bằng/chứng chỉ hàng thứ ${item.row.join(', ')} có lỗi: "${item.title}"`)
          })
          return
        }

        formatData.error.forEach((item) => {
          showNotification('error', `Văn bằng/chứng chỉ hàng thứ ${item.row.join(', ')} có lỗi: "${item.title}" `)
        })

        showNotification(
          'success',
          `Văn bằng/chứng chỉ hàng thứ ${formatData.success.join(', ')} đã được thêm thành công`
        )
        queryCertificates.mutate()
      },
      onError: (error) => {
        showNotification('error', error.message || 'Lỗi khi nhập tệp excel')
      }
    }
  )
  const mutatePushCertificatesBlockchain = useSWRMutation(
    'push-certificates-blockchain',
    async (_key, { arg }: { arg: any }) => {
      const facultyId = facultyOptions.find((item) => item.code === arg.faculty)?.id

      const res = await uploadCertificatesBlockchain(facultyId as string, arg.certificateType, arg.course)
      queryCertificates.mutate()

      return res
    },
    {
      onError: (error) => {
        showNotification('error', error.message || 'Lỗi khi đẩy cả khóa lên Blockchain')
      },
      onSuccess: () => {
        showNotification('success', 'Đẩy cả khóa lên Blockchain thành công')
      }
    }
  )

  const handleUpload = useCallback(
    (file: FormData) => {
      if (typeUpload === 'degree') {
        mutateUploadFile.trigger(file)
      } else {
        mutateUploadCertificateFile.trigger(file)
      }
    },
    [mutateUploadCertificateFile, mutateUploadFile, typeUpload]
  )

  const handleImportCertificateExcel = useCallback(
    (file: FormData) => {
      mutateImportCertificateExcel.trigger(file)
    },
    [mutateImportCertificateExcel]
  )

  const handleCreateCertificate = useCallback(
    (data: any) => {
      mutateCreateCertificate.trigger(data)
    },
    [mutateCreateCertificate]
  )

  const handleCreateDegree = useCallback(
    (data: any) => {
      mutateCreateDegree.trigger(data)
    },
    [mutateCreateDegree]
  )

  const encodeCertificateData = (data: any): string => {
    return (
      encodeJSON({
        university_id: data.universityId,
        university_code: data.universityCode,
        faculty_id: facultyOptions.find((faculty) => faculty.code === data.faculty)?.id,
        certificate_type: data.certificateType,
        course: data.course,
        certificate_id: data.id
      }) ?? ''
    )
  }

  return (
    <>
      <PageHeader
        title='Văn bằng & Chứng chỉ'
        extra={[
          <UploadButton
            key='upload-excel'
            handleUpload={handleImportCertificateExcel}
            loading={mutateImportCertificateExcel.isMutating}
            title={'Tải Excel'}
            icon={<FileUpIcon />}
          />,
          <Button key='create-new-degree' onClick={() => setOpenCreateDegreeDialog(true)}>
            <PlusIcon />
            <span className='hidden md:block'>Cấp văn bằng</span>
          </Button>,
          <Button
            variant={'secondary'}
            key='create-new-certificate'
            onClick={() => setOpenCreateCertificateDialog(true)}
          >
            <PlusIcon />
            <span className='hidden md:block'>Cấp chứng chỉ</span>
          </Button>,
          <Dialog key='upload-pdf' open={openUploadDialog} onOpenChange={setOpenUploadDialog}>
            <DialogTrigger asChild>
              <Button variant={'outline'} title='Có hỗ trợ tải nhiều tệp cùng lúc'>
                <FileUpIcon />
                <span className='hidden md:block'>Tải PDF</span>
              </Button>
            </DialogTrigger>
            <DialogContent>
              <DialogHeader>
                <DialogTitle>Tải tệp PDF chứng chỉ/văn bằng</DialogTitle>
                <DialogDescription>
                  Nếu tải văn bằng thì tên tệp là <strong>số hiệu văn bằng</strong>, nếu tải chứng chỉ thì tên tệp là{' '}
                  <strong>mã sinh viên</strong>
                </DialogDescription>
              </DialogHeader>
              <Label>Chọn loại</Label>
              <Select defaultValue='degree' onValueChange={(value) => setTypeUpload(value as 'degree' | 'certificate')}>
                <SelectTrigger>
                  <SelectValue placeholder='Chọn loại' />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value='degree'>Văn bằng</SelectItem>
                  <SelectItem value='certificate'>Chứng chỉ</SelectItem>
                </SelectContent>
              </Select>
              {typeUpload === 'certificate' && (
                <>
                  <Label>Tên tệp</Label>
                  <Input
                    value={certificateName}
                    onChange={(e) => setCertificateName(e.target.value)}
                    placeholder='Nhập tên tệp'
                  />
                </>
              )}
              <DialogFooter>
                <DialogClose asChild>
                  <Button variant={'outline'}>Hủy bỏ</Button>
                </DialogClose>
                <Button
                  onClick={handleUploadPDF}
                  disabled={mutateUploadFile.isMutating || mutateUploadCertificateFile.isMutating}
                >
                  Tải tệp
                </Button>
              </DialogFooter>
            </DialogContent>
          </Dialog>,
          <AlertDialog key='blockchain-push-degrees'>
            <AlertDialogTrigger asChild>
              <Button title='Đẩy cả khóa lên Blockchain'>
                <Grid2X2Plus />
                <span className='hidden md:block'>{'Blockchain'}</span>
              </Button>
            </AlertDialogTrigger>
            <AlertDialogContent>
              <AlertDialogHeader>
                <AlertDialogTitle>Xác nhận đẩy cả khóa lên Blockchain</AlertDialogTitle>
              </AlertDialogHeader>
              {filter.faculty ? (
                <Alert variant={'success'}>
                  <CheckCircle2Icon />
                  <AlertTitle>Sẵn sàng</AlertTitle>
                  <AlertDescription>
                    <ul className='list-inside list-disc'>
                      <li>Chuyên ngành: {findLabel(filter.faculty, formatFacultyOptions(facultyOptions))}</li>
                      {filter.certificateType && <li>Loại bằng: {filter.certificateType}</li>}
                      {filter.course && <li>Khóa học: {filter.course}</li>}
                    </ul>
                  </AlertDescription>
                </Alert>
              ) : (
                <Alert variant={'warning'}>
                  <AlertCircleIcon />
                  <AlertTitle>Cảnh báo</AlertTitle>
                  <AlertDescription>
                    Vui lòng chọn chuyên ngành trong <strong>phần tìm kiếm</strong> để tiến hành cấp bằng số.
                  </AlertDescription>
                </Alert>
              )}
              <AlertDialogFooter>
                <AlertDialogCancel>Hủy bỏ</AlertDialogCancel>
                <AlertDialogAction
                  disabled={!filter.faculty}
                  onClick={() => mutatePushCertificatesBlockchain.trigger(filter)}
                >
                  Xác nhận
                </AlertDialogAction>
              </AlertDialogFooter>
            </AlertDialogContent>
          </AlertDialog>
        ]}
      />
      <div className='hidden'>
        <UploadButton
          handleUpload={handleUpload}
          loading={mutateUploadFile.isMutating || mutateUploadCertificateFile.isMutating}
          ref={uploadButtonRef}
        />
      </div>
      <Filter
        items={[
          {
            type: 'query_select',
            placeholder: 'Nhập và chọn MSV',
            name: 'studentCode',
            setting: STUDENT_CODE_SEARCH_SETTING
          },
          {
            type: 'select',
            name: 'faculty',
            placeholder: 'Chọn chuyên ngành',
            setting: {
              select: {
                groups: [
                  {
                    label: 'Chuyên ngành',
                    options: formatFacultyOptions(facultyOptions)
                  }
                ]
              }
            }
          },
          {
            type: 'select',
            placeholder: 'Chọn loại bằng',
            name: 'certificateType',
            setting: {
              select: {
                groups: [
                  {
                    label: 'Bằng tốt nghiệp',
                    options: CERTIFICATE_TYPE_OPTIONS
                  }
                ]
              }
            }
          },
          {
            type: 'input',
            name: 'course',
            placeholder: 'Nhập khóa'
          }
        ]}
        handleSetFilter={setFilter}
      />
      <TableList
        items={[
          { header: 'Mã SV', value: 'studentCode', className: 'min-w-[80px] font-semibold text-blue-500' },
          { header: 'Họ và tên', value: 'studentName', className: 'min-w-[200px]' },
          { header: 'Chuyên ngành', value: 'facultyName', className: 'min-w-[150px]' },
          {
            header: 'Phân loại',
            value: 'isDegree',
            render: (item) => {
              return item.isDegree ? (
                <div className='flex items-center gap-2'>
                  <Badge>Văn bằng</Badge>
                  <Badge className='bg-blue-500 text-white hover:bg-blue-400'>{item.certificateType}</Badge>
                </div>
              ) : (
                <Badge variant='outline'>Chứng chỉ</Badge>
              )
            }
          },
          { header: 'Tên văn bằng/chứng chỉ', value: 'name', className: 'min-w-[100px]' },
          { header: 'Ngày cấp', value: 'date', className: 'min-w-[100px]' },
          {
            header: 'Blockchain',
            value: 'onBlockchain',

            render: (item) => (
              <Badge variant={item.onBlockchain ? 'default' : 'outline'}>
                {item.onBlockchain ? 'Đã đẩy' : 'Chưa đẩy'}
              </Badge>
            )
          },
          {
            header: 'Hành động',
            value: 'action',

            render: (item) => (
              <div className='flex gap-2'>
                <Link href={`/education-admin/certificate-management/${item.id}`}>
                  <Button size={'icon'} variant={'outline'} title='Xem dữ liệu trên cơ sở dữ liệu'>
                    <EyeIcon />
                  </Button>
                </Link>
                <Link
                  href={`/education-admin/certificate-management/${encodeCertificateData(item)}/blockchain`}
                  onClick={(e) => {
                    if (!item.onBlockchain) {
                      e.preventDefault()
                      showNotification('error', (item.isDegree ? 'Văn bằng' : 'Chứng chỉ') + ' chưa đẩy lên blockchain')
                      return
                    }
                  }}
                >
                  <Button size={'icon'} title='Xem dữ liệu trên blockchain' disabled={!item.onBlockchain}>
                    <Blocks />
                  </Button>
                </Link>
                <CertificateQRCode id={encodeCertificateData(item)} isIcon={true} disable={!item.onBlockchain} />
              </div>
            )
          }
        ]}
        data={queryCertificates.data?.data || []}
        page={queryCertificates.data?.page || 1}
        pageSize={queryCertificates.data?.page_size || PAGE_SIZE}
      />
      <CommonPagination
        page={queryCertificates.data?.page || 1}
        totalPage={queryCertificates.data?.total_page || 1}
        handleChangePage={(page) => {
          setFilter({ ...filter, page })
        }}
      />
      <DetailDialog
        title='Cấp văn bằng'
        items={[
          {
            type: 'query_select',
            placeholder: 'Nhập và chọn MSV',
            name: 'studentCode',
            setting: STUDENT_CODE_SEARCH_SETTING,
            label: 'Mã sinh viên',
            validator: validateNoEmpty('Mã sinh viên')
          },
          {
            type: 'input',
            name: 'major',
            label: 'Chuyên ngành',
            placeholder: 'Nhập chuyên ngành',
            validator: validateNoEmpty('Chuyên ngành')
          },
          {
            type: 'select',
            placeholder: 'Chọn loại bằng',
            name: 'certificateType',
            setting: {
              select: {
                groups: [
                  {
                    label: undefined,
                    options: CERTIFICATE_TYPE_OPTIONS
                  }
                ]
              }
            },
            label: 'Loại bằng',
            validator: validateNoEmpty('Loại bằng')
          },
          {
            type: 'input',
            placeholder: 'Nhập tên bằng',
            name: 'name',
            label: 'Tên bằng',
            validator: validateNoEmpty('Tên bằng')
          },
          {
            type: 'select',
            placeholder: 'Chọn xếp loại',
            name: 'graduationRank',
            label: 'Xếp loại',
            setting: {
              select: {
                groups: [{ label: undefined, options: GRADUATION_RANK_OPTIONS }]
              }
            }
          },
          {
            type: 'input',
            placeholder: 'Nhập khóa',
            name: 'course',
            label: 'Khóa',
            validator: validateNoEmpty('Khóa')
          },
          {
            type: 'input',
            name: 'serialNumber',
            placeholder: 'Nhập số hiệu',
            label: 'Số hiệu',
            validator: validateNoEmpty('Số hiệu')
          },
          {
            type: 'input',
            name: 'regNo',
            placeholder: 'Nhập số vào sổ gốc cấp văn bằng',
            label: 'Số vào sổ gốc cấp văn bằng',
            validator: validateNoEmpty('Số vào sổ gốc cấp văn bằng')
          },
          {
            type: 'input',
            name: 'date',
            placeholder: 'Nhập ngày cấp',
            label: 'Ngày cấp',
            validator: validateNoEmpty('Ngày cấp'),
            setting: {
              input: {
                type: 'date'
              }
            }
          },
          {
            type: 'input',
            name: 'gpa',
            placeholder: 'Nhập điểm GPA',
            label: 'Điểm GPA',
            validator: validateGPA
          },
          {
            type: 'textarea',
            name: 'description',
            label: 'Mô tả',
            placeholder: 'Nhập mô tả'
          }
        ]}
        data={[]}
        mode={openCreateDegreeDialog ? 'create' : undefined}
        handleSubmit={handleCreateDegree}
        handleClose={() => setOpenCreateDegreeDialog(false)}
      />
      <DetailDialog
        title='Cấp chứng chỉ'
        items={[
          {
            type: 'query_select',
            placeholder: 'Nhập và chọn MSV',
            name: 'studentCode',
            setting: STUDENT_CODE_SEARCH_SETTING,
            label: 'Mã sinh viên',
            validator: validateNoEmpty('Mã sinh viên')
          },
          {
            type: 'input',
            placeholder: 'Nhập tên chứng chỉ',
            name: 'name',
            label: 'Tên chứng chỉ',
            validator: validateNoEmpty('Tên chứng chỉ')
          },
          {
            type: 'input',
            name: 'serialNumber',
            placeholder: 'Nhập số hiệu',
            label: 'Số hiệu',
            validator: validateNoEmpty('Số hiệu')
          },
          {
            type: 'input',
            name: 'regNo',
            placeholder: 'Nhập số vào sổ gốc cấp văn bằng',
            label: 'Số vào sổ gốc cấp văn bằng',
            validator: validateNoEmpty('Số vào sổ gốc cấp văn bằng')
          },
          {
            type: 'input',
            name: 'date',
            placeholder: 'Nhập ngày cấp',
            label: 'Ngày cấp',
            validator: validateNoEmpty('Ngày cấp'),
            setting: {
              input: {
                type: 'date'
              }
            }
          },
          {
            type: 'textarea',
            name: 'description',
            label: 'Mô tả',
            placeholder: 'Nhập mô tả'
          }
        ]}
        data={[]}
        mode={openCreateCertificateDialog ? 'create' : undefined}
        handleSubmit={handleCreateCertificate}
        handleClose={() => setOpenCreateCertificateDialog(false)}
      />
    </>
  )
}

export default CertificateManagementPage
