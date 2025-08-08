import HtmlView from '@/components/common/html-view'
import PageHeader from '@/components/common/page-header'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'

const DigitalCertificateManagementPage = () => {
  return (
    <Tabs defaultValue='certificate'>
      <TabsList>
        <TabsTrigger value='certificate'>Văn bằng số</TabsTrigger>
        <TabsTrigger value='template'>Mẫu bằng số</TabsTrigger>
      </TabsList>
      <TabsContent value='certificate'>
        <PageHeader title='Quản lý văn bằng số' />
      </TabsContent>
      <TabsContent value='template'>
        <PageHeader title='Quản lý mẫu bằng số' />
      </TabsContent>
    </Tabs>
  )
}

export default DigitalCertificateManagementPage
