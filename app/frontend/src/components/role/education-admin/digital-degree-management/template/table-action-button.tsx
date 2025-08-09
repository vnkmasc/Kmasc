'use client'

import HtmlView from '@/components/common/html-view'
import { Button } from '@/components/ui/button'
import { Sheet, SheetContent, SheetHeader, SheetTitle, SheetTrigger } from '@/components/ui/sheet'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { getDegreeTemplateView } from '@/lib/api/degree'
import { formatDegreeTemplateHTML } from '@/lib/utils/format-api'
import { CodeXml, KeyRound, PencilIcon } from 'lucide-react'
import { Dispatch, SetStateAction } from 'react'
import useSWRMutation from 'swr/mutation'

interface Props {
  id: string
  handleSetIdDetail: Dispatch<SetStateAction<string | null | undefined>>
  canEdit: boolean
  canSign: boolean
  handleSign: () => void
}

const TableActionButton: React.FC<Props> = (props) => {
  const queryTemplateView = useSWRMutation(props.id ? `templates/view/${props.id}` : null, async () => {
    const res = await getDegreeTemplateView(props.id)
    return formatDegreeTemplateHTML(res)
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
          <Tabs defaultValue='template' className='mt-4'>
            <TabsList>
              <TabsTrigger value='template'>Mẫu văn bằng</TabsTrigger>
              <TabsTrigger value='code'>Code</TabsTrigger>
            </TabsList>
            <TabsContent value='template'>
              <HtmlView html={queryTemplateView.data} loading={queryTemplateView.isMutating} />
            </TabsContent>
          </Tabs>
        </SheetContent>
      </Sheet>
      <Button size='icon' variant={'secondary'} onClick={() => props.handleSign()} disabled={!props.canSign}>
        <KeyRound />
      </Button>
    </div>
  )
}

export default TableActionButton
