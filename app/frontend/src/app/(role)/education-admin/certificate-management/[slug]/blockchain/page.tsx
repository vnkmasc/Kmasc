import Back from '@/components/common/back'
import CertificateView from '@/components/common/certificate-view'
import { decodeJSON } from '@/lib/utils/lz-string'

interface Props {
  params: Promise<{ slug: string }>
}

const CertificateBlockchainDetailPage = async ({ params }: Props) => {
  const { slug } = await params
  const decodeCertificateData = decodeJSON(slug)

  return (
    <>
      <div className='mb-4 flex items-center gap-2'>
        <Back />
        <h2>Chi tiết thông tin trên blockchain</h2>
      </div>
      <CertificateView
        id={decodeCertificateData.certificate_id}
        certificateType={decodeCertificateData.certificate_type}
        course={decodeCertificateData.course}
        facultyId={decodeCertificateData.faculty_id}
        universityId={decodeCertificateData.university_id}
        universityCode={decodeCertificateData.university_code}
        isBlockchain={true}
      />
    </>
  )
}

export default CertificateBlockchainDetailPage
