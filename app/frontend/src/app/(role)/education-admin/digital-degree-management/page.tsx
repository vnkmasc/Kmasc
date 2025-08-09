import DegreeTemplate from '@/components/role/education-admin/digital-degree-management/template/degree-template'
import DigitalDegreeView from '@/components/role/education-admin/digital-degree-management/degree/digital-degree-view'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'

const DigitalDegreeManagementPage = () => {
  return (
    <Tabs defaultValue='degree'>
      <TabsList>
        <TabsTrigger value='degree'>Văn bằng số</TabsTrigger>
        <TabsTrigger value='template'>Mẫu bằng số</TabsTrigger>
      </TabsList>
      <TabsContent value='degree'>
        <DigitalDegreeView />
      </TabsContent>
      <TabsContent value='template'>
        <DegreeTemplate />
      </TabsContent>
    </Tabs>
  )
}

export default DigitalDegreeManagementPage
