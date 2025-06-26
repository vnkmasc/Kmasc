'use client'
import { Badge } from '@/components/ui/badge'
import { Book, Mail, Calendar, ChartAreaIcon, AwardIcon, User, School, Library } from 'lucide-react'
import FastView from '@/components/common/fast-view'
import DecriptionView from '@/components/common/description-view'
import useSWR from 'swr'
import { getStudentInformation } from '@/lib/api/student'

export default function StudentDashboard() {
  const fastViewData = {
    gpa: 8.5,
    numberOfSubjects: 59
  }

  const queryData = useSWR('student-information', () => getStudentInformation())

  const personalInfoItems = [
    {
      icon: <School className='h-5 w-5 text-gray-500' />,
      title: 'Trường/Học viện',
      value: `${queryData.data?.universityCode} - ${queryData.data?.univeristyName}`
    },
    {
      icon: <User className='h-5 w-5 text-gray-500' />,
      title: 'Họ và tên',
      value: queryData.data?.name
    },
    {
      icon: <Book className='h-5 w-5 text-gray-500' />,
      title: 'Mã sinh viên',
      value: queryData.data?.code
    },
    {
      icon: <Mail className='h-5 w-5 text-gray-500' />,
      title: 'Email',
      value: queryData.data?.email
    },
    {
      icon: <Library className='h-5 w-5 text-gray-500' />,
      title: 'Ngành học',
      value: `${queryData.data?.facultyCode} - ${queryData.data?.facultyName}`
    },
    {
      icon: <Calendar className='h-5 w-5 text-gray-500' />,
      title: 'Năm nhập học',
      value: queryData.data?.year
    }
  ]

  return (
    <div>
      {/* Statistics Section */}
      <h2>Thông tin cá nhân</h2>

      <div className='my-4 flex flex-col gap-4 sm:flex-row'>
        <FastView
          title='GPA'
          value={fastViewData.gpa}
          icon={<ChartAreaIcon className='text-blue-500' />}
          color='text-blue-500'
        />
        <FastView
          title='Số môn đã học'
          value={fastViewData.numberOfSubjects}
          icon={<Book className='text-green-500' />}
          color='text-green-500'
        />
        <FastView
          title='Trạng thái'
          value={
            <Badge variant={queryData.data?.status === 'true' ? 'default' : 'outline'}>
              {queryData.data?.status === 'true' ? 'Đã tốt nghiệp' : 'Đang học'}
            </Badge>
          }
          icon={<AwardIcon />}
          color='text-green-500'
        />
      </div>

      <DecriptionView
        title='Thông tin cá nhân'
        description='Thông tin chi tiết về hồ sơ sinh viên'
        items={personalInfoItems}
      />
    </div>
  )
}
