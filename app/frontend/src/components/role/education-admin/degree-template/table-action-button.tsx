'use client'

import HtmlView from '@/components/common/html-view'
import { Button } from '@/components/ui/button'
import { Sheet, SheetContent, SheetHeader, SheetTitle, SheetTrigger } from '@/components/ui/sheet'
import { getDegreeTemplateView } from '@/lib/api/degree'
import { CodeXml, PencilIcon } from 'lucide-react'
import { Dispatch, SetStateAction } from 'react'
import useSWR from 'swr'

interface Props {
  id: string
  handleSetIdDetail: Dispatch<SetStateAction<string | null | undefined>>
  canEdit: boolean
}

const TableActionButton: React.FC<Props> = (props) => {
  const queryTemplateView = useSWR(props.id ? `templates/view/${props.id}` : null, () =>
    getDegreeTemplateView(props.id)
  )

  return (
    <div className='flex gap-2'>
      <Button variant='outline' size='icon' onClick={() => props.handleSetIdDetail(props.id)} disabled={!props.canEdit}>
        <PencilIcon />
      </Button>
      <Sheet>
        <SheetTrigger asChild>
          <Button size='icon'>
            <CodeXml />
          </Button>
        </SheetTrigger>
        <SheetContent className='min-w-full max-w-full md:min-w-[900px]'>
          <SheetHeader>
            <SheetTitle>Mẫu văn bằng</SheetTitle>
          </SheetHeader>
          <HtmlView html={queryTemplateView.data?.data} loading={queryTemplateView.isLoading} />
        </SheetContent>
      </Sheet>
    </div>
  )
}

export default TableActionButton
