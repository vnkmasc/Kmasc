import CodeView from '@/components/common/code-view'
import HtmlView from '@/components/common/html-view'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { formatExampleTemplateHTML } from '@/lib/utils/format-api'

interface Props {
  baseHtml: string | undefined
  htmlLoading?: boolean
}

const TemplateView: React.FC<Props> = (props) => {
  return (
    <Tabs defaultValue='base-template'>
      <TabsList>
        <TabsTrigger value='base-template'>{'Mẫu văn bằng (Gốc)'}</TabsTrigger>
        <TabsTrigger value='example-template'>{'Mẫu văn bằng (Ví dụ)'}</TabsTrigger>
        <TabsTrigger value='code'>Code</TabsTrigger>
      </TabsList>
      <TabsContent value='base-template'>
        <HtmlView html={props.baseHtml} loading={props.htmlLoading} />
      </TabsContent>
      <TabsContent value='example-template'>
        <HtmlView html={formatExampleTemplateHTML(props.baseHtml || '')} loading={props.htmlLoading} />
      </TabsContent>
      <TabsContent value='code'>
        <CodeView code={props.baseHtml || ''} />
      </TabsContent>
    </Tabs>
  )
}

export default TemplateView
