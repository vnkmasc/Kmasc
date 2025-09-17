'use client'

import { School, User, Book, Mail, Library, Calendar } from 'lucide-react'
import DecriptionView from '@/components/common/description-view'
import { Separator } from '@/components/ui/separator'
import { getStudentById2 } from '@/lib/api/student'
import useSWR from 'swr'

interface Props {
  id: string
}

const getStudentInfoItems = (data: any) => [
  {
    icon: <School className='h-5 w-5 text-gray-500' />,
    title: 'Trường/Học viện',
    value: `${data?.universityCode} - ${data?.univeristyName}`
  },
  {
    icon: <User className='h-5 w-5 text-gray-500' />,
    title: 'Họ và tên',
    value: data?.name
  },
  {
    icon: <Book className='h-5 w-5 text-gray-500' />,
    title: 'Mã sinh viên',
    value: data?.code
  },
  {
    icon: <Mail className='h-5 w-5 text-gray-500' />,
    title: 'Email',
    value: data?.email
  },
  {
    icon: <Library className='h-5 w-5 text-gray-500' />,
    title: 'Ngành học',
    value: `${data?.facultyCode} - ${data?.facultyName}`
  },
  {
    icon: <Calendar className='h-5 w-5 text-gray-500' />,
    title: 'Năm nhập học',
    value: data?.year
  }
]

const StudentView: React.FC<Props> = (props) => {
  const queryStudentInfo = useSWR(props.id, () => getStudentById2(props.id))

  return (
    <>
      <DecriptionView
        title='Thông tin cá nhân'
        description='Thông tin chi tiết về hồ sơ sinh viên'
        items={getStudentInfoItems(queryStudentInfo.data || {})}
      />
      <Separator className='my-4' />
    </>
  )
}

export default StudentView
