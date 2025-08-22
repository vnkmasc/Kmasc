'use client'

import { Blocks, EyeIcon, Grid2X2Plus } from 'lucide-react'
import { Button } from '../../../ui/button'
import Link from 'next/link'
import useSWRMutation from 'swr/mutation'
import { pushCertificateIntoBlockchain } from '@/lib/api/certificate'
import { showNotification } from '@/lib/utils/common'
import CertificateQrCode from '../../../common/certificate-qr-code'

interface Props {
  id: string
  onBlockchain: boolean
}

const CertificateActionButton: React.FC<Props> = (props) => {
  const mutatePushCertificateIntoBlockchain = useSWRMutation(
    'push-blockchain',
    (_, { arg }: { arg: string }) => pushCertificateIntoBlockchain(arg),
    {
      onSuccess: () => {
        showNotification('success', 'Đẩy dữ liệu lên blockchain thành công')
      },
      onError: (error) => {
        showNotification('error', error.message || 'Lỗi khi đẩy dữ liệu lên blockchain')
      }
    }
  )

  return (
    <div className='flex gap-2'>
      <Link href={`/education-admin/certificate-management/${props.id}`} target='_blank'>
        <Button size={'icon'} variant={'outline'} title='Xem dữ liệu trên cơ sở dữ liệu'>
          <EyeIcon />
        </Button>
      </Link>
      <Link href={`/education-admin/certificate-management/${props.id}/blockchain`} target='_blank'>
        <Button size={'icon'} variant={'secondary'} title='Xem dữ liệu trên blockchain'>
          <Blocks />
        </Button>
      </Link>
      <CertificateQrCode id={props.id} isIcon={true} />

      <Button
        disabled={props.onBlockchain}
        size={'icon'}
        title='Đẩy dữ liệu lên blockchain'
        onClick={() => mutatePushCertificateIntoBlockchain.trigger(props.id)}
        isLoading={mutatePushCertificateIntoBlockchain.isMutating}
      >
        <Grid2X2Plus />
      </Button>
    </div>
  )
}

export default CertificateActionButton
