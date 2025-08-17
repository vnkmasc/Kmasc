'use client'

import {
  AlertDialog,
  AlertDialogDescription,
  AlertDialogContent,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
  AlertDialogFooter,
  AlertDialogCancel,
  AlertDialogAction
} from '@/components/ui/alert-dialog'
import { Button } from '@/components/ui/button'
import { deleteDegreeTemplate, signDegreeTemplateById } from '@/lib/api/digital-degree'
import { showMessage, showNotification } from '@/lib/utils/common'
import { CodeXml, KeyRound, PencilIcon, TrashIcon } from 'lucide-react'
import { Dispatch, SetStateAction } from 'react'
import useSWRMutation from 'swr/mutation'
import Link from 'next/link'
import { getSignDegreeConfig } from '@/lib/utils/handle-storage'
import { signDigitalSignature } from '@/lib/utils/handle-vgca'

interface Props {
  id: string
  handleSetIdDetail: Dispatch<SetStateAction<string | null | undefined>>
  canEdit: boolean
  canSign: boolean
  refetch: () => void
  hashTemplate: string
  templateSampleId: string
}

const TableActionButton: React.FC<Props> = (props) => {
  const signDegreeConfig = getSignDegreeConfig()
  const mutateSignDegreeTemplateById = useSWRMutation(
    'sign-degree-template-by-id',
    (_, { arg }: { arg: string }) => signDegreeTemplateById(props.id, arg),
    {
      onSuccess: () => {
        showNotification('success', 'Ký mẫu bằng số thành công')
        props.refetch()
      },
      onError: (error) => {
        showNotification('error', error.message || 'Lỗi khi ký mẫu bằng số')
      }
    }
  )

  const mutateDeleteDegreeTemplate = useSWRMutation('delete-degree-template', () => deleteDegreeTemplate(props.id), {
    onSuccess: () => {
      showNotification('success', 'Xóa mẫu bằng số thành công')
      props.refetch()
    },
    onError: (error) => {
      showNotification('error', error.message || 'Lỗi khi xóa mẫu bằng số')
    }
  })

  return (
    <div className='flex gap-2'>
      <Button variant='outline' size='icon' onClick={() => props.handleSetIdDetail(props.id)} disabled={!props.canEdit}>
        <PencilIcon />
      </Button>
      <Link
        href={`/education-admin/digital-degree-management?tab=template-interface&id=${props.templateSampleId}`}
        target='_blank'
      >
        <Button size='icon'>
          <CodeXml />
        </Button>
      </Link>
      <Button
        size='icon'
        variant={'secondary'}
        onClick={async () => {
          if (signDegreeConfig?.signService === '') {
            showMessage('Vui lòng cấu hình số cho link server ký số')
            return
          }
          // *@*
          const signature = await signDigitalSignature(props.hashTemplate)

          if (!signature) {
            showMessage('Ký số không thành công')
            return
          }

          mutateSignDegreeTemplateById.trigger(signature)
        }}
        disabled={!props.canSign}
        isLoading={mutateSignDegreeTemplateById.isMutating}
      >
        <KeyRound />
      </Button>
      <AlertDialog>
        <AlertDialogTrigger asChild>
          <Button variant='destructive' size='icon'>
            <TrashIcon />
          </Button>
        </AlertDialogTrigger>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Xóa mẫu bằng số</AlertDialogTitle>
            <AlertDialogDescription>
              Mẫu bằng số có ID <b>{props.id}</b> sẽ bị xóa khỏi hệ thống.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Hủy bỏ</AlertDialogCancel>
            <AlertDialogAction onClick={() => mutateDeleteDegreeTemplate.trigger()}>Xóa</AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  )
}

export default TableActionButton
