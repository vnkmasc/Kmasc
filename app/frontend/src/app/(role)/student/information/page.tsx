'use client'
import { Badge } from '@/components/ui/badge'
import { Book, ChartAreaIcon, AwardIcon } from 'lucide-react'
import FastView from '@/components/common/fast-view'
import DecriptionView from '@/components/common/description-view'
import useSWR from 'swr'
import { getStudentInformation } from '@/lib/api/student'
import { getStudentInfoItems } from '@/lib/utils/render-ui'

export default function StudentDashboard() {
  const fastViewData = {
    gpa: 8.5,
    numberOfSubjects: 59
  }

  const queryData = useSWR('student-information', () => getStudentInformation())

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
        items={getStudentInfoItems(queryData.data || {})}
      />
    </div>
  )
}
