import CertificateView from '@/components/common/certificate-view'

interface Props {
  params: Promise<{ slug: string }>
}

const CertificateDetailPage = async ({ params }: Props) => {
  const { slug } = await params

  return (
    <>
      <h2 className='mb-4'>Chi tiết thông tin</h2>
      <CertificateView id={slug} />
    </>
  )
}

export default CertificateDetailPage
