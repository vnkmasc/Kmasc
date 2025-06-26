import Back from '@/components/common/back'
import CertificateView from '@/components/common/certificate-view'

interface Props {
  params: Promise<{ slug: string }>
}

const CertificateDetailPage = async ({ params }: Props) => {
  const { slug } = await params

  return (
    <>
      <div className='mb-4 flex items-center gap-2'>
        <Back />
        <h2>Chi tiết thông tin</h2>
      </div>
      <CertificateView id={slug} />
    </>
  )
}

export default CertificateDetailPage
