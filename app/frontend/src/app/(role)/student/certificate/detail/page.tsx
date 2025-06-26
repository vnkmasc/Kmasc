'use client'

import CertificateView from '@/components/common/certificate-view'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { getCertificatesNameByStudent } from '@/lib/api/certificate'
import useSWR from 'swr'

const StudentCertificateDetailPage = () => {
  const queryListName = useSWR('list-name', getCertificatesNameByStudent)

  return (
    <>
      <h2>Chi tiết các văn bằng</h2>
      <Tabs className='mt-4'>
        <TabsList>
          {queryListName.data?.map((item: any) => (
            <TabsTrigger key={item.id} value={item.id}>
              {item.certificate_name}
            </TabsTrigger>
          ))}
        </TabsList>
        {queryListName.data?.map((item: any) => (
          <TabsContent key={item.id} value={item.id}>
            <CertificateView id={item.id} />
          </TabsContent>
        ))}
      </Tabs>
    </>
  )
}

export default StudentCertificateDetailPage
