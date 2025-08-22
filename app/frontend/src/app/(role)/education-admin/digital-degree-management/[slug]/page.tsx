import Back from '@/components/common/back'
import DigitalDegreeView from '@/components/common/digital-degree-view'

interface Props {
  params: Promise<{ slug: string }>
}

const DigitalDegreeDetailPage = async ({ params }: Props) => {
  const { slug } = await params

  return (
    <>
      <div className='mb-4 flex items-center gap-2'>
        <Back />
        <h2>Chi tiết bằng số trên cơ sở dữ liệu</h2>
      </div>
      <DigitalDegreeView id={slug} isBlockchain={false} />
    </>
  )
}

export default DigitalDegreeDetailPage
