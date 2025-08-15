import Back from '@/components/common/back'
import StudentView from '@/components/role/education-admin/student-management/student-view'

interface Props {
  params: Promise<{ slug: string }>
}

const StudentDetailPage = async ({ params }: Props) => {
  const { slug } = await params

  return (
    <>
      <div className='mb-4 flex items-center gap-2'>
        <Back />
        <h2>Chi tiết thông tin sinh viên</h2>
      </div>
      <StudentView id={slug} />
    </>
  )
}

export default StudentDetailPage
