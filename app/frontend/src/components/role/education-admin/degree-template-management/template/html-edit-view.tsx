import HtmlView from '@/components/common/html-view'
import { Label } from '@/components/ui/label'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { ReactNode } from 'react'

interface Props {
  textarea: ReactNode
  html: string
}

const HtmlEditView: React.FC<Props> = (props) => {
  return (
    <Tabs defaultValue='code'>
      <TabsList>
        <TabsTrigger value='code'>Code</TabsTrigger>
        <TabsTrigger value='preview'>Xem trước</TabsTrigger>
      </TabsList>
      <TabsContent value='code'>{props.textarea}</TabsContent>
      <TabsContent value='preview'>
        <Label>Mẫu bằng số</Label>
        <HtmlView html={props.html} />
      </TabsContent>
    </Tabs>
  )
}

export default HtmlEditView
