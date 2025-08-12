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
import { Sheet, SheetContent, SheetHeader, SheetTitle, SheetTrigger } from '@/components/ui/sheet'
import { deleteDegreeTemplate, getDegreeTemplateView, signDegreeTemplateById } from '@/lib/api/degree'
import { showNotification } from '@/lib/utils/common'
import { CodeXml, KeyRound, PencilIcon, TrashIcon } from 'lucide-react'
import { Dispatch, SetStateAction } from 'react'
import useSWRMutation from 'swr/mutation'
import TemplateView from '../template-view'

interface Props {
  id: string
  handleSetIdDetail: Dispatch<SetStateAction<string | null | undefined>>
  canEdit: boolean
  canSign: boolean
  refetch: () => void
}

const TableActionButton: React.FC<Props> = (props) => {
  const queryTemplateView = useSWRMutation(props.id ? `templates/view/${props.id}` : null, async () => {
    const res = await getDegreeTemplateView(props.id)
    return res
  })

  const mutateSignDegreeTemplateById = useSWRMutation(
    'sign-degree-template-by-id',
    () => signDegreeTemplateById(props.id),
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
      <Sheet>
        <SheetTrigger asChild>
          <Button size='icon' onClick={() => queryTemplateView.trigger()}>
            <CodeXml />
          </Button>
        </SheetTrigger>
        <SheetContent className='min-w-full max-w-full md:min-w-[900px]'>
          <SheetHeader>
            <SheetTitle>Mẫu văn bằng số</SheetTitle>
          </SheetHeader>
          <TemplateView baseHtml={queryTemplateView.data} htmlLoading={queryTemplateView.isMutating} />
        </SheetContent>
      </Sheet>
      <Button
        size='icon'
        variant={'secondary'}
        onClick={() => mutateSignDegreeTemplateById.trigger()}
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
