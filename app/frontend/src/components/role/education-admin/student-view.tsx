'use client'

import { getStudentInfoItems } from '@/app/(role)/student/information/page'
import DecriptionView from '@/components/common/description-view'
import { Separator } from '@/components/ui/separator'
import { getStudentById2 } from '@/lib/api/student'
import useSWR from 'swr'

interface Props {
  id: string
}

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
