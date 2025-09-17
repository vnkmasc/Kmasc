'use client'

import DigitalDegreeView from '@/components/common/digital-degree-view'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { getDigitalDegreesByStudent } from '@/lib/api/digital-degree'
import { clearFalsyValueObject } from '@/lib/utils/common'
import { Suspense, useState } from 'react'
import useSWR from 'swr'
import { useRouter, useSearchParams } from 'next/navigation'
import SuspendPage from '@/components/common/suspend-page'

const StudentDigitalDegreeDetail = () => {
  const [tab, setTab] = useState<string | undefined>(undefined)
  const history = useRouter()
  const searchParams = useSearchParams()
  const queryListName = useSWR('list-name-digital-degree', getDigitalDegreesByStudent, {
    onSuccess: (data) => {
      const tab = searchParams.get('tab')
      setTab(tab ?? data?.[0]?.id)
    }
  })

  const getPropsDigitalDegreeView = (item: any) => {
    return {
      id: item.id,
      isBlockchain: true,
      ...clearFalsyValueObject({
        facultyId: item.on_blockchain_verify?.faculty_id,
        universityId: item.on_blockchain_verify?.university_id,
        certificateType: item.on_blockchain_verify?.certificate_type,
        course: item.on_blockchain_verify?.course,
        universityCode: item.on_blockchain_verify?.university_code ?? 'KMA'
      })
    }
  }

  return (
    <>
      <h2>Danh sách văn bằng số</h2>
      <Tabs
        className='mt-4'
        value={tab}
        onValueChange={(value) => {
          setTab(value)
          history.replace(`/student/digital-degree?tab=${value}`)
        }}
      >
        <TabsList>
          {queryListName.data?.map(
            (item: any) =>
              item.on_blockchain_verify && (
                <TabsTrigger key={item.id} value={item.id}>
                  {item.name}
                </TabsTrigger>
              )
          )}
        </TabsList>
        {queryListName.data?.map((item: any) => {
          const props = getPropsDigitalDegreeView(item)

          return (
            <TabsContent key={item.id} value={item.id}>
              <DigitalDegreeView {...props} />
            </TabsContent>
          )
        })}
      </Tabs>
    </>
  )
}

const StudentDigitalDegreeDetailPage = () => (
  <Suspense fallback={<SuspendPage />}>
    <StudentDigitalDegreeDetail />
  </Suspense>
)

export default StudentDigitalDegreeDetailPage
